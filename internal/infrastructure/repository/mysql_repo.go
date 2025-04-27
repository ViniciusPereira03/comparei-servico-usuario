package repository

import (
	"comparei-servico-usuario/internal/domain/user"
	"database/sql"
	"fmt"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Cadastrar usuário
func (r *MySQLRepository) CreateUser(user *user.User) (*user.User, error) {
	user.Status = 1
	user.Level = 1

	result, err := r.db.Exec("INSERT INTO user (name, username, email, password, status, level) VALUES (?, ?, ?, ?, ?, ?)",
		user.Name, user.Username, user.Email, user.Password, user.Status, user.Level)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)
	return user, nil
}

// Buscar usuário por ID
func (r *MySQLRepository) GetUser(id int) (*user.User, error) {
	var user user.User
	err := r.db.QueryRow("SELECT id, name, username, email, status FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Buscar lista de usuários
func (r *MySQLRepository) GetUsers() ([]*user.User, error) {
	var users []user.User
	rows, err := r.db.Query("SELECT id, name, username, email, status FROM user WHERE status = 1 AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Status); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	userPtrs := make([]*user.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	return userPtrs, nil
}

// Atualizar usuário
func (r *MySQLRepository) UpdateUser(user *user.User) error {
	_, err := r.db.Exec("UPDATE user SET name = ?, username = ?, email = ?, status = ? WHERE id = ?",
		user.Name, user.Username, user.Email, user.Status, user.ID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// Deletar usuário
func (r *MySQLRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("UPDATE user SET status = 0, deleted_at = NOW() WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// GetUserByUsername busca um usuário pelo nome de usuário
func (repo *MySQLRepository) GetUserByUsername(username string) (*user.User, error) {
	var user user.User
	query := "SELECT id, username, email, password FROM user WHERE username = ?"
	err := repo.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
