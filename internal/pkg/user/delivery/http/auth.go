package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
)

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.User)
	var p []byte
	p, err := ioutil.ReadAll(r.Body)
	err = userForm.UnmarshalJSON(p)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	xss.SanitizeUser(userForm)

	us, err := handler.UserUsecase.Register(userForm)
	if err == customErrors.ErrorAlreadyExists {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusConflict}
		resp.Write(w)
		return
	}
	if err == customErrors.ErrorBadInput {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	rq := cast.CookieToRq(r, us.ID)
	result, err := handler.AuthClient.StartSession(r.Context(), &rq)
	http.SetCookie(w, sessions.NewCookie(result.Name, result.Value, &sessions.Options{
		Path:     result.Path,
		Domain:   rq.Domain,
		MaxAge:   int(result.MaxAge),
		Secure:   result.Secure,
		HttpOnly: result.HttpOnly,
		SameSite: http.SameSite(result.SameSite),
	}))
	if err != nil && us.ID != 0 {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	x, err := us.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusCreated,
	}
	resp.Write(w)
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.UserToLogin)
	p, err := ioutil.ReadAll(r.Body)
	err = userForm.UnmarshalJSON(p)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	rq := cast.CookieToRq(r, 0)
	_, err = handler.AuthClient.CheckSession(r.Context(), &rq)
	st, _ := status.FromError(err)
	s2, _ := status.FromError(customErrors.ErrorUserNotLoggedIn)
	if st.Message() != s2.Message() {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorAlreadyIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}

	us, err := handler.UserUsecase.Login(userForm)
	if err == customErrors.ErrorBadInput {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err == customErrors.ErrorBadCredentials {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	rq.ID = us.ID
	result, err := handler.AuthClient.StartSession(r.Context(), &rq)
	http.SetCookie(w, sessions.NewCookie(result.Name, result.Value, &sessions.Options{
		Path:     result.Path,
		Domain:   rq.Domain,
		MaxAge:   int(result.MaxAge),
		Secure:   result.Secure,
		HttpOnly: result.HttpOnly,
		SameSite: http.SameSite(result.SameSite),
	}))
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
	}

	x, err := us.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	rq := cast.CookieToRq(r, 0)
	idMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusForbidden}
		resp.Write(w)
		return
	}
	rq.ID = idMessage.ID
	result, err := handler.AuthClient.EndSession(r.Context(), &rq)
	http.SetCookie(w, sessions.NewCookie(result.Name, result.Value, &sessions.Options{
		Path:     result.Path,
		Domain:   "/",
		MaxAge:   int(result.MaxAge),
		Secure:   result.Secure,
		HttpOnly: result.HttpOnly,
		SameSite: http.SameSite(result.SameSite),
	}))
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	rq := cast.CookieToRq(r, 0)
	idMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusForbidden}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	userID := idMessage.ID

	us, err := handler.UserUsecase.CheckAuth(userID)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	x, err := us.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
