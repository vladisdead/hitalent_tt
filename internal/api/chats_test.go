package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hitalent_tt/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_api "hitalent_tt/internal/api/mocks"
)

func TestHandler_CreateChatHandler(t *testing.T) {
	type mockBehavior func(s *mock_api.MockСhatStorage, chat model.Chat)

	var buff bytes.Buffer
	mockLogger := slog.New(slog.NewTextHandler(&buff, &slog.HandlerOptions{}))

	inputChat := model.Chat{
		Tittle: "Test",
	}
	createdAt := time.Now()

	testChat := []struct {
		testName       string
		inputBody      string
		inputChat      model.Chat
		mockBehavior   mockBehavior
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:  "success create",
			inputBody: `{"tittle":"Test"}`,
			inputChat: inputChat,
			mockBehavior: func(s *mock_api.MockСhatStorage, chat model.Chat) {
				s.EXPECT().CreateChat(inputChat.Tittle).Return(&model.Chat{Id: 0, Tittle: inputChat.Tittle, CreatedAt: createdAt}, nil)
			},
			expectedBody:   fmt.Sprintf("{\"id\":0,\"tittle\":\"Test\",\"created_at\":\"%s\"}\n", createdAt.Format(time.RFC3339Nano)),
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "tittle is empty",
			inputBody:      `{"tittle":"     "}`,
			inputChat:      inputChat,
			mockBehavior:   func(s *mock_api.MockСhatStorage, chat model.Chat) {},
			expectedBody:   TittleError + "\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range testChat {
		t.Run(test.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			chat := mock_api.NewMockСhatStorage(c)
			test.mockBehavior(chat, test.inputChat)

			srv := &Server{provider: chat, log: mockLogger}

			handlerFunc := http.HandlerFunc(srv.CreateChatHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewBufferString(test.inputBody))

			handlerFunc.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
