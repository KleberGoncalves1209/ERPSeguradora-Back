package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// TipoPerfil representa um tipo de perfil de usuário no sistema
type TipoPerfil struct {
	ID        int64     `json:"id"`
	Perfil    string    `json:"perfil"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Ativo     bool      `json:"ativo"`
}

// TipoPerfilRepository gerencia operações de banco de dados para tipos de perfil
type TipoPerfilRepository struct {
	DB *sql.DB
}

// NewTipoPerfilRepository cria um novo repositório de tipos de perfil
func NewTipoPerfilRepository(db *sql.DB) *TipoPerfilRepository {
	return &TipoPerfilRepository{DB: db}
}

// Create insere um novo tipo de perfil no banco de dados
func (r *TipoPerfilRepository) Create(tipoPerfil *TipoPerfil) error {
	// Validar dados do tipo de perfil
	if err := validateTipoPerfil(tipoPerfil); err != nil {
		return err
	}
	
	// Sanitizar dados
	tipoPerfil.Perfil = utils.SanitizeString(tipoPerfil.Perfil)
	
	query := `
	INSERT INTO tipo_perfil 
	(perfil, ativo) 
	VALUES (?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		tipoPerfil.Perfil, 
		tipoPerfil.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar tipo de perfil: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do tipo de perfil: %v", err)
	}
	
	tipoPerfil.ID = id
	return nil
}

// GetAll retorna todos os tipos de perfil do banco de dados
func (r *TipoPerfilRepository) GetAll() ([]TipoPerfil, error) {
	query := `
	SELECT 
		id_tipo_perfil, perfil, created_at, updated_at, ativo 
	FROM tipo_perfil 
	ORDER BY id_tipo_perfil DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tipos de perfil: %v", err)
	}
	defer rows.Close()
	
	var tiposPerfil []TipoPerfil
	
	for rows.Next() {
		var tp TipoPerfil
		if err := rows.Scan(
			&tp.ID, 
			&tp.Perfil, 
			&tp.CreatedAt, 
			&tp.UpdatedAt, 
			&tp.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler tipo de perfil: %v", err)
		}
		tiposPerfil = append(tiposPerfil, tp)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre tipos de perfil: %v", err)
	}
	
	return tiposPerfil, nil
}

// GetByID busca um tipo de perfil pelo ID
func (r *TipoPerfilRepository) GetByID(id int64) (*TipoPerfil, error) {
	query := `
	SELECT 
		id_tipo_perfil, perfil, created_at, updated_at, ativo 
	FROM tipo_perfil 
	WHERE id_tipo_perfil = ?`
	
	var tp TipoPerfil
	err := r.DB.QueryRow(query, id).Scan(
		&tp.ID, 
		&tp.Perfil, 
		&tp.CreatedAt, 
		&tp.UpdatedAt, 
		&tp.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tipo de perfil não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar tipo de perfil: %v", err)
	}
	
	return &tp, nil
}

// Update atualiza os dados de um tipo de perfil existente
func (r *TipoPerfilRepository) Update(tipoPerfil *TipoPerfil) error {
	// Validar dados do tipo de perfil
	if err := validateTipoPerfil(tipoPerfil); err != nil {
		return err
	}
	
	// Sanitizar dados
	tipoPerfil.Perfil = utils.SanitizeString(tipoPerfil.Perfil)
	
	query := `
	UPDATE tipo_perfil 
	SET perfil = ?, ativo = ? 
	WHERE id_tipo_perfil = ?`
	
	_, err := r.DB.Exec(
		query, 
		tipoPerfil.Perfil, 
		tipoPerfil.Ativo, 
		tipoPerfil.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar tipo de perfil: %v", err)
	}
	
	return nil
}

// Delete remove um tipo de perfil do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *TipoPerfilRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM tipo_perfil WHERE id_tipo_perfil = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE tipo_perfil SET ativo = false WHERE id_tipo_perfil = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir tipo de perfil: %v", err)
	}
	
	return nil
}

// validateTipoPerfil valida os dados de um tipo de perfil
func validateTipoPerfil(tp *TipoPerfil) error {
	// Validar nome do perfil
	if err := utils.ValidateRequired("perfil", tp.Perfil); err != nil {
		return err
	}
	if err := utils.ValidateLength("perfil", tp.Perfil, 3, 100); err != nil {
		return err
	}
	
	return nil
}
