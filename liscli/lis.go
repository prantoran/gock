package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type SocListener struct {
	Token string
	Path  string
	Host  string
}

func NewSocListener(host, path, tok string) *SocListener {
	s := SocListener{}
	s.Token = tok
	s.Host = host
	s.Path = path
	return &s
}

type MsgBody struct {
	Type  string `json:"type,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}

// Run creates websocket url, dials and listens to websock
func (s *SocListener) Run() error {

	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	defer c.Close()

	body := MsgBody{
		Token: s.Token,
		Type:  "register",
	}

	// registering token, useless
	if err := c.WriteJSON(body); err != nil {
		return fmt.Errorf("listener reg: %v", err)
	}

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-done:
			return nil
		// case t := <-ticker.C:

		// 	body := MsgBody{
		// 		Msg: "Token:" + s.Token + " time:" + t.String(),
		// 	}
		// 	if err := c.WriteJSON(body); err != nil {
		// 		return fmt.Errorf("write: %v", err)
		// 	}
		// 	// err := c.WriteMessage(websocket.TextMessage, []byte("Token:"+s.Token+" time:"+t.String()))
		// 	// if err != nil {
		// 	// 	return fmt.Errorf("write: %v", err)
		// 	// }
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
