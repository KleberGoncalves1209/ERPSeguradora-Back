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

// TipoPerfilHandler gerencia requisições relacionadas a tipos de perfil
type TipoPerfilHandler struct {
	repo *models.TipoPerfilRepository
}

// NewTipoPerfilHandler cria um novo handler de tipos de perfil
func NewTipoPerfilHandler(db *sql.DB) *TipoPerfilHandler {
	return &TipoPerfilHandler{
		repo: models.NewTipoPerfilRepository(db),
	}
}

// HandleTipoPerfil gerencia todas as requisições relacionadas a tipos de perfil
func (h *TipoPerfilHandler) HandleTipoPerfil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "tipos-perfil" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getTipoPerfilByID(w, id)
		case http.MethodPut:
			h.updateTipoPerfil(w, r, id)
		case http.MethodDelete:
			h.deleteTipoPerfil(w, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getTiposPerfil(w, r)
	case http.MethodPost:
		h.createTipoPerfil(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getTiposPerfil retorna todos os tipos de perfil
func (h *TipoPerfilHandler) getTiposPerfil(w http.ResponseWriter, r *http.Request) {
	tiposPerfil, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar tipos de perfil: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tiposPerfil)
}

// getTipoPerfilByID retorna um tipo de perfil específico pelo ID
func (h *TipoPerfilHandler) getTipoPerfilByID(w http.ResponseWriter, id int64) {
	tipoPerfil, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(tipoPerfil)
}

// createTipoPerfil cria um novo tipo de perfil
func (h *TipoPerfilHandler) createTipoPerfil(w http.ResponseWriter, r *http.Request) {
	var tipoPerfil models.TipoPerfil
	if err := json.NewDecoder(r.Body).Decode(&tipoPerfil); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if tipoPerfil.Perfil == "" {
		http.Error(w, "Nome do perfil é obrigatório", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, o tipo de perfil é ativo
	if !tipoPerfil.Ativo {
		tipoPerfil.Ativo = true
	}

	if err := h.repo.Create(&tipoPerfil); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar tipo de perfil: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tipoPerfil)
}

// updateTipoPerfil atualiza um tipo de perfil existente
func (h *TipoPerfilHandler) updateTipoPerfil(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o tipo de perfil existe
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
	var tipoPerfil models.TipoPerfil
	if err := json.NewDecoder(r.Body).Decode(&tipoPerfil); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	tipoPerfil.ID = id

	// Atualizar o tipo de perfil
	if err := h.repo.Update(&tipoPerfil); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar tipo de perfil: %v", err), http.StatusInternalServerError)
		return
	}

	// Buscar o tipo de perfil atualizado
	updatedTipoPerfil, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar tipo de perfil atualizado: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTipoPerfil)
}

// deleteTipoPerfil remove um tipo de perfil
func (h *TipoPerfilHandler) deleteTipoPerfil(w http.ResponseWriter, id int64) {
	// Verificar se o tipo de perfil existe
	_, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir o tipo de perfil
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir tipo de perfil: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
