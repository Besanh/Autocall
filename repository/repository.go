package repository

import (
	sqlclient "autocall/internal/sql-client"

	"github.com/uptrace/bun"
)

var FusionSqlClient sqlclient.ISqlClientConn
var FreeswitchSqlClient sqlclient.ISqlClientConn
var MySqlClient sqlclient.ISqlClientConn

var Db *bun.DB
