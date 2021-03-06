package musiclibrary

import (
	"database/sql"
	"fmt"
	"github.com/kisasexypantera94/khalzam/fingerprint"
	"github.com/remeh/sizedwaitgroup"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Config holds the configuration used for database
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBname   string
}

// MusicLibrary is the central structure of the algorithm.
// It is the link for fingerprinting and repository interaction.
type MusicLibrary struct {
	db *sql.DB
}

// Open connects to existing audio repository
func Open(cfg *Config) (*MusicLibrary, error) {
	fmt.Printf("Initializing library...\n\n")
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return &MusicLibrary{db}, err
}

// Close closes library
func (lib *MusicLibrary) Close() error {
	err := lib.db.Close()
	return err
}

// Index inserts song into library
func (lib *MusicLibrary) Index(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("Index: file not found")
	}

	fmt.Printf("Indexing '%s'...\n", filename)
	dotIdx := strings.LastIndex(filename, ".")
	slashIdx := strings.LastIndex(filename, "/")
	if dotIdx == -1 {
		return fmt.Errorf("Index: invalid file '%s'", filename)
	}
	songName := filename[slashIdx+1 : dotIdx]
	var songID int
	err := lib.db.QueryRow("INSERT INTO songs(song) VALUES($1) returning sid;", songName).Scan(&songID)
	if err != nil {
		return err
	}

	hashArray, err := fingerprint.Fingerprint(filename)
	if err != nil {
		return err
	}

	stmt, _ := lib.db.Prepare("INSERT INTO hashes(hash, time, sid) VALUES($1, $2, $3) returning hid;")
	for time, hash := range hashArray {
		stmt.Exec(hash, time, songID)
	}

	return nil
}

// IndexDir indexes whole directory
func (lib *MusicLibrary) IndexDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("IndexDir: invalid directory '%s'", path)
	}

	wg := sizedwaitgroup.New(8)
	for _, f := range files {
		filename := path + "/" + f.Name()
		if filepath.Ext(f.Name()) == ".mp3" {
			wg.Add()
			go func() {
				defer wg.Done()
				lib.Index(filename)
			}()
		}
	}
	wg.Wait()

	return nil
}

type table struct {
	absoluteBest  uint          // highest number of matches among timedeltaBest
	timedeltaBest map[uint]uint // highest number of matches for every timedelta
}

type candidate struct {
	songID   uint
	matchNum uint
}

// Recognize searches library and returns table
func (lib *MusicLibrary) Recognize(filename string) (result string, err error) {
	fmt.Printf("Recognizing '%s'...\n", filename)

	hashArray, err := fingerprint.ParallelFingerprint(filename)
	if err != nil {
		return "", err
	}

	cnt := make(map[uint]*table)
	for t, h := range hashArray {
		rows, err := lib.db.Query("SELECT * FROM hashes WHERE hash=$1;", h)
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
				cnt[sid].timedeltaBest = make(map[uint]uint)
			}

			cnt[sid].timedeltaBest[time-uint(t)]++
			if cnt[sid].timedeltaBest[time-uint(t)] > cnt[sid].absoluteBest {
				cnt[sid].absoluteBest = cnt[sid].timedeltaBest[time-uint(t)]
			}
		}
	}

	matchings := make([]*candidate, 0)
	for song, table := range cnt {
		matchings = append(matchings, &candidate{song, table.absoluteBest})
	}

	sort.Slice(matchings, func(i, j int) bool {
		return matchings[i].matchNum > matchings[j].matchNum
	})

	var songName string
	lib.db.QueryRow("SELECT song FROM songs WHERE sid=$1;", matchings[0].songID).Scan(&songName)

	result = fmt.Sprintf("Best match: %s (%d%% matched)\n", songName, int(100*float64(matchings[0].matchNum)/float64(len(hashArray))))
	return
}

// RecognizeDir recognizes whole directory
func (lib *MusicLibrary) RecognizeDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("IndexDir: invalid directory '%s'", path)
	}

	for _, f := range files {
		filename := path + "/" + f.Name()
		if filepath.Ext(f.Name()) == ".mp3" {
			res, err := lib.Recognize(filename)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(res)
		}
	}

	return nil
}

// Delete deletes song from library
func (lib *MusicLibrary) Delete(song string) (affected int64, err error) {
	fmt.Printf("Deleting '%s'...\n", song)
	statement, err := lib.db.Prepare("DELETE FROM songs WHERE song=$1;")
	if err != nil {
		return 0, err
	}
	res, err := statement.Exec(song)
	if err != nil {
		return 0, err
	}
	affected, err = res.RowsAffected()
	return affected, err
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
