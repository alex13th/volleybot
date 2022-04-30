package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
	"sync"
)

type Request interface {
	GetParams() (url.Values, string, error)
}

type UpdatesRequest struct {
	mu             sync.RWMutex
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func (req *UpdatesRequest) GetParams() (val url.Values, method string, err error) {
	req.mu.RLock()
	method = "getUpdates"
	val = url.Values{}
	if req.Offset != 0 {
		val.Add("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit > 0 {
		val.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Timeout > 0 {
		val.Add("timeout", strconv.Itoa(req.Timeout))
	}
	for _, au := range req.AllowedUpdates {
		val.Add("allowed_updates", au)
	}
	defer req.mu.RUnlock()
	return
}

type MessageRequest struct {
	mu                       sync.RWMutex
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

func (req *MessageRequest) GetParams() (val url.Values, method string, err error) {
	req.mu.RLock()
	method = "sendMessage"
	val = url.Values{}
	val.Add("chat_id", strconv.Itoa(req.ChatId))
	val.Add("text", req.Text)
	if req.ParseMode != "" {
		val.Add("parse_mode", req.ParseMode)
	}
	if req.DisableWebPagePreview {
		val.Add("disable_web_page_preview", strconv.FormatBool(req.DisableWebPagePreview))
	}
	if req.DisableNotification {
		val.Add("disable_notification", strconv.FormatBool(req.DisableNotification))
	}
	if req.ReplyToMessageId > 0 {
		val.Add("reply_to_message_id", strconv.Itoa(req.ReplyToMessageId))
	}
	if req.AllowSendingWithoutReply {
		val.Add("allow_sending_without_reply", strconv.FormatBool(req.AllowSendingWithoutReply))
	}
	if len(req.Entities) > 0 {
		data, err := json.Marshal(req.Entities)
		if err != nil {
			return nil, "", err
		}
		val.Add("entities", string(data))
	}

	if req.ReplyMarkup != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	defer req.mu.RUnlock()
	return
}

type EditMessageTextRequest struct {
	mu                    sync.RWMutex
	ChatId                int             `json:"chat_id"`
	MessageId             int             `json:"message_id"`
	InlineMessageId       int             `json:"inline_message_id"`
	Text                  string          `json:"text"`
	ParseMode             string          `json:"parse_mode"`
	Entities              []MessageEntity `json:"entities"`
	DisableWebPagePreview bool            `json:"disable_web_page_preview"`
	ReplyMarkup           interface{}     `json:"reply_markup"`
}

func (req *EditMessageTextRequest) GetParams() (val url.Values, method string, err error) {
	method = "editMessageText"
	val = url.Values{}
	req.mu.RLock()
	val.Add("chat_id", strconv.Itoa(req.ChatId))
	val.Add("message_id", strconv.Itoa(req.MessageId))
	val.Add("inline_message_id", strconv.Itoa(req.InlineMessageId))
	val.Add("text", req.Text)
	if req.ParseMode != "" {
		val.Add("parse_mode", req.ParseMode)
	}
	if req.DisableWebPagePreview {
		val.Add("disable_web_page_preview", strconv.FormatBool(req.DisableWebPagePreview))
	}
	if len(req.Entities) > 0 {
		data, err := json.Marshal(req.Entities)
		if err != nil {
			return nil, "", err
		}
		val.Add("entities", string(data))
	}

	if req.ReplyMarkup != nil {
		data, err := json.Marshal(req.ReplyMarkup)
		if err != nil {
			return nil, "", err
		}
		val.Add("reply_markup", string(data))
	}
	defer req.mu.RUnlock()
	return
}

type AnswerCallbackQueryRequest struct {
	mu              sync.RWMutex
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
	URL             string `json:"url"`
	CacheTime       int    `json:"cache_time"`
}

func (req *AnswerCallbackQueryRequest) GetParams() (val url.Values, method string, err error) {
	method = "answerCallbackQuery"
	val = url.Values{}
	req.mu.RLock()
	val.Add("callback_query_id", req.CallbackQueryId)
	val.Add("text", req.Text)
	val.Add("show_alert", strconv.FormatBool(req.ShowAlert))
	val.Add("url", req.URL)
	val.Add("cache_time", strconv.Itoa(req.CacheTime))
	defer req.mu.RUnlock()
	return
}
