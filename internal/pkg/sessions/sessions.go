package sessions

import (
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

var ErrUserNotLoggedIn = errors.New("user not logged in")

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// MaxAge=0 means no Max-Age attribute specified and the cookie will be
// deleted after the browser session ends.
// MaxAge<0 means delete cookie immediately.
// MaxAge>0 means Max-Age attribute present and given in seconds.

func StartSession(w http.ResponseWriter, r *http.Request, id uint) error {
	session, _ := store.Get(r, "session-name")
	session.Values["id"] = id
	session.Options = &sessions.Options{MaxAge: 100000}
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func EndSession(w http.ResponseWriter, r *http.Request, id uint) error {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return err
	}
	// Get() always returns a session, even if empty, so check isIn
	sessionId, isIn := session.Values["id"]
	if isIn && id == sessionId {
		// deleting a session may only happen at maxage < 0
		session.Options = &sessions.Options{MaxAge: -1}
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

func CheckSession(r *http.Request) (uint, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return 0, err
	}
	id, isIn := session.Values["id"]
	if !isIn {
		return 0, ErrUserNotLoggedIn
	}
	return id.(uint), nil
}
