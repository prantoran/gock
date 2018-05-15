package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

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
type PubMsg struct {
	Type   string `json:"type,omitempty"`
	Msg    string `json:"msg,omitempty"`
	UserID int32  `json:"user_id,omitempty"`
}

// Run creates websocket url, dials and listens to websock
func (s *SocListener) Run() error {

	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}

	body := MsgBody{
		Token: s.Token,
		Type:  "register",
	}

	// registering token, useless
	if err := c.WriteJSON(body); err != nil {
		return fmt.Errorf("listener reg: %v", err)
	}

	go func() {
		defer func() {
			drpCnt++
			fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		}()
		crtCnt++
		fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		for {
			msg := PubMsg{}
			err := websocket.ReadJSON(c, msg)

			if err != nil {
				log.Println("listener:", s.Token, " read:", err)
				return
			}
			
			uid := strconv.Itoa(int(msg.UserID))

			if uid != s.Token {
				incorrect++
				log.Println("Incorrect cnt:", incorrect, " target user:", uid, " received user:", s.Token)
			}
			log.Printf("recv: %v", msg)
		}
	}()

	return nil

}
