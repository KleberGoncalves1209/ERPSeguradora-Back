package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// SistemaContabil representa um sistema contábil no sistema
type SistemaContabil struct {
	ID              int64     `json:"idSistemaContabil"`
	SistemaContabil string    `json:"sistemaContabil"`
	IdSeguradora    int64     `json:"idSeguradora"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Ativo           bool      `json:"ativo"`
}

// SistemaContabilRepository gerencia operações de banco de dados para sistemas contábeis
type SistemaContabilRepository struct {
	DB *sql.DB
}

// NewSistemaContabilRepository cria um novo repositório de sistemas contábeis
func NewSistemaContabilRepository(db *sql.DB) *SistemaContabilRepository {
	return &SistemaContabilRepository{DB: db}
}

// Create insere um novo sistema contábil no banco de dados
func (r *SistemaContabilRepository) Create(sistema *SistemaContabil) error {
	// Validar dados do sistema
	if err := validateSistemaContabil(sistema); err != nil {
		return err
	}
	
	// Sanitizar dados
	sistema.SistemaContabil = utils.SanitizeString(sistema.SistemaContabil)
	
	query := `
	INSERT INTO sistema_contabil 
	(SistemaContabil, idSeguradora, ativo) 
	VALUES (?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		sistema.SistemaContabil, 
		sistema.IdSeguradora, 
		sistema.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar sistema contábil: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do sistema contábil: %v", err)
	}
	
	sistema.ID = id
	return nil
}

// GetAll retorna todos os sistemas contábeis do banco de dados
func (r *SistemaContabilRepository) GetAll() ([]SistemaContabil, error) {
	query := `
	SELECT 
		idSistemaContabil, SistemaContabil, idSeguradora, 
		created_at, updated_at, ativo 
	FROM sistema_contabil 
	ORDER BY idSistemaContabil DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sistemas contábeis: %v", err)
	}
	defer rows.Close()
	
	var sistemas []SistemaContabil
	
	for rows.Next() {
		var s SistemaContabil
		if err := rows.Scan(
			&s.ID, 
			&s.SistemaContabil, 
			&s.IdSeguradora, 
			&s.CreatedAt, 
			&s.UpdatedAt, 
			&s.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler sistema contábil: %v", err)
		}
		sistemas = append(sistemas, s)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre sistemas contábeis: %v", err)
	}
	
	return sistemas, nil
}

// GetByID busca um sistema contábil pelo ID
func (r *SistemaContabilRepository) GetByID(id int64) (*SistemaContabil, error) {
	query := `
	SELECT 
		idSistemaContabil, SistemaContabil, idSeguradora, 
		created_at, updated_at, ativo 
	FROM sistema_contabil 
	WHERE idSistemaContabil = ?`
	
	var s SistemaContabil
	err := r.DB.QueryRow(query, id).Scan(
		&s.ID, 
		&s.SistemaContabil, 
		&s.IdSeguradora, 
		&s.CreatedAt, 
		&s.UpdatedAt, 
		&s.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sistema contábil não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar sistema contábil: %v", err)
	}
	
	return &s, nil
}

// GetBySeguradora busca sistemas contábeis por seguradora
func (r *SistemaContabilRepository) GetBySeguradora(idSeguradora int64) ([]SistemaContabil, error) {
	query := `
	SELECT 
		idSistemaContabil, SistemaContabil, idSeguradora, 
		created_at, updated_at, ativo 
	FROM sistema_contabil 
	WHERE idSeguradora = ? 
	ORDER BY idSistemaContabil DESC`
	
	rows, err := r.DB.Query(query, idSeguradora)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sistemas contábeis por seguradora: %v", err)
	}
	defer rows.Close()
	
	var sistemas []SistemaContabil
	
	for rows.Next() {
		var s SistemaContabil
		if err := rows.Scan(
			&s.ID, 
			&s.SistemaContabil, 
			&s.IdSeguradora, 
			&s.CreatedAt, 
			&s.UpdatedAt, 
			&s.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler sistema contábil: %v", err)
		}
		sistemas = append(sistemas, s)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre sistemas contábeis: %v", err)
	}
	
	return sistemas, nil
}

// Update atualiza os dados de um sistema contábil existente
func (r *SistemaContabilRepository) Update(sistema *SistemaContabil) error {
	// Validar dados do sistema
	if err := validateSistemaContabil(sistema); err != nil {
		return err
	}
	
	// Sanitizar dados
	sistema.SistemaContabil = utils.SanitizeString(sistema.SistemaContabil)
	
	query := `
	UPDATE sistema_contabil 
	SET SistemaContabil = ?, idSeguradora = ?, ativo = ? 
	WHERE idSistemaContabil = ?`
	
	_, err := r.DB.Exec(
		query, 
		sistema.SistemaContabil, 
		sistema.IdSeguradora, 
		sistema.Ativo, 
		sistema.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar sistema contábil: %v", err)
	}
	
	return nil
}

// Delete remove um sistema contábil do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *SistemaContabilRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM sistema_contabil WHERE idSistemaContabil = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE sistema_contabil SET ativo = false WHERE idSistemaContabil = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir sistema contábil: %v", err)
	}
	
	return nil
}

// validateSistemaContabil valida os dados de um sistema contábil
func validateSistemaContabil(s *SistemaContabil) error {
	// Validar nome do sistema contábil
	if err := utils.ValidateRequired("sistemaContabil", s.SistemaContabil); err != nil {
		return err
	}
	if err := utils.ValidateLength("sistemaContabil", s.SistemaContabil, 3, 100); err != nil {
		return err
	}
	
	// Validar seguradora
	if s.IdSeguradora <= 0 {
		return utils.ValidationError{
			Field:   "idSeguradora",
			Message: "deve ser um número positivo",
		}
	}
	
	return nil
}
