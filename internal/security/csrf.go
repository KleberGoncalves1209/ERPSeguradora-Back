package security

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

// CSRFToken representa um token CSRF
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFProtection implementa proteção contra CSRF
type CSRFProtection struct {
	tokens      map[string]CSRFToken
	mu          sync.Mutex
	tokenExpiry time.Duration
}

// NewCSRFProtection cria uma nova proteção CSRF
func NewCSRFProtection(tokenExpiry time.Duration) *CSRFProtection {
	return &CSRFProtection{
		tokens:      make(map[string]CSRFToken),
		tokenExpiry: tokenExpiry,
	}
}

// GenerateToken gera um novo token CSRF
func (c *CSRFProtection) GenerateToken() (string, error) {
	// Gerar 32 bytes aleatórios
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	
	// Codificar em base64
	token := base64.StdEncoding.EncodeToString(b)
	
	// Armazenar o token
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.tokens[token] = CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(c.tokenExpiry),
	}
	
	// Limpar tokens expirados
	c.cleanExpiredTokens()
	
	return token, nil
}

// ValidateToken valida um token CSRF
func (c *CSRFProtection) ValidateToken(token string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Verificar se o token existe e não expirou
	if csrfToken, exists := c.tokens[token]; exists {
		if time.Now().Before(csrfToken.ExpiresAt) {
			return true
		}
		// Remover token expirado
		delete(c.tokens, token)
	}
	
	return false
}

// cleanExpiredTokens remove tokens expirados
func (c *CSRFProtection) cleanExpiredTokens() {
	now := time.Now()
	for token, csrfToken := range c.tokens {
		if now.After(csrfToken.ExpiresAt) {
			delete(c.tokens, token)
		}
	}
}

// Middleware cria um middleware HTTP para proteção CSRF
func (c *CSRFProtection) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignorar métodos seguros (GET, HEAD, OPTIONS)
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		
		// Verificar o token CSRF
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			http.Error(w, "Token CSRF ausente", http.StatusForbidden)
			return
		}
		
		if !c.ValidateToken(token) {
			http.Error(w, "Token CSRF inválido ou expirado", http.StatusForbidden)
			return
		}
		
		// Continuar com a requisição
		next.ServeHTTP(w, r)
	})
}

// GetTokenHandler cria um handler para obter um token CSRF
func (c *CSRFProtection) GetTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := c.GenerateToken()
		if err != nil {
			http.Error(w, "Erro ao gerar token CSRF", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"csrf_token":"` + token + `"}`))
	}
}
