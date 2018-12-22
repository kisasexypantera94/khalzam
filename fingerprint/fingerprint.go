package fingerprint

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bobertlo/go-mpg123/mpg123"
	"log"
	"math"
	"os/exec"
	"strings"
)

// Decode returns float32 slice of samples
func Decode(filename string) []float32 {
	decoder, err := mpg123.NewDecoder("")
	checkErr(err)

	err = decoder.Open(filename)
	checkErr(err)
	defer decoder.Close()

	rate, channels, _ := decoder.GetFormat()
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	var rawData []float32
	tmp := make([]int16, 8192)

	for {
		buf := make([]byte, 2*len(tmp))
		_, err := decoder.Read(buf)

		if err != nil {
			break
		}

		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
		for i := 1; i < len(tmp); i += 2 {
			rawData = append(rawData, (float32)(tmp[i-1]+tmp[i])/2/math.MaxInt16)
		}
	}

	decoder.Delete()

	return rawData
}

// StereoToMono converts file to mono using ffmpeg
func StereoToMono(filename string) string {
	dotIdx := strings.LastIndex(filename, ".")
	monoFilename := filename[:dotIdx] + "_mono"
	if dotIdx != -1 {
		monoFilename += filename[dotIdx:]
	}
	fmt.Println(monoFilename)
	err := exec.Command("/usr/local/bin/ffmpeg", "-i", filename, "-ac", "1", monoFilename).Run()
	checkErr(err)
	return monoFilename
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
