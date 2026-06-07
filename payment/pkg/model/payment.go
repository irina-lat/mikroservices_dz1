package model

type PaymentMethod string

const (
	PaymentMethodUnknown      PaymentMethod = "UNKNOWN"
	PaymentMethodCard         PaymentMethod = "CARD"
	PaymentMethodSBP          PaymentMethod = "SBP"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PayOrderRequest struct {
	OrderUUID     string        `json:"order_uuid"`
	UserUUID      string        `json:"user_uuid"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}