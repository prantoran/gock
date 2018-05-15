package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	Register = "reg"
	Admin    = "admin"
	Send     = "relay"
)

type SocSender struct {
	Token string
	Path  string
	Host  string
	MaxID int
	Turns int
}

func NewSocSender(host, path, tok string, n, m int) *SocSender {
	s := SocSender{}
	s.Token = tok
	s.Host = host
	s.Path = path
	s.MaxID = n
	s.Turns = m
	return &s
}

type MsgBody struct {
	Type      string  `json:"type,omitempty"`
	Msg       string  `json:"msg,omitempty"`
	UserID    int32   `json:"user_id,omitempty"`
	DriverID  int32   `json:"driver_id,omitempty"`
	CurLat    float32 `json:"current_latitude,omitempty"`
	CurLong   float32 `json:"current_longitude,omitempty"`
	PromptAct string  `json:"prompt_action,omitempty"`
}

// Run creates websocket url, dials and listens to websock
func (s *SocSender) Run() error {

	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	r := rand.New(rand.NewSource(99))

	_ := `{
		"type": "relay",
		"driver_id": 12,
		"user_id": 43,
		"current_latitude": 23.89778678,
		"current_longitude": 90.87764554,
		"prompt_action": "update_ride"
	  }
	`
	
	fixedObj := MsgBody{
		UserID:    43,
		DriverID:  12,
		CurLat:    23.89778678,
		CurLong:   90.87764554,
		Type:      "relay",
		PromptAct: "update_ride",
	}

	for i := 1; i <= s.Turns; i++ {
		go func(n, msgID int) {
			u := r.Intn(n) + 1
			body := MsgBody{}
			body.Type = "relay"
			body.UserID = int32(u)
			body.Msg = fmt.Sprintf("msg%d", msgID)

			// fmt.Println("body tp:", body.Type, " userid:", body.UserID, " msg:", body.Msg)
			fmt.Println("Sent to uid:", u, " msg:", msgID)
			if err := c.WriteJSON(body); err != nil {
				log.Printf("write: %v", err)
			}
			if err := c.WriteJSON(fixedObj); err != nil {
				log.Printf("write: %v", err)
			}
		}(s.MaxID, i)
	}
	return nil
}
