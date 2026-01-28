package api

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_api "hitalent_tt/internal/api/mocks"
	"hitalent_tt/model"

	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_SendMessage(t *testing.T) {
	type mockBehavior func(s *mock_api.MockСhatStorage, message string)

	var buff bytes.Buffer
	mockLogger := slog.New(slog.NewTextHandler(&buff, &slog.HandlerOptions{}))

	createdAt := time.Now()

	expectedBodyJSON := &model.Chat{
		Id:        1,
		Tittle:    "test chat",
		CreatedAt: time.Now(),
		Messages: []model.Message{
			{
				Id:        1,
				ChatID:    1,
				Text:      "test message",
				CreatedAt: createdAt,
			},
		},
	}

	testMessage := []struct {
		testName       string
		inputBody      string
		message        string
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:  "success message",
			inputBody: `{"text": "test message"}`,
			message:   "test message",
			mockBehavior: func(s *mock_api.MockСhatStorage, message string) {
				s.EXPECT().SendMessage(1, message).Return(expectedBodyJSON, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf("{\"id\":1,\"tittle\":\"test chat\",\"created_at\":\"%s\",\"messages\":[{\"id\":1,\"chat_id\":1,\"text\":\"test message\",\"created_at\":\"%s\"}]}\n",
				createdAt.Format(time.RFC3339Nano),
				createdAt.Format(time.RFC3339Nano)),
		},
		{
			testName:  "chat not found",
			inputBody: `{"text": "test message"}`,
			message:   "test message",
			mockBehavior: func(s *mock_api.MockСhatStorage, message string) {
				s.EXPECT().SendMessage(1, message).Return(nil, nil)
			},
			expectedBody:   ChatNotFound + "\n",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range testMessage {
		t.Run(test.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mock_api.NewMockСhatStorage(c)
			test.mockBehavior(chat, test.message)

			srv := &Server{provider: chat, log: mockLogger}

			handlerFunc := http.HandlerFunc(srv.SendMessage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/chats/1/messages", bytes.NewBufferString(test.inputBody))

			handlerFunc.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
