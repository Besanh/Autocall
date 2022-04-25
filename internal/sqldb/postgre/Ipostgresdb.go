package postgre

import (
	"context"

	pg "github.com/go-pg/pg/v10"
)

type IPgSqlConnector interface {
	GetConn() *pg.DB
	Ping() error
	QueryOne(model, query interface{}, params ...interface{}) (interface{}, error)
	Query(model, query interface{}, params ...interface{}) (interface{}, error)
	QueryOneContext(c context.Context, model, query interface{}, params ...interface{}) (interface{}, error)
	QueryContext(c context.Context, model, query interface{}, params ...interface{}) (interface{}, error)
	HandleError(err error) (int, error)
}

var FusionSqlConnector IPgSqlConnector
var FreeswitchSqlConnector IPgSqlConnector
