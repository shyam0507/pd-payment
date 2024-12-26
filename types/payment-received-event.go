package types

type PaymentReceivedEvent struct {
	SpecVersion     string      `json:"specversion"`
	Type            string      `json:"type"`
	Source          string      `json:"source"`
	Subject         string      `json:"subject"`
	Id              string      `json:"id"`
	Time            string      `json:"time"`
	DataContentType string      `json:"datacontenttype"`
	Data            PaymentInfo `json:"data"`
}

type PaymentInfo struct {
	OrderId   string  `json:"order_id"`
	PaymentId string  `json:"payment_id"`
	Total     float64 `json:"amount"`
}
