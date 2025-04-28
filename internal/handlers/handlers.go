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

// HomeHandler gerencia requisições para a página inicial
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Go com MySQL - Use /usuarios para acessar a API")
}

// UserHandler gerencia requisições relacionadas a usuários
type UserHandler struct {
	repo         *models.UsuarioRepository
	auditService *services.AuditService
}

// NewUserHandler cria um novo handler de usuários
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		repo:         models.NewUsuarioRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleUsers gerencia todas as requisições relacionadas a usuários
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar se há um ID na URL para operações específicas
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 && parts[1] == "usuarios" && parts[2] != "" {
		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getUserByID(w, r, id)
		case http.MethodPut:
			h.updateUser(w, r, id)
		case http.MethodDelete:
			h.deleteUser(w, r, id)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Operações que não requerem ID específico
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// getUsers retorna todos os usuários
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar usuários: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LIST",
		"USUARIOS",
		"",
		fmt.Sprintf("Listados %d usuários", len(usuarios)),
	)

	json.NewEncoder(w).Encode(usuarios)
}

// getUserByID retorna um usuário específico pelo ID
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	usuario, err := h.repo.GetByID(id)
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
		"USUARIO",
		fmt.Sprintf("%d", id),
		"Consulta de usuário",
	)

	json.NewEncoder(w).Encode(usuario)
}

// createUser cria um novo usuário
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if usuario.Nome == "" || usuario.Email == "" || usuario.Login == "" || usuario.Senha == "" {
		http.Error(w, "Nome, email, login e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	// Por padrão, se não for especificado, o usuário é ativo
	if !usuario.Ativo {
		usuario.Ativo = true
	}

	if err := h.repo.Create(&usuario); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"CREATE",
		"USUARIO",
		fmt.Sprintf("%d", usuario.ID),
		fmt.Sprintf("Criado usuário: %s (%s)", usuario.Nome, usuario.Email),
	)

	// Não retornar a senha na resposta
	usuario.Senha = ""

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// updateUser atualiza um usuário existente
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o usuário existe
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
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garantir que o ID seja o mesmo
	usuario.ID = id

	// Atualizar o usuário
	if err := h.repo.Update(&usuario); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Se a senha foi fornecida, atualizá-la separadamente
	senhaAlterada := false
	if usuario.Senha != "" {
		if err := h.repo.UpdatePassword(id, usuario.Senha); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao atualizar senha: %v", err), http.StatusInternalServerError)
			return
		}
		senhaAlterada = true
	}

	// Registrar na auditoria
	detalhes := fmt.Sprintf("Atualizado usuário: %s (%s)", usuario.Nome, usuario.Email)
	if senhaAlterada {
		detalhes += " - Senha alterada"
	}
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"UPDATE",
		"USUARIO",
		fmt.Sprintf("%d", id),
		detalhes,
	)

	// Buscar o usuário atualizado
	updatedUser, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar usuário atualizado: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
}

// deleteUser remove um usuário
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int64) {
	// Verificar se o usuário existe
	usuario, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Excluir o usuário
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao excluir usuário: %v", err), http.StatusInternalServerError)
		return
	}

	// Registrar na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"DELETE",
		"USUARIO",
		fmt.Sprintf("%d", id),
		fmt.Sprintf("Desativado usuário: %s (%s)", usuario.Nome, usuario.Email),
	)

	// Responder com sucesso
	w.WriteHeader(http.StatusNoContent)
}
