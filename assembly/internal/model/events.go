package model

// OrderPaidEvent — входящее событие об оплате заказа
// Контракт: AssemblyService получает это событие из Kafka
type OrderPaidEvent struct {
	EventUUID       string `json:"event_uuid"`        // Уникальный ID события (идемпотентность)
	OrderUUID       string `json:"order_uuid"`        // ID оплаченного заказа
	UserUUID        string `json:"user_uuid"`         // ID пользователя
	PaymentMethod   string `json:"payment_method"`    // Способ оплаты (CARD, SBP, ...)
	TransactionUUID string `json:"transaction_uuid"`  // ID транзакции из PaymentService
}

// ShipAssembledEvent — исходящее событие о сборке заказа
// Контракт: AssemblyService отправляет это событие в OrderService
type ShipAssembledEvent struct {
	EventUUID    string `json:"event_uuid"`     // Уникальный ID события (идемпотентность)
	OrderUUID    string `json:"order_uuid"`     // ID собранного заказа
	UserUUID     string `json:"user_uuid"`      // ID пользователя
	BuildTimeSec int64  `json:"build_time_sec"` // Время сборки в секундах
}