package http
import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type PersonHandler struct {
	PersonUsecase domain.PersonUsecase
}

func NewHandlers(router *mux.Router, uc domain.PersonUsecase) {
	handler := &PersonHandler{
		PersonUsecase: uc,
	}

	router.HandleFunc("/person/getPerson", handler.GetPerson).Methods("GET", "OPTIONS")
	router.HandleFunc("/person/getPersonFilms", handler.GetPersonFilms).Methods("GET", "OPTIONS")

}