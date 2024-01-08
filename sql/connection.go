package sql

import "database/sql"

type Connection struct {
	db *sql.DB
}

func NewConnection(db *sql.DB) Connection {
	return Connection{
		db: db,
	}
}

func (c *Connection) AddTable(query string) error {
	_, err := c.db.Exec(query)
	return err
}

func (c *Connection) Execute(query string) {
	_, err := c.db.Exec(query)
	if err != nil {
		return
	}

}
