package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// variables used to count active and dropped conns,
	drpCnt = 0
	crtCnt = 0
	// to count the number of times wrong clients received msgs
	incorrect = 0
)

// Listener stores the vars req by a websocket listener
type Listener struct {
	mu    sync.Mutex
	Token string
	Path  string
	Host  string
	Ack   map[string]chan bool
	Conn  *websocket.Conn
}

// NewListener returns a new Listener
func NewListener(host, path, tok string) *Listener {
	s := Listener{}
	s.Token = tok
	s.Host = host
	s.Path = path
	s.Ack = make(map[string]chan bool)
	return &s
}

// Msg represents tokens received by Listener
type LisMsg struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}

// Connect creates a websocket conn
func (s *Listener) Connect() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	s.Conn = c
	return err
}

// Register registers listener to websocket server
func (s *Listener) Register() error {
	// registering token, useless
	return s.Conn.WriteJSON(LisMsg{
		Token: s.Token,
		Type:  "register",
	})
}

// Run creates websocket url, dials and listens to websock
func (s *Listener) Run(wg *sync.WaitGroup) error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("listener Run Connect err: %v", err)
	}
	if err := s.Register(); err != nil {
		return fmt.Errorf("listener register err: %v", err)
	}

	wg.Done()

	go func() {
		defer func() {
			drpCnt++
			fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		}()
		crtCnt++
		fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		for {

			msg := PubMsg{}
			err := websocket.ReadJSON(s.Conn, &msg)

			if err != nil {
				log.Println("listener:", s.Token, " read:", err)
				return
			}

			uid := strconv.Itoa(int(msg.UserID))

			log.Println("uid:", uid, "s.Token:", s.Token)
			if uid != s.Token {
				incorrect++
				log.Println("Incorrect cnt:", incorrect, " target user:", uid, " received user:", s.Token)
			}
			log.Printf("recv: %v", msg)
			s.Ack[msg.Msg] <- true
		}
	}()
	return nil
}
