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
	songName := filename[slashIdx+1 : dotIdx]
	var songID int
	err := lib.db.QueryRow("INSERT INTO songs(song) VALUES($1) returning sid;", songName).Scan(&songID)
	if err != nil {
		return err
	}

	hashArray := fingerprint.Fingerprint(filename)
	for time, hash := range hashArray {
		row := lib.db.QueryRow("INSERT INTO hashes(hash, time, sid) VALUES($1, $2, $3) returning hid;", hash, time, songID)
		var lastID int
		row.Scan(&lastID)
	}

	return err
}

type table struct {
	songID  uint
	curBest map[uint]uint
}

type part struct {
	songID   uint
	matchNum uint
}

// RecogniseSong searches library and returns table
func (lib *MusicLibrary) RecogniseSong(filename string) {
	cnt := make(map[uint]*table)

	hashArray := fingerprint.Fingerprint(filename)
	for t, h := range hashArray {
		rows, err := lib.db.Query("SELECT * FROM hashes WHERE hash=$1", h)
		checkErr(err)

		for rows.Next() {
			var hid int
			var hash string
			var sid uint
			var time uint
			err = rows.Scan(&hid, &hash, &time, &sid)
			checkErr(err)
			if cnt[sid] == nil {
				cnt[sid] = &table{}
				cnt[sid].curBest = make(map[uint]uint)
			}

			cnt[sid].curBest[time-(uint)(t)]++
			if cnt[sid].curBest[time-(uint)(t)] > cnt[sid].songID {
				cnt[sid].songID = cnt[sid].curBest[time-(uint)(t)]
			}
		}
	}

	matchings := make([]*part, 0)
	fmt.Println("Recognising", filename, "...")
	fmt.Println("Number of samples:", len(hashArray))
	for key, val := range cnt {
		matchings = append(matchings, &part{val.songID, key})
		fmt.Println(matchings[len(matchings)-1])
	}
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
