package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KleberGoncalves1209/EstudoGo/internal/models"
)

// SeguradoraHandler gerencia requisições relacionadas a seguradoras
type SeguradoraHandler struct {
	repo *models.SeguradoraRepository
}

// NewSeguradoraHandler cria um novo handler de seguradoras
func NewSeguradoraHandler(db *sql.DB) *SeguradoraHandler {
	return &SeguradoraHandler{
		repo: models.NewSeguradoraRepository(db),
	}
}

// HandleSeguradora gerencia todas as requisições relacionadas a seguradoras
func (h *SeguradoraHandler) HandleSeguradora(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "seguradoras" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getSeguradoraByID(w, id)
		case http.MethodPut:
			h.updateSeguradora(w, r, id)
		case http.MethodDelete:
			h.deleteSeguradora(w, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getSeguradoras(w, r)
	case http.MethodPost:
		h.createSeguradora(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getSeguradoras retorna todas as seguradoras
func (h *SeguradoraHandler) getSeguradoras(w http.ResponseWriter, r *http.Request) {
	seguradoras, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar seguradoras: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(seguradoras)
}

// getSeguradoraByID retorna uma seguradora específica pelo ID
func (h *SeguradoraHandler) getSeguradoraByID(w http.ResponseWriter, id int64) {
	seguradora, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(seguradora)
}

// createSeguradora cria uma nova seguradora
func (h *SeguradoraHandler) createSeguradora(w http.ResponseWriter, r *http.Request) {
	var seguradora models.Seguradora
	if err := json.NewDecoder(r.Body).Decode(&seguradora); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if seguradora.Nome == "" {
		http.Error(w, "Nome da seguradora é obrigatório", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, a seguradora é ativa
	if !seguradora.Ativo {
		seguradora.Ativo = true
	}

	if err := h.repo.Create(&seguradora); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(seguradora)
}

// updateSeguradora atualiza uma seguradora existente
func (h *SeguradoraHandler) updateSeguradora(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se a seguradora existe
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
	var seguradora models.Seguradora
	if err := json.NewDecoder(r.Body).Decode(&seguradora); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	seguradora.ID = id

	// Atualizar a seguradora
	if err := h.repo.Update(&seguradora); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Buscar a seguradora atualizada
	updatedSeguradora, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar seguradora atualizada: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedSeguradora)
}

// deleteSeguradora remove uma seguradora
func (h *SeguradoraHandler) deleteSeguradora(w http.ResponseWriter, id int64) {
	// Verificar se a seguradora existe
	_, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir a seguradora
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir seguradora: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
