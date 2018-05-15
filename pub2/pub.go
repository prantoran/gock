package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"

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
	Type   string `json:"type,omitempty"`
	Msg    string `json:"msg,omitempty"`
	UserID int32  `json:"user_id,omitempty"`
}

// Run creates websocket url, dials and listens to websock
func (s *SocSender) Run() error {

	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func(n, m int) {
		r := rand.New(rand.NewSource(99))

		for i := 1; i <= m; i++ {
			u := r.Intn(n) + 1
			body := MsgBody{}
			body.Type = "relay"
			body.UserID = int32(u)
			body.Msg = fmt.Sprintf("msg%d", i)

			// fmt.Println("body tp:", body.Type, " userid:", body.UserID, " msg:", body.Msg)
			if err := c.WriteJSON(body); err != nil {
				log.Printf("write: %v", err)
			}
		}
	}(s.MaxID, s.Turns)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %v", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
