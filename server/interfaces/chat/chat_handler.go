package chat

import (
	"encoding/json"
	"net/http"
	"pinterest/application"
	"pinterest/domain/entity"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ChatInfo struct {
	chatApp application.ChatAppInterface
	userApp application.UserAppInterface
	logger  *zap.Logger
}

func NewChatnfo(chatApp application.ChatAppInterface, userApp application.UserAppInterface,
	logger *zap.Logger) *ChatInfo {
	return &ChatInfo{
		chatApp: chatApp,
		userApp: userApp,
		logger:  logger,
	}
}

func (chatInfo *ChatInfo) HandleAddMessage(w http.ResponseWriter, r *http.Request) {
	otherUserID := 0
	var err error
	vars := mux.Vars(r)
	otherIDStr, passedID := vars[string(entity.IDKey)]

	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID

	switch passedID {
	case true:
		{
			otherUserID, _ = strconv.Atoi(otherIDStr)

		}
	case false: // ID was not passed
		{
			otherUsername := vars[string(entity.UsernameKey)]
			otherUser, err := chatInfo.userApp.GetUserByUsername(otherUsername)
			if err != nil {
				chatInfo.logger.Info(err.Error(),
					zap.String("url", r.RequestURI),
					zap.String("method", r.Method))
				switch err {
				case entity.UserNotFoundError:
					w.WriteHeader(http.StatusNotFound)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			otherUserID = otherUser.UserID
		}
	}

	chatID, err := chatInfo.chatApp.GetChatIDByUsers(userID, otherUserID)
	if err != nil {
		if err != entity.ChatNotFoundError {
			chatInfo.logger.Info(err.Error(),
				zap.String("url", r.RequestURI),
				zap.String("method", r.Method))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		chatID, err = chatInfo.chatApp.CreateChat(userID, otherUserID)
		if err != nil {
			chatInfo.logger.Info(err.Error(),
				zap.String("url", r.RequestURI),
				zap.String("method", r.Method))
			switch err {
			case entity.UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	messageInput := new(entity.MessageInput)
	err = json.NewDecoder(r.Body).Decode(messageInput)
	if err != nil {
		chatInfo.logger.Info(err.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if messageInput.MessageText == "" {
		chatInfo.logger.Info("Passed message is empty",
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := entity.Message{
		MessageID:      0,
		AuthorID:       userID,
		Text:           messageInput.MessageText,
		TimeOfCreation: "", // TODO: add datetime here
	}

	messageID, err := chatInfo.chatApp.AddMessage(chatID, &message)
	if err != nil {
		chatInfo.logger.Info(err.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = chatInfo.chatApp.SendMessage(chatID, messageID, otherUserID)
	if err != nil {
		if err != entity.ClientNotSetError {
			chatInfo.logger.Info(err.Error(),
				zap.String("url", r.RequestURI),
				zap.String("method", r.Method))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (chatInfo *ChatInfo) HandleReadChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr, _ := vars[string(entity.IDKey)]
	chatID, _ := strconv.Atoi(chatIDStr)

	userID := r.Context().Value(entity.CookieInfoKey).(*entity.CookieInfo).UserID

	err := chatInfo.chatApp.ReadChat(chatID, userID)
	if err != nil {
		chatInfo.logger.Info(err.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		switch err {
		case entity.ChatNotFoundError:
			w.WriteHeader(http.StatusNotFound)
		case entity.UserNotInChatError:
			w.WriteHeader(http.StatusForbidden)
		case entity.ChatAlreadyReadError:
			w.WriteHeader(http.StatusConflict)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
