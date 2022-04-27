package main

import (
	"autocall/api"
	apiv1 "autocall/api/v1"
	commonAuth "autocall/common/auth"
	IRedis "autocall/internal/redis"
	redis "autocall/internal/redis/driver"
	imysql "autocall/internal/sqldb/mysql"
	mysql "autocall/internal/sqldb/mysql/driver"
	"autocall/middleware/auth"
	"autocall/service"
	"io"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port     string
	LogType  string
	LogLevel string
	LogFile  string
	LogAddr  string
	DB       string
	Redis    string
	Auth     string
	AuthUrl  string
	DBConfig
}

type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	SSLMode         string
	Timeout         int
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime int
}

var config Config

func init() {
	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Dir:      config.Dir,
		Port:     viper.GetString(`main.port`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile:  viper.GetString(`main.log_file`),
		LogAddr:  viper.GetString(`main.log_addr`),
		DB:       viper.GetString(`main.db`),
		Redis:    viper.GetString(`main.redis`),
		Auth:     viper.GetString(`main.oauth`),
		AuthUrl:  viper.GetString(`main.auth_url`),
	}
	if cfg.DB == "enabled" {
		cfg.DBConfig = DBConfig{
			Driver:          viper.GetString(`db.driver`),
			Host:            viper.GetString(`db.host`),
			Port:            viper.GetString(`db.port`),
			Username:        viper.GetString(`db.username`),
			Password:        viper.GetString(`db.password`),
			Database:        viper.GetString(`db.name`),
			SSLMode:         viper.GetString(`db.disable`),
			Timeout:         viper.GetInt(`db.timeout`),
			MaxOpenConns:    viper.GetInt(`db.max_open_conns`),
			MaxIdleConns:    viper.GetInt(`db.max_idle_conns`),
			MaxConnLifetime: viper.GetInt(`db.conn_max_lifetime`),
		}
	}
	if cfg.Redis == "enabled" {
		var err error
		IRedis.Redis, err = redis.NewRedis(redis.Config{
			Addr:         viper.GetString(`redis.address`),
			Password:     viper.GetString(`redis.password`),
			DB:           viper.GetInt(`redis.database`),
			PoolSize:     30,
			PoolTimeout:  20,
			IdleTimeout:  10,
			ReadTimeout:  20,
			WriteTimeout: 15,
		})
		if err != nil {
			panic(err)
		}
	}
	switch cfg.Auth {
	case "oauth":
		auth.AuthUrl = cfg.AuthUrl
	}
	commonAuth.NewAuthUtil(commonAuth.Config{
		ExpiredTime: viper.GetInt(`oauth.expired_in`),
		TokenType:   viper.GetString(`oauth.tokenType`),
	})
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	mysqlconfig := mysql.MySqlConfig{
		Host:         config.DBConfig.Host,
		Database:     config.DBConfig.Database,
		User:         config.DBConfig.Username,
		Password:     config.DBConfig.Password,
		Port:         config.DBConfig.Port,
		Charset:      "utf8",
		PingInterval: config.DBConfig.MaxConnLifetime,
		MaxOpenConns: config.DBConfig.MaxOpenConns,
		MaxIdleConns: config.DBConfig.MaxIdleConns,
	}
	imysql.MySqlConnector = mysql.NewMySqlConnector(mysqlconfig)
	imysql.MySqlConnector.Ping()

	server := api.NewServer()

	apiv1.NewAuthHandler(server.Engine, service.NewUserService())
	apiv1.NewUserHandler(server.Engine, service.NewUserService())
	apiv1.NewContactCallHandler(server.Engine, service.NewContactService())

	server.Start(config.Port)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "GRAYLOG":
		gelfWriter, err := gelf.NewUDPWriter(cfg.LogAddr)
		if err != nil {
			log.Error("main", "setAppLogger", err.Error())
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(io.MultiWriter(os.Stdout, gelfWriter))
		}
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(os.Stdout)
		}
	case "GRAYLOG_FILE":
		gelfWriter, err := gelf.NewUDPWriter(cfg.LogAddr)
		if file != nil && err == nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file, gelfWriter))
		} else if file == nil && err == nil {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(io.MultiWriter(os.Stdout, gelfWriter))
		} else if file != nil && err != nil {
			log.Error("main", "setAppLogger", err.Error())
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}
