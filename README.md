# Khalzam
## About
Khalzam is a simple audio recognition program. Its algrorithm is based on
[this article by Jovan Jovanovic](https://www.toptal.com/algorithms/shazam-it-music-processing-fingerprinting-and-recognition)

## Setup
You need to create and initialize database:
```
➜  khalzam git:(master) ✗ createdb -O user databasename
➜  khalzam git:(master) ✗ psql -f createdb.sql databasename
```

## Usage
### Shell mode
```
➜  khalzam git:(master) ✗ DBUSER=kisasexypantera94 DBNAME=khalzam go run shell.go
Initializing library...

MusicLibrary interactive shell
>>> help

Commands:
  clear          clear the screen
  delete         delete audio from database
  exit           exit the program
  help           display help
  index          index audiofile
  indexdir       index directory
  recognize      recognize audiofile


>>> index resources/Modjo\ -\ Lady\ \(Hear\ Me\ Tonight\).mp3
Indexing 'resources/Modjo - Lady (Hear Me Tonight).mp3'
>>> recognize resources/modjo_lady_sample.mp3
Recognizing 'resources/modjo_lady_sample.mp3'...
Best match: Modjo - Lady (Hear Me Tonight) (11% matched)
```

### API
```golang
package main

import (
	"fmt"
	"github.com/kisasexypantera94/khalzam/musiclibrary"
	_ "github.com/lib/pq"
)

func main() {
	cfg := &musiclibrary.Config{
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASSWORD"),
		DBname:   os.Getenv("DBNAME"),
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
	}

	musicLib, _ := musiclibrary.Open(cfg)
	defer musicLib.Close()

	musicLib.Index("resources/Modjo - Lady (Hear Me Tonight).mp3")
	result := musicLib.Recognize("resources/modjo_lady_sample.mp3")
	fmt.Println(result)
}
```