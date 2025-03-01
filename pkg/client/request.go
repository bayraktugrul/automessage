package client

type Request struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
