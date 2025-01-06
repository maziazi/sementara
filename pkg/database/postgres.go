package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"project_sprint/pkg/config"
)

var DB *pgxpool.Pool

// InitDB menginisialisasi koneksi database menggunakan connection pooling
func InitDB() {
	// Memuat konfigurasi dari .env
	cfg := config.LoadEnv()

	// Buat connection string PostgreSQL
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	// Koneksi ke database dengan connection pool
	var err error
	DB, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")
}

// CloseDB menutup koneksi database saat aplikasi berhenti
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
