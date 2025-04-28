package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// ObjetoContabilizacaoEvento representa uma relação entre objeto de contabilização e evento
type ObjetoContabilizacaoEvento struct {
	ID                        int64     `json:"idObjetoContabilizacaoEvento"`
	IdObjetoContabilizacao    int64     `json:"idObjetoContabilizacao"`
	IdCodigoEvento            int64     `json:"idCodigoEvento"`
	IdSeguradora              int64     `json:"idSeguradora"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Ativo                     bool      `json:"ativo"`
	// Campos para exibição de informações relacionadas
	ObjetoContabilizacaoNome  string    `json:"objetoContabilizacaoNome,omitempty"`
	EventoNumero              int       `json:"eventoNumero,omitempty"`
	EventoDescricao           string    `json:"eventoDescricao,omitempty"`
}

// ObjetoContabilizacaoEventoRepository gerencia operações de banco de dados para relações entre objetos de contabilização e eventos
type ObjetoContabilizacaoEventoRepository struct {
	DB *sql.DB
}

// NewObjetoContabilizacaoEventoRepository cria um novo repositório de relações entre objetos de contabilização e eventos
func NewObjetoContabilizacaoEventoRepository(db *sql.DB) *ObjetoContabilizacaoEventoRepository {
	return &ObjetoContabilizacaoEventoRepository{DB: db}
}

// Create insere uma nova relação entre objeto de contabilização e evento no banco de dados
func (r *ObjetoContabilizacaoEventoRepository) Create(relacao *ObjetoContabilizacaoEvento) error {
	// Validar dados da relação
	if err := validateObjetoContabilizacaoEvento(relacao); err != nil {
		return err
	}
	
	query := `
	INSERT INTO objeto_contabilizacao_evento 
	(idObjetoContabilizacao, idCodigoEvento, idSeguradora, ativo) 
	VALUES (?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		relacao.IdObjetoContabilizacao, 
		relacao.IdCodigoEvento, 
		relacao.IdSeguradora, 
		relacao.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar relação entre objeto de contabilização e evento: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID da relação: %v", err)
	}
	
	relacao.ID = id
	return nil
}

