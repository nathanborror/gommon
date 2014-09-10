// Package spokes is a basic WebSocket publish / subscribe toolkit. It builds
// off of gorilla/websocket's chat example. Clients can subscribe to URLs and
// receive updates when they change.
package spokes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type hub struct {
	Connections   map[*Connection]bool // Registered connections
	Incoming      chan wsRequest       // Inbound messages from the connections
	Register      chan *Connection     // Register requests from the connections
	Unregister    chan *Connection     // Unregister requests from connections
	Subscriptions Subscriptions
}

// Hub maintains the set of active connections and broadcasts messages to
// subscribers. Use the Run function to start the Hub (i.e go spokes.Hub.Run())
var Hub = hub{
	Connections:   make(map[*Connection]bool),
	Incoming:      make(chan wsRequest),
	Register:      make(chan *Connection),
	Unregister:    make(chan *Connection),
	Subscriptions: Subscriptions{},
}

// Subscriptions holds all the subscribed clients
type Subscriptions map[string]map[*Connection]bool

func (s Subscriptions) add(channel string, c *Connection) {
	if _, ok := s[channel]; !ok {
		s[channel] = make(map[*Connection]bool)
	}
	s[channel][c] = true
}

func (s Subscriptions) remove(c *Connection) {
	for _, subs := range s {
		delete(subs, c)
	}
}

type action string

// Subscribe subscribes someone to receive updates for
// a particular url. Request ...
const (
	Request   action = "request"
	Subscribe        = "subscribe"
)

type message struct {
	Port   string
	URL    string
	Action action
}

type wsResponse struct {
	Channel string
	Data    *json.RawMessage
}

func (h *hub) sendToChannel(channel string, message []byte) {
	// Only send to clients subscribed to the request URL
	subscriptions, ok := h.Subscriptions[channel]
	if !ok {
		log.Printf("[Spokes]: No subscriptions for '%s'\n", channel)
		return
	}
	toWrap := json.RawMessage(message)
	m := wsResponse{
		Channel: channel,
		Data:    &toWrap,
	}
	body, err := json.Marshal(m)
	if err != nil {
		log.Println("[Spokes]: Error marshalling channel message:", err)
	}
	for c := range subscriptions {
		log.Println("[Spokes]: Sending to ", c.User.Hash)

		select {
		case c.send <- body:
		default:
			close(c.send)
			delete(subscriptions, c)
			delete(h.Connections, c)
		}
	}
}

func (h *hub) handleMessage(conn *Connection, m message) {
	switch m.Action {
	case Subscribe:
		h.Subscriptions.add(m.URL, conn)
		log.Println("[Spokes]: Adding subscription to", m.URL, "for", conn.User.Hash)
		log.Println("[Spokes]:", len(h.Subscriptions), "subscriptions.")

	// Request actions perform an internal GET request and send the results to
	// all subscribed clients
	case Request:
		req, err := http.NewRequest("GET", "http://localhost:"+m.Port+m.URL, nil)
		if err != nil {
			log.Println(err)
			return
		}
		// Add cookie
		req.AddCookie(conn.Cookie)

		// TODO: Replace this with a "the better way" referenced in
		// render.go and update here
		req.Header.Add("X-Requested-With", "XMLHttpRequest")

		// Make request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		// Get the body of the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		h.sendToChannel(m.URL, body)
	}
}

// Run starts the WebSocket.
func (h *hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
		case c := <-h.Unregister:
			delete(h.Connections, c)
			close(c.send)
			h.Subscriptions.remove(c)
		case r := <-h.Incoming:
			var m message
			// TODO: handle unmarshalling errors
			json.Unmarshal(r.message, &m)
			h.handleMessage(r.connection, m)
		}
	}
}
