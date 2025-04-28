package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/auth"
	"github.com/KleberGoncalves1209/EstudoGo/internal/models"
	"github.com/KleberGoncalves1209/EstudoGo/internal/services"
)

// LoginRequest representa os dados de requisição de login
type LoginRequest struct {
	Login string `json:"login"`
	Senha string `json:"senha"`
}

// LoginResponse representa a resposta de login bem-sucedido
type LoginResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	Usuario      models.Usuario `json:"usuario"`
	ExpiresIn    int            `json:"expires_in"` // Tempo de expiração em segundos
}

// RefreshRequest representa os dados de requisição de refresh de token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse representa a resposta de refresh de token bem-sucedido
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Tempo de expiração em segundos
}

// AuthHandler gerencia requisições relacionadas a autenticação
type AuthHandler struct {
	repo        *models.UsuarioRepository
	auditService *services.AuditService
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		repo:        models.NewUsuarioRepository(db),
		auditService: services.NewAuditService(db),
	}
}

// HandleLogin processa requisições de login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	// Decodificar os dados da requisição
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	
	// Verificar se o login e senha foram fornecidos
	if loginReq.Login == "" || loginReq.Senha == "" {
		http.Error(w, "Login e senha são obrigatórios", http.StatusBadRequest)
		return
	}
	
	// Verificar se a conta está bloqueada
	locked, blockedUntil, err := h.auditService.IsAccountLocked(loginReq.Login)
	if err != nil {
		// Registrar erro, mas continuar para verificar as credenciais
		fmt.Printf("Erro ao verificar bloqueio de conta: %v\n", err)
	}
	
	if locked {
		// Registrar tentativa de login em conta bloqueada
		_ = h.auditService.LogLoginAttempt(r, loginReq.Login, false)
		
		// Calcular tempo restante de bloqueio
		remainingTime := blockedUntil.Sub(time.Now())
		minutes := int(remainingTime.Minutes()) + 1 // Arredondar para cima
		
		// Responder com erro
		errorMsg := fmt.Sprintf("Conta bloqueada temporariamente. Tente novamente em %d minutos.", minutes)
		http.Error(w, errorMsg, http.StatusTooManyRequests)
		
		// Registrar na auditoria
		_ = h.auditService.LogAction(
			r.Context(),
			r,
			"LOGIN_BLOCKED",
			"USER",
			loginReq.Login,
			"Tentativa de login em conta bloqueada",
		)
		
		return
	}
	
	// Verificar se excedeu o limite de tentativas de login
	exceeded, blockedUntil, err := h.auditService.CheckLoginAttempts(r, loginReq.Login)
	if err != nil {
		// Registrar erro, mas continuar para verificar as credenciais
		fmt.Printf("Erro ao verificar tentativas de login: %v\n", err)
	}
	
	if exceeded {
		// Registrar tentativa de login que excedeu o limite
		_ = h.auditService.LogLoginAttempt(r, loginReq.Login, false)
		
		// Calcular tempo de bloqueio
		remainingTime := blockedUntil.Sub(time.Now())
		minutes := int(remainingTime.Minutes()) + 1 // Arredondar para cima
		
		// Responder com erro
		errorMsg := fmt.Sprintf("Muitas tentativas de login. Conta bloqueada por %d minutos.", minutes)
		http.Error(w, errorMsg, http.StatusTooManyRequests)
		
		// Registrar na auditoria
		_ = h.auditService.LogAction(
			r.Context(),
			r,
			"LOGIN_ATTEMPTS_EXCEEDED",
			"USER",
			loginReq.Login,
			"Excedeu o limite de tentativas de login",
		)
		
		return
	}
	
	// Verificar as credenciais
	usuario, err := h.repo.VerifyPassword(loginReq.Login, loginReq.Senha)
	
	// Registrar tentativa de login
	loginSuccess := err == nil
	_ = h.auditService.LogLoginAttempt(r, loginReq.Login, loginSuccess)
	
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		
		// Registrar na auditoria
		_ = h.auditService.LogAction(
			r.Context(),
			r,
			"LOGIN_FAILED",
			"USER",
			loginReq.Login,
			"Tentativa de login com credenciais inválidas",
		)
		
		return
	}
	
	// Gerar token JWT
	accessToken, err := auth.GenerateToken(usuario.ID, usuario.Login, usuario.IdTipoPerfil)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}
	
	// Gerar refresh token
	refreshToken, err := auth.GenerateRefreshToken(usuario.ID, usuario.Login, usuario.IdTipoPerfil)
	if err != nil {
		http.Error(w, "Erro ao gerar refresh token", http.StatusInternalServerError)
		return
	}
	
	// Registrar login bem-sucedido na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"LOGIN_SUCCESS",
		"USER",
		fmt.Sprintf("%d", usuario.ID),
		"Login bem-sucedido",
	)
	
	// Preparar resposta
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Usuario:      *usuario,
		ExpiresIn:    int(auth.TokenExpiration.Seconds()),
	}
	
	// Definir cabeçalho de resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Enviar resposta
	json.NewEncoder(w).Encode(response)
}

// HandleRefresh processa requisições de refresh de token
func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	
	// Decodificar os dados da requisição
	var refreshReq RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	
	// Verificar se o token foi fornecido
	if refreshReq.RefreshToken == "" {
		http.Error(w, "Refresh token é obrigatório", http.StatusBadRequest)
		return
	}
	
	// Validar o refresh token e gerar novos tokens
	accessToken, newRefreshToken, err := auth.RefreshToken(refreshReq.RefreshToken)
	if err != nil {
		http.Error(w, "Refresh token inválido ou expirado", http.StatusUnauthorized)
		
		// Registrar na auditoria
		_ = h.auditService.LogAction(
			r.Context(),
			r,
			"TOKEN_REFRESH_FAILED",
			"AUTH",
			"",
			"Falha ao renovar token: " + err.Error(),
		)
		
		return
	}
	
	// Extrair informações do usuário do token para auditoria
	claims, _ := auth.ValidateRefreshToken(refreshReq.RefreshToken)
	
	// Registrar refresh de token na auditoria
	_ = h.auditService.LogAction(
		r.Context(),
		r,
		"TOKEN_REFRESHED",
		"AUTH",
		fmt.Sprintf("%d", claims.UserID),
		"Token renovado com sucesso",
	)
	
	// Preparar resposta
	response := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(auth.TokenExpiration.Seconds()),
	}
	
	// Definir cabeçalho de resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Enviar resposta
	json.NewEncoder(w).Encode(response)
}