// GetAll retorna todas as relações entre objetos de contabilização e eventos do banco de dados
func (r *ObjetoContabilizacaoEventoRepository) GetAll() ([]ObjetoContabilizacaoEvento, error) {
	query := `
	SELECT 
		oce.idObjetoContabilizacaoEvento, oce.idObjetoContabilizacao, oce.idCodigoEvento, 
		oce.idSeguradora, oce.created_at, oce.updated_at, oce.ativo,
		oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM objeto_contabilizacao_evento oce
	JOIN objeto_contabilizacao oc ON oce.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON oce.idCodigoEvento = e.idCodigoEvento
	ORDER BY oce.idObjetoContabilizacaoEvento DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações: %v", err)
	}
	defer rows.Close()
	
	var relacoes []ObjetoContabilizacaoEvento
	
	for rows.Next() {
		var r ObjetoContabilizacaoEvento
		if err := rows.Scan(
			&r.ID, 
			&r.IdObjetoContabilizacao, 
			&r.IdCodigoEvento, 
			&r.IdSeguradora, 
			&r.CreatedAt, 
			&r.UpdatedAt, 
			&r.Ativo,
			&r.ObjetoContabilizacaoNome,
			&r.EventoNumero,
			&r.EventoDescricao,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler relação: %v", err)
		}
		relacoes = append(relacoes, r)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre relações: %v", err)
	}
	
	return relacoes, nil
}

// GetByID busca uma relação pelo ID
func (r *ObjetoContabilizacaoEventoRepository) GetByID(id int64) (*ObjetoContabilizacaoEvento, error) {
	query := `
	SELECT 
		oce.idObjetoContabilizacaoEvento, oce.idObjetoContabilizacao, oce.idCodigoEvento, 
		oce.idSeguradora, oce.created_at, oce.updated_at, oce.ativo,
		oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM objeto_contabilizacao_evento oce
	JOIN objeto_contabilizacao oc ON oce.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON oce.idCodigoEvento = e.idCodigoEvento
	WHERE oce.idObjetoContabilizacaoEvento = ?`
	
	var rel ObjetoContabilizacaoEvento
	err := r.DB.QueryRow(query, id).Scan(
		&rel.ID, 
		&rel.IdObjetoContabilizacao, 
		&rel.IdCodigoEvento, 
		&rel.IdSeguradora, 
		&rel.CreatedAt, 
		&rel.UpdatedAt, 
		&rel.Ativo,
		&rel.ObjetoContabilizacaoNome,
		&rel.EventoNumero,
		&rel.EventoDescricao,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("relação não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar relação: %v", err)
	}
	
	return &rel, nil
}

// GetBySeguradora busca relações por seguradora
func (r *ObjetoContabilizacaoEventoRepository) GetBySeguradora(idSeguradora int64) ([]ObjetoContabilizacaoEvento, error) {
	query := `
	SELECT 
		oce.idObjetoContabilizacaoEvento, oce.idObjetoContabilizacao, oce.idCodigoEvento, 
		oce.idSeguradora, oce.created_at, oce.updated_at, oce.ativo,
		oc.ObjetoContabilizacao, e.Evento, e.Descricao
	FROM objeto_contabilizacao_evento oce
	JOIN objeto_contabilizacao oc ON oce.idObjetoContabilizacao = oc.idObjetoContabilizacao
	JOIN eventos e ON oce.idCodigoEvento = e.idCodigoEvento
	WHERE oce.idSeguradora = ? 
	ORDER BY oce.idObjetoContabilizacaoEvento DESC`
	
	rows, err := r.DB.Query(query, idSeguradora)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações por seguradora: %v", err)
	}
	defer rows.Close()
	
	var relacoes []ObjetoContabilizacaoEvento
	
	for rows.Next() {
		var r ObjetoContabilizacaoEvento
		if err := rows.Scan(
			&r.ID, 
			&r.IdObjetoContabilizacao, 
			&r.IdCodigoEvento, 
			&r.IdSeguradora, 
			&r.CreatedAt, 
			&r.UpdatedAt, 
			&r.Ativo,
			&r.ObjetoContabilizacaoNome,
			&r.EventoNumero,
			&r.EventoDescricao,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler relação: %v", err)
		}
		relacoes = append(relacoes, r)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre relações: %v", err)
	}
	
	return relacoes, nil
}

// Update atualiza os dados de uma relação existente
func (r *ObjetoContabilizacaoEventoRepository) Update(relacao *ObjetoContabilizacaoEvento) error {
	// Validar dados da relação
	if err := validateObjetoContabilizacaoEvento(relacao); err != nil {
		return err
	}
	
	query := `
	UPDATE objeto_contabilizacao_evento 
	SET idObjetoContabilizacao = ?, idCodigoEvento = ?, idSeguradora = ?, ativo = ? 
	WHERE idObjetoContabilizacaoEvento = ?`
	
	_, err := r.DB.Exec(
		query, 
		relacao.IdObjetoContabilizacao, 
		relacao.IdCodigoEvento, 
		relacao.IdSeguradora, 
		relacao.Ativo, 
		relacao.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar relação: %v", err)
	}
	
	return nil
}

// Delete remove uma relação do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *ObjetoContabilizacaoEventoRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM objeto_contabilizacao_evento WHERE idObjetoContabilizacaoEvento = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE objeto_contabilizacao_evento SET ativo = false WHERE idObjetoContabilizacaoEvento = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir relação: %v", err)
	}
	
	return nil
}

// validateObjetoContabilizacaoEvento valida os dados de uma relação
func validateObjetoContabilizacaoEvento(r *ObjetoContabilizacaoEvento) error {
	// Validar objeto de contabilização
	if r.IdObjetoContabilizacao <= 0 {
		return utils.ValidationError{
			Field:   "idObjetoContabilizacao",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar evento
	if r.IdCodigoEvento <= 0 {
		return utils.ValidationError{
			Field:   "idCodigoEvento",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar seguradora
	if r.IdSeguradora <= 0 {
		return utils.ValidationError{
			Field:   "idSeguradora",
			Message: "deve ser um número positivo",
		}
	}
	
	return nil
}
