package main

import (
	"database/sql"
	. "fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

	args := []string{
		"user=postgres", is("password", os.Getenv("JADE_DB_PASSWORD")),
		"host=localhost", is("port", os.Getenv("JADE_DB_PORT")),
		"dbname=postgres",
		"connect_timeout=5",
		"sslmode=disable",
	}

	db, err := sql.Open("postgres", strings.Join(args, " "))
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)
	initDatabase(db)
}

func startServer() {
	err := exec.Command("docker", "run",
		"-d",
		// remove container entirely when closed. Only for testing.
		"--rm",
		// Bind host machine port to container port.
		"-p", Sprintf("%s:5432", os.Getenv("JANE_DB_PORT")),
		// Set the environment variable that determines the user's password.
		"-e", is("POSTGRES_PASSWORD", os.Getenv("JADE_DB_PASSWORD")),
	).Run()
	if err != nil {
		log.Fatal(err)
	}
}

// is is a convenience function for assembling a key-value pair.
func is(key string, val string) string {
	return Sprintf("%s=%s", key, val)
}
