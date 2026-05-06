package messages

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
)

type Hub struct {
	mu          sync.RWMutex
	connections map[uint]map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[uint]map[*websocket.Conn]bool),
	}
}

func (h *Hub) Add(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		h.connections[userID] = make(map[*websocket.Conn]bool)
	}

	h.connections[userID][conn] = true
}

func (h *Hub) Remove(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		return
	}

	delete(h.connections[userID], conn)

	if len(h.connections[userID]) == 0 {
		delete(h.connections, userID)
	}
}

func (h *Hub) SendToUser(ctx context.Context, userID uint, payload any) error {
	h.mu.RLock()
	userConnections := h.connections[userID]
	h.mu.RUnlock()

	if len(userConnections) == 0 {
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	for conn := range userConnections {
		if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
			h.Remove(userID, conn)
			_ = conn.Close(websocket.StatusInternalError, "write failed")
		}
	}

	return nil
}
