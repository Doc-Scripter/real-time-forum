package models

type Message struct {
	Type string      `json:"type"`
	Data any `json:"data"`
	From string      `json:"from,omitempty"`
	To   string      `json:"to,omitempty"`
}