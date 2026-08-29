package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"notification/internal/client/http/telegram"
	"notification/internal/model"
	"platform/pkg/logger"
)

type NotificationService struct {
	telegramClient *telegram.Client
}

func NewNotificationService(telegramClient *telegram.Client) *NotificationService {
	return &NotificationService{
		telegramClient: telegramClient,
	}
}

// SendOrderPaidNotification отправляет уведомление об оплате
func (s *NotificationService) SendOrderPaidNotification(ctx context.Context, event *model.OrderPaidEvent) error {
	log := logger.Logger()

	text := fmt.Sprintf(
		"🛒 <b>Заказ оплачен</b>\n\n"+
			"Номер заказа: <code>%s</code>\n"+
			"Пользователь: <code>%s</code>\n"+
			"Способ оплаты: %s\n"+
			"Транзакция: <code>%s</code>\n\n"+
			"⏳ Ожидайте сборку корабля...",
		event.OrderUUID,
		event.UserUUID,
		event.PaymentMethod,
		event.TransactionUUID,
	)

	if err := s.telegramClient.SendMessage(ctx, text); err != nil {
		log.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}

	log.Info(ctx, "Order paid notification sent", zap.String("order_uuid", event.OrderUUID))
	return nil
}

// SendOrderAssembledNotification отправляет уведомление о сборке
func (s *NotificationService) SendOrderAssembledNotification(ctx context.Context, event *model.OrderAssembledEvent) error {
	log := logger.Logger()

	text := fmt.Sprintf(
		"🚀 <b>Корабль собран!</b>\n\n"+
			"Номер заказа: <code>%s</code>\n"+
			"Пользователь: <code>%s</code>\n"+
			"⏱ Время сборки: %d секунд\n\n"+
			"✅ Заказ готов к отправке!",
		event.OrderUUID,
		event.UserUUID,
		event.BuildTimeSec,
	)

	if err := s.telegramClient.SendMessage(ctx, text); err != nil {
		log.Error(ctx, "Failed to send order assembled notification", zap.Error(err))
		return err
	}

	log.Info(ctx, "Order assembled notification sent", zap.String("order_uuid", event.OrderUUID))
	return nil
}