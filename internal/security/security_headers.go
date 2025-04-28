package security

import (
	"net/http"
)

// SecurityHeaders adiciona cabeçalhos de segurança às respostas HTTP
type SecurityHeaders struct {
	// Configurações personalizáveis
	HSTS              bool
	HSTSMaxAge        int
	HSTSIncludeSubdomains bool
	HSTSPreload       bool
	
	ContentSecurityPolicy string
	XFrameOptions     string
	XContentTypeOptions bool
	ReferrerPolicy    string
	PermissionsPolicy string
}

// NewSecurityHeaders cria uma nova configuração de cabeçalhos de segurança
func NewSecurityHeaders() *SecurityHeaders {
	return &SecurityHeaders{
		HSTS:                 true,
		HSTSMaxAge:           31536000, // 1 ano
		HSTSIncludeSubdomains: true,
		HSTSPreload:          false,
		
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline';",
		XFrameOptions:         "DENY",
		XContentTypeOptions:   true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		PermissionsPolicy:     "camera=(), microphone=(), geolocation=(), interest-cohort=()",
	}
}

// Middleware cria um middleware HTTP para adicionar cabeçalhos de segurança
func (sh *SecurityHeaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strict-Transport-Security (HSTS)
		if sh.HSTS {
			hstsValue := "max-age=" + string(sh.HSTSMaxAge)
			if sh.HSTSIncludeSubdomains {
				hstsValue += "; includeSubDomains"
			}
			if sh.HSTSPreload {
				hstsValue += "; preload"
			}
			w.Header().Set("Strict-Transport-Security", hstsValue)
		}
		
		// Content-Security-Policy
		if sh.ContentSecurityPolicy != "" {
			w.Header().Set("Content-Security-Policy", sh.ContentSecurityPolicy)
		}
		
		// X-Frame-Options
		if sh.XFrameOptions != "" {
			w.Header().Set("X-Frame-Options", sh.XFrameOptions)
		}
		
		// X-Content-Type-Options
		if sh.XContentTypeOptions {
			w.Header().Set("X-Content-Type-Options", "nosniff")
		}
		
		// Referrer-Policy
		if sh.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", sh.ReferrerPolicy)
		}
		
		// Permissions-Policy
		if sh.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", sh.PermissionsPolicy)
		}
		
		// Cache-Control para APIs
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		
		// Continuar com a requisição
		next.ServeHTTP(w, r)
	})
}
