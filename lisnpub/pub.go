package main

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Register = "reg"
	Admin    = "admin"
	Send     = "relay"
)

type Publisher struct {
	Token string
	Path  string
	Host  string
	Conn  *websocket.Conn
}

func NewPublisher(host, path, tok string) *Publisher {
	s := Publisher{}
	s.Token = tok
	s.Host = host
	s.Path = path
	return &s
}

// PubMsg is a msg sent by the publisher
type PubMsg struct {
	Type      string  `json:"type,omitempty"`
	Msg       string  `json:"msg,omitempty"`
	UserID    int32   `json:"user_id,omitempty"`
	DriverID  int32   `json:"driver_id,omitempty"`
	CurLat    float32 `json:"current_latitude,omitempty"`
	CurLong   float32 `json:"current_longitude,omitempty"`
	PromptAct string  `json:"prompt_action,omitempty"`
}

func (s *Publisher) Connect() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	s.Conn = c
	return err
}

var FixedObj = PubMsg{
	UserID:    43,
	DriverID:  12,
	CurLat:    23.89778678,
	CurLong:   90.87764554,
	Type:      "relay",
	PromptAct: "update_ride",
}

func insertDuration(t time.Duration, a []time.Duration) {
	l := len(a)
	cur := t
	for i := 0; i < l && i < 10; i++ {
		u := a[i]
		if u < cur {
			for j := i; j < 10 && j < l; j++ {
				u = a[j]
				a[j] = cur
				cur = u
			}
			if l < 10 {
				a = append(a, cur)
			}
			return
		}
	}
	if l < 10 {
		maxTimes = append(maxTimes, t)
	}
}
