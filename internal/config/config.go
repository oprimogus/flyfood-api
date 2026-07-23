package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/subosito/gotenv"
)

var (
	cfg   *Config
	cOnce sync.Once
)

type Config struct {
	API              *API
	Postgres         *Postgres
	Zitadel          *Zitadel
	AWS              *AWS
	ExternalServices *ExternalServices
}

func Get() *Config {
	cOnce.Do(func() {
		cfg = newConfig()
	})

	return cfg
}

type API struct {
	Port        string
	BasePath    string
	Name        string
	Environment string
	SQLCDebug   string
}

type Postgres struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
}

func (d *Postgres) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Host", "redacted"),
		slog.String("Port", "redacted"),
		slog.String("Name", "redacted"),
		slog.String("User", "redacted"),
		slog.String("Password", "redacted"),
	)
}

type Zitadel struct {
	Domain                string
	Port                  string
	ClientID              string
	ProjectID             string
	KeyPath               string
	ServiceAccountKeyPath string
}

type AWS struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionKey      string
}

func (a *AWS) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("AccessKeyID", "redacted"),
		slog.String("SecretAccessKey", "redacted"),
		slog.String("SessionKey", "redacted"),
	)
}

type ExternalService struct {
	BaseURL string
	ApiKey  string
}

func (a *ExternalService) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ApiKey", "redacted"),
	)
}

type ExternalServices struct {
	Resend    *ExternalService
	Nominatim *ExternalService
}

func newConfig() *Config {
	err := gotenv.Load(".env")
	if err != nil {
		slog.Info("arquivo .env não encontrado, usando variáveis de ambiente")
	}
	cfg := &Config{
		API: &API{
			Port:        os.Getenv("API_PORT"),
			BasePath:    os.Getenv("API_BASE_PATH"),
			Name:        os.Getenv("API_NAME"),
			Environment: os.Getenv("API_ENVIRONMENT"),
			SQLCDebug:   os.Getenv("API_SQLC_DEBUG"),
		},
		Postgres: &Postgres{
			Host:         os.Getenv("POSTGRES_HOST"),
			Port:         os.Getenv("POSTGRES_PORT"),
			UserName:     os.Getenv("POSTGRES_USERNAME"),
			Password:     os.Getenv("POSTGRES_PASSWORD"),
			DatabaseName: os.Getenv("POSTGRES_DATABASE_NAME"),
		},
		Zitadel: &Zitadel{
			Domain:                os.Getenv("ZITADEL_DOMAIN"),
			Port:                  os.Getenv("ZITADEL_PORT"),
			ClientID:              os.Getenv("ZITADEL_CLIENT_ID"),
			ProjectID:             os.Getenv("ZITADEL_PROJECT_ID"),
			KeyPath:               os.Getenv("ZITADEL_KEY_PATH"),
			ServiceAccountKeyPath: os.Getenv("ZITADEL_SERVICE_ACCOUNT_KEY_PATH"),
		},
		AWS: &AWS{
			Region:          os.Getenv("AWS_REGION"),
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		ExternalServices: &ExternalServices{
			Resend: &ExternalService{
				BaseURL: os.Getenv("RESEND_BASE_URL"),
				ApiKey:  os.Getenv("RESEND_API_KEY"),
			},
			Nominatim: &ExternalService{
				BaseURL: os.Getenv("NOMINATIM_BASE_URL"),
				ApiKey:  os.Getenv("NOMINATIM_API_KEY"),
			},
		},
	}
	return cfg
}
