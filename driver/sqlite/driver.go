package mysql

import (
	"database/sql"
	"log"

	"github.com/elwinar/rambler/driver"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	driver.Register("sqlite", Driver{})
}

type Driver struct{}

func (d Driver) New(dsn, schema string) (driver.Conn, error) {
	log.Println("opening", dsn)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return Conn{
		db:     db,
		schema: schema,
	}, nil
}

type Conn struct {
	db     *sql.DB
	schema string
}

func (c Conn) HasTable() (bool, error) {
	var table string
	err := c.db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, "migrations").Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (c Conn) CreateTable() error {
	_, err := c.db.Exec(`CREATE TABLE migrations ( migration VARCHAR(255) NOT NULL );`)
	return err
}

func (c Conn) GetApplied() ([]string, error) {
	rows, err := c.db.Query(`SELECT migration FROM migrations ORDER BY migration ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (c Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (migration) VALUES (?)`, migration)
	return err
}

func (c Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE migration = ?`, migration)
	return err
}

func (c Conn) Execute(statement string) error {
	_, err := c.db.Exec(statement)
	return err
}
