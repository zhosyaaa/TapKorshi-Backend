package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute
	defaultVerificationCodeLength = 8

	EnvLocal = "env"
	Prod     = "prod"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Database DatabaseConfig
	}

	HTTPConfig struct {
		Port           string        `mapstructure:"port"`
		MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
		ReadTimeout    time.Duration `mapstructure:"readTimeout"`
		WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	}

	DatabaseConfig struct {
		Name     string
		Host     string
		Port     string
		Sslmode  string
		User     string
		Password string
	}
	EmailConfig struct {
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	EmailTemplates struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	EmailSubjects struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}
)

// Init initializes HTTP and Database configurations from environment variables.
func Init(configsDir string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	populateDefaults()
	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)

	fmt.Println(cfg)
	return &cfg, nil
}

func unmarshal(cfg *Config) error {

	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) {
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Sslmode = os.Getenv("DB_SSLMODE")

	fmt.Println(os.Getenv("DB_NAME"))
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")
}
func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
