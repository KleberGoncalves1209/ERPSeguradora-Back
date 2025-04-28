package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Connect estabelece uma conexão com o banco de dados MySQL
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %v", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// CreateTables cria as tabelas necessárias se não existirem
func CreateTables(db *sql.DB) error {
	// Criar tabela de tipos de perfil
	tipoPerfilQuery := `
	CREATE TABLE IF NOT EXISTS tipo_perfil (
		id_tipo_perfil INT AUTO_INCREMENT PRIMARY KEY,
		perfil VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE
	);`

	_, err := db.Exec(tipoPerfilQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela tipo_perfil: %v", err)
	}

	// Criar tabela de seguradoras
	seguradoraQuery := `
	CREATE TABLE IF NOT EXISTS seguradoras (
		id_seguradora INT AUTO_INCREMENT PRIMARY KEY,
		seguradora VARCHAR(100) NOT NULL,
		nome_abreviado VARCHAR(50),
		codigo_susep VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE
	);`

	_, err = db.Exec(seguradoraQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela seguradoras: %v", err)
	}

	// Criar tabela de usuários
	usuarioQuery := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INT AUTO_INCREMENT PRIMARY KEY,
		nome VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		login VARCHAR(50) NOT NULL UNIQUE,
		senha VARCHAR(255) NOT NULL,
		idTipoPerfil INT NOT NULL,
		idSeguradora INT NOT NULL,
		AdminERP BOOLEAN DEFAULT FALSE,
		bloqueado BOOLEAN DEFAULT FALSE,
		bloqueado_ate DATETIME NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idTipoPerfil) REFERENCES tipo_perfil(id_tipo_perfil),
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(usuarioQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela usuarios: %v", err)
	}

	// Criar tabela de tentativas de login
	loginAttemptsQuery := `
	CREATE TABLE IF NOT EXISTS login_attempts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		login VARCHAR(50) NOT NULL,
		ip_address VARCHAR(45) NOT NULL,
		success BOOLEAN NOT NULL,
		attempt_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_login (login),
		INDEX idx_ip_address (ip_address),
		INDEX idx_attempt_time (attempt_time)
	);`

	_, err = db.Exec(loginAttemptsQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela login_attempts: %v", err)
	}

	// Criar tabela de auditoria
	auditLogQuery := `
	CREATE TABLE IF NOT EXISTS audit_log (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NULL,
		username VARCHAR(50) NULL,
		action VARCHAR(100) NOT NULL,
		entity_type VARCHAR(50) NOT NULL,
		entity_id VARCHAR(50) NULL,
		details TEXT NULL,
		ip_address VARCHAR(45) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_action (action),
		INDEX idx_entity_type (entity_type),
		INDEX idx_created_at (created_at)
	);`

	_, err = db.Exec(auditLogQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela audit_log: %v", err)
	}

	// Criar tabela de eventos
	eventosQuery := `
	CREATE TABLE IF NOT EXISTS eventos (
		idCodigoEvento INT AUTO_INCREMENT PRIMARY KEY,
		Evento INT NOT NULL,
		Descricao VARCHAR(255) NOT NULL,
		idSeguradora INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(eventosQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela eventos: %v", err)
	}

	// Criar tabela de objeto contabilização
	objetoContabilizacaoQuery := `
	CREATE TABLE IF NOT EXISTS objeto_contabilizacao (
		idObjetoContabilizacao INT AUTO_INCREMENT PRIMARY KEY,
		ObjetoContabilizacao VARCHAR(100) NOT NULL,
		Descricao VARCHAR(255) NOT NULL,
		idSeguradora INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(objetoContabilizacaoQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela objeto_contabilizacao: %v", err)
	}

	// Criar tabela de objeto contabilização evento
	objetoContabilizacaoEventoQuery := `
	CREATE TABLE IF NOT EXISTS objeto_contabilizacao_evento (
		idObjetoContabilizacaoEvento INT AUTO_INCREMENT PRIMARY KEY,
		idObjetoContabilizacao INT NOT NULL,
		idCodigoEvento INT NOT NULL,
		idSeguradora INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idObjetoContabilizacao) REFERENCES objeto_contabilizacao(idObjetoContabilizacao),
		FOREIGN KEY (idCodigoEvento) REFERENCES eventos(idCodigoEvento),
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(objetoContabilizacaoEventoQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela objeto_contabilizacao_evento: %v", err)
	}

	// Criar tabela de sistema contábil
	sistemaContabilQuery := `
	CREATE TABLE IF NOT EXISTS sistema_contabil (
		idSistemaContabil INT AUTO_INCREMENT PRIMARY KEY,
		SistemaContabil VARCHAR(100) NOT NULL,
		idSeguradora INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(sistemaContabilQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela sistema_contabil: %v", err)
	}

	// Criar tabela de configuração de sistema contábil
	sistemaContabilConfigQuery := `
	CREATE TABLE IF NOT EXISTS sistema_contabil_config (
		idSistemaContabilConfig INT AUTO_INCREMENT PRIMARY KEY,
		idSistemaContabil INT NOT NULL,
		idObjetoContabilizacao INT NOT NULL,
		idCodigoEvento INT NOT NULL,
		idSeguradora INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		ativo BOOLEAN DEFAULT TRUE,
		FOREIGN KEY (idSistemaContabil) REFERENCES sistema_contabil(idSistemaContabil),
		FOREIGN KEY (idObjetoContabilizacao) REFERENCES objeto_contabilizacao(idObjetoContabilizacao),
		FOREIGN KEY (idCodigoEvento) REFERENCES eventos(idCodigoEvento),
		FOREIGN KEY (idSeguradora) REFERENCES seguradoras(id_seguradora)
	);`

	_, err = db.Exec(sistemaContabilConfigQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela sistema_contabil_config: %v", err)
	}

	return nil
}
