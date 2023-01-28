package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/chat"
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
	router.HandlerFunc(http.MethodPost, createChat, c.createChat)
	router.HandlerFunc(http.MethodPost, findUserChats, c.findUserChats)
}

func (c ChatHandler) createChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Accept", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	d := new(dto.CreateChatDTO)
	if err := json.Unmarshal(body, d); err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	chatID, err := c.service.Create(r.Context(), chat.CreateChatDTO(*d))
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	response := dto.ShowChatIdDTO{
		ID: chatID,
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

func (c ChatHandler) findUserChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Accept", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	d := new(dto.GetUserDTO)
	if err := json.Unmarshal(body, d); err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	chats, err := c.service.FindUserChats(r.Context(), d.UserID)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	marshal, err := json.Marshal(chats)
	if err != nil {
		// TODO: Error handling

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}
