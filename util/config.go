package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment             string        `mapstructure:"ENVIRONMENT"`
	DBDriver                string        `mapstructure:"DB_DRIVER"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	MigrationURL            string        `mapstructure:"MIGRATION_URL"`
	RedisAddress            string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress       string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress       string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey       string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	URL_LOCALHOST           string        `mapstructure:"URL_LOCALHOST"`
	GIN_MODE                string        `mapstructure:"GIN_MODE"`
	MINIO_ENDPOINT          string        `mapstructure:"MINIO_ENDPOINT"`
	MINIO_ACCESS_KEY_ID     string        `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MINIO_SECRET_ACCESS_KEY string        `mapstructure:"MINIO_SECRET_ACCESS_KEY"`
	MINIO_USE_SSL           bool          `mapstructure:"MINIO_USE_SSL"`
	MINIO_BUCKET_NAME       string        `mapstructure:"MINIO_BUCKET_NAME"`
	MINIO_URL_RESULT        string        `mapstructure:"MINIO_URL_RESULT"`
	EmailSenderName         string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress      string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	config.Environment = os.Getenv("ENVIRONMENT")
	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DBSource = os.Getenv("DB_SOURCE")
	config.MigrationURL = os.Getenv("MIGRATION_URL")

	config.HTTPServerAddress = os.Getenv("HTTP_SERVER_ADDRESS")
	config.GIN_MODE = os.Getenv("GIN_MODE")
	config.URL_LOCALHOST = os.Getenv("URL_LOCALHOST")
	config.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	config.MINIO_ENDPOINT = os.Getenv("MINIO_ENDPOINT")
	config.MINIO_ACCESS_KEY_ID = os.Getenv("MINIO_ACCESS_KEY_ID")
	config.MINIO_SECRET_ACCESS_KEY = os.Getenv("MINIO_SECRET_ACCESS_KEY")
	config.MINIO_USE_SSL = os.Getenv("MINIO_USE_SSL") == "true"
	config.MINIO_BUCKET_NAME = os.Getenv("MINIO_BUCKET_NAME")
	config.MINIO_URL_RESULT = os.Getenv("MINIO_URL_RESULT")
	config.EmailSenderName = os.Getenv("EMAIL_SENDER_NAME")
	config.EmailSenderAddress = os.Getenv("EMAIL_SENDER_ADDRESS")
	config.EmailSenderPassword = os.Getenv("EMAIL_SENDER_PASSWORD")

	// accessTokenDuration, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	// refreshTokenDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))

	// config.AccessTokenDuration = time.Duration(accessTokenDuration)
	// config.RefreshTokenDuration = time.Duration(refreshTokenDuration)
	return
}
