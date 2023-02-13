package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/chat"
	"messenger-rest-api/app/internal/errors"
	"net/http"
)

var (
	createChat    = "/chats/add/"
	findUserChats = "/chats/get/"
)

type ChatHandler struct {
	service *chat.Service
}

func NewChatHandler(service *chat.Service) *ChatHandler {
	return &ChatHandler{service: service}
}

func (c ChatHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createChat, custom_error.Middleware(c.createChat))
	router.HandlerFunc(http.MethodPost, findUserChats, custom_error.Middleware(c.findUserChats))
}

// Create Chat godoc
// @Summary      Create Chat
// @Description  create chat with specific name and users
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        chat   body  dto.CreateChatDTO  true  "Chat"
// @Success      201  {object}  dto.ShowChatIdDTO
// @Failure      418
// @Failure      400
// @Failure      500
// @Router       /chats/add [post]
func (c ChatHandler) createChat(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	d := new(dto.CreateChatDTO)
	if err := json.Unmarshal(body, d); err != nil {
		return err
	}

	chatID, err := c.service.Create(r.Context(), chat.CreateChatDTO(*d))
	if err != nil {
		return err
	}

	response := dto.ShowChatIdDTO{
		ID: chatID,
	}

	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaledResponse)
	return nil
}

// Find User Chats godoc
// @Summary      Find User Chats
// @Description  find all chats where our user is participated
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        user_id   body  dto.GetUserDTO  true  "User ID"
// @Success      200  {array}  entities.Chat
// @Failure      418
// @Failure      400
// @Failure      500
// @Router       /chats/get/ [post]
func (c ChatHandler) findUserChats(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	d := new(dto.GetUserDTO)
	if err := json.Unmarshal(body, d); err != nil {
		return err
	}

	chats, err := c.service.FindUserChats(r.Context(), d.UserID)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(chats)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
	return nil
}
