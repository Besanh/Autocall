package postgre

import (
	"context"
	"time"

	postgre "autocall/internal/sqldb/postgre"

	pg "github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

type PgConfig struct {
	Host                  string
	Port                  string
	Database              string
	User                  string
	Password              string
	MaxRetries            int
	RetryStatementTimeout bool
	DialTimeout           int
	ReadTimeout           int
	WriteTimeout          int
	PoolSize              int
	PoolTimeout           int
}

type FusionSqlConnector struct {
	Host                  string
	Port                  string
	Database              string
	User                  string
	Password              string
	MaxRetries            int
	RetryStatementTimeout bool
	DialTimeout           int
	ReadTimeout           int
	WriteTimeout          int
	PoolSize              int
	PoolTimeout           int
	Sqldb                 *pg.DB
}

func NewPgSqlConnector(config PgConfig) postgre.IPgSqlConnector {
	Pq := &FusionSqlConnector{
		Host:                  config.Host,
		Port:                  config.Port,
		User:                  config.User,
		Password:              config.Password,
		Database:              config.Database,
		MaxRetries:            config.MaxRetries,
		RetryStatementTimeout: config.RetryStatementTimeout,
		DialTimeout:           config.DialTimeout,
		ReadTimeout:           config.ReadTimeout,
		WriteTimeout:          config.WriteTimeout,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout,
	}
	Pq.Connect()
	return Pq
}
func (conn *FusionSqlConnector) GetConn() *pg.DB {
	return conn.Sqldb
}

func (conn *FusionSqlConnector) Connect() {
	opts := pg.Options{
		User:                  conn.User,
		Password:              conn.Password,
		Database:              conn.Database,
		Addr:                  conn.Host + ":" + conn.Port,
		MaxRetries:            conn.MaxRetries,
		RetryStatementTimeout: conn.RetryStatementTimeout,
		DialTimeout:           time.Duration(conn.DialTimeout) * time.Second,
		ReadTimeout:           time.Duration(conn.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(conn.WriteTimeout) * time.Second,
		PoolSize:              conn.PoolSize,
		PoolTimeout:           time.Duration(conn.PoolTimeout) * time.Second,
	}
	conn.Sqldb = pg.Connect(&opts)
}

func (conn *FusionSqlConnector) Ping() error {
	_, err := conn.Sqldb.Exec("SELECT 1")
	if err == nil {
		log.Info("Connect database postgre successful :", conn.Host)
	} else {
		log.Error("Connect database postgre fail :", err)
	}
	return err
}

func (conn *FusionSqlConnector) QueryOne(model, query interface{}, params ...interface{}) (interface{}, error) {
	result, err := conn.Sqldb.QueryOne(model, query, params...)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (conn *FusionSqlConnector) Query(model, query interface{}, params ...interface{}) (interface{}, error) {
	result, err := conn.Sqldb.Query(model, query, params...)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (conn *FusionSqlConnector) QueryOneContext(c context.Context, model, query interface{}, params ...interface{}) (interface{}, error) {
	result, err := conn.Sqldb.QueryOneContext(c, model, query, params...)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (conn *FusionSqlConnector) QueryContext(c context.Context, model, query interface{}, params ...interface{}) (interface{}, error) {
	result, err := conn.Sqldb.QueryContext(c, model, query, params...)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (conn *FusionSqlConnector) HandleError(err error) (int, error) {
	if err == pg.ErrNoRows {
		return 404, nil
	}
	return 0, err
}
