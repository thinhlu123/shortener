package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	Log    Logger
	Url    UrlConfig
}

type UrlConfig struct {
	Host string
}

type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

type DBConfig struct {
	DBHost    []string `json:"db_host"` // host1:port1,host2:port2 for multiple address
	DBUser    string   `json:"db_user"`
	DBPass    string   `json:"db_pass"`
	DBName    string   `json:"db_name"`
	DBAuth    string   `json:"db_auth"`
	EnableSSL bool     `json:"ssl"` // true | false
}

type RedisConfig struct {
	RedisAddr      string
	RedisPort      string
	RedisUser      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Load config file from given path
func loadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

func GetConfig(cfgPath string) error {
	if Conf == nil {
		cfgFile, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		Conf, err = parseConfig(cfgFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetConfigPath(configPath string) string {
	if configPath == "prd" {
		return ""
	} else if configPath == "docker" {
		return "config/config_docker"
	}
	return "config/config_local"
}
