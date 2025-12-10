package node

import (
	"log"

	"github.com/gorilla/websocket"
)

type SubController struct {
	node *Node
	subscribers map[string]([]*websocket.Conn)
}

func NewSubController(node *Node) *SubController {
	return &SubController{
		subscribers: make(map[string]([]*websocket.Conn)),
		node: node,
	}
}

func (sc *SubController) Node() *Node {
	return sc.node
}

func (sc *SubController) SetNode(node *Node) {
	sc.node = node
}

func (sc *SubController) AddSubscriber(listID string, conn *websocket.Conn) {
	if sc.subscribers[listID] == nil {
		sc.subscribers[listID] = []*websocket.Conn{}
	}

	for _, existingConn := range sc.subscribers[listID] {
		if existingConn == conn {
			return // Already subscribed
		}
	}

	sc.subscribers[listID] = append(sc.subscribers[listID], conn)
}

func (sc *SubController) RemoveSubscriber(listID string, conn *websocket.Conn) {
	subscribers := sc.subscribers[listID]
	if subscribers == nil {
		return
	}

	for i, existingConn := range subscribers {
		if existingConn == conn {
			sc.subscribers[listID] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}

func (sc *SubController) NotifySubscribers(listID string, message []byte) {
	subscribers := sc.subscribers[listID]
	if subscribers == nil {
		return
	}

	for _, conn := range subscribers {
		if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			// Handle error, possibly remove subscriber
			log.Println("Error notifying subscriber:", err)
			continue
		}
	}
}