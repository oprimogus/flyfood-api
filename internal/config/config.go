package config

import (
	"log/slog"
	"os"

	"github.com/subosito/gotenv"
)

var (
	conf *Config
)

type DBConf struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (d *DBConf) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Host", "redacted"),
		slog.String("Port", "redacted"),
		slog.String("Name", "redacted"),
		slog.String("User", "redacted"),
		slog.String("Password", "redacted"),
	)
}

type APIConf struct {
	ServiceName string
	BasePath    string
	Port        string
	GinMode     string
	Environment string
	sqlcDebug   string
	Consts      map[string]string
}

type zitadelConfig struct {
	Issuer                string
	Api                   string
	Domain                string
	Port                  string
	ClientID              string
	KeyPath               string
	ServiceAccountKeyPath string
	ProjectID             string
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

type ExternalService struct {
	BaseURL string
}

type Nominatim struct {
	ExternalService
}

type Config struct {
	Database  *DBConf
	Api       *APIConf
	Zitadel   *zitadelConfig
	Resend    *resendConfig
	Aws       *aws
	Nominatim *Nominatim
}

func newConfig() *Config {
	err := gotenv.Load(".env")
	if err != nil {
        slog.Info("arquivo .env não encontrado, usando variáveis de ambiente")
    }
	return &Config{
		Database: &DBConf{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Api: &APIConf{
			ServiceName: os.Getenv("API_SERVICE_NAME"),
			BasePath:    os.Getenv("API_BASE_PATH"),
			Port:        os.Getenv("API_PORT"),
			GinMode:     os.Getenv("GIN_MODE"),
			Environment: os.Getenv("ENVIRONMENT"),
			sqlcDebug:   os.Getenv("SQLCDEBUG"),
		},
		Zitadel: &zitadelConfig{
			Issuer:                os.Getenv("ZITADEL_ISSUER"),
			Api:                   os.Getenv("ZITADEL_API"),
			Domain:                os.Getenv("ZITADEL_DOMAIN"),
			Port:                  os.Getenv("ZITADEL_PORT"),
			ClientID:              os.Getenv("ZITADEL_CLIENT_ID"),
			KeyPath:               os.Getenv("ZITADEL_KEY"),
			ServiceAccountKeyPath: os.Getenv("ZITADEL_SERVICE_ACCOUNT_KEY"),
			ProjectID:             os.Getenv("ZITADEL_PROJECT_ID"),
		},
		Resend: &resendConfig{
			apiKey: os.Getenv("RESEND_API_KEY"),
		},
		Aws: &aws{
			region:          os.Getenv("AWS_REGION"),
			accessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			secretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		Nominatim: &Nominatim{
			ExternalService: ExternalService{
				BaseURL: os.Getenv("NOMINATIM_BASE_URL"),
			},
		},
	}
}

func GetInstance() *Config {
	if conf == nil {
		conf = newConfig()
	}
	return conf
}
