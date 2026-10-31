package main

import (
	"database/sql"
	"log"
)

func initDatabase(db *sql.DB) {
	tx := newTx(db)
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

func insert(db *sql.DB, folders []Folder, files []File) (rows int64) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range folders {
		statementText := insertRootFolder
		args := []string{f.uuid, f.name}
		if f.parent != nil {
			statementText = insertFolder
			args = append(args, *f.parent)
		}
		rows += Exec(tx, statementText, args)
		if f.parent != nil {
			rows += Exec(tx, insertFolderRelation, []string{*f.parent, f.uuid})
		}
	}

	for _, f := range files {
		args := []string{f.uuid, f.name, f.parent}
		rows += Exec(tx, insertFile, args)
		rows += Exec(tx, insertChildRelation, []string{f.parent, f.uuid})
	}
	commit(tx)
	return rows
}

func Exec(tx *sql.Tx, statementText string, args []string) (rows int64) {
	stmt, err := tx.Prepare(statementText)
	if err != nil {
		rollback(tx)
		log.Fatal(err)
	}
	res, err := stmt.Exec(args)
	if err != nil {
		closeStmt(stmt)
		rollback(tx)
		log.Fatal(err)
	}
	rows, err = res.RowsAffected()
	if err != nil {
		closeStmt(stmt)
		rollback(tx)
		log.Fatal(err)
	}
	err = stmt.Close()
	if err != nil {
		rollback(tx)
		log.Fatal(err)
	}
	return rows
}
