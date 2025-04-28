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

// ObjetoContabilizacaoEventoHandler gerencia requisições relacionadas a relações entre objetos de contabilização e eventos
type ObjetoContabilizacaoEventoHandler struct {
	repo         *models.ObjetoContabilizacaoEventoRepository
	auditService *services.AuditService
}

// NewObjetoContabilizacaoEventoHandler cria um novo handler de relações entre objetos de contabilização e eventos
func NewObjetoContabilizacaoEventoHandler(db *sql.DB) *ObjetoContabilizacaoEventoHandler {
	return &ObjetoContabilizacaoEventoHandler{
		repo:         models.NewObjetoContabilizacaoEventoRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleObjetoContabilizacaoEvento gerencia todas as requisições relacionadas a relações entre objetos de contabilização e eventos
func (h *ObjetoContabilizacaoEventoHandler) HandleObjetoContabilizacaoEvento(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "objetos-contabilizacao-eventos" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getObjetoContabilizacaoEventoByID(w, r, id)
		case http.MethodPut:
			h.updateObjetoContabilizacaoEvento(w, r, id)
		case http.MethodDelete:
			h.deleteObjetoContabilizacaoEvento(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Verificar se há um parâmetro de seguradora na URL
	if len(parts) > 3 && parts[1] == "objetos-contabilizacao-eventos" && parts[2] == "seguradora" && parts[3] != "" {
		idSeguradora, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de seguradora inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getObjetosContabilizacaoEventosBySeguradora(w, r, idSeguradora)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getObjetosContabilizacaoEventos(w, r)
	case http.MethodPost:
		h.createObjetoContabilizacaoEvento(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getObjetosContabilizacaoEventos retorna todas as relações entre objetos de contabilização e eventos
func (h *ObjetoContabilizacaoEventoHandler) getObjetosContabilizacaoEventos(w http.ResponseWriter, r *http.Request) {
	relacoes, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar relações: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"OBJETOS_CONTABILIZACAO_EVENTOS",
		"",
		fmt.Sprintf("Listadas %d relações entre objetos de contabilização e eventos", len(relacoes)),
	)

	json.NewEncoder(w).Encode(relacoes)
}

// getObjetoContabilizacaoEventoByID retorna uma relação específica pelo ID
func (h *ObjetoContabilizacaoEventoHandler) getObjetoContabilizacaoEventoByID(w http.ResponseWriter, r *http.Request, id int64) {
	relacao, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
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
		"OBJETO_CONTABILIZACAO_EVENTO",
		fmt.Sprintf("%d", id),
		"Consulta de relação entre objeto de contabilização e evento",
	)

	json.NewEncoder(w).Encode(relacao)
}

// getObjetosContabilizacaoEventosBySeguradora retorna relações de uma seguradora específica
func (h *ObjetoContabilizacaoEventoHandler) getObjetosContabilizacaoEventosBySeguradora(w http.ResponseWriter, r *http.Request, idSeguradora int64) {
	relacoes, err := h.repo.GetBySeguradora(idSeguradora)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar relações por seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"OBJETOS_CONTABILIZACAO_EVENTOS",
		fmt.Sprintf("seguradora/%d", idSeguradora),
		fmt.Sprintf("Listadas %d relações da seguradora %d", len(relacoes), idSeguradora),
	)

	json.NewEncoder(w).Encode(relacoes)
}

// createObjetoContabilizacaoEvento cria uma nova relação entre objeto de contabilização e evento
func (h *ObjetoContabilizacaoEventoHandler) createObjetoContabilizacaoEvento(w http.ResponseWriter, r *http.Request) {
	var relacao models.ObjetoContabilizacaoEvento
	if err := json.NewDecoder(r.Body).Decode(&relacao); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if relacao.IdObjetoContabilizacao <= 0 || relacao.IdCodigoEvento <= 0 || relacao.IdSeguradora <= 0 {
		http.Error(w, "ID do objeto de contabilização, ID do evento e ID da seguradora são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, a relação é ativa
	if !relacao.Ativo {
		relacao.Ativo = true
	}

	if err := h.repo.Create(&relacao); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar relação: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"OBJETO_CONTABILIZACAO_EVENTO",
		fmt.Sprintf("%d", relacao.ID),
		fmt.Sprintf("Criada relação entre objeto de contabilização %d e evento %d", relacao.IdObjetoContabilizacao, relacao.IdCodigoEvento),
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(relacao)
}

// updateObjetoContabilizacaoEvento atualiza uma relação existente
func (h *ObjetoContabilizacaoEventoHandler) updateObjetoContabilizacaoEvento(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se a relação existe
	_, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Decodificar os dados da requisição
	var relacao models.ObjetoContabilizacaoEvento
	if err := json.NewDecoder(r.Body).Decode(&relacao); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	relacao.ID = id

	// Atualizar a relação
	if err := h.repo.Update(&relacao); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar relação: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"OBJETO_CONTABILIZACAO_EVENTO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Atualizada relação entre objeto de contabilização %d e evento %d", relacao.IdObjetoContabilizacao, relacao.IdCodigoEvento),
	)

	// Buscar a relação atualizada
	updatedRelacao, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar relação atualizada: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedRelacao)
}

// deleteObjetoContabilizacaoEvento remove uma relação
func (h *ObjetoContabilizacaoEventoHandler) deleteObjetoContabilizacaoEvento(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se a relação existe
	relacao, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir a relação
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir relação: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"OBJETO_CONTABILIZACAO_EVENTO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativada relação entre objeto de contabilização %d e evento %d", relacao.IdObjetoContabilizacao, relacao.IdCodigoEvento),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
