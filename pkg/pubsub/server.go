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
	mutex    sync.Mutex
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

	srv.mutex.Lock()
	for ws, wg := range srv.conns {
		if err := ws.WriteMessage(1, body); err != nil {
			log.Println(err)
			wg.Done()
			delete(srv.conns, ws)
		}
	}
	srv.mutex.Unlock()
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
	log.Printf("Registered client %s", ws.RemoteAddr().String())
	var wg sync.WaitGroup
	wg.Add(1)
	srv.mutex.Lock()
	srv.conns[ws] = &wg
	srv.mutex.Unlock()
	wg.Wait()
}
