package model

// OrderPaidEvent — событие об оплате заказа
type OrderPaidEvent struct {
	EventUUID       string `json:"event_uuid"`
	OrderUUID       string `json:"order_uuid"`
	UserUUID        string `json:"user_uuid"`
	PaymentMethod   string `json:"payment_method"`
	TransactionUUID string `json:"transaction_uuid"`
}

// OrderAssembledEvent — событие о сборке заказа
type OrderAssembledEvent struct {
	EventUUID    string `json:"event_uuid"`
	OrderUUID    string `json:"order_uuid"`
	UserUUID     string `json:"user_uuid"`
	BuildTimeSec int64  `json:"build_time_sec"`
}