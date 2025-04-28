package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// SistemaContabilConfig representa uma configuração de sistema contábil no sistema
type SistemaContabilConfig struct {
	ID                     int64     `json:"idSistemaContabilConfig"`
	IdSistemaContabil      int64     `json:"idSistemaContabil"`
	IdObjetoContabilizacao int64     `json:"idObjetoContabilizacao"`
	IdCodigoEvento         int64     `json:"idCodigoEvento"`
	IdSeguradora           int64     `json:"idSeguradora"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	Ativo                  bool      `json:"ativo"`
	// Campos para exibição de informações relacionadas
	SistemaContabilNome       string    `json:"sistemaContabilNome,omitempty"`
	ObjetoContabilizacaoNome  string    `json:"objetoContabilizacaoNome,omitempty"`
	EventoNumero              int       `json:"eventoNumero,omitempty"`
	EventoDescricao           string    `json:"eventoDescricao,omitempty"`
}

// SistemaContabilConfigRepository gerencia operações de banco de dados para configurações de sistema contábil
type SistemaContabilConfigRepository struct {
	DB *sql.DB
}

// NewSistemaContabilConfigRepository cria um novo repositório de configurações de sistema contábil
func NewSistemaContabilConfigRepository(db *sql.DB) *SistemaContabilConfigRepository {
	return &SistemaContabilConfigRepository{DB: db}
}

// Create insere uma nova configuração de sistema contábil no banco de dados
func (r *SistemaContabilConfigRepository) Create(config *SistemaContabilConfig) error {
	// Validar dados da configuração
	if err := validateSistemaContabilConfig(config); err != nil {
		return err
	}
	
	query := `
	INSERT INTO sistema_contabil_config 
	(idSistemaContabil, idObjetoContabilizacao, idCodigoEvento, idSeguradora, ativo) 
	VALUES (?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		config.IdSistemaContabil, 
		config.IdObjetoContabilizacao, 
		config.IdCodigoEvento, 
		config.IdSeguradora, 
		config.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar configuração de sistema contábil: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da configuração: %v", err)
	}
	
	config.ID = id
	return nil
}

