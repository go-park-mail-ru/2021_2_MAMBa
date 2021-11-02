package http

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (handler *UserHandler) GetBasicInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	u64, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	us, err := handler.UserUsecase.GetBasicInfo(u64)
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
