package main

import (
	"database/sql"
	"log"
)

func initDatabase(db *sql.DB) {
	tx := newTx(db)
	defer rollback(tx)
	execStmt(dropTables, tx)
	execStmt(createTables, tx)
	commit(tx)
}

func execStmt(query string, tx *sql.Tx) {
	_, err := tx.Exec(query)
	if err != nil {
		rollback(tx)
		log.Fatal(err)
	}
}

func closeDB(db *sql.DB) {
	if db != nil {
		_ = db.Close()
	}
}

func rollback(tx *sql.Tx) {
	if tx != nil {
		log.Println("Rolling back transaction.")
		_ = tx.Rollback()
	}
}

func closeStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
}

func newTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func commit(tx *sql.Tx) {
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
