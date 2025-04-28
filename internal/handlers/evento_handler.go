package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KleberGoncalves1209/EstudoGo/internal/models"
	"github.com/KleberGoncalves1209/EstudoGo/internal/services"
)

// EventoHandler gerencia requisições relacionadas a eventos
type EventoHandler struct {
	repo         *models.EventoRepository
	auditService *services.AuditService
}

// NewEventoHandler cria um novo handler de eventos
func NewEventoHandler(db *sql.DB) *EventoHandler {
	return &EventoHandler{
		repo:         models.NewEventoRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleEvento gerencia todas as requisições relacionadas a eventos
func (h *EventoHandler) HandleEvento(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "eventos" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getEventoByID(w, r, id)
		case http.MethodPut:
			h.updateEvento(w, r, id)
		case http.MethodDelete:
			h.deleteEvento(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Verificar se há um parâmetro de seguradora na URL
	if len(parts) > 3 && parts[1] == "eventos" && parts[2] == "seguradora" && parts[3] != "" {
		idSeguradora, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de seguradora inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getEventosBySeguradora(w, r, idSeguradora)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getEventos(w, r)
	case http.MethodPost:
		h.createEvento(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getEventos retorna todos os eventos
func (h *EventoHandler) getEventos(w http.ResponseWriter, r *http.Request) {
	eventos, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar eventos: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"EVENTOS",
		"",
		fmt.Sprintf("Listados %d eventos", len(eventos)),
	)

	json.NewEncoder(w).Encode(eventos)
}

// getEventoByID retorna um evento específico pelo ID
func (h *EventoHandler) getEventoByID(w http.ResponseWriter, r *http.Request, id int64) {
	evento, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"READ",
		"EVENTO",
		fmt.Sprintf("%d", id),
		"Consulta de evento",
	)

	json.NewEncoder(w).Encode(evento)
}

// getEventosBySeguradora retorna eventos de uma seguradora específica
func (h *EventoHandler) getEventosBySeguradora(w http.ResponseWriter, r *http.Request, idSeguradora int64) {
	eventos, err := h.repo.GetBySeguradora(idSeguradora)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar eventos por seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"EVENTOS",
		fmt.Sprintf("seguradora/%d", idSeguradora),
		fmt.Sprintf("Listados %d eventos da seguradora %d", len(eventos), idSeguradora),
	)

	json.NewEncoder(w).Encode(eventos)
}

// createEvento cria um novo evento
func (h *EventoHandler) createEvento(w http.ResponseWriter, r *http.Request) {
	var evento models.Evento
	if err := json.NewDecoder(r.Body).Decode(&evento); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if evento.Evento <= 0 || evento.Descricao == "" || evento.IdSeguradora <= 0 {
		http.Error(w, "Evento, descrição e ID da seguradora são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, o evento é ativo
	if !evento.Ativo {
		evento.Ativo = true
	}

	if err := h.repo.Create(&evento); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar evento: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"EVENTO",
		fmt.Sprintf("%d", evento.ID),
		fmt.Sprintf("Criado evento: %d (%s)", evento.Evento, evento.Descricao),
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(evento)
}

// updateEvento atualiza um evento existente
func (h *EventoHandler) updateEvento(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o evento existe
	_, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Decodificar os dados da requisição
	var evento models.Evento
	if err := json.NewDecoder(r.Body).Decode(&evento); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	evento.ID = id

	// Atualizar o evento
	if err := h.repo.Update(&evento); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar evento: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"EVENTO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Atualizado evento: %d (%s)", evento.Evento, evento.Descricao),
	)

	// Buscar o evento atualizado
	updatedEvento, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar evento atualizado: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedEvento)
}

// deleteEvento remove um evento
func (h *EventoHandler) deleteEvento(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o evento existe
	evento, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir o evento
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir evento: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"EVENTO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativado evento: %d (%s)", evento.Evento, evento.Descricao),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
