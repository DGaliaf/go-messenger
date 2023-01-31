package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	d "messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/message"
	"net/http"
)

var (
	sendMessage = "/messages/add"
	getMessages = "/messages/get"
)

type MessageHandler struct {
	service *message.Service
}

func NewMessageHandler(service *message.Service) *MessageHandler {
	return &MessageHandler{service: service}
}

func (m MessageHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, getMessages, m.getMessageFromChatByID)
	router.HandlerFunc(http.MethodPost, sendMessage, m.sendMessage)
}

func (m MessageHandler) getMessageFromChatByID(w http.ResponseWriter, r *http.Request) {
	var dto d.GetMessagesFromChatDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	if err := json.Unmarshal(body, &dto); err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	messages, err := m.service.GetMessagesFromChatByID(r.Context(), message.GetMessagesFromChatDTO(dto))
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	response, err := json.Marshal(messages)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (m MessageHandler) sendMessage(w http.ResponseWriter, r *http.Request) {
	var dto d.CreateMessageDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	if err := json.Unmarshal(body, &dto); err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	id, err := m.service.SendMessage(r.Context(), message.CreateMessageDTO(dto))
	if err != nil {
		// TODO: Error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}

	response := d.ShowMessageIdDTO{
		ID: id,
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
