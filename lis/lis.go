package lis

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/prantoran/gock/pub"
)

var (
	// variables used to count active and dropped conns,
	drpCnt = 0
	crtCnt = 0
	// to count the number of times wrong clients received msgs
	incorrect = 0
)

// SocListener stores the vars req by a websocket listener
type SocListener struct {
	Token string
	Path  string
	Host  string
	Conn  *websocket.Conn
}

// NewSocListener returns a new SocListener
func NewSocListener(host, path, tok string) *SocListener {
	s := SocListener{}
	s.Token = tok
	s.Host = host
	s.Path = path
	return &s
}

// LisMsg represents tokens received by Listener
type LisMsg struct {
	Type  string `json:"type,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}

// Connect creates a websocket conn
func (s *SocListener) Connect() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Token}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	s.Conn = c
	return err
}

// Register registers listener to websocket server
func (s *SocListener) Register() error {
	// registering token, useless
	return s.Conn.WriteJSON(LisMsg{
		Token: s.Token,
		Type:  "register",
	})
}

// Run creates websocket url, dials and listens to websock
func (s *SocListener) Run() error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("listener Run Connect err: %v", err)
	}
	if err := s.Register(); err != nil {
		return fmt.Errorf("listener register err: %v", err)
	}
	go func() {
		defer func() {
			drpCnt++
			fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		}()
		crtCnt++
		fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		for {
			msg := pub.PubMsg{}
			err := websocket.ReadJSON(s.Conn, msg)

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

// CreateListeners create (endID-startID+1) listeners by calling NewSocListener.
// NewSocListener calls a goroutine, so calling CreateListeners as a gorouting is unadvised.
func CreateListeners(addr string, startID, endID int) error {
	if endID < startID {
		return fmt.Errorf("StartID greater than endID")
	}
	for i := startID; i <= endID; i++ {
		listener := NewSocListener(addr, "/ws", strconv.Itoa(i))
		listener.Run()
	}
	return nil
}
