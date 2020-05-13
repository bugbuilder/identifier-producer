package types

type Identifier struct {
	TransactionId string `json:"transactionId, omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
