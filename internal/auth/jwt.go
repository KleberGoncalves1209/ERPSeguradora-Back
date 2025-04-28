package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta para assinar os tokens JWT
// Em produção, isso deve ser uma variável de ambiente ou configuração segura
var jwtKey = []byte("sua_chave_secreta_muito_segura_e_longa")

// Constantes para configuração de tokens
const (
	TokenExpiration     = 24 * time.Hour    // Tempo de expiração do token
	TokenRefreshBefore  = 30 * time.Minute  // Renovar token se faltar menos de 30 minutos para expirar
	RefreshTokenExpiration = 7 * 24 * time.Hour // Tempo de expiração do refresh token (7 dias)
)

// Claims representa as claims do JWT
type Claims struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	TipoPerfilID int    `json:"tipo_perfil_id"`
	TokenType    string `json:"token_type"` // "access" ou "refresh"
	jwt.RegisteredClaims
}

// GenerateToken gera um novo token JWT para um usuário
func GenerateToken(userID int64, username string, tipoPerfilID int) (string, error) {
	return generateTokenWithType(userID, username, tipoPerfilID, "access", TokenExpiration)
}

// GenerateRefreshToken gera um novo refresh token JWT para um usuário
func GenerateRefreshToken(userID int64, username string, tipoPerfilID int) (string, error) {
	return generateTokenWithType(userID, username, tipoPerfilID, "refresh", RefreshTokenExpiration)
}

// generateTokenWithType gera um token com tipo e duração específicos
func generateTokenWithType(userID int64, username string, tipoPerfilID int, tokenType string, expiration time.Duration) (string, error) {
	// Define o tempo de expiração do token
	expirationTime := time.Now().Add(expiration)
	
	// Cria as claims
	claims := &Claims{
		UserID:       userID,
		Username:     username,
		TipoPerfilID: tipoPerfilID,
		TokenType:    tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api-seguradoras",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	
	// Cria o token com as claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna as claims
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse do token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Verifica se o token é válido
	if !token.Valid {
		return nil, errors.New("token inválido")
	}
	
	// Verifica se é um token de acesso
	if claims.TokenType != "access" {
		return nil, errors.New("tipo de token inválido")
	}
	
	return claims, nil
}

// ValidateRefreshToken valida um refresh token JWT e retorna as claims
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	// Parse do token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Verifica se o token é válido
	if !token.Valid {
		return nil, errors.New("token inválido")
	}
	
	// Verifica se é um refresh token
	if claims.TokenType != "refresh" {
		return nil, errors.New("tipo de token inválido")
	}
	
	return claims, nil
}

// ShouldRefreshToken verifica se um token deve ser renovado
func ShouldRefreshToken(claims *Claims) bool {
	// Verificar se o token expira em menos de 30 minutos
	if claims.ExpiresAt != nil {
		expiresAt := claims.ExpiresAt.Time
		refreshBefore := time.Now().Add(TokenRefreshBefore)
		return expiresAt.Before(refreshBefore)
	}
	return false
}

// RefreshToken gera um novo token a partir de um refresh token válido
func RefreshToken(refreshTokenString string) (string, string, error) {
	// Valida o refresh token
	claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}
	
	// Gera um novo token de acesso
	newAccessToken, err := GenerateToken(claims.UserID, claims.Username, claims.TipoPerfilID)
	if err != nil {
		return "", "", err
	}
	
	// Gera um novo refresh token
	newRefreshToken, err := GenerateRefreshToken(claims.UserID, claims.Username, claims.TipoPerfilID)
	if err != nil {
		return "", "", err
	}
	
	return newAccessToken, newRefreshToken, nil
}
