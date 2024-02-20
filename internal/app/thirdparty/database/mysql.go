package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // needed
)

const (
	defaultRetry           = 3
	ErrMysqlRecordNotFound = "Record not found"
	ErrMysqlFailed         = "SQL failed"
	ErrRecordAlreadyExists = "Record already exists."
)

// A MysqlService represents the MySql client
type MysqlService struct {
	db           *sql.DB
	DatabaseName string
}

// type MySQLError struct {
// 	OrigErr error                  `json:"trace_stack,omitempty"`
// 	Details map[string]interface{} `json:"error,omitempty"`
// }


// NewMysqlService a singleton function that returns only one instance of the Mysql database client
func NewMysqlService(ctx context.Context, conf *Config) (*MysqlService, error) {
	var dbh *sql.DB
	var err error

	svc := &MysqlService{}

	if conf.Retry <= 0 {
		conf.Retry = defaultRetry
	}

	for retry := 1; retry <= conf.Retry; retry++ {

		//connect
		dbh, err = sql.Open("mysql", conf.DSN)
		if err != nil {
			time.Sleep(time.Millisecond * 1000)
			continue
		}

		err = dbh.Ping()
		if err != nil {
			time.Sleep(time.Millisecond * 1000)
			continue
		}
		// most important tweak is here :-)
		// https://www.alexedwards.net/blog/configuring-sqldb

		// MySQL's wait_timeout setting will automatically close any connections
		// that haven't been used for 8 hours (by default).

		// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
		// means that there is no maximum lifetime and the connection is reused
		// forever (which is the default behavior).

		dbh.SetConnMaxLifetime(12 * time.Hour)

		// Set the maximum number of concurrently idle connections to 5. Setting this
		// to less than or equal to 0 will mean that no idle connections are retained.

		dbh.SetMaxIdleConns(5)

		// Set the maximum number of concurrently open connections to 5. Setting this
		// to less than or equal to 0 will mean there is no maximum limit (which
		// is also the default setting).

		dbh.SetMaxOpenConns(5)

		// Set the number of open and idle connection
		// to a maximum total of (idle:2 + open:3) = 5

		break
	}

	if dbh != nil {
		rows, err := dbh.QueryContext(ctx, "select version()")
		if err != nil {
			return nil, svc.ConvertError(ctx, err)
		}
		defer rows.Close()
		for rows.Next() {
			var ver string
			if err := rows.Scan(&ver); err != nil {
				return nil, svc.ConvertError(ctx, err)
			}
		}
		// Check for errors from iterating over rows.
		if err := rows.Err(); err != nil {
			return nil, svc.ConvertError(ctx, err)
		}

	}

	svc.db = dbh
	svc.DatabaseName = conf.DatabaseName

	//oks
	return svc, nil
}

func (s *MysqlService) Client() *sql.DB {
	return s.db
}

// Get retrieves one record from the db table specified given its primary key and value
func (s *MysqlService) Get(ctx context.Context, tbl, id string, value interface{}, model interface{}) (interface{}, error) {
	rows, err := s.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s=?", tbl, id), model, value)
	if err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		v := reflect.ValueOf(model).Elem()
		v.Set(reflect.ValueOf(rows[0]).Elem())
		return model, nil
	}

	return nil,sql.ErrNoRows
}

// Insert creates a new database record of a table given a structure
func (s *MysqlService) Insert(ctx context.Context, tbl string, obj interface{}) (interface{}, error) {
	var tx *sql.Tx
	var pErr error

	isNewTx := false
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		tx, pErr = s.BeginTx(ctx)
		if pErr != nil {
			return nil, pErr
		}
		isNewTx = true
	}

	names, values := s.marshalAttributes(obj)

	sqlStr := fmt.Sprintf("INSERT %s SET %s", tbl, strings.Join(names, ", "))

	res, err := tx.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		if isNewTx {
			tx.Rollback()
		}
		return nil, s.ConvertError(ctx, err, sqlStr)
	}
	if isNewTx {
		tx.Commit()
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, s.ConvertError(ctx, err, sqlStr)
	}
	return id, nil
}

// Update updates a record given its primary key and value, the of the object passed will set/update the data of that record
func (s *MysqlService) Update(ctx context.Context, tbl, id string, value interface{}, obj interface{}) error {
	var tx *sql.Tx
	var pErr error

	isNewTx := false
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		tx, pErr = s.BeginTx(ctx)
		if pErr != nil {
			return pErr
		}
		isNewTx = true
	}

	names, values := s.marshalAttributes(obj)
	values = append(values, value)
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", tbl, strings.Join(names, ", "), id)
	_, err := tx.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		if isNewTx {
			tx.Rollback()
		}
		return s.ConvertError(ctx, err, sqlStr)
	}
	if isNewTx {
		tx.Commit()
	}
	return nil
}

// Delete deletes a database record of a particular table given its promary key value
func (s *MysqlService) Delete(ctx context.Context, tbl, id string, value interface{}) error {
	var tx *sql.Tx
	var pErr error

	isNewTx := false
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		tx, pErr = s.BeginTx(ctx)
		if pErr != nil {
			return pErr
		}
		isNewTx = true
	}

	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s=?", tbl, id)
	res, err := tx.ExecContext(ctx, sqlStr, value)
	if err != nil {
		if isNewTx {
			tx.Rollback()
		}
		return s.ConvertError(ctx, err, sqlStr)
	}
	if isNewTx {
		tx.Commit()
	}

	_ = res

	return s.ConvertError(ctx, err, sqlStr)
}

