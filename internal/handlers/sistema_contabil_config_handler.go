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

// SistemaContabilConfigHandler gerencia requisições relacionadas a configurações de sistema contábil
type SistemaContabilConfigHandler struct {
	repo         *models.SistemaContabilConfigRepository
	auditService *services.AuditService
}

// NewSistemaContabilConfigHandler cria um novo handler de configurações de sistema contábil
func NewSistemaContabilConfigHandler(db *sql.DB) *SistemaContabilConfigHandler {
	return &SistemaContabilConfigHandler{
		repo:         models.NewSistemaContabilConfigRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleSistemaContabilConfig gerencia todas as requisições relacionadas a configurações de sistema contábil
func (h *SistemaContabilConfigHandler) HandleSistemaContabilConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "sistemas-contabeis-config" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getSistemaContabilConfigByID(w, r, id)
		case http.MethodPut:
			h.updateSistemaContabilConfig(w, r, id)
		case http.MethodDelete:
			h.deleteSistemaContabilConfig(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Verificar se há um parâmetro de seguradora na URL
	if len(parts) > 3 && parts[1] == "sistemas-contabeis-config" && parts[2] == "seguradora" && parts[3] != "" {
		idSeguradora, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de seguradora inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getSistemasContabeisConfigBySeguradora(w, r, idSeguradora)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar se há um parâmetro de sistema contábil na URL
	if len(parts) > 3 && parts[1] == "sistemas-contabeis-config" && parts[2] == "sistema" && parts[3] != "" {
		idSistemaContabil, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "ID de sistema contábil inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.getSistemasContabeisConfigBySistemaContabil(w, r, idSistemaContabil)
			return
		}

		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getSistemasContabeisConfig(w, r)
	case http.MethodPost:
		h.createSistemaContabilConfig(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getSistemasContabeisConfig retorna todas as configurações de sistema contábil
func (h *SistemaContabilConfigHandler) getSistemasContabeisConfig(w http.ResponseWriter, r *http.Request) {
	configs, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar configurações de sistema contábil: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"SISTEMAS_CONTABEIS_CONFIG",
		"",
		fmt.Sprintf("Listadas %d configurações de sistema contábil", len(configs)),
	)

	json.NewEncoder(w).Encode(configs)
}

// getSistemaContabilConfigByID retorna uma configuração específica pelo ID
func (h *SistemaContabilConfigHandler) getSistemaContabilConfigByID(w http.ResponseWriter, r *http.Request, id int64) {
	config, err := h.repo.GetByID(id)
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
		"SISTEMA_CONTABIL_CONFIG",
		fmt.Sprintf("%d", id),
		"Consulta de configuração de sistema contábil",
	)

	json.NewEncoder(w).Encode(config)
}

// getSistemasContabeisConfigBySeguradora retorna configurações de uma seguradora específica
func (h *SistemaContabilConfigHandler) getSistemasContabeisConfigBySeguradora(w http.ResponseWriter, r *http.Request, idSeguradora int64) {
	configs, err := h.repo.GetBySeguradora(idSeguradora)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar configurações por seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"SISTEMAS_CONTABEIS_CONFIG",
		fmt.Sprintf("seguradora/%d", idSeguradora),
		fmt.Sprintf("Listadas %d configurações da seguradora %d", len(configs), idSeguradora),
	)

	json.NewEncoder(w).Encode(configs)
}

// getSistemasContabeisConfigBySistemaContabil retorna configurações de um sistema contábil específico
func (h *SistemaContabilConfigHandler) getSistemasContabeisConfigBySistemaContabil(w http.ResponseWriter, r *http.Request, idSistemaContabil int64) {
	configs, err := h.repo.GetBySistemaContabil(idSistemaContabil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar configurações por sistema contábil: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"SISTEMAS_CONTABEIS_CONFIG",
		fmt.Sprintf("sistema/%d", idSistemaContabil),
		fmt.Sprintf("Listadas %d configurações do sistema contábil %d", len(configs), idSistemaContabil),
	)

	json.NewEncoder(w).Encode(configs)
}

// createSistemaContabilConfig cria uma nova configuração de sistema contábil
func (h *SistemaContabilConfigHandler) createSistemaContabilConfig(w http.ResponseWriter, r *http.Request) {
	var config models.SistemaContabilConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if config.IdSistemaContabil <= 0 || config.IdObjetoContabilizacao <= 0 || config.IdCodigoEvento <= 0 || config.IdSeguradora <= 0 {
		http.Error(w, "ID do sistema contábil, ID do objeto de contabilização, ID do evento e ID da seguradora são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, a configuração é ativa
	if !config.Ativo {
		config.Ativo = true
	}

	if err := h.repo.Create(&config); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar configuração: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"SISTEMA_CONTABIL_CONFIG",
		fmt.Sprintf("%d", config.ID),
		fmt.Sprintf("Criada configuração para sistema contábil %d, objeto %d e evento %d", config.IdSistemaContabil, config.IdObjetoContabilizacao, config.IdCodigoEvento),
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

// updateSistemaContabilConfig atualiza uma configuração existente
func (h *SistemaContabilConfigHandler) updateSistemaContabilConfig(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se a configuração existe
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
	var config models.SistemaContabilConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	config.ID = id

	// Atualizar a configuração
	if err := h.repo.Update(&config); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar configuração: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"SISTEMA_CONTABIL_CONFIG",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Atualizada configuração para sistema contábil %d, objeto %d e evento %d", config.IdSistemaContabil, config.IdObjetoContabilizacao, config.IdCodigoEvento),
	)

	// Buscar a configuração atualizada
	updatedConfig, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar configuração atualizada: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedConfig)
}

// deleteSistemaContabilConfig remove uma configuração
func (h *SistemaContabilConfigHandler) deleteSistemaContabilConfig(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se a configuração existe
	config, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir a configuração
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir configuração: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"SISTEMA_CONTABIL_CONFIG",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativada configuração para sistema contábil %d, objeto %d e evento %d", config.IdSistemaContabil, config.IdObjetoContabilizacao, config.IdCodigoEvento),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
