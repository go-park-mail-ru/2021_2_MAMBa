package database

import (
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbUser     = "dev"
	dbPassword = "1234"
	dbHost     = "localhost"
	dbPport    = 5432
	dbName     = "film4u"
)

type ConnectionPool interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type DBManager struct {
	Pool ConnectionPool
}

func Connect() *DBManager {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", dbUser, dbPassword, dbHost, dbPport, dbName)
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		mylog.Warn("Postgres error")
		mylog.Error(err)
		return nil
	}
	mylog.Info("Successful connection to postgres")
	return &DBManager{Pool: pool}
}

func (dbm *DBManager) Disconnect() {
	dbm.Pool.Close()
	mylog.Info("Postgres disconnected")
}

func (dbm *DBManager) Query(queryString string, params ...interface{}) ([][][]byte, error) {
	transactionContext := context.Background()
	tx, err := dbm.Pool.Begin(transactionContext)
	if err != nil {
		mylog.Warn("Error connecting to a pool")
		mylog.Error(err)
		return nil, err
	}
	// ВАЖНО - Rollback не проходит после commit
	defer tx.Rollback(transactionContext)

	rows, err := tx.Query(transactionContext, queryString, params...)
	if err != nil {
		mylog.Warn("Error in query")
		mylog.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([][][]byte, 0)
	for rows.Next() {
		rowBuffer := make([][]byte, 0)
		rowBuffer = append(rowBuffer, rows.RawValues()...)
		result = append(result, rowBuffer)
	}

	err = tx.Commit(transactionContext)
	if err != nil {
		mylog.Warn("Error committing")
		mylog.Error(err)
		return nil, err
	}
	return result, nil
}

func (dbm *DBManager) Execute(queryString string, params ...interface{}) error {
	transactionContext := context.Background()
	tx, err := dbm.Pool.Begin(transactionContext)
	if err != nil {
		mylog.Warn("Error connecting to a pool")
		mylog.Error(err)
		return err
	}
	// ВАЖНО - Rollback не проходит после commit
	defer tx.Rollback(transactionContext)

	_, err = tx.Exec(transactionContext, queryString, params...)
	if err != nil {
		mylog.Warn("Error in query")
		mylog.Error(err)
		return err
	}

	err = tx.Commit(transactionContext)
	if err != nil {
		mylog.Warn("Error committing")
		mylog.Error(err)
		return err
	}
	return nil
}