// GetAll retorna todas as configurações de sistema contábil do banco de dados
func (r *SistemaContabilConfigRepository) GetAll() ([]SistemaContabilConfig, error) {
	query := `
	SELECT 
		scc.idSistemaContabilConfig, scc.idSistemaContabil, scc.idObjetoContabilizacao, 
		scc.idCodigoEvento, scc.idSeguradora, scc.created_at, scc.updated_at, scc.ativo,
		sc.SistemaContabil, oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM sistema_contabil_config scc
	JOIN sistema_contabil sc ON scc.idSistemaContabil = sc.idSistemaContabil
	JOIN objeto_contabilizacao oc ON scc.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON scc.idCodigoEvento = e.idCodigoEvento
	ORDER BY scc.idSistemaContabilConfig DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configurações: %v", err)
	}
	defer rows.Close()
	
	var configs []SistemaContabilConfig
	
	for rows.Next() {
		var c SistemaContabilConfig
		if err := rows.Scan(
			&c.ID, 
			&c.IdSistemaContabil, 
			&c.IdObjetoContabilizacao, 
			&c.IdCodigoEvento, 
			&c.IdSeguradora, 
			&c.CreatedAt, 
			&c.UpdatedAt, 
			&c.Ativo,
			&c.SistemaContabilNome,
			&c.ObjetoContabilizacaoNome,
			&c.EventoNumero,
			&c.EventoDescricao,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler configuração: %v", err)
		}
		configs = append(configs, c)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre configurações: %v", err)
	}
	
	return configs, nil
}

// GetByID busca uma configuração pelo ID
func (r *SistemaContabilConfigRepository) GetByID(id int64) (*SistemaContabilConfig, error) {
	query := `
	SELECT 
		scc.idSistemaContabilConfig, scc.idSistemaContabil, scc.idObjetoContabilizacao, 
		scc.idCodigoEvento, scc.idSeguradora, scc.created_at, scc.updated_at, scc.ativo,
		sc.SistemaContabil, oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM sistema_contabil_config scc
	JOIN sistema_contabil sc ON scc.idSistemaContabil = sc.idSistemaContabil
	JOIN objeto_contabilizacao oc ON scc.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON scc.idCodigoEvento = e.idCodigoEvento
	WHERE scc.idSistemaContabilConfig = ?`
	
	var c SistemaContabilConfig
	err := r.DB.QueryRow(query, id).Scan(
		&c.ID, 
		&c.IdSistemaContabil, 
		&c.IdObjetoContabilizacao, 
		&c.IdCodigoEvento, 
		&c.IdSeguradora, 
		&c.CreatedAt, 
		&c.UpdatedAt, 
		&c.Ativo,
		&c.SistemaContabilNome,
		&c.ObjetoContabilizacaoNome,
		&c.EventoNumero,
		&c.EventoDescricao,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("configuração não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar configuração: %v", err)
	}
	
	return &c, nil
}

// GetBySeguradora busca configurações por seguradora
func (r *SistemaContabilConfigRepository) GetBySeguradora(idSeguradora int64) ([]SistemaContabilConfig, error) {
	query := `
	SELECT 
		scc.idSistemaContabilConfig, scc.idSistemaContabil, scc.idObjetoContabilizacao, 
		scc.idCodigoEvento, scc.idSeguradora, scc.created_at, scc.updated_at, scc.ativo,
		sc.SistemaContabil, oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM sistema_contabil_config scc
	JOIN sistema_contabil sc ON scc.idSistemaContabil = sc.idSistemaContabil
	JOIN objeto_contabilizacao oc ON scc.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON scc.idCodigoEvento = e.idCodigoEvento
	WHERE scc.idSeguradora = ? 
	ORDER BY scc.idSistemaContabilConfig DESC`
	
	rows, err := r.DB.Query(query, idSeguradora)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configurações por seguradora: %v", err)
	}
	defer rows.Close()
	
	var configs []SistemaContabilConfig
	
	for rows.Next() {
		var c SistemaContabilConfig
		if err := rows.Scan(
			&c.ID, 
			&c.IdSistemaContabil, 
			&c.IdObjetoContabilizacao, 
			&c.IdCodigoEvento, 
			&c.IdSeguradora, 
			&c.CreatedAt, 
			&c.UpdatedAt, 
			&c.Ativo,
			&c.SistemaContabilNome,
			&c.ObjetoContabilizacaoNome,
			&c.EventoNumero,
			&c.EventoDescricao,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler configuração: %v", err)
		}
		configs = append(configs, c)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre configurações: %v", err)
	}
	
	return configs, nil
}

// GetBySistemaContabil busca configurações por sistema contábil
func (r *SistemaContabilConfigRepository) GetBySistemaContabil(idSistemaContabil int64) ([]SistemaContabilConfig, error) {
	query := `
	SELECT 
		scc.idSistemaContabilConfig, scc.idSistemaContabil, scc.idObjetoContabilizacao, 
		scc.idCodigoEvento, scc.idSeguradora, scc.created_at, scc.updated_at, scc.ativo,
		sc.SistemaContabil, oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM sistema_contabil_config scc
	JOIN sistema_contabil sc ON scc.idSistemaContabil = sc.idSistemaContabil
	JOIN objeto_contabilizacao oc ON scc.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON scc.idCodigoEvento = e.idCodigoEvento
	WHERE scc.idSistemaContabil = ? 
	ORDER BY scc.idSistemaContabilConfig DESC`
	
	rows, err := r.DB.Query(query, idSistemaContabil)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configurações por sistema contábil: %v", err)
	}
	defer rows.Close()
	
	var configs []SistemaContabilConfig
	
	for rows.Next() {
		var c SistemaContabilConfig
		if err := rows.Scan(
			&c.ID, 
			&c.IdSistemaContabil, 
			&c.IdObjetoContabilizacao, 
			&c.IdCodigoEvento, 
			&c.IdSeguradora, 
			&c.CreatedAt, 
			&c.UpdatedAt, 
			&c.Ativo,
			&c.SistemaContabilNome,
			&c.ObjetoContabilizacaoNome,
			&c.EventoNumero,
			&c.EventoDescricao,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler configuração: %v", err)
		}
		configs = append(configs, c)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre configurações: %v", err)
	}
	
	return configs, nil
}

// Update atualiza os dados de uma configuração existente
func (r *SistemaContabilConfigRepository) Update(config *SistemaContabilConfig) error {
	// Validar dados da configuração
	if err := validateSistemaContabilConfig(config); err != nil {
		return err
	}
	
	query := `
	UPDATE sistema_contabil_config 
	SET idSistemaContabil = ?, idObjetoContabilizacao = ?, idCodigoEvento = ?, idSeguradora = ?, ativo = ? 
	WHERE idSistemaContabilConfig = ?`
	
	_, err := r.DB.Exec(
		query, 
		config.IdSistemaContabil, 
		config.IdObjetoContabilizacao, 
		config.IdCodigoEvento, 
		config.IdSeguradora, 
		config.Ativo, 
		config.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar configuração: %v", err)
	}
	
	return nil
}

// Delete remove uma configuração do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *SistemaContabilConfigRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM sistema_contabil_config WHERE idSistemaContabilConfig = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE sistema_contabil_config SET ativo = false WHERE idSistemaContabilConfig = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir configuração: %v", err)
	}
	
	return nil
}

// validateSistemaContabilConfig valida os dados de uma configuração
func validateSistemaContabilConfig(c *SistemaContabilConfig) error {
	// Validar sistema contábil
	if c.IdSistemaContabil <= 0 {
		return utils.ValidationError{
			Field:   "idSistemaContabil",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar objeto de contabilização
	if c.IdObjetoContabilizacao <= 0 {
		return utils.ValidationError{
			Field:   "idObjetoContabilizacao",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar evento
	if c.IdCodigoEvento <= 0 {
		return utils.ValidationError{
			Field:   "idCodigoEvento",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar seguradora
	if c.IdSeguradora <= 0 {
		return utils.ValidationError{
			Field:   "idSeguradora",
			Message: "deve ser um número positivo",
		}
	}
	
	return nil
}
