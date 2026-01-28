package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"hitalent_tt/model"
)

//go:generate mockgen -source=chats.go -destination=mocks/mock.go
type СhatStorage interface {
	CreateChat(tittle string) (*model.Chat, error)
	GetChatByID(chatId int, messageCount int) (*model.Chat, error)
	DeleteChat(id int) error
	SendMessage(chatID int, text string) (*model.Chat, error)
}

func (s *Server) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		s.log.Error(MethodNotAllowed)
		return
	}

	var chat *model.Chat

	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error(ParsingBodyError)
		return
	}

	tittle := strings.TrimSpace(chat.Tittle)

	if tittle == "" || utf8.RuneCountInString(tittle) > 200 {
		http.Error(w, TittleError, http.StatusBadRequest)
		s.log.Error(TittleError)
		return
	}

	chat, err = s.provider.CreateChat(tittle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error("Error creating chat: " + err.Error())
		return
	}

	replyReponse(w, chat)
}

func (s *Server) GetChatByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		s.log.Error(MethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/chats/"))
	if err != nil {
		http.Error(w, ParsingBodyError, http.StatusBadRequest)
		s.log.Error(ParsingBodyError)
		return
	}

	var limit int

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, ParsingBodyError, http.StatusBadRequest)
			s.log.Error(ParsingBodyError)
			return
		}
	}

	if limit > 100 {
		http.Error(w, LimitError, http.StatusBadRequest)
		s.log.Error(LimitError)
		return
	}

	if limit == 0 {
		limit = DefaultLimit
	}

	chat, err := s.provider.GetChatByID(id, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error(err.Error())
		return
	}

	if chat == nil {
		http.Error(w, ChatNotFound, http.StatusNotFound)
		s.log.Error(ChatNotFound)
		return
	}

	replyReponse(w, chat)
}

func (s *Server) DeleteChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		s.log.Error(MethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/chats/"))
	if err != nil {
		http.Error(w, ParsingBodyError, http.StatusBadRequest)
		s.log.Error(ParsingBodyError)
		return
	}

	err = s.provider.DeleteChat(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
