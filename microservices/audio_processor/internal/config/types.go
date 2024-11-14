package config

import (
	"strconv"
	"time"
)

type (
	// Config es la estructura principal que agrupa todas las configuraciones
	Config struct {
		Environment string
		Service     ServiceConfig
		AWS         AWSConfig
		Messaging   MessagingConfig
		Storage     StorageConfig
		Database    DatabaseConfig
		API         APIConfig
		GinConfig   GinConfig
	}

	// ServiceConfig contiene configuración general del servicio
	ServiceConfig struct {
		MaxAttempts int
		Timeout     time.Duration
	}

	GinConfig struct {
		Mode string
	}

	// AWSConfig contiene todas las configuraciones relacionadas con AWS
	AWSConfig struct {
		Region string
	}

	// MessagingConfig maneja la configuración de mensajería
	MessagingConfig struct {
		Type  string // kafka o sqs
		Kafka *KafkaConfig
		SQS   *SQSConfig
	}

	// KafkaConfig configuración específica de Kafka
	KafkaConfig struct {
		Brokers []string
		Topic   string
	}

	// SQSConfig configuración específica de SQS
	SQSConfig struct {
		QueueURL string
	}

	// StorageConfig maneja la configuración de almacenamiento
	StorageConfig struct {
		Type        string // s3 o local
		S3Config    *S3Config
		LocalConfig *LocalConfig
	}

	LocalConfig struct {
		BasePath string
	}

	// S3Config configuración específica de S3
	S3Config struct {
		BucketName string
	}

	// DatabaseConfig maneja la configuración de base de datos
	DatabaseConfig struct {
		Type     string
		Mongo    *MongoConfig
		DynamoDB *DynamoDBConfig
	}

	// MongoConfig configuración específica de MongoDB
	MongoConfig struct {
		User        string
		Password    string
		Port        string
		Host        []string
		Database    string
		Collections Collections
	}

	// DynamoDBConfig configuración específica de DynamoDB
	DynamoDBConfig struct {
		Tables Tables
	}

	// Collections nombres de colecciones para MongoDB
	Collections struct {
		Songs      string
		Operations string
	}

	// Tables nombres de tablas para DynamoDB
	Tables struct {
		Songs      string
		Operations string
	}

	// APIConfig maneja la configuración de APIs externas
	APIConfig struct {
		YouTube YouTubeConfig
		OAuth2  OAuth2Config
	}

	// YouTubeConfig configuración específica de YouTube API
	YouTubeConfig struct {
		ApiKey string
	}

	// OAuth2Config configuración de OAuth2
	OAuth2Config struct {
		Enabled string
	}
)

func (c OAuth2Config) ParseBool() bool {
	oAuthEnabled, err := strconv.ParseBool(c.Enabled)
	if err != nil {
		return false
	}
	return oAuthEnabled
}
