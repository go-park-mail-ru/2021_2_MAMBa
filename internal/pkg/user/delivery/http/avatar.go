package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/filesaver"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

const (
	root = "."
	path = "/static/media/img/users/"
)

func (handler *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	rq :=  cast.CookieToRq(r,0)
	clientIDMessage, err := handler.AuthClient.CheckSession(r.Context(),&rq)
	if err != nil || err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	clientID := clientIDMessage.ID

	err = r.ParseMultipartForm(10 * 1024 * 1024) // лимит 10МБ
	if err != nil {
		log.Warn(fmt.Sprintf("parse multipart form error: %s", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Warn(fmt.Sprintf("error while parsing file: %s", err))
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
	}
	defer uploaded.Close()

	filename, err := filesaver.UploadFile(uploaded, root, path, filepath.Ext(header.Filename))
	if err != nil {
		log.Error(err)
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
	}

	us, err := handler.UserUsecase.UpdateAvatar(clientID, filename)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
