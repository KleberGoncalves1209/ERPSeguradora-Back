package security

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implementa limitação de taxa de requisições por IP
type RateLimiter struct {
	requests     map[string][]time.Time
	mu           sync.Mutex
	maxRequests  int           // Número máximo de requisições permitidas
	timeWindow   time.Duration // Janela de tempo para contagem de requisições
	blockDuration time.Duration // Duração do bloqueio quando o limite é excedido
	blockedIPs    map[string]time.Time
}

// NewRateLimiter cria um novo limitador de taxa
func NewRateLimiter(maxRequests int, timeWindow, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:      make(map[string][]time.Time),
		maxRequests:   maxRequests,
		timeWindow:    timeWindow,
		blockDuration: blockDuration,
		blockedIPs:    make(map[string]time.Time),
	}
}

// IsAllowed verifica se um IP está permitido a fazer uma requisição
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Verificar se o IP está bloqueado
	if blockedUntil, blocked := rl.blockedIPs[ip]; blocked {
		// Se o tempo de bloqueio já passou, remover o bloqueio
		if time.Now().After(blockedUntil) {
			delete(rl.blockedIPs, ip)
		} else {
			return false
		}
	}
	
	// Obter o tempo atual
	now := time.Now()
	
	// Limpar requisições antigas
	if times, exists := rl.requests[ip]; exists {
		var validTimes []time.Time
		for _, t := range times {
			if now.Sub(t) <= rl.timeWindow {
				validTimes = append(validTimes, t)
			}
		}
		rl.requests[ip] = validTimes
	}
	
	// Verificar se o IP excedeu o limite
	if len(rl.requests[ip]) >= rl.maxRequests {
		// Bloquear o IP
		rl.blockedIPs[ip] = now.Add(rl.blockDuration)
		return false
	}
	
	// Adicionar a requisição atual
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// GetRemainingTime retorna o tempo restante de bloqueio para um IP
func (rl *RateLimiter) GetRemainingTime(ip string) time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if blockedUntil, blocked := rl.blockedIPs[ip]; blocked {
		if time.Now().After(blockedUntil) {
			delete(rl.blockedIPs, ip)
			return 0
		}
		return blockedUntil.Sub(time.Now())
	}
	
	return 0
}

// Middleware cria um middleware HTTP para limitação de taxa
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter o IP do cliente
		ip := getClientIP(r)
		
		// Verificar se o IP está permitido
		if !rl.IsAllowed(ip) {
			// Calcular tempo restante de bloqueio
			remainingTime := rl.GetRemainingTime(ip)
			
			// Configurar cabeçalhos de resposta
			w.Header().Set("Retry-After", remainingTime.String())
			http.Error(w, "Muitas requisições. Tente novamente mais tarde.", http.StatusTooManyRequests)
			return
		}
		
		// Continuar com a requisição
		next.ServeHTTP(w, r)
	})
}

// getClientIP obtém o endereço IP do cliente
func getClientIP(r *http.Request) string {
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
