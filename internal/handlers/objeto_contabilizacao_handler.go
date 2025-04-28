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

// ObjetoContabilizacaoHandler gerencia requisições relacionadas a objetos de contabilização
type ObjetoContabilizacaoHandler struct {
	repo         *models.ObjetoContabilizacaoRepository
	auditService *services.AuditService
}

// NewObjetoContabilizacaoHandler cria um novo handler de objetos de contabilização
func NewObjetoContabilizacaoHandler(db *sql.DB) *ObjetoContabilizacaoHandler {
	return &ObjetoContabilizacaoHandler{
		repo:         models.NewObjetoContabilizacaoRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleObjetoContabilizacao gerencia todas as requisições relacionadas a objetos de contabilização
func (h *ObjetoContabilizacaoHandler) HandleObjetoContabilizacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "objetos-contabilizacao" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getObjetoContabilizacaoByID(w, r, id)
		case http.MethodPut:
			h.updateObjetoContabilizacao(w, r, id)
		case http.MethodDelete:
			h.deleteObjetoContabilizacao(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Verificar se há um parâmetro de seguradora na URL
	if len(parts) > 3 && parts[1] == "objetos-contabilizacao" && parts[2] == "seguradora" && parts[3] != "" {
		idSeguradora, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de seguradora inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getObjetosContabilizacaoBySeguradora(w, r, idSeguradora)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getObjetosContabilizacao(w, r)
	case http.MethodPost:
		h.createObjetoContabilizacao(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getObjetosContabilizacao retorna todos os objetos de contabilização
func (h *ObjetoContabilizacaoHandler) getObjetosContabilizacao(w http.ResponseWriter, r *http.Request) {
	objetos, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar objetos de contabilização: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"OBJETOS_CONTABILIZACAO",
		"",
		fmt.Sprintf("Listados %d objetos de contabilização", len(objetos)),
	)

	json.NewEncoder(w).Encode(objetos)
}

// getObjetoContabilizacaoByID retorna um objeto de contabilização específico pelo ID
func (h *ObjetoContabilizacaoHandler) getObjetoContabilizacaoByID(w http.ResponseWriter, r *http.Request, id int64) {
	objeto, err := h.repo.GetByID(id)
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
		"OBJETO_CONTABILIZACAO",
		fmt.Sprintf("%d", id),
		"Consulta de objeto de contabilização",
	)

	json.NewEncoder(w).Encode(objeto)
}

// getObjetosContabilizacaoBySeguradora retorna objetos de contabilização de uma seguradora específica
func (h *ObjetoContabilizacaoHandler) getObjetosContabilizacaoBySeguradora(w http.ResponseWriter, r *http.Request, idSeguradora int64) {
	objetos, err := h.repo.GetBySeguradora(idSeguradora)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar objetos de contabilização por seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"OBJETOS_CONTABILIZACAO",
		fmt.Sprintf("seguradora/%d", idSeguradora),
		fmt.Sprintf("Listados %d objetos de contabilização da seguradora %d", len(objetos), idSeguradora),
	)

	json.NewEncoder(w).Encode(objetos)
}

// createObjetoContabilizacao cria um novo objeto de contabilização
func (h *ObjetoContabilizacaoHandler) createObjetoContabilizacao(w http.ResponseWriter, r *http.Request) {
	var objeto models.ObjetoContabilizacao
	if err := json.NewDecoder(r.Body).Decode(&objeto); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if objeto.ObjetoContabilizacao == "" || objeto.Descricao == "" || objeto.IdSeguradora <= 0 {
		http.Error(w, "Objeto de contabilização, descrição e ID da seguradora são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, o objeto é ativo
	if !objeto.Ativo {
		objeto.Ativo = true
	}

	if err := h.repo.Create(&objeto); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar objeto de contabilização: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"OBJETO_CONTABILIZACAO",
		fmt.Sprintf("%d", objeto.ID),
		fmt.Sprintf("Criado objeto de contabilização: %s (%s)", objeto.ObjetoContabilizacao, objeto.Descricao),
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(objeto)
}

// updateObjetoContabilizacao atualiza um objeto de contabilização existente
func (h *ObjetoContabilizacaoHandler) updateObjetoContabilizacao(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o objeto existe
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
	var objeto models.ObjetoContabilizacao
	if err := json.NewDecoder(r.Body).Decode(&objeto); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	objeto.ID = id

	// Atualizar o objeto
	if err := h.repo.Update(&objeto); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar objeto de contabilização: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"OBJETO_CONTABILIZACAO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Atualizado objeto de contabilização: %s (%s)", objeto.ObjetoContabilizacao, objeto.Descricao),
	)

	// Buscar o objeto atualizado
	updatedObjeto, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar objeto de contabilização atualizado: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedObjeto)
}

// deleteObjetoContabilizacao remove um objeto de contabilização
func (h *ObjetoContabilizacaoHandler) deleteObjetoContabilizacao(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o objeto existe
	objeto, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir o objeto
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir objeto de contabilização: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"OBJETO_CONTABILIZACAO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativado objeto de contabilização: %s (%s)", objeto.ObjetoContabilizacao, objeto.Descricao),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
