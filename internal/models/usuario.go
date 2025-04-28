package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KleberGoncalves1209/EstudoGo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// Usuario representa um usuário no sistema
type Usuario struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Email        string    `json:"email"`
	Login        string    `json:"login"`
	Senha        string    `json:"senha,omitempty"` // omitempty para não retornar a senha em respostas JSON
	IdTipoPerfil int       `json:"idTipoPerfil"`
	IdSeguradora int       `json:"idSeguradora"`
	AdminERP     bool      `json:"adminERP"`
	Bloqueado    bool      `json:"bloqueado"`
	BloqueadoAte *time.Time `json:"bloqueado_ate,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Ativo        bool      `json:"ativo"`
}

// UsuarioRepository gerencia operações de banco de dados para usuários
type UsuarioRepository struct {
	DB *sql.DB
}

// NewUsuarioRepository cria um novo repositório de usuários
func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{DB: db}
}

// Create insere um novo usuário no banco de dados
func (r *UsuarioRepository) Create(usuario *Usuario) error {
	// Validar dados do usuário
	if err := validateUsuario(usuario); err != nil {
		return err
	}
	
	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}
	
	query := `
	INSERT INTO usuarios 
	(nome, email, login, senha, idTipoPerfil, idSeguradora, AdminERP, bloqueado, ativo) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		usuario.Nome, 
		usuario.Email, 
		usuario.Login, 
		string(hashedPassword), 
		usuario.IdTipoPerfil, 
		usuario.IdSeguradora, 
		usuario.AdminERP,
		usuario.Bloqueado,
		usuario.Ativo,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar usuário: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter ID do usuário: %v", err)
	}
	
	usuario.ID = id
	return nil
}

// GetAll retorna todos os usuários do banco de dados
func (r *UsuarioRepository) GetAll() ([]Usuario, error) {
	query := `
	SELECT 
		id, nome, email, login, idTipoPerfil, idSeguradora, 
		AdminERP, bloqueado, bloqueado_ate, created_at, updated_at, ativo 
	FROM usuarios 
	ORDER BY id DESC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %v", err)
	}
	defer rows.Close()
	
	var usuarios []Usuario
	
	for rows.Next() {
		var u Usuario
		var bloqueadoAte sql.NullTime
		
		if err := rows.Scan(
			&u.ID, 
			&u.Nome, 
			&u.Email, 
			&u.Login, 
			&u.IdTipoPerfil, 
			&u.IdSeguradora, 
			&u.AdminERP,
			&u.Bloqueado,
			&bloqueadoAte,
			&u.CreatedAt, 
			&u.UpdatedAt, 
			&u.Ativo,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler usuário: %v", err)
		}
		
		if bloqueadoAte.Valid {
			u.BloqueadoAte = &bloqueadoAte.Time
		}
		
		usuarios = append(usuarios, u)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre usuários: %v", err)
	}
	
	return usuarios, nil
}

// GetByID busca um usuário pelo ID
func (r *UsuarioRepository) GetByID(id int64) (*Usuario, error) {
	query := `
	SELECT 
		id, nome, email, login, idTipoPerfil, idSeguradora, 
		AdminERP, bloqueado, bloqueado_ate, created_at, updated_at, ativo 
	FROM usuarios 
	WHERE id = ?`
	
	var u Usuario
	var bloqueadoAte sql.NullTime
	
	err := r.DB.QueryRow(query, id).Scan(
		&u.ID, 
		&u.Nome, 
		&u.Email, 
		&u.Login, 
		&u.IdTipoPerfil, 
		&u.IdSeguradora, 
		&u.AdminERP,
		&u.Bloqueado,
		&bloqueadoAte,
		&u.CreatedAt, 
		&u.UpdatedAt, 
		&u.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}
	
	if bloqueadoAte.Valid {
		u.BloqueadoAte = &bloqueadoAte.Time
	}
	
	return &u, nil
}

// GetByLogin busca um usuário pelo login
func (r *UsuarioRepository) GetByLogin(login string) (*Usuario, error) {
	query := `
	SELECT 
		id, nome, email, login, senha, idTipoPerfil, idSeguradora, 
		AdminERP, bloqueado, bloqueado_ate, created_at, updated_at, ativo 
	FROM usuarios 
	WHERE login = ?`
	
	var u Usuario
	var bloqueadoAte sql.NullTime
	
	err := r.DB.QueryRow(query, login).Scan(
		&u.ID, 
		&u.Nome, 
		&u.Email, 
		&u.Login, 
		&u.Senha,
		&u.IdTipoPerfil, 
		&u.IdSeguradora, 
		&u.AdminERP,
		&u.Bloqueado,
		&bloqueadoAte,
		&u.CreatedAt, 
		&u.UpdatedAt, 
		&u.Ativo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %v", err)
	}
	
	if bloqueadoAte.Valid {
		u.BloqueadoAte = &bloqueadoAte.Time
	}
	
	return &u, nil
}

// Update atualiza os dados de um usuário existente
func (r *UsuarioRepository) Update(usuario *Usuario) error {
	// Validar dados do usuário (exceto senha que é tratada separadamente)
	if err := validateUsuarioUpdate(usuario); err != nil {
		return err
	}
	
	query := `
	UPDATE usuarios 
	SET nome = ?, email = ?, login = ?, idTipoPerfil = ?, 
		idSeguradora = ?, AdminERP = ?, bloqueado = ?, ativo = ? 
	WHERE id = ?`
	
	_, err := r.DB.Exec(
		query, 
		usuario.Nome, 
		usuario.Email, 
		usuario.Login, 
		usuario.IdTipoPerfil, 
		usuario.IdSeguradora, 
		usuario.AdminERP,
		usuario.Bloqueado,
		usuario.Ativo, 
		usuario.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %v", err)
	}
	
	return nil
}

// UpdatePassword atualiza apenas a senha do usuário
func (r *UsuarioRepository) UpdatePassword(id int64, novaSenha string) error {
	// Validar a nova senha
	if err := utils.ValidatePassword(novaSenha); err != nil {
		return err
	}
	
	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(novaSenha), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}
	
	query := `UPDATE usuarios SET senha = ? WHERE id = ?`
	
	_, err = r.DB.Exec(query, string(hashedPassword), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar senha: %v", err)
	}
	
	return nil
}

// Delete remove um usuário do banco de dados (ou desativa, dependendo da regra de negócio)
func (r *UsuarioRepository) Delete(id int64) error {
	// Opção 1: Exclusão física
	// query := `DELETE FROM usuarios WHERE id = ?`
	
	// Opção 2: Exclusão lógica (recomendada)
	query := `UPDATE usuarios SET ativo = false WHERE id = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir usuário: %v", err)
	}
	
	return nil
}

