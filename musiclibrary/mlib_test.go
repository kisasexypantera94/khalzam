package musiclibrary

import (
	"fmt"
	_ "github.com/lib/pq"
	// "sync"
	"testing"
)

// func TestIndexing(t *testing.T) {
// 	musicLib, err := Open()
// 	checkErr(err)          //
// 	wg := sync.WaitGroup{} //
// 	defer musicLib.Close() //
// 	wg.Add(1)              //
// 	go func() {            //
// 		defer wg.Done() //
// 		musicLib.Delete("Modjo - Lady (Hear Me Tonight)")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Delete("Beastie Boys - Intergalactic")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Delete("Mogwai - Travel Is Dangerous")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Delete("journeydontstop")
// 	}()
// 	wg.Wait()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Index("../resources/Modjo - Lady (Hear Me Tonight).mp3")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Index("../resources/Beastie Boys - Intergalactic.mp3")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Index("../resources/Mogwai - Travel Is Dangerous.mp3")
// 	}()
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		musicLib.Index("../resources/journeydontstop.mp3")
// 	}()
// 	wg.Wait()
// }

// func TestIndexDir(t *testing.T) {
// 	musicLib, err := Open()
// 	checkErr(err)
// 	defer musicLib.Close()
// 	// musicLib.IndexDirectory("../resources")
// 	musicLib.Index("../resources/Арсений Креститель - Мой Вейп.mp3")
// 	musicLib.Index("../resources/пасош - мандельштам.mp3")
// 	musicLib.Index("../resources/Хаски - Пуля-дура.mp3")
// }

func TestRecogniseOnly(t *testing.T) {
	musicLib, err := Open()
	checkErr(err)
	defer musicLib.Close()
	fmt.Println(musicLib.Recognize("../resources/intergalactic_sample.mp3"))
	fmt.Println(musicLib.Recognize("../resources/travel_chorus_sample.mp3"))
	fmt.Println(musicLib.Recognize("../resources/travel_sample.mp3"))
	fmt.Println(musicLib.Recognize("../resources/modjo_lady_sample.mp3"))
	fmt.Println(musicLib.Recognize("../resources/intergalacticnew.mp3"))
	fmt.Println(musicLib.Recognize("../resources/bloodorangegoodenough.mp3"))
	fmt.Println(musicLib.Recognize("../resources/xtal_sample.mp3"))
	fmt.Println(musicLib.Recognize("../resources/disorderlive.mp3"))
	fmt.Println(musicLib.Recognize("../resources/journeylive.mp3"))
}

// func TestWavDecode(t *testing.T) {
// 	mlib, err := Open()
// 	checkErr(err)
// 	defer mlib.Close()
// 	// mlib.Index("../resources/fr.wav")
// 	mlib.Delete("blood")
// 	mlib.Index("../resources/blood.wav")
// 	mlib.Recognize("../resources/bloodourangegoodenough.wav")
// }
