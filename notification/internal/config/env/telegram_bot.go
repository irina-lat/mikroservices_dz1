package env

import (
	"os"
)

type TelegramBotConfig struct {
	token  string
	chatID string
}

func LoadTelegramBotConfig() (*TelegramBotConfig, error) {
	token := os.Getenv("NOTIFICATION_TELEGRAM_BOT_TOKEN")
	if token == "" {
		token = "your_bot_token_here"
	}

	chatID := os.Getenv("NOTIFICATION_TELEGRAM_CHAT_ID")
	if chatID == "" {
		chatID = "your_chat_id_here"
	}

	return &TelegramBotConfig{
		token:  token,
		chatID: chatID,
	}, nil
}

func (c *TelegramBotConfig) Token() string  { return c.token }
func (c *TelegramBotConfig) ChatID() string { return c.chatID }