package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rsterbin/rambler/driver"
	"strings"
)

func init() {
	driver.Register("postgresql", Driver{})
}

type Driver struct{}

func (d Driver) New(dsn, dbname, schema string) (driver.Conn, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(schema) != "" {
		_, err := db.Exec(fmt.Sprintf(`SET search_path="%s";`, schema))
		if err != nil {
			return nil, err
		}
	} else {
		schema = "public"
	}

	c := &Conn{
		db:     db,
		dbname: dbname,
		schema: schema,
	}

	return c, nil
}

type Conn struct {
	db     *sql.DB
	dbname string
	schema string
}

func (c *Conn) HasTable() (bool, error) {
	var name string
	err := c.db.QueryRow(fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_catalog = '%s' AND table_schema = '%s' AND table_name = 'migrations'`, c.dbname, c.schema)).Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (c *Conn) CreateTable() error {
	_, err := c.db.Exec(fmt.Sprintf(`CREATE TABLE "%s".migrations ( migration VARCHAR(255) NOT NULL );`, c.schema))
	return err
}

func (c *Conn) GetApplied() ([]string, error) {
	rows, err := c.db.Query(fmt.Sprintf(`SELECT migration FROM "%s".migrations ORDER BY migration ASC`, c.schema))
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

func (c *Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`INSERT INTO "%s".migrations (migration) VALUES ($1)`, c.schema), migration)
	return err
}

func (c *Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`DELETE FROM "%s".migrations WHERE migration = $1`, c.schema), migration)
	return err
}

func (c *Conn) Execute(query string) error {
	_, err := c.db.Exec(query)
	return err
}
