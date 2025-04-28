package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

// Error implementa a interface error
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateEmail valida o formato de um email
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email não pode ser vazio")
	}

	// Expressão regular para validação básica de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("formato de email inválido")
	}

	return nil
}

// ValidatePassword verifica se a senha atende aos requisitos mínimos
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("senha deve ter pelo menos 8 caracteres")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	}

	return nil
}

// SanitizeString remove espaços extras e caracteres indesejados
func SanitizeString(s string) string {
	// Remover espaços extras no início e fim
	s = strings.TrimSpace(s)
	
	// Remover múltiplos espaços
	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")
	
	return s
}

// ValidateLength verifica se uma string tem o comprimento dentro dos limites
func ValidateLength(field, value string, min, max int) error {
	if len(value) < min {
		return ValidationError{
			Field:   field,
			Message: fmt.Sprintf("deve ter pelo menos %d caracteres", min),
		}
	}
	
	if max > 0 && len(value) > max {
		return ValidationError{
			Field:   field,
			Message: fmt.Sprintf("não pode ter mais que %d caracteres", max),
		}
	}
	
	return nil
}

// ValidateRequired verifica se um campo obrigatório está preenchido
func ValidateRequired(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return ValidationError{
			Field:   field,
			Message: "campo obrigatório",
		}
	}
	
	return nil
}

// ValidateNumericRange verifica se um número está dentro de um intervalo
func ValidateNumericRange(field string, value, min, max int) error {
	if value < min {
		return ValidationError{
			Field:   field,
			Message: fmt.Sprintf("deve ser pelo menos %d", min),
		}
	}
	
	if max > 0 && value > max {
		return ValidationError{
			Field:   field,
			Message: fmt.Sprintf("não pode ser maior que %d", max),
		}
	}
	
	return nil
}
