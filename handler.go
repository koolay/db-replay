package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-mysql/mysql"
)

type MysqlHandler struct {
	db *sql.DB
}

func (h *MysqlHandler) UseDB(dbName string) error {
	return nil
}

func (h *MysqlHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Println("HandleQuery", query)
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}

	result := &mysql.Result{
		Status: mysql.SERVER_STATUS_AUTOCOMMIT,
	}
	return result, nil
	// return nil, fmt.Errorf("not supported now")
}

func (h *MysqlHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	log.Println("HandleFieldList")
	return nil, fmt.Errorf("not supported now")
}

func (h *MysqlHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	log.Println("HandleStmtPrepare")
	return 0, 0, nil, fmt.Errorf("not supported now")
}

func (h *MysqlHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	log.Println("HandleStmtExecute")
	return nil, fmt.Errorf("not supported now")
}

func (h *MysqlHandler) HandleStmtClose(context interface{}) error {
	log.Println("HandleStmtClose")
	return nil
}

func (h *MysqlHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Println("HandleOtherCommand")
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}
