package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
)

const (
	numberDuplicateKey = 1062
)

func IsDuplicateKeyError(err error) bool {
	mysqlErr, ok := lo.ErrorsAs[*mysql.MySQLError](err)
	if ok {
		return mysqlErr.Number == numberDuplicateKey
	}

	return false
}
