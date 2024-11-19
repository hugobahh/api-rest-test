package repository

import (
	"api-rest-test/internal/config"
	mysql "api-rest-test/pkg/database/mysql"
	"api-rest-test/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type IUserRepository interface {
	Register(context.Context, string, string, string, string) error
	Search(context.Context, string, string) ([]map[string]interface{}, error)
	ExistUser(context.Context, string, string) (bool, error)
}

type UserRepository struct {
	config         *config.Configuration
	Log            logger.Logger
	mysqlConnector *mysql.MySQLConnector
	statements     map[string]*sql.Stmt
	timeout        time.Duration
	Id_Pay         uint64     `json:"id_pay" gorm:"primaryKey"`
	Id_Res         int32      `json:"id_res"`
	Id_Inm         int32      `json:"id_inmueble"`
	Name           string     `json:"name"`
	Monto          float32    `json:"monto"`
	Date_Reg       *time.Time `json:"date_reg"`
	Id_Client      int32      `json:"id_client"`
}

func NewUserRepository(config *config.Configuration, log *logger.Log) *UserRepository {
	sqlConfig := mysql.DataSource{
		Host:         config.MySQLDataSource.Host,
		User:         config.MySQLDataSource.User,
		Password:     config.MySQLDataSource.Password,
		Name:         config.MySQLDataSource.Name,
		Port:         config.MySQLDataSource.Port,
		ReadOnly:     config.MySQLDataSource.ReadOnly,
		Timeout:      config.MySQLDataSource.Timeout,
		TimeoutQuery: config.MySQLDataSource.TimeoutQuery,
	}
	statements := make(map[string]*sql.Stmt)
	return &UserRepository{
		Log:            log,
		mysqlConnector: mysql.GetMySQLConnector(&sqlConfig, log),
		statements:     statements,
		timeout:        config.MySQLDataSource.Timeout,
	}
}

func (ur *UserRepository) Register(ctx context.Context, User string, Mail string, Tel string, Pwd string) error {
	ctx, cancel := context.WithTimeout(ctx, ur.timeout)
	defer cancel()
	stmt, ok := ur.mysqlConnector.GetStatements("Register")
	if !ok {
		return fmt.Errorf("statement 'Register' not found")
	}

	_, err := stmt.QueryContext(ctx, User, Mail, Tel, Pwd)
	if err != nil {
		logger.GetLogger().Error("RegisterEntrance", "Error executing query", err)
		return err
	}
	return nil
	//return nil
	ur.Log.Infof("Repository-RegisterEntrance...")
	return nil
}

func (dr *UserRepository) Search(ctxNo context.Context, Mail string, Tel string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	var bOK = false

	ctx, errCtx := context.WithTimeout(context.Background(), time.Duration(dr.timeout)*time.Second)
	log.Println(errCtx)

	stmt, ok := dr.mysqlConnector.GetStatements("SearchMailTel")
	if !ok || stmt == nil {
		return nil, fmt.Errorf("statement 'Select' not found or is nil")
	}
	//rows, err := stmt.Query(Id, IdRes)
	rows, err := stmt.QueryContext(ctx, Mail, Tel)
	if err != nil {
		logger.GetLogger().Error("SelectRepository", "Get", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		mapJson := make(map[string]interface{})
		var id uint64
		var idClient, idHouse int32
		var incident sql.NullString
		var report sql.NullString
		var path sql.NullString
		var dateReg sql.NullString
		var st string

		if err := rows.Scan(&id, &idClient, &idHouse, &incident, &report, &path, &dateReg, &st); err != nil {
			logger.GetLogger().Error("SelectRepository", "Scan", err)
			return nil, err
		}

		mapJson["id_Incident"] = id
		mapJson["id_Client"] = idClient
		mapJson["id_House"] = idHouse
		mapJson["Incident"] = incident.String
		mapJson["Report"] = report.String
		mapJson["Path"] = path.String
		mapJson["St"] = st
		mapJson["Date_Reg"] = dateReg

		results = append(results, mapJson)
	}

	if len(results) == 0 {
		return nil, errors.New("No encontrado")
	} else {
		bOK = true
	}

	if bOK == true {
		return results, nil
	} else {
		mapJson := make(map[string]interface{})
		mapJson["Process"] = "No hay datos"
		mapJson["Datos"] = "No en contrado"
		err := errors.New("No encontrado")
		return nil, err
	}
}

func (ur *UserRepository) UserExist(ctx context.Context, Mail string, Tel string) (bool, error) {
	var results []map[string]interface{}

	//ctx, errCtx := context.WithTimeout(context.Background(), time.Duration(ur.timeout)*time.Second)
	ctx, _ = context.WithTimeout(ctx, ur.timeout)
	log.Println(ctx)

	stmt, ok := ur.mysqlConnector.GetStatements("ExistUsr")
	if !ok || stmt == nil {
		return false, fmt.Errorf("statement 'Select' not found or is nil")
	}

	rows, err := stmt.QueryContext(ctx, Mail, Tel)
	if err != nil {
		logger.GetLogger().Error("SelectRepository", "Get", err)
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		mapJson := make(map[string]interface{})
		var id uint64
		var st string

		if err := rows.Scan(&id, &st); err != nil {
			logger.GetLogger().Error("ExistUser", "Scan", err)
			return false, err
		}

		mapJson["id_User"] = id
		mapJson["St"] = st

		results = append(results, mapJson)
	}

	if len(results) == 0 {
		return false, errors.New("No encontrado")
	} else {
		return true, nil
	}
}
