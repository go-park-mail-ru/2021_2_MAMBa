package sessions

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	sGrpc "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"net/http/httptest"
)

var store = sessions.NewFilesystemStore("", securecookie.GenerateRandomKey(32))
var sessionName = "session-name"

// MaxAge=0 means no Max-Age attribute specified and the cookie will be
// deleted after the browser session ends.
// MaxAge<0 means delete cookie immediately.
// MaxAge>0 means Max-Age attribute present and given in seconds.

type SessionManager struct{
	secure bool
}

func NewSessionManager(secure bool) sGrpc.SessionRPCServer {
	return &SessionManager{secure: secure}
}

func (sm *SessionManager) StartSession(ctx context.Context, rq *sGrpc.Request) (*sGrpc.Session, error) {
	rd := new(io.Reader)
	r, err := http.NewRequest("", "", *rd)
	r.AddCookie(&http.Cookie{
		Name:     rq.Name,
		Value:    rq.Value,
		Path:     rq.Path,
		Domain:   rq.Domain,
		MaxAge:   int(rq.MaxAge),
		Secure:   rq.Secure,
		HttpOnly: rq.HttpOnly,
		SameSite: http.SameSite(rq.SameSite),
		Raw:      rq.Raw,
		Unparsed: rq.Unparsed,
	})
	session, _ := store.Get(r, sessionName)
	session.Values["id"] = rq.ID
	session.Options = &sessions.Options{
		MaxAge:   100000, // ~27 hours
		Secure:   sm.secure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	w := httptest.NewRecorder()
	err = session.Save(r, w)
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		store.Codecs...)
	if err != nil {
		return &sGrpc.Session{}, err
	}
	return &sGrpc.Session{
		Name:     sessionName,
		Path:     "/",
		MaxAge:   int64(session.Options.MaxAge),
		Secure:   session.Options.Secure,
		HttpOnly: session.Options.HttpOnly,
		SameSite: int64(session.Options.SameSite),
		Value:    encoded,
	}, nil
}

func (sm *SessionManager) EndSession(ctx context.Context, rq *sGrpc.Request) (*sGrpc.Session, error) {
	rd := new(io.Reader)
	r, err := http.NewRequest("", "", *rd)
	r.WithContext(ctx)
	r.AddCookie(&http.Cookie{
		Name:     rq.Name,
		Value:    rq.Value,
		Path:     rq.Path,
		Domain:   rq.Domain,
		MaxAge:   int(rq.MaxAge),
		Secure:   rq.Secure,
		HttpOnly: rq.HttpOnly,
		SameSite: http.SameSite(rq.SameSite),
		Raw:      rq.Raw,
		Unparsed: rq.Unparsed,
	})
	session, err := store.Get(r, sessionName)
	if err != nil {
		return &sGrpc.Session{}, err
	}
	// Get() always returns a session, even if empty, so check isIn
	sessionId, isIn := session.Values["id"]
	if isIn && rq.ID == sessionId {
		// deleting a session may only happen at maxage < 0
		session.Options.MaxAge = -1
		w := httptest.NewRecorder()
		err = session.Save(r, w)
		if err != nil {
			return &sGrpc.Session{}, err
		}
	}
	return &sGrpc.Session{
		Name:     session.Name(),
		Path:     session.Options.Path,
		MaxAge:   int64(session.Options.MaxAge),
		Secure:   session.Options.Secure,
		HttpOnly: session.Options.HttpOnly,
		SameSite: int64(session.Options.SameSite),
		Value:    "",
	}, nil
}

func (sm *SessionManager) CheckSession(ctx context.Context, rq *sGrpc.Request) (*sGrpc.ID, error) {
	rd := new(io.Reader)
	r, err := http.NewRequest("", "", *rd)
	r.WithContext(ctx)
	r.AddCookie(&http.Cookie{
		Name:     rq.Name,
		Value:    rq.Value,
		Path:     rq.Path,
		Domain:   rq.Domain,
		MaxAge:   int(rq.MaxAge),
		Secure:   rq.Secure,
		HttpOnly: rq.HttpOnly,
		SameSite: http.SameSite(rq.SameSite),
		Raw:      rq.Raw,
		Unparsed: rq.Unparsed,
	})
	session, err := store.Get(r, sessionName)
	if err != nil && !session.IsNew {
		return &sGrpc.ID{ID: 0}, err
	}
	id, isIn := session.Values["id"]
	if !isIn || session.IsNew {
		return &sGrpc.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn
	}
	idCasted, ok := id.(uint64)
	if !ok {
		return &sGrpc.ID{ID: 0}, customErrors.ErrorUint64Cast
	}
	return &sGrpc.ID{ID: idCasted}, nil
}
