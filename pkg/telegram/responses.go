package telegram

import (
	"encoding/json"
	"io"
)

type UpdateResponse struct {
	Ok          bool               `json:"ok"`
	Result      []Update           `json:"result"`
	Description string             `json:"description"`
	ErrorCode   int                `json:"error_code"`
	Parameters  ResponseParameters `json:"parameters"`
}

type MessageResponse struct {
	Ok         bool               `json:"ok"`
	Result     Message            `json:"result"`
	ErrorCode  int                `json:"error_code"`
	Parameters ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
}

func ParseJson(i interface{}, reader io.Reader) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(i)
}

func (update *UpdateResponse) Parse(reader io.Reader) error {
	return ParseJson(update, reader)
}

func (message *MessageResponse) Parse(reader io.Reader) error {
	return ParseJson(message, reader)
}
