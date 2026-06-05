package handler

import "github.com/ready4god2513/requestbin/store"

type Handler struct {
	store *store.Store
	hub   *Hub
}

func New(s *store.Store, hub *Hub) *Handler {
	return &Handler{store: s, hub: hub}
}
