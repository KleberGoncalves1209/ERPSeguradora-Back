package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config armazena as configurações da aplicação
type Config struct {
	DatabaseURL string
	ServerPort  int
}

// Load carrega as configurações da aplicação
func Load() (*Config, error) {
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configuração do banco de dados
	dbUser := getEnv("DB_USER", "klebe351_kleberGo")
	dbPass := getEnv("DB_PASS", "D05m09@123")
	dbHost := getEnv("DB_HOST", "br38.hostgator.com.br")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "klebe351_portalSeguradora")

	// Montar string de conexão MySQL
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Porta do servidor
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("porta do servidor inválida: %v", err)
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  serverPort,
	}, nil
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
