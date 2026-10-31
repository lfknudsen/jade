package main

import (
	"database/sql"
	. "fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"uuid"

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

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)

	switch os.Args[1] {
	case "init":
		initDatabase(db)
	case "load":
		var dir = "."
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}

		abs, err := filepath.Abs(dir)
		if err != nil {
			log.Fatal(err)
		}
		folders, files, err := loadDir(abs, nil)
		if err != nil {
			log.Fatal(err)
		}
		insert(db, folders, files)
	}
}

func connect() (*sql.DB, error) {
	args := []string{
		"user=postgres", is("password", os.Getenv("JADE_DB_PASSWORD")),
		"host=localhost", is("port", os.Getenv("JADE_DB_PORT")),
		"dbname=postgres",
		"connect_timeout=5",
		"sslmode=disable",
	}
	return sql.Open("postgres", strings.Join(args, " "))
}

func startServer() {
	err := exec.Command("docker", "run",
		"-d",
		// remove container entirely when closed. For testing.
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

func loadDir(dir string, parent *string) ([]Folder, []File, error) {
	// Compile list of directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []Folder{}, []File{}, err
	}

	// Create output slices.
	folders := make([]Folder, 0, len(entries)+1)
	files := make([]File, 0, len(entries))

	// Create this folder's struct
	this := Folder{
		uuid:   uuid.New().String(),
		name:   filepath.Base(dir),
		parent: parent,
	}
	folders = append(folders, this)

	// Walk
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return []Folder{}, []File{}, err
		}
		if info.IsDir() {
			subfolders, subfiles, err := loadDir(filepath.Join(dir, entry.Name()), &this.uuid)
			if err != nil {
				return folders, files, err
			}
			folders = append(folders, subfolders...)
			files = append(files, subfiles...)
		} else if info.Mode().IsRegular() {
			files = append(files, File{
				uuid:   uuid.New().String(),
				name:   entry.Name(),
				parent: this.uuid,
			})
		}
	}
	return folders, files, nil
}

type Folder struct {
	uuid   string
	name   string
	parent *string
}

type File struct {
	uuid   string
	name   string
	parent string
}
