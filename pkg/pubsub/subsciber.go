package pubsub

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Subscriber starts subscriber loop
func Subscriber(addr string, data chan<- string) error {
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/subscribe", addr), nil)
	if err != nil {
		return err
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return err
		}
		data <- string(message)
	}
}
