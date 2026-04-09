package nats

import (
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"log/slog"
	"strings"
	"time"
)

type NATS struct {
	Conn *nats.Conn
}

func NewNATS(url string) (*NATS, error) {
	var nc *nats.Conn
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(url, nats.Name("food-delivery-nats"))
		if err == nil {
			return &NATS{Conn: nc}, nil
		}

		log.Println("Retrying NATS connection...", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Error connecting to NATS:: %s", err)
	return nil, err
}

func (n *NATS) Pub(topic string, message []byte) error {
	err := n.Conn.Publish(topic, message)
	if err != nil {
		return err
	}
	return nil
}

func (n *NATS) Sub(topic string, clients map[string]*websocket.Conn) error {

	_, err := n.Conn.Subscribe(topic, func(msg *nats.Msg) {
		message := string(msg.Data)
		slog.Info("MESSAGE_REPLY_FROM_NATS", "RECEIVED_MESSAGE", message)
		userId, messageData := n.formatMessage(message)
		if conn, ok := clients[userId]; ok {
			err := conn.WriteMessage(websocket.TextMessage, []byte(messageData))
			if err != nil {
				log.Println("Error sending message to client:", err)
				conn.Close()
				delete(clients, userId)
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *NATS) formatMessage(message string) (userId string, messageData string) {
	parts := strings.Split(message, "|")
	result := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2) // Split into key and value
		if len(kv) == 2 {
			result[kv[0]] = kv[1] // Store in a map
		}
	}
	return result["USER_ID"], result["MESSAGE"]
}
