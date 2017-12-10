package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"

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
	Title   string `json:"title"`
	Message string `json:"message"`
}

func main() {
	fmt.Printf("Start\n")

	url := "ws://localhost:8080/ws"
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

		var res Notification
		json.Unmarshal(msg, &res)

		notify(res.Message, res.Title, "", "")

		log.Printf("received: %s\n", msg)
	}
}
