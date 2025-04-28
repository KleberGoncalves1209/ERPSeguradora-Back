package handlers

import (
	"net/http"
	"strings"

	"github.com/KleberGoncalves1209/EstudoGo/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger configura as rotas para a documentação Swagger
func SetupSwagger() http.Handler {
	// Inicializar a documentação Swagger
	docs.SwaggerInfo.Title = "API de Gerenciamento de Seguradoras"
	docs.SwaggerInfo.Description = "API para gerenciamento de usuários, tipos de perfil, seguradoras e outros recursos"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Adicionar um parâmetro de versão para evitar problemas de cache
	return httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json?v=1.0"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	)
}

// SwaggerHandler retorna o handler para a documentação Swagger
func SwaggerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir acesso ao arquivo doc.json
		if strings.HasSuffix(r.URL.Path, "doc.json") {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		// Servir a documentação Swagger
		SetupSwagger().ServeHTTP(w, r)
	})
}
