package presentation

type AmqpLogResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
