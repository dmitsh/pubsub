package pubsub

import (
	"bytes"
	"fmt"
	"net/http"
)

// Publish sends a message to broadcast
func Publish(addr, data string) error {
	resp, err := http.Post(fmt.Sprintf("http://%s/publish", addr), "application/text", bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status %d : %s", resp.StatusCode, resp.Status)
	}
	return nil
}
