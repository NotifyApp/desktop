package main

import (
	"runtime"
	"fmt"
	"os/exec"
	"github.com/gorilla/websocket"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"encoding/json"
	"io"
	"strings"
	//"time"
)

func notify(text string, title string, sound string, icon string) {
	switch runtime.GOOS {
	case "darwin":
		notification := fmt.Sprintf("display notification \"%s\" with title \"%s\" sound name \"%s\"", text, title, sound)
		cmd := exec.Command("osascript", "-e", notification)
		cmd.Start()
	case "linux":
		cmd := exec.Command("notify-send", "-i", icon, title, text)
		cmd.Start()
	case "windows":
		fmt.Printf("Soon\n")
	}
}

var addr = flag.String("addr", "localhost:5000", "http service address")

type Notification struct {
	title, text string
}

func main() {
	fmt.Printf("Start\n")

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			dec := json.NewDecoder(strings.NewReader(string(message)))
			var m Notification
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			if (m.text != "") {
				notify(m.text, m.title, "Glass", "")
			}
		}
	}()

	for {
		
	}
}