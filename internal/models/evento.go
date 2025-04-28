package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
)

// Evento representa um evento no sistema
type Evento struct {
	ID           int64     `json:"idCodigoEvento"`
	Evento       int       `json:"evento"`
	Descricao    string    `json:"descricao"`
	IdSeguradora int64     `json:"idSeguradora"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Ativo        bool      `json:"ativo"`
}

// EventoRepository gerencia operações de banco de dados para eventos
type EventoRepository struct {
	DB *sql.DB
}

// NewEventoRepository cria um novo repositório de eventos
func NewEventoRepository(db *sql.DB) *EventoRepository {
	return &EventoRepository{DB: db}
}

// Create insere um novo evento no banco de dados
func (r *EventoRepository) Create(evento *Evento) error {
	// Validar dados do evento
	if err := validateEvento(evento); err != nil {
		return err
	}
	
	// Sanitizar dados
	evento.Descricao = utils.SanitizeString(evento.Descricao)
	
	query := `
	INSERT INTO eventos 
	(Evento, Descricao, idSeguradora, ativo) 
	VALUES (?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		evento.Evento, 
		evento.Descricao, 
		evento.IdSeguradora, 
		evento.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar evento: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do evento: %v", err)
	}
	
	evento.ID = id
	return nil
}

// GetAll retorna todos os eventos do banco de dados
func (r *EventoRepository) GetAll() ([]Evento, error) {
	query := `
	SELECT 
		idCodigoEvento, Evento, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM eventos 
	ORDER BY idCodigoEvento DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %v", err)
	}
	defer rows.Close()
	
	var eventos []Evento
	
	for rows.Next() {
		var e Evento
		if err := rows.Scan(
			&e.ID, 
			&e.Evento, 
			&e.Descricao, 
			&e.IdSeguradora, 
			&e.CreatedAt, 
			&e.UpdatedAt, 
			&e.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler evento: %v", err)
		}
		eventos = append(eventos, e)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre eventos: %v", err)
	}
	
	return eventos, nil
}

// GetByID busca um evento pelo ID
func (r *EventoRepository) GetByID(id int64) (*Evento, error) {
	query := `
	SELECT 
		idCodigoEvento, Evento, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM eventos 
	WHERE idCodigoEvento = ?`
	
	var e Evento
	err := r.DB.QueryRow(query, id).Scan(
		&e.ID, 
		&e.Evento, 
		&e.Descricao, 
		&e.IdSeguradora, 
		&e.CreatedAt, 
		&e.UpdatedAt, 
		&e.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("evento não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar evento: %v", err)
	}
	
	return &e, nil
}

// GetBySeguradora busca eventos por seguradora
func (r *EventoRepository) GetBySeguradora(idSeguradora int64) ([]Evento, error) {
	query := `
	SELECT 
		idCodigoEvento, Evento, Descricao, idSeguradora, 
		created_at, updated_at, ativo 
	FROM eventos 
	WHERE idSeguradora = ? 
	ORDER BY idCodigoEvento DESC`
	
	rows, err := r.DB.Query(query, idSeguradora)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos por seguradora: %v", err)
	}
	defer rows.Close()
	
	var eventos []Evento
	
	for rows.Next() {
		var e Evento
		if err := rows.Scan(
			&e.ID, 
			&e.Evento, 
			&e.Descricao, 
			&e.IdSeguradora, 
			&e.CreatedAt, 
			&e.UpdatedAt, 
			&e.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler evento: %v", err)
		}
		eventos = append(eventos, e)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre eventos: %v", err)
	}
	
	return eventos, nil
}

// Update atualiza os dados de um evento existente
func (r *EventoRepository) Update(evento *Evento) error {
	// Validar dados do evento
	if err := validateEvento(evento); err != nil {
		return err
	}
	
	// Sanitizar dados
	evento.Descricao = utils.SanitizeString(evento.Descricao)
	
	query := `
	UPDATE eventos 
	SET Evento = ?, Descricao = ?, idSeguradora = ?, ativo = ? 
	WHERE idCodigoEvento = ?`
	
	_, err := r.DB.Exec(
		query, 
		evento.Evento, 
		evento.Descricao, 
		evento.IdSeguradora, 
		evento.Ativo, 
		evento.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar evento: %v", err)
	}
	
	return nil
}

// Delete remove um evento do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *EventoRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM eventos WHERE idCodigoEvento = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE eventos SET ativo = false WHERE idCodigoEvento = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir evento: %v", err)
	}
	
	return nil
}

// validateEvento valida os dados de um evento
func validateEvento(e *Evento) error {
	// Validar evento
	if e.Evento <= 0 {
		return utils.ValidationError{
			Field:   "evento",
			Message: "deve ser um número positivo",
		}
	}
	
	// Validar descrição
	if err := utils.ValidateRequired("descricao", e.Descricao); err != nil {
		return err
	}
	if err := utils.ValidateLength("descricao", e.Descricao, 3, 255); err != nil {
		return err
	}
	
	// Validar seguradora
	if e.IdSeguradora <= 0 {
		return utils.ValidationError{
			Field:   "idSeguradora",
			Message: "deve ser um número positivo",
		}
	}
	
	return nil
}
