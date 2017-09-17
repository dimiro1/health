package db

import (
	"database/sql"

	"github.com/dhiemaz/health"
)

// Checker is a checker that check a database connection
type Checker struct {
	CheckSQL   string
	VersionSQL string
	DB         *sql.DB
}

// NewChecker returns a new db.Checker with the given URL
func NewChecker(checkSQL, versionSQL string, db *sql.DB) Checker {
	return Checker{CheckSQL: checkSQL, VersionSQL: versionSQL, DB: db}
}

// NewMySQLChecker returns a new db.Checker configured for use in MySQL
func NewMySQLChecker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT VERSION()", db)
}

// NewPostgreSQLChecker returns a new db.Checker configured for use in PostgreSQL
func NewPostgreSQLChecker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT VERSION()", db)
}

// NewSqlite3Checker returns a new db.Checker configured for use in Sqlite3
func NewSqlite3Checker(db *sql.DB) Checker {
	return NewChecker("SELECT 1", "SELECT sqlite_version()", db)
}

// Check execute two queries in the database
// The first is a simple one used to verify if the database is up
// If is Up then another query is executed, querying for the database version
func (c Checker) Check() health.Health {
	var (
		version string
		ok      string
	)

	health := health.NewHealth()

	if c.DB == nil {
		health.Down().AddInfo("error", "Empty resource")
		return health
	}

	err := c.DB.QueryRow(c.CheckSQL).Scan(&ok)

	if err != nil {
		health.Down().AddInfo("error", err.Error())
		return health
	}

	err = c.DB.QueryRow(c.VersionSQL).Scan(&version)

	if err != nil {
		health.Down().AddInfo("error", err.Error())
		return health
	}

	health.Up().AddInfo("version", version)

	return health
}
