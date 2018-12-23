package musiclibrary

import (
	// "fmt"
	_ "github.com/lib/pq"
	"testing"
)

func TestRecogniseSong(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	defer musicLib.Close()
	err = musicLib.DeleteSong("Modjo - Lady (Hear Me Tonight)")
	err = musicLib.DeleteSong("Beastie Boys - Intergalactic")
	err = musicLib.DeleteSong("Mogwai - Travel Is Dangerous")

	err = musicLib.InsertSong("../resources/Modjo - Lady (Hear Me Tonight).mp3")
	err = musicLib.InsertSong("../resources/Beastie Boys - Intergalactic.mp3")
	err = musicLib.InsertSong("../resources/Mogwai - Travel Is Dangerous.mp3")

	musicLib.RecogniseSong("../resources/intergalactic_sample.mp3")
	musicLib.RecogniseSong("../resources/modjo_lady_sample.mp3")
	musicLib.RecogniseSong("../resources/travel_sample.mp3")
	musicLib.RecogniseSong("../resources/travel_chorus_sample.mp3")
}