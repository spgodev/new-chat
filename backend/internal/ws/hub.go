package ws

import (
	"New_Chat/backend/internal/domain"
	"sync"
)

type Hub struct {
	clients map[int64]*Client
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int64]*Client)}
}

func (h *Hub) Add(userID int64, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[userID] = client
}

func (h *Hub) Remove(userID int64, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userID] == client {
		delete(h.clients, userID)
	}
}

func (h *Hub) Get(userID int64) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[userID]
	if !ok {
		return nil, false
	}
	return client, true
}

func (h *Hub) Send(userID int64, data any) error {
	client, ok := h.Get(userID)
	if !ok {
		return domain.ErrConnectionNotFound
	}

	if err := client.WriteJSON(data); err != nil {
		h.Remove(userID, client)
		_ = client.Close()
		return err
	}

	return nil
}

func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	if !ok {
		return false
	}
	return true
}
