package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/KleberGoncalves1209/EstudoGo/internal/auth"
)

// Chaves para o contexto
type contextKey string

const (
	UserIDKey       contextKey = "user_id"
	UsernameKey     contextKey = "username"
	TipoPerfilIDKey contextKey = "tipo_perfil_id"
	
	// Cabeçalhos para rotação de token
	HeaderNewToken        = "X-New-Access-Token"
	HeaderNewRefreshToken = "X-New-Refresh-Token"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter o token do cabeçalho Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Autorização necessária", http.StatusUnauthorized)
			return
		}
		
		// O token deve estar no formato "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de autorização inválido", http.StatusUnauthorized)
			return
		}
		
		tokenString := parts[1]
		
		// Validar o token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
			return
		}
		
		// Adicionar informações do usuário ao contexto da requisição
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, TipoPerfilIDKey, claims.TipoPerfilID)
		
		// Verificar se o token precisa ser renovado
		if auth.ShouldRefreshToken(claims) {
			// Gerar um novo token
			newToken, err := auth.GenerateToken(claims.UserID, claims.Username, claims.TipoPerfilID)
			if err == nil {
				// Adicionar o novo token ao cabeçalho da resposta
				w.Header().Add(HeaderNewToken, newToken)
			}
		}
		
		// Chamar o próximo handler com o contexto atualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin verifica se o usuário tem perfil de administrador
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter o tipo de perfil do contexto
		tipoPerfilID, ok := r.Context().Value(TipoPerfilIDKey).(int)
		if !ok {
			http.Error(w, "Erro ao obter tipo de perfil", http.StatusInternalServerError)
			return
		}
		
		// Verificar se é administrador (assumindo que o ID 1 é o perfil de administrador)
		// Em uma implementação real, você pode ter uma lista de perfis com permissões
		if tipoPerfilID != 1 {
			http.Error(w, "Acesso negado: permissão de administrador necessária", http.StatusForbidden)
			return
		}
		
		// Chamar o próximo handler
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext obtém o ID do usuário do contexto
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUsernameFromContext obtém o nome de usuário do contexto
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// GetTipoPerfilIDFromContext obtém o ID do tipo de perfil do contexto
func GetTipoPerfilIDFromContext(ctx context.Context) (int, bool) {
	tipoPerfilID, ok := ctx.Value(TipoPerfilIDKey).(int)
	return tipoPerfilID, ok
}

// GetClientIP obtém o endereço IP do cliente
func GetClientIP(r *http.Request) string {
	// Tentar obter o IP real se estiver atrás de um proxy
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
