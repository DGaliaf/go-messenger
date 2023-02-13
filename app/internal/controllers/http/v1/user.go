package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/user"
	"messenger-rest-api/app/internal/errors"
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
	router.HandlerFunc(http.MethodPost, createUser, custom_error.Middleware(u.createUser))
}

// Create User godoc
// @Summary      Create User
// @Description  create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        username   body  dto.CreateUserDTO  true  "Username"
// @Success      201  {object}  dto.ShowUserIdDTO
// @Failure      418
// @Failure      400
// @Failure      500
// @Router       /users/add/ [post]
func (u UserHandler) createUser(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	d := new(dto.CreateUserDTO)
	if err := json.Unmarshal(body, d); err != nil {
		return err
	}

	userId, err := u.service.Create(r.Context(), user.CreateUserDTO(*d))
	if err != nil {
		return err
	}

	response := dto.ShowUserIdDTO{
		UserID: userId,
	}

	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaledResponse)
	return nil
}
