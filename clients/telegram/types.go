package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type BotMessageResponse struct {
	Ok     bool       `json:"ok"`
	Result BotMessage `json:"result"`
}

type BotMessage struct {
	ID   int64 `json:"message_id"`
	Chat Chat  `json:"chat"`
}

type From struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

type Chat struct {
	ID int64 `json:"id"`
}
