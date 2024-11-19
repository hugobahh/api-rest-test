package mysql

import (
	"api-rest-test/internal/app/queries"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Importa el controlador MySQL
)

type MySQLConnector struct {
	DBClient     *sql.DB
	log          Logger
	shutdownOnce sync.Once
	timeout      time.Duration
	statements   map[string]*sql.Stmt
}

var (
	mySQLConnector *MySQLConnector
	once           sync.Once
)

func GetMySQLConnector(cnnDataSource *DataSource, log Logger) *MySQLConnector {
	once.Do(func() {
		mySQLConnector, _ = NewMySQLConnector(cnnDataSource, log)
	})
	return mySQLConnector
}

func NewMySQLConnector(cnnDataSource *DataSource, log Logger) (*MySQLConnector, error) {
	cnnDB := &MySQLConnector{log: log, timeout: cnnDataSource.Timeout}
	cnnDB.setDBConnection(cnnDataSource)
	cnnDB.setupDBConnection(cnnDataSource)
	var err error
	statements, err := cnnDB.InitStatements(time.Duration(cnnDataSource.TimeoutQuery) * time.Second)
	if err != nil {
		return nil, err
	}

	cnnDB.statements = statements
	return cnnDB, nil
}

func (c *MySQLConnector) setDBConnection(cnnDataSource *DataSource) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&readTimeout=%ds",
		cnnDataSource.User, url.QueryEscape(cnnDataSource.Password), cnnDataSource.Host, cnnDataSource.Port, cnnDataSource.Name,
		int(cnnDataSource.Timeout.Seconds()), int(cnnDataSource.Timeout.Seconds()))
	//if cnnDataSource.ReadOnly {
	//	dsn += "&ApplicationIntent=ReadOnly"
	//}
	//log.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		c.log.Fatal("MySQLConnector", "setDBConnection", err)
	}
	c.DBClient = db
}

func (c *MySQLConnector) setupDBConnection(cnnConfig *DataSource) {
	c.DBClient.SetMaxOpenConns(cnnConfig.MaxOpenConnections)
	c.DBClient.SetMaxIdleConns(cnnConfig.MaxIdleConnections)
	c.DBClient.SetConnMaxLifetime(cnnConfig.MaxLifetimeConnections)
}

func (c *MySQLConnector) InitStatements(timeout time.Duration) (map[string]*sql.Stmt, error) {
	statements := make(map[string]*sql.Stmt)
	preparedQueries := map[string]string{
		"Register": queries.Register,
		"ExistUsr": queries.ExistUsr,
		"Cancel":   queries.Cancel,
	}
	for name, query := range preparedQueries {
		stmt, err := c.prepareContext(timeout, query)
		if err != nil {
			c.log.Error("SQLConnector", "InitializeStatements", err)
			return nil, err
		}
		statements[name] = stmt
	}
	return statements, nil
}

func (c *MySQLConnector) prepareContext(timeout time.Duration, query string) (*sql.Stmt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stmt, err := c.DBClient.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (c *MySQLConnector) GetStatements(name string) (*sql.Stmt, bool) {
	stmt, ok := c.statements[name]
	return stmt, ok
}

// ============================================
func (m *MySQLConnector) HealthCheck() error {
	if m.DBClient != nil {
		if err := m.DBClient.Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (m *MySQLConnector) Shutdown() error {
	var err error
	m.shutdownOnce.Do(func() {
		log.Println("Shutting down MySQL client...")
		if err = m.DBClient.Close(); err != nil {
			log.Println("MySQLConnector", "Shutdown", err)
		} else {
			log.Println("MySQL client is now disconnected")
		}
	})
	return err
}
