package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/config"
	"github.com/KleberGoncalves1209/EstudoGo/internal/database"
	"github.com/KleberGoncalves1209/EstudoGo/internal/handlers"
	"github.com/KleberGoncalves1209/EstudoGo/internal/middleware"
	"github.com/KleberGoncalves1209/EstudoGo/internal/security"
	"github.com/KleberGoncalves1209/EstudoGo/internal/services"
)

// @title API de Gerenciamento de Seguradoras
// @version 1.0
// @description API para gerenciamento de usuários, tipos de perfil, seguradoras, eventos, objetos de contabilização e sistemas contábeis com autenticação JWT, limite de tentativas de login, rotação de tokens, auditoria e diversas medidas de segurança
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Autenticação usando JWT. Exemplo: "Bearer {token}"
func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializar conexão com o banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar conexão com o banco
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao verificar conexão com o banco: %v", err)
	}
	log.Println("Conexão com o banco de dados Hostgator estabelecida com sucesso!")

	// Criar tabelas se não existirem
	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}
	log.Println("Tabelas verificadas/criadas com sucesso!")
	
	// Inserir dados iniciais se necessário
	if err := database.SeedInitialData(db); err != nil {
		log.Fatalf("Erro ao inserir dados iniciais: %v", err)
	}
	
	// Inicializar serviço de auditoria
	auditService := services.NewAuditService(db)
	
	// Inicializar componentes de segurança
	rateLimiter := security.NewRateLimiter(60, time.Minute, 5*time.Minute)
	csrfProtection := security.NewCSRFProtection(time.Hour)
	securityHeaders := security.NewSecurityHeaders()
	
	// Criar mux para rotas
	mux := http.NewServeMux()
	
	// Configurar handlers públicos
	mux.HandleFunc("/", handlers.HomeHandler)
	
	// Rotas de autenticação (públicas, mas com rate limiting)
	authHandler := handlers.NewAuthHandler(db)
	mux.Handle("/auth/login", rateLimiter.Middleware(http.HandlerFunc(authHandler.HandleLogin)))
	mux.Handle("/auth/refresh", rateLimiter.Middleware(http.HandlerFunc(authHandler.HandleRefresh)))
	
	// Rota para obter token CSRF (protegida)
	mux.Handle("/csrf/token", middleware.AuthMiddleware(csrfProtection.GetTokenHandler()))
	
	// Rota para a documentação Swagger (pública)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Desabilitar temporariamente os cabeçalhos de segurança para o Swagger
		if strings.Contains(r.URL.Path, "swagger") {
			handlers.SwaggerHandler().ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})))
	
	// Handlers para rotas protegidas
	userHandler := handlers.NewUserHandler(db)
	tipoPerfilHandler := handlers.NewTipoPerfilHandler(db)
	seguradoraHandler := handlers.NewSeguradoraHandler(db)
	
	// Novos handlers para as novas entidades
	eventoHandler := handlers.NewEventoHandler(db)
	objetoContabilizacaoHandler := handlers.NewObjetoContabilizacaoHandler(db)
	objetoContabilizacaoEventoHandler := handlers.NewObjetoContabilizacaoEventoHandler(db)
	sistemaContabilHandler := handlers.NewSistemaContabilHandler(db)
	sistemaContabilConfigHandler := handlers.NewSistemaContabilConfigHandler(db)
	
	// Middleware para registrar todas as requisições na auditoria
	auditMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Registrar a requisição na auditoria (exceto para rotas públicas como Swagger)
			if !strings.HasPrefix(r.URL.Path, "/swagger/") {
				_ = auditService.LogAction(
					r.Context(),
					r,
					"API_REQUEST",
					"ENDPOINT",
					r.URL.Path,
					fmt.Sprintf("Método: %s", r.Method),
				)
			}
			
			// Continuar com a requisição
			next.ServeHTTP(w, r)
		})
	}
	
	// Aplicar middlewares de segurança às rotas protegidas
	secureMiddleware := func(next http.Handler) http.Handler {
		// Aplicar middlewares na ordem correta
		handler := next
		handler = csrfProtection.Middleware(handler)
		handler = middleware.AuthMiddleware(handler)
		handler = rateLimiter.Middleware(handler)
		handler = securityHeaders.Middleware(handler)
		handler = auditMiddleware(handler)
		return handler
	}
	
	// Rotas para usuários (protegidas)
	mux.Handle("/usuarios/", secureMiddleware(http.HandlerFunc(userHandler.HandleUsers)))
	mux.Handle("/usuarios", secureMiddleware(http.HandlerFunc(userHandler.HandleUsers)))
	
	// Rotas para tipos de perfil (protegidas)
	mux.Handle("/tipos-perfil/", secureMiddleware(http.HandlerFunc(tipoPerfilHandler.HandleTipoPerfil)))
	mux.Handle("/tipos-perfil", secureMiddleware(http.HandlerFunc(tipoPerfilHandler.HandleTipoPerfil)))
	
	// Rotas para seguradoras (protegidas)
	mux.Handle("/seguradoras/", secureMiddleware(http.HandlerFunc(seguradoraHandler.HandleSeguradora)))
	mux.Handle("/seguradoras", secureMiddleware(http.HandlerFunc(seguradoraHandler.HandleSeguradora)))
	
	// Rotas para eventos (protegidas)
	mux.Handle("/eventos/", secureMiddleware(http.HandlerFunc(eventoHandler.HandleEvento)))
	mux.Handle("/eventos", secureMiddleware(http.HandlerFunc(eventoHandler.HandleEvento)))
	
	// Rotas para objetos de contabilização (protegidas)
	mux.Handle("/objetos-contabilizacao/", secureMiddleware(http.HandlerFunc(objetoContabilizacaoHandler.HandleObjetoContabilizacao)))
	mux.Handle("/objetos-contabilizacao", secureMiddleware(http.HandlerFunc(objetoContabilizacaoHandler.HandleObjetoContabilizacao)))
	
	// Rotas para relações entre objetos de contabilização e eventos (protegidas)
	mux.Handle("/objetos-contabilizacao-eventos/", secureMiddleware(http.HandlerFunc(objetoContabilizacaoEventoHandler.HandleObjetoContabilizacaoEvento)))
	mux.Handle("/objetos-contabilizacao-eventos", secureMiddleware(http.HandlerFunc(objetoContabilizacaoEventoHandler.HandleObjetoContabilizacaoEvento)))
	
	// Rotas para sistemas contábeis (protegidas)
	mux.Handle("/sistemas-contabeis/", secureMiddleware(http.HandlerFunc(sistemaContabilHandler.HandleSistemaContabil)))
	mux.Handle("/sistemas-contabeis", secureMiddleware(http.HandlerFunc(sistemaContabilHandler.HandleSistemaContabil)))
	
	// Rotas para configurações de sistema contábil (protegidas)
	mux.Handle("/sistemas-contabeis-config/", secureMiddleware(http.HandlerFunc(sistemaContabilConfigHandler.HandleSistemaContabilConfig)))
	mux.Handle("/sistemas-contabeis-config", secureMiddleware(http.HandlerFunc(sistemaContabilConfigHandler.HandleSistemaContabilConfig)))
	
	// Iniciar servidor HTTP
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Servidor iniciado em http://localhost%s", serverAddr)
	log.Printf("Documentação Swagger disponível em http://localhost%s/swagger/index.html", serverAddr)
	
	// Aplicar headers de segurança a todas as respostas
	secureServer := securityHeaders.Middleware(mux)
	
	log.Fatal(http.ListenAndServe(serverAddr, secureServer))
}
