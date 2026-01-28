package psql

import (
	"time"

	"gorm.io/gorm"
	"hitalent_tt/model"
)

func (s *Storage) CreateChat(tittle string) (*model.Chat, error) {
	var chat *model.Chat

	chat = &model.Chat{
		Tittle:    tittle,
		CreatedAt: time.Now(),
	}

	chatResult := s.db.Create(&chat)

	if chatResult.Error != nil {
		return nil, chatResult.Error
	}

	return chat, nil
}

func (s *Storage) GetChatByID(chatId int, messageCount int) (*model.Chat, error) {

	var chat *model.Chat
	var message []model.Message

	chatResult := s.db.First(&chat, chatId)

	if chatResult.Error != nil {
		switch chatResult.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, chatResult.Error
		}
	}

	messageResult := s.db.Order("id desc").Limit(messageCount).Find(&message, model.Message{ChatID: chatId})
	if messageResult.Error != nil {
		return nil, messageResult.Error
	}

	chat.Messages = message

	return chat, nil
}

func (s *Storage) DeleteChat(id int) error {
	return s.db.Delete(&model.Chat{}, id).Error
}
