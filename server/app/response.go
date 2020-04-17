package app

type Response struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Channel  string `json:"channel"`
	SendUser string `json:"send_user"`
}
