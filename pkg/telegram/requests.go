package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Request interface {
	GetParams() (url.Values, string, error)
}

type UpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func (req *UpdatesRequest) GetParams() (url.Values, string, error) {
	values := url.Values{}
	if req.Offset != 0 {
		values.Add("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		values.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Timeout > 0 {
		values.Add("timeout", strconv.Itoa(req.Timeout))
	}
	for _, au := range req.AllowedUpdates {
		values.Add("allowed_updates", au)
	}
	return values, "getUpdates", nil
}

type MessageRequest struct {
	ChatId                   int             `json:"chat_id"`
	Text                     string          `json:"text"`
	ParseMode                string          `json:"parse_mode"`
	Entities                 []MessageEntity `json:"entities"`
	DisableWebPagePreview    bool            `json:"disable_web_page_preview"`
	DisableNotification      bool            `json:"disable_notification"`
	ReplyToMessageId         int             `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`
	ReplyMarkup              interface{}     `json:"reply_markup"`
}

func (req *MessageRequest) GetParams() (url.Values, string, error) {
	values := url.Values{}
	values.Add("chat_id", strconv.Itoa(req.ChatId))
	values.Add("text", req.Text)
	if req.ParseMode != "" {
		values.Add("parse_mode", req.ParseMode)
	}
	if req.DisableWebPagePreview {
		values.Add("disable_web_page_preview", strconv.FormatBool(req.DisableWebPagePreview))
	}
	if req.DisableNotification {
		values.Add("disable_notification", strconv.FormatBool(req.DisableNotification))
	}
	if req.ReplyToMessageId > 0 {
		values.Add("reply_to_message_id", strconv.Itoa(req.ReplyToMessageId))
	}
	if req.AllowSendingWithoutReply {
		values.Add("allow_sending_without_reply", strconv.FormatBool(req.AllowSendingWithoutReply))
	}
	if len(req.Entities) > 0 {
		data, err := json.Marshal(req.Entities)
		if err != nil {
			return nil, "", err
		}
		values.Add("entities", string(data))
	}

	if req.ReplyMarkup != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		values.Add("reply_markup", string(data))
	}

	return values, "sendMessage", nil
}
