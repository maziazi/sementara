package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"project_sprint/pkg/config"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

// InitDB menginisialisasi koneksi database menggunakan connection pooling (hanya sekali)
func InitDB() {
	once.Do(func() {
		// Memuat konfigurasi dari .env
		cfg := config.LoadEnv()

		// Buat connection string PostgreSQL dengan konfigurasi optimal
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

		// Konfigurasi pool dengan opsi tambahan (timeout, max connections, dll.)
		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			log.Fatalf("Error parsing database config: %v", err)
		}

		// Atur parameter koneksi (sesuaikan dengan kebutuhan aplikasi)
		poolConfig.MaxConns = 10                       // Maksimal 10 koneksi
		poolConfig.MinConns = 2                        // Minimal 2 koneksi
		poolConfig.MaxConnLifetime = 30 * time.Minute  // Maksimal umur koneksi 30 menit
		poolConfig.MaxConnIdleTime = 5 * time.Minute   // Koneksi idle selama 5 menit akan ditutup
		poolConfig.HealthCheckPeriod = 1 * time.Minute // Cek kesehatan koneksi tiap 1 menit

		// Buat connection pool
		dbPool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		log.Println("✅ Connected to database successfully")
	})
}

// CloseDB menutup koneksi database saat aplikasi berhenti
func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("🛑 Database connection closed")
	}
}

// GetDBPool mengembalikan instance database pool
func GetDBPool() *pgxpool.Pool {
	if dbPool == nil {
		log.Println("⚠️ Database connection is not initialized, calling InitDB()")
		InitDB()
	}
	return dbPool
}
