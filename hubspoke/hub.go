// Package hubspoke is a basic WebSocket publish / subscribe toolkit. It builds
// off of gorilla/websocket's chat example. Clients can subscribe to URLs and
// receive updates when they change.
package hubspoke

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type hub struct {
	connections   map[*connection]bool // Registered connections
	incoming      chan wsRequest       // Inbound messages from the connections
	register      chan *connection     // Register requests from the connections
	unregister    chan *connection     // Unregister requests from connections
	subscriptions subscriptions
}

// Hub maintains the set of active connections and broadcasts messages to
// subscribers. Use the Run function to start the Hub (i.e go hubspoke.Hub.Run())
var Hub = hub{
	connections:   make(map[*connection]bool),
	incoming:      make(chan wsRequest),
	register:      make(chan *connection),
	unregister:    make(chan *connection),
	subscriptions: subscriptions{},
}

type subscriptions map[string]map[*connection]bool

func (s subscriptions) add(channel string, c *connection) {
	if _, ok := s[channel]; !ok {
		s[channel] = make(map[*connection]bool)
	}
	s[channel][c] = true
}

func (s subscriptions) remove(c *connection) {
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
	url    string
	Action action
}

type wsResponse struct {
	Channel string
	Data    *json.RawMessage
}

func (h *hub) sendToChannel(channel string, message []byte) {
	// Only send to clients subscribed to the request url
	subscriptions, ok := h.subscriptions[channel]
	if !ok {
		log.Printf("No subscriptions for '%s'\n", channel)
		return
	}
	toWrap := json.RawMessage(message)
	m := wsResponse{
		Channel: channel,
		Data:    &toWrap,
	}
	body, err := json.Marshal(m)
	if err != nil {
		log.Println("Error marshalling channel message:", err)
	}
	for c := range subscriptions {
		select {
		case c.send <- body:
		default:
			close(c.send)
			delete(subscriptions, c)
			delete(h.connections, c)
		}
	}
}

func (h *hub) handleMessage(conn *connection, m message) {
	switch m.Action {
	case Subscribe:
		h.subscriptions.add(m.url, conn)
		log.Println("Adding subscription to", m.url, "for", conn.User.Hash)

	// Request actions perform an internal GET request and send the results to
	// all subscribed clients
	case Request:
		req, err := http.NewRequest("GET", "http://localhost:8080"+m.url, nil)
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
		h.sendToChannel(m.url, body)
	}
}

// Run starts the WebSocket.
func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
			h.subscriptions.remove(c)
		case r := <-h.incoming:
			var m message
			// TODO: handle unmarshalling errors
			json.Unmarshal(r.message, &m)
			h.handleMessage(r.connection, m)
		}
	}
}
