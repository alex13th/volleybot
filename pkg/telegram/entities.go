package telegram

import "regexp"

type User struct {
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	UserName                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

type Location struct {
	Longitude            float32 `json:"longitude"`
	Latitude             float32 `json:"latitude"`
	HorizontalAccuracy   float32 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}

type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

type Chat struct {
	Id                    int             `json:"id"`
	Type                  string          `json:"type"`
	Title                 string          `json:"title"`
	Username              string          `json:"username"`
	FirstName             string          `json:"first_name"`
	LastName              string          `json:"last_name"`
	Photo                 ChatPhoto       `json:"photo"`
	Bio                   string          `json:"bio"`
	Description           string          `json:"description"`
	InviteLink            string          `json:"invite_link"`
	PinnedMessage         Message         `json:"pinned_message"`
	Permissions           ChatPermissions `json:"permissions"`
	SlowModeDelay         int             `json:"slow_mode_delay"`
	MessageAutoDeleteTime int             `json:"message_auto_delete_time"`
	StickerSetName        string          `json:"sticker_set_name"`
	CanSetStickerSet      bool            `json:"can_set_sticker_set"`
	LinkedChatId          int             `json:"linked_chat_id"`
	Location              ChatLocation    `json:"location"`
}

type ChatShared struct {
	RequestId int `json:"request_id"`
	ChatId    int `json:"chat_id"`
}

type KeyboardButtonRequestChat struct {
	RequestId   int  `json:"request_id"`
	BotIsMember bool `json:"bot_is_member"`
}

type KeyboardButton struct {
	Text        string                    `json:"text"`
	RequestChat KeyboardButtonRequestChat `json:"request_chat"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	Url                          string `json:"url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
	Pay                          bool   `json:"pay"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	Url      string `json:"url"`
	User     *User  `json:"user"`
	Language string `json:"language"`
}

type Message struct {
	MessageId             int             `json:"message_id"`
	From                  *User           `json:"from"`
	SenderChat            *Chat           `json:"sender_chat"`
	Date                  int             `json:"date"`
	Chat                  *Chat           `json:"chat"`
	ForwardFrom           *User           `json:"forward_from"`
	ForwardFromChat       *Chat           `json:"forward_from_chat"`
	ForwardFromMessage_id int             `json:"forward_from_message_id"`
	ForwardSignature      string          `json:"forward_signature"`
	ForwardSenderName     string          `json:"forward_sender_name"`
	ForwardDate           int             `json:"forward_date"`
	ReplyToMessage        *Message        `json:"reply_to_message"`
	ChatShared            *ChatShared     `json:"chat_shared"`
	ViaBot                *User           `json:"via_bot"`
	EditDate              int             `json:"edit_date"`
	MediaGroupId          int             `json:"media_group_id"`
	AuthorSignature       string          `json:"author_signature"`
	Text                  string          `json:"text"`
	Entities              []MessageEntity `json:"entities"`
	Caption               string          `json:"caption"`
	CaptionEentities      []MessageEntity `json:"caption_entities"`
	ReplyMarkup           interface{}     `json:"reply_markup"`
}

func (msg Message) GetCommand() string {
	re, err := regexp.Compile(`^/([a-zA-Z0-9_]*)`)
	if err != nil {
		return ""
	}
	matches := re.FindStringSubmatch(msg.Text)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (msg Message) IsCommand() bool {
	return msg.GetCommand() != ""
}

func (msg Message) SendMessage(tb Bot, Text string, mr MessageRequest) (*MessageResponse, error) {
	return tb.SendMessage(msg.CreateMessageRequest(Text, mr))
}

func (msg Message) DeleteMessage(tb Bot) {
	tb.SendMessage(&DeleteMessageRequest{ChatId: msg.Chat.Id, MessageId: msg.MessageId})
}

func (msg Message) Reply(tb Bot, Text string, mr MessageRequest) (*MessageResponse, error) {
	return tb.SendMessage(msg.CreateReplyRequest(Text, mr))
}

func (msg Message) CreateReplyRequest(Text string, mr MessageRequest) (result MessageRequest) {
	result = msg.CreateMessageRequest(Text, mr)
	result.ReplyToMessageId = msg.MessageId
	return
}

func (msg *Message) CreateMessageRequest(Text string, mr MessageRequest) (result MessageRequest) {
	result.ChatId = msg.Chat.Id
	result.Text += Text

	return
}

func (msg Message) EditText(tb Bot, Text string, mer EditMessageTextRequest) (*MessageResponse, error) {
	return tb.SendMessage(msg.CreateEditTextRequest(Text, mer))
}

func (msg Message) CreateEditTextRequest(Text string, Request EditMessageTextRequest) EditMessageTextRequest {
	return EditMessageTextRequest{ChatId: msg.Chat.Id, MessageId: msg.MessageId, Text: Text}
}

type CallbackQuery struct {
	Id              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageId string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

func (cq CallbackQuery) Answer(tb Bot, Text string, req AnswerCallbackQueryRequest) (*MessageResponse, error) {
	req.CallbackQueryId = cq.Id
	req.Text = Text
	return tb.SendMessage(req)
}

type Update struct {
	UpdateId          int            `json:"update_id"`
	Message           *Message       `json:"message"`
	EditedMessage     *Message       `json:"edited_message"`
	ChannelPost       *Message       `json:"channel_post"`
	EditedChannelPost *Message       `json:"edited_channel_post"`
	CallbackQuery     *CallbackQuery `json:"callback_query"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type BotCommandScope struct {
	Type string `json:"type"`
}

type BotCommandScopeChat struct {
	Type   string      `json:"type"`
	ChatId interface{} `json:"chat_id"`
}
