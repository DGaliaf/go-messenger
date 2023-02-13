package v1

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	d "messenger-rest-api/app/internal/controllers/http/dto"
	"messenger-rest-api/app/internal/domain/service/message"
	"messenger-rest-api/app/internal/errors"
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
	router.HandlerFunc(http.MethodPost, getMessages, custom_error.Middleware(m.getMessageFromChatByID))
	router.HandlerFunc(http.MethodPost, sendMessage, custom_error.Middleware(m.sendMessage))
}

// Get Message godoc
// @Summary      Get Message
// @Description  get message from chat
// @Tags         message
// @Accept       json
// @Produce      json
// @Param        chat_id   body  dto.GetMessagesFromChatDTO  true  "Chat ID"
// @Success      200  {array}  entities.Message
// @Failure      418
// @Failure      400
// @Failure      500
// @Router       /messages/get [post]
func (m MessageHandler) getMessageFromChatByID(w http.ResponseWriter, r *http.Request) error {
	var dto d.GetMessagesFromChatDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &dto); err != nil {
		return err
	}

	messages, err := m.service.GetMessagesFromChatByID(r.Context(), message.GetMessagesFromChatDTO(dto))
	if err != nil {
		return err
	}

	response, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return nil
}

// Send Message godoc
// @Summary      Send Message
// @Description  send message to chat
// @Tags         message
// @Accept       json
// @Produce      json
// @Param        message   body  dto.CreateMessageDTO  true  "Message"
// @Success      201  {object}  dto.ShowMessageIdDTO
// @Failure      418
// @Failure      400
// @Failure      500
// @Router       /messages/add [post]
func (m MessageHandler) sendMessage(w http.ResponseWriter, r *http.Request) error {
	var dto d.CreateMessageDTO

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &dto); err != nil {
		return err
	}

	id, err := m.service.SendMessage(r.Context(), message.CreateMessageDTO(dto))
	if err != nil {
		return err
	}

	response := d.ShowMessageIdDTO{
		ID: id,
	}

	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaledResponse)
	return nil
}
