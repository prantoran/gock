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
	Conn  *websocket.Conn
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

// LisErr stores error generated inside Listen's goroutine
type LisErr struct {
	v string
}

func (e *LisErr) Error() string {
	return e.v
}

// Connect dials a websocket url and stores the connection
func (s *SocListener) Connect() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	s.Conn = c
	return nil
}

// Listen listens to websocket
func (s *SocListener) Listen() error {
	var errL LisErr
	go func() {
		for {
			_, message, err := s.Conn.ReadMessage()
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

// Run creates websocket url, dials and listens to websock
func (s *SocListener) Run() error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("Connect err: %v", err)
	}
	defer s.Conn.Close()

	body := MsgBody{
		Token: s.Token,
		Type:  "register",
	}

	// registering token, useless
	if err := s.Conn.WriteJSON(body); err != nil {
		return fmt.Errorf("listener reg: %v", err)
	}

	if err := s.Listen(); err != nil {
		return fmt.Errorf("Listen err: %v", err)
	}

	return nil
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
