package fingerprint

import (
	"testing"
)

func TestDecodeMpg123(t *testing.T) {
	Decode("../resources/Metallica - Master Of Puppets.mp3")
}

// func TestDecodeGo(t *testing.T) {
// 	f, err := os.Open("../resources/journeydontstop.mp3")
// 	checkErr(err)
// 	defer f.Close()

// 	d, err := mp3.NewDecoder(f)
// 	checkErr(err)
// 	defer d.Close()

// 	var rawData []float32
// 	tmp := make([]int16, 4096)

// 	// decode mp3 file and dump output
// 	for {
// 		buf := make([]byte, 2*len(tmp))
// 		_, err := d.Read(buf)

// 		if err != nil {
// 			break
// 		}

// 		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
// 		for i := 1; i < len(tmp); i += 2 {
// 			rawData = append(rawData, (float32)(tmp[i-1]+tmp[i])/2/math.MaxInt16)
// 		}
// 	}
// }

// func TestConvertStereoToMono(t *testing.T) {
// 	StereoToMono("/Users/chingachgook/dev/gocode/src/github.com/kisasexypantera94/khalzam/resources/journeydontstop.mp3")
// }
