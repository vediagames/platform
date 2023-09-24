package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/vediagames/platform/session/domain"
)

type CreateResponse struct {
	ID         string    `json:"id"`
	IP         string    `json:"ip"`
	Device     string    `json:"device"`
	PageURL    string    `json:"page_url"`
	CreatedAt  time.Time `json:"created_at"`
	InsertedAt time.Time `json:"inserted_at"`
}

type CreateRequest struct {
	IP        string `json:"ip"`
	Device    string `json:"device"`
	PageURL   string `json:"page_url"`
	CreatedAt string `json:"created_at"`
}

func CreateHandler(s domain.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to decode: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to parse: %s", err)
			http.Error(w, domain.ErrInvalidCreatedAt.Error(), http.StatusBadRequest)
			return
		}

		res, err := s.Create(r.Context(), domain.CreateRequest{
			IP:        domain.IP(req.IP),
			Device:    domain.Device(req.Device),
			PageURL:   req.PageURL,
			CreatedAt: createdAt,
		})
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to create: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonRes, err := json.Marshal(CreateResponse{
			ID:         res.Session.ID,
			IP:         res.Session.IP.String(),
			Device:     res.Session.Device.String(),
			PageURL:    res.Session.PageURL,
			CreatedAt:  res.Session.CreatedAt,
			InsertedAt: res.Session.InsertedAt,
		})
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to marshal: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write(jsonRes); err != nil {
			zerolog.Ctx(r.Context()).Error().Msgf("failed to write: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
