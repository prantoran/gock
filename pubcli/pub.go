package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Register = "reg"
	Admin    = "admin"
	Send     = "relay"
)

/* Function to run the groutine to run for stdin read */
func read(r io.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- scan.Text()
		}
	}()
	return lines
}

type SocSender struct {
	Token string
	Path  string
	Host  string
}

func NewSocSender(host, path, tok string) *SocSender {
	s := SocSender{}
	s.Token = tok
	s.Host = host
	s.Path = path
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
	log.Printf("connecting to %s", u.String())

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

	mes := read(os.Stdin) //Reading from Stdin

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case msg := <-mes:

			w := strings.Fields(msg)
			tp := w[0]
			uid, _ := strconv.Atoi(w[1])

			body := MsgBody{}
			body.Type = tp
			body.UserID = int32(uid)
			if tp == Send {
				body.Msg = w[2]
			}
			fmt.Printf("w: %v\n", w)
			fmt.Println("body tp:", body.Type, " userid:", body.UserID, " msg:", body.Msg)
			if err := c.WriteJSON(body); err != nil {
				return fmt.Errorf("write: %v", err)
			}

		case <-done:
			return nil
		// case t := <-ticker.C:
		// 	body := MsgBody{
		// 		Msg: "Token:" + s.Token + " time:" + t.String(),
		// 	}
		// 	if err := c.WriteJSON(body); err != nil {
		// 		return fmt.Errorf("write: %v", err)
		// 	}
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
