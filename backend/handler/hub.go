package handler

import "sync"

type Hub struct {
	mu          sync.RWMutex
	clients     map[string]map[chan []byte]struct{}
	broadcast   chan broadcastMsg
	subscribe   chan subscription
	unsubscribe chan subscription
}

type broadcastMsg struct {
	binID string
	data  []byte
}

type subscription struct {
	binID string
	ch    chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]map[chan []byte]struct{}),
		broadcast:   make(chan broadcastMsg, 64),
		subscribe:   make(chan subscription, 16),
		unsubscribe: make(chan subscription, 16),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.subscribe:
			h.mu.Lock()
			if h.clients[sub.binID] == nil {
				h.clients[sub.binID] = make(map[chan []byte]struct{})
			}
			h.clients[sub.binID][sub.ch] = struct{}{}
			h.mu.Unlock()

		case sub := <-h.unsubscribe:
			h.mu.Lock()
			if clients, ok := h.clients[sub.binID]; ok {
				delete(clients, sub.ch)
				if len(clients) == 0 {
					delete(h.clients, sub.binID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			targets := make([]chan []byte, 0, len(h.clients[msg.binID]))
			for ch := range h.clients[msg.binID] {
				targets = append(targets, ch)
			}
			h.mu.RUnlock()
			for _, ch := range targets {
				select {
				case ch <- msg.data:
				default:
				}
			}
		}
	}
}

func (h *Hub) Subscribe(binID string) chan []byte {
	ch := make(chan []byte, 16)
	h.subscribe <- subscription{binID: binID, ch: ch}
	return ch
}

func (h *Hub) Unsubscribe(binID string, ch chan []byte) {
	h.unsubscribe <- subscription{binID: binID, ch: ch}
}

func (h *Hub) Broadcast(binID string, data []byte) {
	select {
	case h.broadcast <- broadcastMsg{binID: binID, data: data}:
	default:
	}
}
