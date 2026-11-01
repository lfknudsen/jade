package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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

// insert performs bulk insertion of folders and files.
// It does so by joining all the different INSERT SQL statements into one per
// table. In one tiny example set of 311 directory entries, this method
// reduced the duration from ~14s to ~245ms.
func insert(db *sql.DB, folders []Folder, files []File) (rows int64) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// ================
	// Folder insertion

	// Contain indices of folders
	nestedFolders := make([]int, 0, len(folders))
	rootFolders := make([]int, 0, len(folders))
	for i := range folders {
		if folders[i].parent == nil {
			rootFolders = append(rootFolders, i)
		} else {
			nestedFolders = append(nestedFolders, i)
		}
	}

	// Insert folders without parent
	columns := 2
	statement := PrepareMultiInsert("folder (uuid, name)", len(rootFolders), columns)
	flatArgs := make([]any, 0, len(rootFolders)*columns)
	for _, i := range rootFolders {
		flatArgs = append(flatArgs, folders[i].uuid, folders[i].name)
	}
	rows += Exec(tx, statement, flatArgs)

	// Insert folders with parent
	columns = 3
	statement = PrepareMultiInsert("folder", len(nestedFolders), columns)
	flatArgs = make([]any, 0, len(nestedFolders)*columns)
	for _, i := range nestedFolders {
		flatArgs = append(flatArgs, folders[i].uuid, folders[i].name, *(folders[i].parent))
	}
	rows += Exec(tx, statement, flatArgs)

	// Insert relations between parent and subdirectory
	columns = 2
	statement = PrepareMultiInsert("subdirectory", len(nestedFolders), columns)
	flatArgs = make([]any, 0, len(nestedFolders)*columns)
	for _, i := range nestedFolders {
		flatArgs = append(flatArgs, *(folders[i].parent), folders[i].uuid)
	}
	rows += Exec(tx, statement, flatArgs)

	// ==============
	// FILE INSERTION

	// Insert files
	columns = 3
	statement = PrepareMultiInsert("file", len(files), columns)
	flatArgs = make([]any, 0, len(files)*columns)
	for _, f := range files {
		flatArgs = append(flatArgs, f.uuid, f.name, f.parent)
	}
	rows += Exec(tx, statement, flatArgs)

	// Insert relation between parent folder and child file
	columns = 2
	statement = PrepareMultiInsert("child", len(files), columns)
	flatArgs = make([]any, 0, len(files)*columns)
	for _, f := range files {
		flatArgs = append(flatArgs, f.parent, f.uuid)
	}
	rows += Exec(tx, statement, flatArgs)

	commit(tx)
	return rows
}

func Exec(tx *sql.Tx, statementText string, args []any) (rows int64) {
	stmt, err := tx.Prepare(statementText)
	if err != nil {
		rollback(tx)
		log.Fatal(err)
	}
	res, err := stmt.Exec(args...)
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

func PrepareMultiInsert(table string, rows int, columns int) string {
	prefix := fmt.Sprintf("INSERT INTO %s VALUES\n", table)
	lines := make([]string, 0, rows)
	for i := range rows {
		row := make([]string, 0, columns)
		for j := 1; j <= columns; j++ {
			row = append(row, fmt.Sprintf("$%d", i*columns+j))
		}
		line := fmt.Sprintf("(%s)", strings.Join(row, ", "))
		lines = append(lines, line)
	}
	return prefix + strings.Join(lines, ",\n") + ";"
}
