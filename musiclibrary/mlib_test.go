package musiclibrary

import (
	"fmt"
	_ "github.com/lib/pq"
	"testing"
)

func TestLibInit(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	defer musicLib.Close()
	err = musicLib.InsertSong("Toto - Africa")
	if err != nil {
		fmt.Println("insert song: ", err)
		return
	}
	err = musicLib.DeleteSong("Toto - Africa")
	checkErr(err)
}

func TestInsertSong(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	defer musicLib.Close()
	err = musicLib.InsertSong("../resources/Metallica - Master Of Puppets.mp3")
}
