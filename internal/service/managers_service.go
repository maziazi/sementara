package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"project_sprint/internal/model"
	"project_sprint/pkg/database"
	"strings"
	"time"
)

func RegisterUser(email, password string) (*model.Managers, error) {
	db := database.GetDBPool()

	var existingUser model.Managers
	err := db.QueryRow(context.Background(), "SELECT email FROM manager WHERE email = $1", email).Scan(&existingUser.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("database error: %v", err)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err = db.Exec(context.Background(), "INSERT INTO manager (email, password, created_at) VALUES ($1, $2, $3)",
		email, string(hashedPassword), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &model.Managers{Email: email, CreatedAt: time.Now()}, nil
}

func AuthenticateManager(email, password string) (*model.Managers, error) {
	db := database.GetDBPool()

	var user model.Managers
	err := db.QueryRow(context.Background(), "SELECT id, email, password FROM manager WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("email not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	_, err = db.Exec(context.Background(), "UPDATE manager SET created_at=$1 WHERE email = $2", time.Now(), email)
	if err != nil {
		return nil, fmt.Errorf("failed to update login timestamp: %v", err)
	}

	return &user, nil
}
func GetUserProfile(userID int) (*model.Managers, error) {
	db := database.GetDBPool()

	var user model.Managers
	err := db.QueryRow(context.Background(), "SELECT id, email, column_name, manager_image_uri, company_name, company_image_uri FROM manager WHERE id = $1", userID).
		Scan(&user.ID, &user.Email, &user.Name, &user.ManagerImageURI, &user.CompanyName, &user.CompanyImageURI)

	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}
func UpdateUserProfile(userID int, email, name, userImageUri, companyName, companyImageUri *string) (*model.Managers, error) {
	db := database.GetDBPool()

	// Ambil data lama dulu
	var oldUser model.Managers
	err := db.QueryRow(context.Background(), "SELECT id, email, column_name, manager_image_uri, company_name, company_image_uri FROM manager WHERE id = $1", userID).
		Scan(&oldUser.ID, &oldUser.Email, &oldUser.Name, &oldUser.ManagerImageURI, &oldUser.CompanyName, &oldUser.CompanyImageURI)

	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Simpan hanya field yang dikirim dalam request
	fields := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1 // Posisi argumen untuk PostgreSQL ($1, $2, ...)

	// Jika parameter tidak nil, tambahkan ke query update
	if email != nil {
		fields = append(fields, fmt.Sprintf("email = $%d", argID))
		args = append(args, *email)
		argID++
	}
	if name != nil {
		fields = append(fields, fmt.Sprintf("column_name = $%d", argID))
		args = append(args, *name)
		argID++
	}
	if userImageUri != nil {
		fields = append(fields, fmt.Sprintf("manager_image_uri = $%d", argID))
		args = append(args, *userImageUri)
		argID++
	}
	if companyName != nil {
		fields = append(fields, fmt.Sprintf("company_name = $%d", argID))
		args = append(args, *companyName)
		argID++
	}
	if companyImageUri != nil {
		fields = append(fields, fmt.Sprintf("company_image_uri = $%d", argID))
		args = append(args, *companyImageUri)
		argID++
	}

	// Jika tidak ada field yang perlu diperbarui, langsung return user lama
	if len(fields) == 0 {
		return &oldUser, nil
	}

	// Bangun query UPDATE
	query := fmt.Sprintf("UPDATE manager SET %s WHERE id = $%d", strings.Join(fields, ", "), argID)
	args = append(args, userID)

	// Eksekusi query update
	_, err = db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Println("Error updating user profile:", err)
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}

	// Ambil kembali data yang sudah diperbarui
	var updatedUser model.Managers
	err = db.QueryRow(context.Background(), "SELECT id, email, column_name, manager_image_uri, company_name, company_image_uri FROM manager WHERE id = $1", userID).
		Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.Name, &updatedUser.ManagerImageURI, &updatedUser.CompanyName, &updatedUser.CompanyImageURI)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated user: %v", err)
	}

	return &updatedUser, nil
}

//func UpdateUserProfile(userID int, email, name, userImageUri, companyName, companyImageUri string) (*model.Managers, error) {
//	db := database.GetDBPool()
//
//	_, err := db.Exec(context.Background(), "UPDATE manager SET email = $1, column_name = $2, manager_image_uri = $3, company_name = $4, company_image_uri = $5 WHERE id = $6",
//		email, name, userImageUri, companyName, companyImageUri, userID)
//
//	if err != nil {
//		return nil, fmt.Errorf("failed to update user profile: %v", err)
//	}
//
//	var user model.Managers
//	err = db.QueryRow(context.Background(), "SELECT id, email, column_name, manager_image_uri, company_name, company_image_uri FROM manager WHERE id = $1", userID).
//		Scan(&user.ID, &user.Email, &user.Name, &user.ManagerImageURI, &user.CompanyName, &user.CompanyImageURI)
//
//	if err != nil {
//		return nil, fmt.Errorf("user not found: %v", err)
//	}
//
//	return &user, nil
//}
