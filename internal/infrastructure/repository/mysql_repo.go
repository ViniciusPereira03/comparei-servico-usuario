package repository

import (
	"comparei-servico-usuario/internal/domain"
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
func (r *MySQLRepository) CreateUser(user *domain.User) error {
	fmt.Println("INSERT USER: ", user)
	_, err := r.db.Exec("INSERT INTO user (name, username, email, password, status) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Username, user.Email, user.Password, user.Status)
	return err
}

// Buscar usuário por ID
func (r *MySQLRepository) GetUser(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, name, username, email, status FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Buscar lista de usuários
func (r *MySQLRepository) GetUsers() (*[]domain.User, error) {
	var users []domain.User
	rows, err := r.db.Query("SELECT id, name, username, email, status FROM user WHERE status = 1 AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Status); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

// Atualizar usuário
func (r *MySQLRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec("UPDATE user SET name = ?, username = ?, email = ?, status = ? WHERE id = ?",
		user.Name, user.Username, user.Email, user.Password, user.Status)
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
func (repo *MySQLRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, username, email, password FROM user WHERE username = ?"
	err := repo.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
