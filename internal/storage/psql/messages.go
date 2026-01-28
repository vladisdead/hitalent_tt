package psql

import (
	"time"

	"gorm.io/gorm"
	"hitalent_tt/model"
)

func (s *Storage) SendMessage(chatID int, text string) (*model.Chat, error) {

	var chat *model.Chat
	var message *model.Message

	chatResult := s.db.First(&chat, chatID)
	if chatResult.Error != nil {
		switch chatResult.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, chatResult.Error
		}
	}

	message = &model.Message{
		ChatID:    chatID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	messageResult := s.db.Create(&message)
	if messageResult.Error != nil {
		return nil, messageResult.Error
	}

	chat.Messages = append(chat.Messages, *message)

	return chat, nil
}
