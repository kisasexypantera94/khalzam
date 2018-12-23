package musiclibrary

import (
	"database/sql"
	"fmt"
	"github.com/kisasexypantera94/khalzam/fingerprint"
	"log"
	"strings"
)

const (
	dbUser     = "kisasexypantera94"
	dbPassword = ""
	dbName     = "khalzam"
)

// MusicLibrary ...
type MusicLibrary struct {
	db *sql.DB
}

// Open return pointer to existing library
func Open() (*MusicLibrary, error) {
	fmt.Println("Initializing library...")
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open("postgres", dbinfo)
	return &MusicLibrary{db}, err
}

// Close closes library
func (lib *MusicLibrary) Close() {
	lib.db.Close()
}

// InsertSong inserts song into library
func (lib *MusicLibrary) InsertSong(filename string) error {
	dotIdx := strings.LastIndex(filename, ".")
	slashIdx := strings.LastIndex(filename, "/")
	songName := filename[slashIdx + 1:dotIdx]
	var songID int
	err := lib.db.QueryRow("INSERT INTO songs(song) VALUES($1) returning sid;", songName).Scan(&songID)
	if err != nil {
		return err
	}

	hashArray := fingerprint.Fingerprint(filename)
	for time, hash := range hashArray {
		lib.db.QueryRow("INSERT INTO hashes(hash, time, sid) VALUES($1, $2, $3)", hash, time, songID)
	}

	return err
}

// DeleteSong deletes song from library
func (lib *MusicLibrary) DeleteSong(song string) error {
	statement, err := lib.db.Prepare("DELETE FROM songs WHERE song=$1")
	checkErr(err)
	_, err = statement.Exec(song)
	return err
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
