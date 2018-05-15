package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Register = "reg"
	Admin    = "admin"
	Send     = "send"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

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

type MsgBody struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Message string `json:"message,omitempty"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	cUserID := flag.Arg(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: "userid=" + cUserID}
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

	for {
		select {
		case msg := <-mes:
			if cUserID != Admin {
				continue
			}
			w := strings.Fields(msg)
			tp := w[0]
			userID := w[1]

			body := MsgBody{}
			body.Type = tp
			body.UserID = userID
			if tp == Send {
				body.Message = w[2]
			}
			fmt.Println("body tp:", body.Type, " userid:", body.UserID, " msg:", body.Message)
			err := c.WriteJSON(body)
			if err != nil {
				log.Println("write:", err)
				return
			}
			// s := userID + " " + msg
			// err := c.WriteMessage(websocket.TextMessage, []byte(s))
			// if err != nil {
			// 	log.Println("write:", err)
			// 	return
			// }
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
