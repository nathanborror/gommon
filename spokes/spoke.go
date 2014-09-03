package spokes

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nathanborror/gommon/auth"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 512                 // Maximum message size allowed from peer.
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	ws     *websocket.Conn // The websocket connection.
	send   chan []byte     // Buffered channel of outbound messages.
	User   *auth.User
	Cookie *http.Cookie
}

// wsRequest
type wsRequest struct {
	connection *Connection
	message    []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump() {
	defer func() {
		Hub.Unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		Hub.Incoming <- wsRequest{c, message}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// SpokeHandler handles webocket requests from the peer.
func SpokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Upgrade request to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	// Check for authenticated user
	if !auth.IsAuthenticated(r) {
		return
	}

	// Grab user making request
	u, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		log.Println(err)
	}

	// Grab cookie
	cookie, _ := r.Cookie("authenticated-user")

	// Create connection
	c := &Connection{send: make(chan []byte, 256), ws: ws, User: u, Cookie: cookie}
	Hub.Register <- c
	go c.writePump()
	c.readPump()
}
