package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"1000words-game/internal/api/responses"
	"1000words-game/internal/config"
)

type HealthHandler struct {
	db  *sql.DB
	cfg config.Config
}

func NewHealthHandler(db *sql.DB, cfg config.Config) *HealthHandler {
	return &HealthHandler{
		db:  db,
		cfg: cfg,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"

	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "error"
	}

	responses.JSON(w, http.StatusOK, map[string]any{
		"app":      h.cfg.AppName,
		"api":      "ok",
		"database": dbStatus,
	})
}
