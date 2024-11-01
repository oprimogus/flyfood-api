package config

import (
	"log/slog"
	"os"

	"github.com/subosito/gotenv"

	"github.com/oprimogus/cardapiogo/internal/utils"
)

var (
	conf *Config
)

type dbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (d *dbConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Host", "redacted"),
		slog.String("Port", "redacted"),
		slog.String("Name", "redacted"),
		slog.String("User", "redacted"),
		slog.String("Password", "redacted"),
	)
}

type apiConfig struct {
	basePath    string
	port        string
	ginMode     string
	Environment string
	sqlcDebug   string
	Consts      map[string]string
}

func (a *apiConfig) BasePath() string {
	return a.basePath
}

func (a *apiConfig) Port() string {
	return a.port
}

func (a *apiConfig) GinMode() string {
	return a.ginMode
}

func (a *apiConfig) SQLCDebug() string {
	return a.sqlcDebug
}

type keycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

func (d *keycloakConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("BaseURL", d.BaseURL),
		slog.String("Realm", d.Realm),
		slog.String("ClientID", "redacted"),
		slog.String("ClientSecret", "redacted"),
	)
}

type resendConfig struct {
	apiKey string
}

func (r *resendConfig) APIKey() string {
	return r.apiKey
}

type aws struct {
	region          string
	accessKeyID     string
	secretAccessKey string
	sessionKey      string
}

func (a *aws) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("accessKeyID", "redacted"),
		slog.String("secretAccessKey", "redacted"),
		slog.String("sessionKey", "redacted"),
	)
}

func (a *aws) Region() string {
	return a.region
}

func (a *aws) AccessKeyID() string {
	return a.accessKeyID
}

func (a *aws) SecretAccessKey() string {
	return a.secretAccessKey
}

func (a *aws) SessionKey() string {
	return a.sessionKey
}

type Config struct {
	Database *dbConfig
	Api      *apiConfig
	Keycloak *keycloakConfig
	Resend   *resendConfig
	Aws      *aws
}

func newConfig() *Config {
	err := utils.SetWorkingDirToProjectRoot()
	if err != nil {
		panic("fail on set project root as workdir")
	}
	err = gotenv.Load(".env")
	if err != nil {
		slog.Error("fail on load env vars: %s", "err", err)
		panic("fail on load env vars")
	}
	return &Config{
		Database: &dbConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Api: &apiConfig{
			basePath:    os.Getenv("API_BASE_PATH"),
			port:        os.Getenv("API_PORT"),
			ginMode:     os.Getenv("GIN_MODE"),
			Environment: os.Getenv("ENVIRONMENT"),
			sqlcDebug:   os.Getenv("SQLCDEBUG"),
		},
		Keycloak: &keycloakConfig{
			BaseURL:      os.Getenv("KEYCLOAK_BASE_URL"),
			Realm:        os.Getenv("KEYCLOAK_REALM"),
			ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
			ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		},
		Resend: &resendConfig{
			apiKey: os.Getenv("RESEND_API_KEY"),
		},
		Aws: &aws{
			region:          os.Getenv("AWS_REGION"),
			accessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			secretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}

func GetInstance() *Config {
	if conf == nil {
		conf = newConfig()
	}
	return conf
}
