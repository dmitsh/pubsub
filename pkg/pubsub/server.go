package pubsub

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server is a PubSub server
type Server struct {
	port     int
	upgrader *websocket.Upgrader
	conns    map[*websocket.Conn]*sync.WaitGroup
}

// NewServer instantiates the server
func NewServer(port int) *Server {
	srv := &Server{
		port: port,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		conns: make(map[*websocket.Conn]*sync.WaitGroup),
	}
	http.HandleFunc("/publish", srv.publishEndpoint)
	http.HandleFunc("/subscribe", srv.subscribeEndpoint)
	return srv
}

// Start starts the server
func (srv *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", srv.port), nil))
}

func (srv *Server) publishEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Printf("Broadcasting %q", string(body))
	for ws, wg := range srv.conns {
		if err := ws.WriteMessage(1, body); err != nil {
			log.Println(err)
			wg.Done()
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) subscribeEndpoint(w http.ResponseWriter, r *http.Request) {
	srv.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Registered client")
	var wg sync.WaitGroup
	wg.Add(1)
	srv.conns[ws] = &wg
	defer delete(srv.conns, ws)
	wg.Wait()
}
