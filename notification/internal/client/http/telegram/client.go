package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"platform/pkg/logger"
)

type Client struct {
	token  string
	chatID string
	client *http.Client
}

func NewClient(token, chatID string) *Client {
	return &Client{
		token:  token,
		chatID: chatID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendMessage отправляет сообщение в Telegram
func (c *Client) SendMessage(ctx context.Context, text string) error {
	log := logger.Logger()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	payload := map[string]interface{}{
		"chat_id": c.chatID,
		"text":    text,
		"parse_mode": "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn(ctx, "Telegram API returned non-200 status", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("telegram API status: %d", resp.StatusCode)
	}

	log.Info(ctx, "Telegram message sent", zap.String("chat_id", c.chatID))
	return nil
}