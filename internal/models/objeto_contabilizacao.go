package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// ObjetoContabilizacao representa um objeto de contabilização no sistema
type ObjetoContabilizacao struct {
	ID                   int64     `json:"idObjetoContabilizacao"`
	ObjetoContabilizacao string    `json:"objetoContabilizacao"`
	Descricao            string    `json:"descricao"`
	IdSeguradora         int64     `json:"idSeguradora"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Ativo                bool      `json:"ativo"`
}

// ObjetoContabilizacaoRepository gerencia operações de banco de dados para objetos de contabilização
type ObjetoContabilizacaoRepository struct {
	DB *sql.DB
}

// NewObjetoContabilizacaoRepository cria um novo repositório de objetos de contabilização
func NewObjetoContabilizacaoRepository(db *sql.DB) *ObjetoContabilizacaoRepository {
	return &ObjetoContabilizacaoRepository{DB: db}
}

// Create insere um novo objeto de contabilização no banco de dados
func (r *ObjetoContabilizacaoRepository) Create(objeto *ObjetoContabilizacao) error {
	// Validar dados do objeto
	if err := validateObjetoContabilizacao(objeto); err != nil {
		return err
	}
	
	// Sanitizar dados
	objeto.ObjetoContabilizacao = utils.SanitizeString(objeto.ObjetoContabilizacao)
	objeto.Descricao = utils.SanitizeString(objeto.Descricao)
	
	query := `
	INSERT INTO objeto_contabilizacao 
	(ObjetoContabilizacao, Descricao, idSeguradora, ativo) 
	VALUES (?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		objeto.ObjetoContabilizacao, 
		objeto.Descricao, 
		objeto.IdSeguradora, 
		objeto.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar objeto de contabilização: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do objeto de contabilização: %v", err)
	}
	
	objeto.ID = id
	return nil
}

// GetAll retorna todos os objetos de contabilização do banco de dados
func (r *ObjetoContabilizacaoRepository) GetAll() ([]ObjetoContabilizacao, error) {
	query := `
	SELECT 
		idObjetoContabilizacao, ObjetoContabilizacao, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM objeto_contabilizacao 
	ORDER BY idObjetoContabilizacao DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar objetos de contabilização: %v", err)
	}
	defer rows.Close()
	
	var objetos []ObjetoContabilizacao
	
	for rows.Next() {
		var o ObjetoContabilizacao
		if err := rows.Scan(
			&o.ID, 
			&o.ObjetoContabilizacao, 
			&o.Descricao, 
			&o.IdSeguradora, 
			&o.CreatedAt, 
			&o.UpdatedAt, 
			&o.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler objeto de contabilização: %v", err)
		}
		objetos = append(objetos, o)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre objetos de contabilização: %v", err)
	}
	
	return objetos, nil
}

// GetByID busca um objeto de contabilização pelo ID
func (r *ObjetoContabilizacaoRepository) GetByID(id int64) (*ObjetoContabilizacao, error) {
	query := `
	SELECT 
		idObjetoContabilizacao, ObjetoContabilizacao, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM objeto_contabilizacao 
	WHERE idObjetoContabilizacao = ?`
	
	var o ObjetoContabilizacao
	err := r.DB.QueryRow(query, id).Scan(
		&o.ID, 
		&o.ObjetoContabilizacao, 
		&o.Descricao, 
		&o.IdSeguradora, 
		&o.CreatedAt, 
		&o.UpdatedAt, 
		&o.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("objeto de contabilização não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar objeto de contabilização: %v", err)
	}
	
	return &o, nil
}

// GetBySeguradora busca objetos de contabilização por seguradora
func (r *ObjetoContabilizacaoRepository) GetBySeguradora(idSeguradora int64) ([]ObjetoContabilizacao, error) {
	query := `
	SELECT 
		idObjetoContabilizacao, ObjetoContabilizacao, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM objeto_contabilizacao 
	WHERE idSeguradora = ? 
	ORDER BY idObjetoContabilizacao DESC`
	
	rows, err := r.DB.Query(query, idSeguradora)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar objetos de contabilização por seguradora: %v", err)
	}
	defer rows.Close()
	
	var objetos []ObjetoContabilizacao
	
	for rows.Next() {
		var o ObjetoContabilizacao
		if err := rows.Scan(
			&o.ID, 
			&o.ObjetoContabilizacao, 
			&o.Descricao, 
			&o.IdSeguradora, 
			&o.CreatedAt, 
			&o.UpdatedAt, 
			&o.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler objeto de contabilização: %v", err)
		}
		objetos = append(objetos, o)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre objetos de contabilização: %v", err)
	}
	
	return objetos, nil
}

// Update atualiza os dados de um objeto de contabilização existente
func (r *ObjetoContabilizacaoRepository) Update(objeto *ObjetoContabilizacao) error {
	// Validar dados do objeto
	if err := validateObjetoContabilizacao(objeto); err != nil {
		return err
	}
	
	// Sanitizar dados
	objeto.ObjetoContabilizacao = utils.SanitizeString(objeto.ObjetoContabilizacao)
	objeto.Descricao = utils.SanitizeString(objeto.Descricao)
	
	query := `
	UPDATE objeto_contabilizacao 
	SET ObjetoContabilizacao = ?, Descricao = ?, idSeguradora = ?, ativo = ? 
	WHERE idObjetoContabilizacao = ?`
	
	_, err := r.DB.Exec(
		query, 
		objeto.ObjetoContabilizacao, 
		objeto.Descricao, 
		objeto.IdSeguradora, 
		objeto.Ativo, 
		objeto.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar objeto de contabilização: %v", err)
	}
	
	return nil
}

// Delete remove um objeto de contabilização do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *ObjetoContabilizacaoRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM objeto_contabilizacao WHERE idObjetoContabilizacao = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE objeto_contabilizacao SET ativo = false WHERE idObjetoContabilizacao = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir objeto de contabilização: %v", err)
	}
	
	return nil
}

// validateObjetoContabilizacao valida os dados de um objeto de contabilização
func validateObjetoContabilizacao(o *ObjetoContabilizacao) error {
	// Validar objeto de contabilização
	if err := utils.ValidateRequired("objetoContabilizacao", o.ObjetoContabilizacao); err != nil {
		return err
	}
	if err := utils.ValidateLength("objetoContabilizacao", o.ObjetoContabilizacao, 3, 100); err != nil {
		return err
	}
	
	// Validar descrição
	if err := utils.ValidateRequired("descricao", o.Descricao); err != nil {
		return err
	}
	if err := utils.ValidateLength("descricao", o.Descricao, 3, 255); err != nil {
		return err
	}
	
	// Validar seguradora
	if o.IdSeguradora <= 0 {
		return utils.ValidationError{
			Field:   "idSeguradora",
			Message: "deve ser um número positivo",
		}
	}
	
	return nil
}
