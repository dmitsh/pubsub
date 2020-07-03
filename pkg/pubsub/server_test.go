package pubsub

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	msg1 = "message 1"
	msg2 = "message 2"
)

func TestPubSub(t *testing.T) {
	// start server
	port, err := getAvailablePort()
	if err != nil {
		t.Fatalf(err.Error())
	}
	srv := NewServer(port)
	go srv.Start()
	// allow to setup
	time.Sleep(time.Second)
	// start 3 subscribers
	addr := fmt.Sprintf("localhost:%d", port)
	data := make(chan string)
	for i := 0; i < 3; i++ {
		go func() {
			if err := Subscriber(addr, data); err != nil {
				t.Fatalf(err.Error())
			}
		}()
	}
	// allow to setup
	time.Sleep(time.Second)
	// publish 2 messages
	Publish(addr, msg1)
	Publish(addr, msg2)

	// verify
	var cnt1, cnt2 int
	timer := time.NewTimer(3 * time.Second)
	for {
		select {
		case msg := <-data:
			switch msg {
			case msg1:
				cnt1++
			case msg2:
				cnt2++
			default:
				t.Fatalf("Unexpected message %q", msg)
			}
		case <-timer.C:
			require.Equal(t, cnt1, 3)
			require.Equal(t, cnt2, 3)
			return
		}
	}
}

func getAvailablePort() (int, error) {
	// find unused port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err = l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}
