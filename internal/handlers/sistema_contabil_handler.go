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

// SistemaContabilHandler gerencia requisições relacionadas a sistemas contábeis
type SistemaContabilHandler struct {
	repo         *models.SistemaContabilRepository
	auditService *services.AuditService
}

// NewSistemaContabilHandler cria um novo handler de sistemas contábeis
func NewSistemaContabilHandler(db *sql.DB) *SistemaContabilHandler {
	return &SistemaContabilHandler{
		repo:         models.NewSistemaContabilRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleSistemaContabil gerencia todas as requisições relacionadas a sistemas contábeis
func (h *SistemaContabilHandler) HandleSistemaContabil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "sistemas-contabeis" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getSistemaContabilByID(w, r, id)
		case http.MethodPut:
			h.updateSistemaContabil(w, r, id)
		case http.MethodDelete:
			h.deleteSistemaContabil(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Verificar se há um parâmetro de seguradora na URL
	if len(parts) > 3 && parts[1] == "sistemas-contabeis" && parts[2] == "seguradora" && parts[3] != "" {
		idSeguradora, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de seguradora inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getSistemasContabeisBySeguradora(w, r, idSeguradora)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getSistemasContabeis(w, r)
	case http.MethodPost:
		h.createSistemaContabil(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getSistemasContabeis retorna todos os sistemas contábeis
func (h *SistemaContabilHandler) getSistemasContabeis(w http.ResponseWriter, r *http.Request) {
	sistemas, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar sistemas contábeis: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"SISTEMAS_CONTABEIS",
		"",
		fmt.Sprintf("Listados %d sistemas contábeis", len(sistemas)),
	)

	json.NewEncoder(w).Encode(sistemas)
}

// getSistemaContabilByID retorna um sistema contábil específico pelo ID
func (h *SistemaContabilHandler) getSistemaContabilByID(w http.ResponseWriter, r *http.Request, id int64) {
	sistema, err := h.repo.GetByID(id)
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
		"SISTEMA_CONTABIL",
		fmt.Sprintf("%d", id),
		"Consulta de sistema contábil",
	)

	json.NewEncoder(w).Encode(sistema)
}

// getSistemasContabeisBySeguradora retorna sistemas contábeis de uma seguradora específica
func (h *SistemaContabilHandler) getSistemasContabeisBySeguradora(w http.ResponseWriter, r *http.Request, idSeguradora int64) {
	sistemas, err := h.repo.GetBySeguradora(idSeguradora)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar sistemas contábeis por seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"SISTEMAS_CONTABEIS",
		fmt.Sprintf("seguradora/%d", idSeguradora),
		fmt.Sprintf("Listados %d sistemas contábeis da seguradora %d", len(sistemas), idSeguradora),
	)

	json.NewEncoder(w).Encode(sistemas)
}

// createSistemaContabil cria um novo sistema contábil
func (h *SistemaContabilHandler) createSistemaContabil(w http.ResponseWriter, r *http.Request) {
	var sistema models.SistemaContabil
	if err := json.NewDecoder(r.Body).Decode(&sistema); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if sistema.SistemaContabil == "" || sistema.IdSeguradora <= 0 {
		http.Error(w, "Nome do sistema contábil e ID da seguradora são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, o sistema é ativo
	if !sistema.Ativo {
		sistema.Ativo = true
	}

	if err := h.repo.Create(&sistema); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar sistema contábil: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"SISTEMA_CONTABIL",
		fmt.Sprintf("%d", sistema.ID),
		fmt.Sprintf("Criado sistema contábil: %s", sistema.SistemaContabil),
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sistema)
}

// updateSistemaContabil atualiza um sistema contábil existente
func (h *SistemaContabilHandler) updateSistemaContabil(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o sistema existe
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
	var sistema models.SistemaContabil
	if err := json.NewDecoder(r.Body).Decode(&sistema); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	sistema.ID = id

	// Atualizar o sistema
	if err := h.repo.Update(&sistema); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar sistema contábil: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"SISTEMA_CONTABIL",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Atualizado sistema contábil: %s", sistema.SistemaContabil),
	)

	// Buscar o sistema atualizado
	updatedSistema, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar sistema contábil atualizado: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedSistema)
}

// deleteSistemaContabil remove um sistema contábil
func (h *SistemaContabilHandler) deleteSistemaContabil(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o sistema existe
	sistema, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir o sistema
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir sistema contábil: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"SISTEMA_CONTABIL",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativado sistema contábil: %s", sistema.SistemaContabil),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
