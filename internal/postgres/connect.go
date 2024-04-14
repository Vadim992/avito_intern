package postgres

import "database/sql"

type DB struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{DB: db}
}

func InitDB(connect string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connect)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
