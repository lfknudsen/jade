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

	db, err := sql.Open("postgres", "host=localhost dbname=postgres connect_timeout=5 port=5432")
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
