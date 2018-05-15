package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	addr = "10.111.105.196:8081"
	// addr = "192.168.0.167:8081"
	// addr = flag.String("addr", "localhost:8080", "http service address")
)

// SocListener listens to websocket messages
type SocListener struct {
	Token string
	Path  string
	Host  string
}

// NewSocListener creates new instances
func NewSocListener(host, path, tok string) *SocListener {
	s := SocListener{}
	s.Token = tok
	s.Host = host
	s.Path = path
	return &s
}

// MsgBody is a websocket message
type MsgBody struct {
	Type  string `json:"type,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}

type LisErr struct {
	v string
}

func (e *LisErr) Error() string {
	return e.v
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

	var errL LisErr
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errL = LisErr{
					v: err.Error(),
				}
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	return &errL
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	token := flag.Arg(0)

	listener := NewSocListener(addr, "/ws", token)
	if err := listener.Run(); err != nil {
		log.Println("NewSocListener err:", err)
	}

	forever := make(chan bool)
	<-forever
}

/*
00wO3ztoWQ16AfsGaCS5ddaUz6pkY2FpeHavP2RB
*/
