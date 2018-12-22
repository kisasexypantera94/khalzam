package musiclibrary

import (
	"fmt"
	_ "github.com/lib/pq"
	"testing"
)

func TestLibInit(t *testing.T) {
	musicLib, err := Open()
	defer musicLib.Close()
	checkErr(err)
	err = musicLib.InsertSong("Toto - Africa")
	if err != nil {
		fmt.Println("insert song: ", err)
		return
	}
	err = musicLib.DeleteSong("Toto - Africa")
	checkErr(err)
}
