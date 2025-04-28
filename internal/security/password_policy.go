package security

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// PasswordPolicy define a política de senhas
type PasswordPolicy struct {
	MinLength           int
	RequireUppercase    bool
	RequireLowercase    bool
	RequireNumbers      bool
	RequireSpecialChars bool
	MaxRepeatedChars    int
	DisallowCommonWords bool
	PasswordHistory     int
	ExpiryDays          int
	commonPasswords     map[string]bool
}

// NewPasswordPolicy cria uma nova política de senhas
func NewPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:           10,
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireNumbers:      true,
		RequireSpecialChars: true,
		MaxRepeatedChars:    3,
		DisallowCommonWords: true,
		PasswordHistory:     5,
		ExpiryDays:          90,
		commonPasswords:     loadCommonPasswords(),
	}
}

// ValidatePassword valida uma senha de acordo com a política
func (pp *PasswordPolicy) ValidatePassword(password string) error {
	// Verificar comprimento mínimo
	if len(password) < pp.MinLength {
		return fmt.Errorf("a senha deve ter pelo menos %d caracteres", pp.MinLength)
	}
	
	// Verificar requisitos de caracteres
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	
	if pp.RequireUppercase && !hasUpper {
		return errors.New("a senha deve conter pelo menos uma letra maiúscula")
	}
	
	if pp.RequireLowercase && !hasLower {
		return errors.New("a senha deve conter pelo menos uma letra minúscula")
	}
	
	if pp.RequireNumbers && !hasNumber {
		return errors.New("a senha deve conter pelo menos um número")
	}
	
	if pp.RequireSpecialChars && !hasSpecial {
		return errors.New("a senha deve conter pelo menos um caractere especial")
	}
	
	// Verificar caracteres repetidos
	if pp.MaxRepeatedChars > 0 {
		for i := 0; i <= len(password)-pp.MaxRepeatedChars; i++ {
			char := password[i]
			repeated := true
			
			for j := 1; j < pp.MaxRepeatedChars; j++ {
				if password[i+j] != char {
					repeated = false
					break
				}
			}
			
			if repeated {
				return fmt.Errorf("a senha não pode conter mais de %d caracteres repetidos consecutivos", pp.MaxRepeatedChars-1)
			}
		}
	}
	
	// Verificar senhas comuns
	if pp.DisallowCommonWords {
		lowerPassword := strings.ToLower(password)
		if pp.commonPasswords[lowerPassword] {
			return errors.New("a senha é muito comum e facilmente adivinhável")
		}
	}
	
	return nil
}

// IsPasswordExpired verifica se uma senha expirou
func (pp *PasswordPolicy) IsPasswordExpired(lastChanged time.Time) bool {
	if pp.ExpiryDays <= 0 {
		return false
	}
	
	expiryDate := lastChanged.AddDate(0, 0, pp.ExpiryDays)
	return time.Now().After(expiryDate)
}

// loadCommonPasswords carrega uma lista de senhas comuns
func loadCommonPasswords() map[string]bool {
	// Em uma implementação real, você carregaria de um arquivo
	// Aqui, apenas incluímos algumas das senhas mais comuns
	commonPasswords := map[string]bool{
		"123456": true,
		"password": true,
		"12345678": true,
		"qwerty": true,
		"123456789": true,
		"12345": true,
		"1234": true,
		"111111": true,
		"1234567": true,
		"dragon": true,
		"123123": true,
		"baseball": true,
		"abc123": true,
		"football": true,
		"monkey": true,
		"letmein": true,
		"shadow": true,
		"master": true,
		"666666": true,
		"qwertyuiop": true,
		"123321": true,
		"mustang": true,
		"1234567890": true,
		"michael": true,
		"654321": true,
		"superman": true,
		"1qaz2wsx": true,
		"7777777": true,
		"121212": true,
		"000000": true,
		"qazwsx": true,
		"123qwe": true,
		"killer": true,
		"trustno1": true,
		"jordan": true,
		"jennifer": true,
		"zxcvbnm": true,
		"asdfgh": true,
		"hunter": true,
		"buster": true,
		"soccer": true,
		"harley": true,
		"batman": true,
		"andrew": true,
		"tigger": true,
		"sunshine": true,
		"iloveyou": true,
		"2000": true,
		"charlie": true,
		"robert": true,
		"thomas": true,
		"hockey": true,
		"ranger": true,
		"daniel": true,
		"starwars": true,
		"klaster": true,
		"112233": true,
		"george": true,
		"computer": true,
		"michelle": true,
		"jessica": true,
		"pepper": true,
		"1111": true,
		"zxcvbn": true,
		"555555": true,
		"11111111": true,
		"131313": true,
		"freedom": true,
		"777777": true,
		"pass": true,
		"maggie": true,
		"159753": true,
		"aaaaaa": true,
		"ginger": true,
		"princess": true,
		"joshua": true,
		"cheese": true,
		"amanda": true,
		"summer": true,
		"love": true,
		"ashley": true,
		"nicole": true,
		"chelsea": true,
		"biteme": true,
		"matthew": true,
		"access": true,
		"yankees": true,
		"987654321": true,
		"dallas": true,
		"austin": true,
		"thunder": true,
		"taylor": true,
		"matrix": true,
		"mobilemail": true,
		"mom": true,
		"monitor": true,
		"monitoring": true,
		"montana": true,
		"moon": true,
		"moscow": true,
	}
	
	return commonPasswords
}
