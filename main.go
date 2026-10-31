package main

import (
	"database/sql"
	"log"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			startServer()
			os.Exit(0)
		}
	}

	db, err := sql.Open("postgres", "host=localhost dbname=jade connect_timeout=5 port=5432")
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)
	initDatabase(db)
}

func startServer() {
	err := exec.Command("docker", "run",
		"-d",
		"--rm",
		"-p", "5432:5432",
		"-e", "POSTGRES_PASSWORD=abc123",
	).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func initDatabase(db *sql.DB) {
	tx := newTx(db)
	defer rollback(tx)
	execStmt(dropTables, tx)
	execStmt(createTables, tx)
	commit(tx)
}

func execStmt(query string, tx *sql.Tx) {
	stmt, err := tx.Prepare(query)
	if err != nil {
		rollback(tx)
		log.Fatal(err)
	}
	if _, err = stmt.Exec(); err != nil {
		closeStmt(stmt)
		rollback(tx)
		log.Fatal(err)
	}
	closeStmt(stmt)
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