// Query accepts any SQL query statement and will return the results in an array with teh same type as the model interface passed
func (s *MysqlService) Query(ctx context.Context, queryStr string, model interface{}, values ...interface{}) ([]interface{}, error) {
	var objects []interface{}

	var rows *sql.Rows
	var err error
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		rows, err = tx.QueryContext(ctx, queryStr, values...)
		if err != nil {
			return objects, s.ConvertError(ctx, err, queryStr)
		}
		defer rows.Close()
	} else {
		rows, err = s.db.QueryContext(ctx, queryStr, values...)
		if err != nil {
			return objects, s.ConvertError(ctx, err, queryStr)
		}
		defer rows.Close()
	}

	// Get reflection value of model passed
	om := reflect.ValueOf(model).Elem()

	for rows.Next() {

		// figure out what columns were returned
		// the column names will be the JSON object field keys
		columns, err := rows.Columns()
		if err != nil {
			return objects, s.ConvertError(ctx, err, queryStr)
		}

		// Create new interface of the same type as the model type, then assign that created instance
		cm := reflect.New(om.Type()).Interface()

		// Returns the value of the returned interface
		v := reflect.ValueOf(cm).Elem()

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		valuesFields := make([]interface{}, len(columns))

		// Get field names from struct and their JSON equivalent
		t := v.Type()
		fieldIdx := map[string]int{}
		for i := 0; i < t.NumField(); i++ {
			jsonName := t.Field(i).Tag.Get("json")
			fieldIdx[jsonName] = i
		}

		for i, name := range columns {
			if _, found := fieldIdx[name]; !found {
				// To avoid the error "...converting NULL to [type] is unsupported"
				// We will set the column that no longer exists in the model to an empty interface
				k := new(interface{})
				valuesFields[i] = &k
			} else {
				// Assign the struct field for the scan value
				valuesFields[i] = v.Field(fieldIdx[name]).Addr().Interface()
			}
		}

		err = rows.Scan(valuesFields...)
		if err != nil {
			return objects, s.ConvertError(ctx, err, queryStr)
		}

		objects = append(objects, cm)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return objects, s.ConvertError(ctx, err, queryStr)
	}

	return objects, nil
}

// marshalAttributes returns the keys and values given a staruct value
func (s *MysqlService) marshalAttributes(obj interface{}) ([]string, []interface{}) {
	var m map[string]interface{}
	b, err := json.Marshal(obj)
	if err == nil {
		_ = json.Unmarshal(b, &m)
	}

	values := []interface{}{}
	names := []string{}
	for k, v := range m {
		names = append(names, fmt.Sprintf("%s=?", k))
		values = append(values, v)
	}

	return names, values
}

func (s *MysqlService) CreateTable(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return s.ConvertError(ctx, err, query)
	}
	return nil
}

func (s *MysqlService) DropTable(ctx context.Context, tableName string) error {
	sqlStr := fmt.Sprintf("DROP TABLE %s", tableName)
	_, err := s.db.ExecContext(ctx, sqlStr)
	if err != nil {
		return s.ConvertError(ctx, err, sqlStr)
	}
	return nil
}

// ConvertError is a helper method that converts mysql error messages to the more generic models.DBError
func (s *MysqlService) ConvertError(ctx context.Context, err error, args ...string) error {
	if err == nil {
		return nil
	}


	if sqlErr, ok := err.(*mysql.MySQLError); ok {

		eMsg := ErrMysqlFailed
		if sqlErr.Number == 1062 {
			eMsg = ErrRecordAlreadyExists
		}

		return &mysql.MySQLError{Number:sqlErr.Number,SQLState: sqlErr.SQLState, Message: fmt.Sprintf("%s; SQL Error: %d", eMsg, sqlErr.Number)}
	}

	log.Println(fmt.Sprintf("Database:: MYSQL::%v - %s", err, strings.Join(args, "; ")))

	return err
}

func (s *MysqlService) BeginTx(ctx context.Context) (*sql.Tx, error) {
	var tx *sql.Tx

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, s.ConvertError(ctx, err)
	}

	return tx, nil
}

// GetAndUpdate one record from the db table specified given its primary key and value and call updateFn within the same DB transaction
func (s *MysqlService) GetAndUpdate(ctx context.Context, tbl, id string, value interface{}, model interface{}, updateFn UpdateFunc) (interface{}, error) {
	tx, pErr := s.BeginTx(ctx)
	if pErr != nil {
		return "", pErr
	}
	defer tx.Rollback()

	tctx := context.WithValue(ctx, "tx", tx)
	qry := fmt.Sprintf("SELECT * FROM %s WHERE %s=? FOR UPDATE", tbl, id)
	rows, pErr := s.Query(tctx, qry, model, value)
	if pErr != nil {
		return nil, pErr
	}
	if len(rows) == 0 {
		return nil,sql.ErrNoRows
	}

	v := reflect.ValueOf(model).Elem()
	v.Set(reflect.ValueOf(rows[0]).Elem())

	updatedDoc, pErr := updateFn(tctx, model)
	if pErr != nil {
		return nil, pErr
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return updatedDoc, nil
}
