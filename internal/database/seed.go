package database

import (
	"database/sql"
	"fmt"
	"log"

	// Removendo a importação não utilizada
	"golang.org/x/crypto/bcrypt"
)

// SeedInitialData insere dados iniciais no banco de dados
func SeedInitialData(db *sql.DB) error {
	// Verificar se já existem dados
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tipo_perfil").Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar dados existentes: %v", err)
	}

	// Se não houver tipos de perfil, criar os padrões
	if count == 0 {
		log.Println("Criando tipos de perfil padrão...")
		
		// Criar tipo de perfil Administrador
		_, err = db.Exec(
			"INSERT INTO tipo_perfil (perfil, ativo) VALUES (?, ?)",
			"Administrador", true,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar tipo de perfil Administrador: %v", err)
		}
		
		// Criar tipo de perfil Usuário
		_, err = db.Exec(
			"INSERT INTO tipo_perfil (perfil, ativo) VALUES (?, ?)",
			"Usuário", true,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar tipo de perfil Usuário: %v", err)
		}
	}

	// Verificar se já existem seguradoras
	err = db.QueryRow("SELECT COUNT(*) FROM seguradoras").Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar seguradoras existentes: %v", err)
	}

	// Se não houver seguradoras, criar uma padrão
	if count == 0 {
		log.Println("Criando seguradora padrão...")
		
		_, err = db.Exec(
			"INSERT INTO seguradoras (seguradora, nome_abreviado, codigo_susep, ativo) VALUES (?, ?, ?, ?)",
			"Seguradora Padrão", "SegPad", "00000", true,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar seguradora padrão: %v", err)
		}
	}

	// Verificar se já existem usuários
	err = db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar usuários existentes: %v", err)
	}

	// Se não houver usuários, criar um administrador padrão
	if count == 0 {
		log.Println("Criando usuário administrador padrão...")
		
		// Hash da senha padrão "Admin@123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash da senha: %v", err)
		}
		
		_, err = db.Exec(
			"INSERT INTO usuarios (nome, email, login, senha, idTipoPerfil, idSeguradora, AdminERP, ativo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"Administrador", "admin@sistema.com", "admin", string(hashedPassword), 1, 1, true, true,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar usuário administrador: %v", err)
		}
		
		log.Println("Usuário administrador criado com sucesso!")
		log.Println("Login: admin")
		log.Println("Senha: Admin@123")
	}

	return nil
}
