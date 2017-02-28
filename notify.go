package main

import (
	"runtime"
	"fmt"
	"os/exec"
	"log"
	"encoding/json"

	"github.com/gorilla/websocket"
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

type Notification struct {
	Message string
}

func main() {
	fmt.Printf("Start\n")

	url := "ws://localhost:5000/"
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		res := Notification{}

		json.Unmarshal(msg, &res)

		notify(res.Message, "test", "", "")

		log.Printf("received: %s\n", msg)
	}
}