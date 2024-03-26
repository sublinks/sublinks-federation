package db

import (
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	Connect() error
	Ping() bool
	RunMigrations()
	Find(interface{}, ...interface{}) error
	Save(interface{}) error
}

func (db *PostgresDB) Find(data interface{}, conds ...interface{}) error {
	db.tx = db.DB.Find(data, conds)
	return db.tx.Error
}

func (db *PostgresDB) Ping() bool {
	return db.DB.Exec("SELECT 1").Error == nil
}

type PostgresDB struct {
	*gorm.DB
	tx *gorm.DB
}

func NewDatabase() Database {
	return &PostgresDB{}
}

func (d *PostgresDB) Save(data interface{}) error {
	db := d.DB.Save(data)
	return db.Error
}

func (d *PostgresDB) Connect() error {
	database, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	maxIdleCons, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		maxIdleCons = 10
	}
	sqlDB.SetMaxIdleConns(maxIdleCons)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	maxOpenCons, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		maxOpenCons = 100
	}
	sqlDB.SetMaxOpenConns(maxOpenCons)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	maxConnLifetime, err := time.ParseDuration(os.Getenv("DB_MAX_CONN_LIFETIME"))
	if err != nil {
		maxConnLifetime = 60
	}
	sqlDB.SetConnMaxLifetime(maxConnLifetime * time.Minute)
	d.DB = database
	return nil
}
