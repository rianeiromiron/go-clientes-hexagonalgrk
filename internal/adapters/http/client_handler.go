package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/grok/crudclienteshex/internal/application"
)

type ClientHandler struct {
	service *application.ClientService
}

func NewClientHandler(service *application.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

type clientRequest struct {
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
}

func (h *ClientHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/clients", h.List)
	mux.HandleFunc("POST /api/clients", h.Create)
	mux.HandleFunc("GET /api/clients/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/clients/{id}", h.Update)
	mux.HandleFunc("DELETE /api/clients/{id}", h.Delete)
}

func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	clients, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	client, err := h.service.Create(r.Context(), req.Nombre, req.Email, req.Telefono, req.Direccion)
	if err != nil {
		if errors.Is(err, application.ErrNombreRequired) || errors.Is(err, application.ErrEmailRequired) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	client, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, application.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	client, err := h.service.Update(r.Context(), id, req.Nombre, req.Email, req.Telefono, req.Direccion)
	if err != nil {
		if errors.Is(err, application.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, application.ErrNombreRequired) || errors.Is(err, application.ErrEmailRequired) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, application.ErrClientNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
