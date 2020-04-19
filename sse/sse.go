package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range clients` started.
const patience time.Duration = time.Second * 1

type uuID []byte

type Client struct {
	UUID           uuID
	ChannelPayload chan []byte
}

type NotificationType string

const (
	BroadcastAll        NotificationType = "BroadcastAll"
	BroadcastByHospital NotificationType = "BroadcastByHospital"
	BroadcastByUUID     NotificationType = "BroadcastByUUID"
)

func (n NotificationType) String() string {
	return string(n)
}

type Notification struct {
	UUID    uuID
	Type    NotificationType
	Payload []byte
}

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan Notification

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan Client

	// Client connections registry
	clients map[string]chan []byte
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan Notification, 1),
		newClients:     make(chan Client),
		closingClients: make(chan Client),
		clients:        make(map[string]chan []byte),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) Hub() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Make sure that the writer supports flushing.
		//
		flusher, ok := w.(http.Flusher)

		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Each connection registers its own message channel with the Broker's connections registry
		channelID := make(chan []byte)
		messageChan := Client{[]byte(vars["uuID"]), channelID}

		// Signal the broker that we have a new connection
		broker.newClients <- messageChan

		// Remove this client from the map of connected clients
		// when this handler exits.
		defer func() {
			broker.closingClients <- messageChan
		}()

		// Listen to connection close and un-register messageChan
		notify := w.(http.CloseNotifier).CloseNotify()

		for {
			select {
			case <-notify:
				return
			default:

				// Write to the ResponseWriter
				// Server Sent Events compatible
				fmt.Fprintf(w, "data: %s\n\n", <-messageChan.ChannelPayload)

				// Flush the data immediatly instead of buffering it for later.
				flusher.Flush()
			}
		}
	}
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[string(s.UUID)] = s.ChannelPayload
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, string(s.UUID))
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			switch notificationType := event.Type; notificationType {
			case BroadcastByHospital:
				fmt.Println("Boardcast by hospital")
			case BroadcastByUUID:
				fmt.Println("Boardcast by UUID")
				if clientMessageChan, ok := broker.clients[string(event.UUID)]; ok {
					clientMessageChan <- event.Payload
				}
			default:
				fmt.Println("Boardcast to all registered client")
				// We got a new event from the outside!
				// Send event to all connected clients
				for _, clientMessageChan := range broker.clients {
					select {
					case clientMessageChan <- event.Payload:
					case <-time.After(patience):
						log.Print("Skipping client.")
					}
				}
			}
		}
	}
}

// SayAll is example function for broadcast to all registered clients
func (broker *Broker) SayAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonStructure, _ := json.Marshal(map[string]string{
			"name":    "exampleName",
			"message": "exampleMessage"})

		broker.Notifier <- Notification{Type: BroadcastAll, Payload: []byte(string(jsonStructure))}

		w.Write([]byte("ok."))
	}
}

// SayByUUID is example function for broadcast by specific UUID registered clients
func (broker *Broker) SayByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		jsonStructure, _ := json.Marshal(map[string]string{
			"name":    "exampleNameByUUID",
			"message": "exampleMessageByUUID"})

		broker.Notifier <- Notification{UUID: []byte(vars["uuID"]), Type: BroadcastByUUID, Payload: []byte(string(jsonStructure))}

		w.Write([]byte("ok."))
	}
}
