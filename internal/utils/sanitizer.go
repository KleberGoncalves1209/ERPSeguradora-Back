package utils

import (
	"html"
	"regexp"
	"strings"
)

// Sanitizer fornece funções para sanitizar entradas de usuário
type Sanitizer struct {
	// Expressões regulares para validação
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
	nameRegex     *regexp.Regexp
	sqlInjectionPatterns *regexp.Regexp
	xssPatterns   *regexp.Regexp
}

// NewSanitizer cria um novo sanitizador
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,50}$`),
		nameRegex:     regexp.MustCompile(`^[a-zA-Z0-9 \-'.]{2,100}$`),
		sqlInjectionPatterns: regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|UNION|INTO|EXEC|EXECUTE|FROM|WHERE|GROUP BY|HAVING|ORDER BY|--|\bOR\b|\bAND\b)`),
		xssPatterns:   regexp.MustCompile(`(?i)(<script|javascript:|on\w+\s*=|<iframe|<object|<embed|<img[^>]+\s+onerror)`),
	}
}

// SanitizeString sanitiza uma string genérica
func (s *Sanitizer) SanitizeString(input string) string {
	// Remover espaços extras no início e fim
	sanitized := strings.TrimSpace(input)
	
	// Remover múltiplos espaços
	spaceRegex := regexp.MustCompile(`\s+`)
	sanitized = spaceRegex.ReplaceAllString(sanitized, " ")
	
	// Escapar HTML
	sanitized = html.EscapeString(sanitized)
	
	return sanitized
}

// SanitizeEmail sanitiza e valida um email
func (s *Sanitizer) SanitizeEmail(email string) (string, bool) {
	// Sanitizar a string
	sanitized := s.SanitizeString(email)
	
	// Validar o formato do email
	isValid := s.emailRegex.MatchString(sanitized)
	
	return sanitized, isValid
}

// SanitizeUsername sanitiza e valida um nome de usuário
func (s *Sanitizer) SanitizeUsername(username string) (string, bool) {
	// Sanitizar a string
	sanitized := s.SanitizeString(username)
	
	// Validar o formato do nome de usuário
	isValid := s.usernameRegex.MatchString(sanitized)
	
	return sanitized, isValid
}

// SanitizeName sanitiza e valida um nome
func (s *Sanitizer) SanitizeName(name string) (string, bool) {
	// Sanitizar a string
	sanitized := s.SanitizeString(name)
	
	// Validar o formato do nome
	isValid := s.nameRegex.MatchString(sanitized)
	
	return sanitized, isValid
}

// DetectSQLInjection detecta possíveis tentativas de injeção SQL
func (s *Sanitizer) DetectSQLInjection(input string) bool {
	return s.sqlInjectionPatterns.MatchString(input)
}

// DetectXSS detecta possíveis tentativas de XSS
func (s *Sanitizer) DetectXSS(input string) bool {
	return s.xssPatterns.MatchString(input)
}

// SanitizeJSON sanitiza valores em uma string JSON
func (s *Sanitizer) SanitizeJSON(jsonStr string) string {
	// Esta é uma implementação simplificada
	// Em um ambiente de produção, você deve usar um parser JSON adequado
	
	// Sanitizar caracteres potencialmente perigosos
	sanitized := strings.ReplaceAll(jsonStr, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#x27;")
	sanitized = strings.ReplaceAll(sanitized, "/", "&#x2F;")
	
	return sanitized
}
