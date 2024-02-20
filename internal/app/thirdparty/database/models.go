package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

// A Manger represents the interface that needs to be implemented for any database service defined
type Manager interface {
	Insert(ctx context.Context, tbl string, obj interface{}) (interface{}, error)
	Update(ctx context.Context, tbl, id string, value interface{}, obj interface{}) error
	Delete(ctx context.Context, tbl, id string, value interface{}) error
	Get(ctx context.Context, tbl, key string, value interface{}, model interface{}) (interface{}, error)
	GetAndUpdate(ctx context.Context, tbl, key string, value interface{}, model interface{}, updateFn UpdateFunc) (interface{}, error)
	Query(ctx context.Context, queryStr string, model interface{}, values ...interface{}) ([]interface{}, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTable(ctx context.Context, query string) error
	DropTable(ctx context.Context, tableName string) error
	Client() *sql.DB
}

// A Config represents the configuration needed to instantiate the database client
type Config struct {
	Name         string
	DSN          string
	Options      map[string]string
	Retry        int
	TLS          bool
	DatabaseName string
}

type UpdateFunc func(context.Context, interface{}) (interface{}, error)


// MysqlConf  ...
type MysqlConf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
}

// ConnStr get the mysql conn-string
func (conf *MysqlConf) ConnStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
	)
}

// Options ...
func (conf *MysqlConf) Options() string {
	dbOptions := map[string]string{
		"charset":          "utf8",
		"tls":              "false",
		"loc":              "Local",
		"maxAllowedPacket": "0",
		"readTimeout":      "1m30s",
		"writeTimeout":     "1m",
		"timeout":          "1m30s", // dial-timeout
	}

	var opts []string
	for k, v := range dbOptions {
		opts = append(opts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	s := fmt.Sprintf("%s", strings.Join(opts, "&"))
	return s
}

// DSN ...
func (conf *MysqlConf) DSN() string {
	return conf.ConnStr() + `?` + conf.Options()
}