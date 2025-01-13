package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"project_sprint/api/v1"
	"project_sprint/internal/handler"
	"project_sprint/pkg/database"
)

func main() {
	// Inisialisasi database
	database.InitDB()
	defer database.CloseDB()

	//cfg := config.LoadEnv()
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-2"),
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})

	if err != nil {
		log.Fatal("Error creating AWS session: ", err)
	}

	s3Client := s3.New(sess)

	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Println("Error listing S3 buckets:", err)
		return
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf(" - %s\n", *bucket.Name)
	}

	router := gin.Default()
	router.POST("/auth", handler.AuthHandler) // Login & Register

	v1Group := router.Group("/v1")
	{
		v1.RegisterDepartmentRoutes(v1Group)
		v1.RegisterManagerRoutes(v1Group)
		v1.RegisterFileRoutes(v1Group)
		v1.RegisterEmployeeRoutes(v1Group)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on http://localhost:%s", port)
	router.Run(":" + port)
}
