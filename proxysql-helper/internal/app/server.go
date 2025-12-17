package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"proxysql-galera-app/internal/httpx"
	"proxysql-galera-app/internal/model"
	"proxysql-galera-app/internal/repository"
)

type Server struct{ writer repository.Writer }

func NewServer(w repository.Writer) *Server { return &Server{writer: w} }

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/todos", s.handleCreateTodo)
	return httpx.LogRequests(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	t := &model.Todo{Title: req.Title, CreatedAt: time.Now()}
	if err := s.writer.CreateTodo(ctx, t); err != nil {
		http.Error(w, "db error", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}
