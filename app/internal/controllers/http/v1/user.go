package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/user"
	"net/http"
)

var (
	createUser = "/users/add"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (u UserHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createUser, u.createUser)
}

func (u UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Accept", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	d := new(dto.CreateUserDTO)
	if err := json.Unmarshal(body, d); err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	userId, err := u.service.Create(r.Context(), user.CreateUserDTO(*d))
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	response := dto.ShowUserIdDTO{
		UserID: userId,
	}

	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaledResponse)
}
