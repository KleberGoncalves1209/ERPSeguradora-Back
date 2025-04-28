package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/middleware"
)

// AuditService gerencia o registro de ações de auditoria
type AuditService struct {
	DB *sql.DB
}

// NewAuditService cria um novo serviço de auditoria
func NewAuditService(db *sql.DB) *AuditService {
	return &AuditService{
		DB: db,
	}
}

// LogAction registra uma ação no log de auditoria
func (s *AuditService) LogAction(ctx context.Context, r *http.Request, action, entityType, entityID string, details string) error {
	// Obter informações do usuário do contexto, se disponíveis
	var userID int64
	var username string
	
	userIDFromCtx, ok := middleware.GetUserIDFromContext(ctx)
	if ok {
		userID = userIDFromCtx
	}
	
	usernameFromCtx, ok := middleware.GetUsernameFromContext(ctx)
	if ok {
		username = usernameFromCtx
	}
	
	// Obter o endereço IP do cliente
	ipAddress := getIPAddress(r)
	
	// Inserir o registro de auditoria
	query := `
	INSERT INTO audit_log 
	(user_id, username, action, entity_type, entity_id, details, ip_address) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	var userIDNullable sql.NullInt64
	if userID > 0 {
		userIDNullable.Int64 = userID
		userIDNullable.Valid = true
	}
	
	var usernameNullable sql.NullString
	if username != "" {
		usernameNullable.String = username
		usernameNullable.Valid = true
	}
	
	var entityIDNullable sql.NullString
	if entityID != "" {
		entityIDNullable.String = entityID
		entityIDNullable.Valid = true
	}
	
	var detailsNullable sql.NullString
	if details != "" {
		detailsNullable.String = details
		detailsNullable.Valid = true
	}
	
	_, err := s.DB.Exec(
		query,
		userIDNullable,
		usernameNullable,
		action,
		entityType,
		entityIDNullable,
		detailsNullable,
		ipAddress,
	)
	
	if err != nil {
		return fmt.Errorf("erro ao registrar ação de auditoria: %v", err)
	}
	
	return nil
}

// LogLoginAttempt registra uma tentativa de login
func (s *AuditService) LogLoginAttempt(r *http.Request, login string, success bool) error {
	// Obter o endereço IP do cliente
	ipAddress := getIPAddress(r)
	
	// Inserir o registro de tentativa de login
	query := `
	INSERT INTO login_attempts 
	(login, ip_address, success) 
	VALUES (?, ?, ?)`
	
	_, err := s.DB.Exec(query, login, ipAddress, success)
	if err != nil {
		return fmt.Errorf("erro ao registrar tentativa de login: %v", err)
	}
	
	return nil
}

// CheckLoginAttempts verifica se um usuário ou IP excedeu o limite de tentativas de login
func (s *AuditService) CheckLoginAttempts(r *http.Request, login string) (bool, time.Time, error) {
	// Obter o endereço IP do cliente
	ipAddress := getIPAddress(r)
	
	// Verificar tentativas de login recentes (últimos 15 minutos)
	query := `
	SELECT COUNT(*) 
	FROM login_attempts 
	WHERE (login = ? OR ip_address = ?) 
	AND success = false 
	AND attempt_time > DATE_SUB(NOW(), INTERVAL 15 MINUTE)`
	
	var count int
	err := s.DB.QueryRow(query, login, ipAddress).Scan(&count)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("erro ao verificar tentativas de login: %v", err)
	}
	
	// Se houver mais de 5 tentativas falhas nos últimos 15 minutos, bloquear a conta
	if count >= 5 {
		// Bloquear a conta por 30 minutos
		blockUntil := time.Now().Add(30 * time.Minute)
		
		// Atualizar o status de bloqueio do usuário
		updateQuery := `
		UPDATE usuarios 
		SET bloqueado = true, bloqueado_ate = ? 
		WHERE login = ?`
		
		_, err := s.DB.Exec(updateQuery, blockUntil, login)
		if err != nil {
			return true, blockUntil, fmt.Errorf("erro ao bloquear conta: %v", err)
		}
		
		return true, blockUntil, nil
	}
	
	return false, time.Time{}, nil
}

// IsAccountLocked verifica se uma conta está bloqueada
func (s *AuditService) IsAccountLocked(login string) (bool, time.Time, error) {
	query := `
	SELECT bloqueado, bloqueado_ate 
	FROM usuarios 
	WHERE login = ?`
	
	var bloqueado bool
	var bloqueadoAte sql.NullTime
	
	err := s.DB.QueryRow(query, login).Scan(&bloqueado, &bloqueadoAte)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, time.Time{}, nil
		}
		return false, time.Time{}, fmt.Errorf("erro ao verificar status de bloqueio: %v", err)
	}
	
	// Se a conta estiver bloqueada, mas o tempo de bloqueio já passou, desbloquear
	if bloqueado && bloqueadoAte.Valid && bloqueadoAte.Time.Before(time.Now()) {
		// Desbloquear a conta
		updateQuery := `
		UPDATE usuarios 
		SET bloqueado = false, bloqueado_ate = NULL 
		WHERE login = ?`
		
		_, err := s.DB.Exec(updateQuery, login)
		if err != nil {
			return true, bloqueadoAte.Time, fmt.Errorf("erro ao desbloquear conta: %v", err)
		}
		
		return false, time.Time{}, nil
	}
	
	return bloqueado, bloqueadoAte.Time, nil
}

// getIPAddress obtém o endereço IP do cliente a partir da requisição
func getIPAddress(r *http.Request) string {
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
