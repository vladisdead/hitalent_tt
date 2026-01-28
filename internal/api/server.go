package api

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"hitalent_tt/internal/config"
)

const (
	TittleError      = "tittle must be lower 200 or not empty"
	MessageTextError = "message must be lower 5000 or not empty"
	MethodNotAllowed = "method not allowed"
	ParsingBodyError = "error parsing body"
	ChatNotFound     = "chat not found"
	DefaultLimit     = 20
)

type Server struct {
	srv      *http.Server
	provider СhatStorage
	log      *slog.Logger
}

func NewAPI(p СhatStorage, cfg *config.Config, log *slog.Logger) *Server {

	srv := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		},
		provider: p,
		log:      log,
	}

	http.HandleFunc("/chats", srv.CreateChatHandler)
	http.HandleFunc("/chats/", func(w http.ResponseWriter, r *http.Request) {
		// проверяем, соответствует ли путь шаблону /chats/{id}/messages/
		if strings.HasSuffix(r.URL.Path, "/messages") {
			srv.SendMessage(w, r)
			return
		} else if r.Method == http.MethodDelete {
			srv.DeleteChatHandler(w, r)
			return
		} else {
			srv.GetChatByIDHandler(w, r)
			return
		}
	})

	return srv
}

func (s *Server) Start() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("error starting server")
		}
	}()

	s.log.Info("server started", "listening on", s.srv.Addr)
	<-stopChan
}

func replyReponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
}
