package http

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
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
	clientID, err := sessions.CheckSession(r)
	if err != nil || err == customErrors.ErrorUserNotLoggedIn {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 * 1024 * 1024) // лимит 10МБ
	if err != nil {
		log.Warn(fmt.Sprintf("parse multipart form error: %s", err))
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Warn(fmt.Sprintf("error while parsing file: %s", err))
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
	}
	defer uploaded.Close()

	filename, err := filesaver.UploadFile(uploaded, root, path, filepath.Ext(header.Filename))
	if err != nil {
		log.Error(err)
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
	}

	us, err := handler.UserUsecase.UpdateAvatar(clientID, filename)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}