// VerifyPassword verifica se a senha fornecida corresponde à senha armazenada
func (r *UsuarioRepository) VerifyPassword(login, senha string) (*Usuario, error) {
	// Buscar o usuário pelo login
	usuario, err := r.GetByLogin(login)
	if err != nil {
		return nil, err
	}
	
	// Verificar se o usuário está ativo
	if !usuario.Ativo {
		return nil, fmt.Errorf("usuário inativo")
	}
	
	// Verificar se o usuário está bloqueado
	if usuario.Bloqueado {
		if usuario.BloqueadoAte != nil && usuario.BloqueadoAte.After(time.Now()) {
			return nil, fmt.Errorf("usuário bloqueado temporariamente")
		} else {
			// Se o tempo de bloqueio já passou, desbloquear o usuário
			usuario.Bloqueado = false
			usuario.BloqueadoAte = nil
			
			// Atualizar o status de bloqueio no banco de dados
			updateQuery := `
			UPDATE usuarios 
			SET bloqueado = false, bloqueado_ate = NULL 
			WHERE id = ?`
			
			_, err := r.DB.Exec(updateQuery, usuario.ID)
			if err != nil {
				return nil, fmt.Errorf("erro ao desbloquear usuário: %v", err)
			}
		}
	}
	
	// Comparar a senha fornecida com o hash armazenado
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha))
	if err != nil {
		return nil, fmt.Errorf("senha incorreta")
	}
	
	// Limpar a senha antes de retornar o usuário
	usuario.Senha = ""
	
	return usuario, nil
}

// Funções de validação

// validateUsuario valida os dados de um novo usuário
func validateUsuario(u *Usuario) error {
	// Validar nome
	if err := utils.ValidateRequired("nome", u.Nome); err != nil {
		return err
	}
	if err := utils.ValidateLength("nome", u.Nome, 3, 100); err != nil {
		return err
	}
	
	// Validar email
	if err := utils.ValidateRequired("email", u.Email); err != nil {
		return err
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return utils.ValidationError{Field: "email", Message: err.Error()}
	}
	
	// Validar login
	if err := utils.ValidateRequired("login", u.Login); err != nil {
		return err
	}
	if err := utils.ValidateLength("login", u.Login, 3, 50); err != nil {
		return err
	}
	
	// Validar senha
	if err := utils.ValidateRequired("senha", u.Senha); err != nil {
		return err
	}
	if err := utils.ValidatePassword(u.Senha); err != nil {
		return utils.ValidationError{Field: "senha", Message: err.Error()}
	}
	
	// Validar tipo de perfil
	if err := utils.ValidateNumericRange("idTipoPerfil", u.IdTipoPerfil, 1, 0); err != nil {
		return err
	}
	
	// Validar seguradora
	if err := utils.ValidateNumericRange("idSeguradora", u.IdSeguradora, 1, 0); err != nil {
		return err
	}
	
	return nil
}

// validateUsuarioUpdate valida os dados de atualização de um usuário (sem senha)
func validateUsuarioUpdate(u *Usuario) error {
	// Validar nome
	if err := utils.ValidateRequired("nome", u.Nome); err != nil {
		return err
	}
	if err := utils.ValidateLength("nome", u.Nome, 3, 100); err != nil {
		return err
	}
	
	// Validar email
	if err := utils.ValidateRequired("email", u.Email); err != nil {
		return err
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return utils.ValidationError{Field: "email", Message: err.Error()}
	}
	
	// Validar login
	if err := utils.ValidateRequired("login", u.Login); err != nil {
		return err
	}
	if err := utils.ValidateLength("login", u.Login, 3, 50); err != nil {
		return err
	}
	
	// Validar tipo de perfil
	if err := utils.ValidateNumericRange("idTipoPerfil", u.IdTipoPerfil, 1, 0); err != nil {
		return err
	}
	
	// Validar seguradora
	if err := utils.ValidateNumericRange("idSeguradora", u.IdSeguradora, 1, 0); err != nil {
		return err
	}
	
	return nil
}
