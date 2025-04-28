package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// Seguradora representa uma seguradora no sistema
type Seguradora struct {
	ID            int64     `json:"id"`
	Nome          string    `json:"nome"`
	NomeAbreviado string    `json:"nome_abreviado"`
	CodigoSusep   string    `json:"codigo_susep"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Ativo         bool      `json:"ativo"`
}

// SeguradoraRepository gerencia operações de banco de dados para seguradoras
type SeguradoraRepository struct {
	DB *sql.DB
}

// NewSeguradoraRepository cria um novo repositório de seguradoras
func NewSeguradoraRepository(db *sql.DB) *SeguradoraRepository {
	return &SeguradoraRepository{DB: db}
}

// Create insere uma nova seguradora no banco de dados
func (r *SeguradoraRepository) Create(seguradora *Seguradora) error {
	// Validar dados da seguradora
	if err := validateSeguradora(seguradora); err != nil {
		return err
	}
	
	// Sanitizar dados
	seguradora.Nome = utils.SanitizeString(seguradora.Nome)
	seguradora.NomeAbreviado = utils.SanitizeString(seguradora.NomeAbreviado)
	seguradora.CodigoSusep = utils.SanitizeString(seguradora.CodigoSusep)
	
	query := `
	INSERT INTO seguradoras 
	(seguradora, nome_abreviado, codigo_susep, ativo) 
	VALUES (?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		seguradora.Nome, 
		seguradora.NomeAbreviado, 
		seguradora.CodigoSusep, 
		seguradora.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar seguradora: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da seguradora: %v", err)
	}
	
	seguradora.ID = id
	return nil
}

// GetAll retorna todas as seguradoras do banco de dados
func (r *SeguradoraRepository) GetAll() ([]Seguradora, error) {
	query := `
	SELECT 
		id_seguradora, seguradora, nome_abreviado, codigo_susep, 
		created_at, updated_at, ativo 
	FROM seguradoras 
	ORDER BY id_seguradora DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar seguradoras: %v", err)
	}
	defer rows.Close()
	
	var seguradoras []Seguradora
	
	for rows.Next() {
		var s Seguradora
		if err := rows.Scan(
			&s.ID, 
			&s.Nome, 
			&s.NomeAbreviado, 
			&s.CodigoSusep, 
			&s.CreatedAt, 
			&s.UpdatedAt, 
			&s.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler seguradora: %v", err)
		}
		seguradoras = append(seguradoras, s)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre seguradoras: %v", err)
	}
	
	return seguradoras, nil
}

// GetByID busca uma seguradora pelo ID
func (r *SeguradoraRepository) GetByID(id int64) (*Seguradora, error) {
	query := `
	SELECT 
		id_seguradora, seguradora, nome_abreviado, codigo_susep, 
		created_at, updated_at, ativo 
	FROM seguradoras 
	WHERE id_seguradora = ?`
	
	var s Seguradora
	err := r.DB.QueryRow(query, id).Scan(
		&s.ID, 
		&s.Nome, 
		&s.NomeAbreviado, 
		&s.CodigoSusep, 
		&s.CreatedAt, 
		&s.UpdatedAt, 
		&s.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("seguradora não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar seguradora: %v", err)
	}
	
	return &s, nil
}

// Update atualiza os dados de uma seguradora existente
func (r *SeguradoraRepository) Update(seguradora *Seguradora) error {
	// Validar dados da seguradora
	if err := validateSeguradora(seguradora); err != nil {
		return err
	}
	
	// Sanitizar dados
	seguradora.Nome = utils.SanitizeString(seguradora.Nome)
	seguradora.NomeAbreviado = utils.SanitizeString(seguradora.NomeAbreviado)
	seguradora.CodigoSusep = utils.SanitizeString(seguradora.CodigoSusep)
	
	query := `
	UPDATE seguradoras 
	SET seguradora = ?, nome_abreviado = ?, codigo_susep = ?, ativo = ? 
	WHERE id_seguradora = ?`
	
	_, err := r.DB.Exec(
		query, 
		seguradora.Nome, 
		seguradora.NomeAbreviado, 
		seguradora.CodigoSusep, 
		seguradora.Ativo, 
		seguradora.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar seguradora: %v", err)
	}
	
	return nil
}

// Delete remove uma seguradora do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *SeguradoraRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM seguradoras WHERE id_seguradora = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE seguradoras SET ativo = false WHERE id_seguradora = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir seguradora: %v", err)
	}
	
	return nil
}

// validateSeguradora valida os dados de uma seguradora
func validateSeguradora(s *Seguradora) error {
	// Validar nome da seguradora
	if err := utils.ValidateRequired("nome", s.Nome); err != nil {
		return err
	}
	if err := utils.ValidateLength("nome", s.Nome, 3, 100); err != nil {
		return err
	}
	
	// Validar nome abreviado (opcional)
	if s.NomeAbreviado != "" {
		if err := utils.ValidateLength("nome_abreviado", s.NomeAbreviado, 2, 50); err != nil {
			return err
		}
	}
	
	// Validar código SUSEP (opcional)
	if s.CodigoSusep != "" {
		if err := utils.ValidateLength("codigo_susep", s.CodigoSusep, 2, 20); err != nil {
			return err
		}
	}
	
	return nil
}
