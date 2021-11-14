package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type IPForwarding struct {
	Enabled bool   `mapstructure:"enabled"`
	IP      string `mapstructure:"ip"`
	Port    string `mapstructure:"port"`
}


type ServerConfig struct {
	Addr            string
	WriteTimeout    int
	ReadTimeout     int
	GraceFulTimeout int
	Registration    bool
	Debug 			bool
}

type ImpersonationConfig struct {
	Password string
	Admin    bool
	User     bool
}

type DBConfig struct {
	Driver          string
	Name            string
	Host            string
	Port 			string
	Username		string
	Password		string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

type DBSecondary struct {
	Driver          string
	Name            string
	Host            string
	Port 			string
	Username		string
	Password		string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

type Config struct {
	Server            ServerConfig
	DB                DBConfig
	DBSecondary                DBSecondary
	Redis             RedisServer
	JWTConfig         JWTConfig
	Mailer            Mailer
	ResetPassword     ResetPassword
	NoSQL 			  NoSQLConfig
	BNIConfig		  BNIConfig
}

type NoSQLConfig struct {
	Driver	string
	Host 	string
	Port 	int
	Username	string
	Password 	string
}

type RedisServer struct {
	Addr     string
	Password string
	Timeout  int
	MaxIdle  int
}

// JWTConfig is JWT configuration object
type JWTConfig struct {
	Issuer            string
	Secret            string
	TokenLifeTimeHour int
}

type Mailer struct {
	Server     string
	Port       int
	Username   string
	Password   string
	UseTls     bool
	Sender     string
	MaxAttempt int
}

type BNIConfig struct {
	Key     string
	Url       string
	ClientID   string
	Cid   	 string
}

type ResetPassword struct {
	UserLink  string
	AdminLink string
}


func InitConfig() Config {
	viper.SetConfigName("config/.env")
	if os.Getenv("ENV") == "staging" {
		viper.SetConfigName(".env.yml-" + "staging")
	}

	if os.Getenv("ENV") == "production" {
		viper.SetConfigName(".env.yml-" + "production")
	}

	viper.AddConfigPath(".")

	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return configuration
}
