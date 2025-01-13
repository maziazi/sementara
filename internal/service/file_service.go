package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"project_sprint/internal/model"
	"project_sprint/pkg/database"
)

const (
	S3BucketName = "penopangsistemuii" // Ganti dengan nama bucket S3 Anda
)

// AddFile untuk menyimpan URI file ke database
//func AddFile(uri string) (*model.File, error) {
//	db := database.GetDBPool()
//
//	var file model.File
//	err := db.QueryRow(context.Background(), "INSERT INTO file (uri) VALUES ($1) RETURNING id, uri", uri).Scan(&file.ID, &file.URI)
//	if err != nil {
//		return nil, err
//	}
//
//	return &file, nil
//}

func AddFile(fileURL string) (*model.File, error) {
	// Simpan langsung URL file ke database tanpa mengunggah ulang
	db := database.GetDBPool()
	var file model.File
	err := db.QueryRow(context.Background(), "INSERT INTO file (uri) VALUES ($1) RETURNING id, uri", fileURL).Scan(&file.ID, &file.URI)
	if err != nil {
		log.Printf("Error inserting file into database: %v", err)
		return nil, err
	}

	log.Printf("File stored in database: ID = %d, URI = %s", file.ID, file.URI)
	return &file, nil
}

//func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
//	// Load AWS configuration dari environment atau shared credentials
//	cfg, err := config.LoadDefaultConfig(context.TODO())
//	if err != nil {
//		log.Fatalf("unable to load SDK config, %v", err)
//		return "", err
//	}
//
//	// Inisialisasi client S3
//	s3Client := s3.NewFromConfig(cfg)
//
//	// Baca file ke dalam buffer
//	buffer := bytes.NewBuffer(nil)
//	if _, err := buffer.ReadFrom(file); err != nil {
//		return "", fmt.Errorf("failed to read file: %v", err)
//	}
//
//	// Generate nama file unik (opsional, bisa diubah sesuai kebutuhan)
//	fileName := fmt.Sprintf("%d%s", os.Getpid(), filepath.Ext(fileHeader.Filename))
//
//	// Kirim file ke S3
//	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
//		Bucket: aws.String(S3BucketName),
//		Key:    aws.String(fileName),
//		Body:   bytes.NewReader(buffer.Bytes()),
//		ACL:    "public-read", // Jika ingin file bisa diakses publik
//	})
//	if err != nil {
//		return "", fmt.Errorf("failed to upload file to S3: %v", err)
//	}
//
//	// URL akses file di S3
//	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", S3BucketName, fileName)
//	return fileURL, nil
//}

// GetFileByID untuk mengambil file berdasarkan ID
func GetFileByID(id int) (*model.File, error) {
	db := database.GetDBPool()

	var file model.File
	err := db.QueryRow(context.Background(), "SELECT id, uri FROM file WHERE id = $1", id).Scan(&file.ID, &file.URI)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return &file, nil
}

// DeleteFile untuk menghapus file berdasarkan ID
func DeleteFile(id int) error {
	db := database.GetDBPool()

	result, err := db.Exec(context.Background(), "DELETE FROM file WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("file not found")
	}

	return nil
}
