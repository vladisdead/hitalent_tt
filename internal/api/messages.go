package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"hitalent_tt/model"
)

func (s *Server) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		if r.Method != "POST" {
			http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
			s.log.Error(MethodNotAllowed)
			return
		}
	}

	trimmed := strings.TrimPrefix(r.URL.Path, "/chats/")
	trimmed = strings.TrimSuffix(trimmed, "/messages")

	chatIDStr := strings.Trim(trimmed, "/")

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, ParsingBodyError, http.StatusBadRequest)
		s.log.Error(ParsingBodyError)
		return
	}

	var message *model.Message

	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error(ParsingBodyError)
		return
	}

	messageText := message.Text

	if strings.TrimSpace(messageText) == "" || utf8.RuneCountInString(message.Text) > 5000 {
		http.Error(w, MessageTextError, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		s.log.Error(MessageTextError)
		return
	}

	chatMessage, err := s.provider.SendMessage(chatID, message.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		s.log.Error(err.Error())
		return
	}

	if chatMessage == nil {
		http.Error(w, ChatNotFound, http.StatusNotFound)
		s.log.Error(ChatNotFound)
		return
	}

	replyReponse(w, chatMessage)
}
