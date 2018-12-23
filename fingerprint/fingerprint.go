package fingerprint

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/mjibson/go-dsp/fft"
	"io"
	"log"
	"math"
	"math/cmplx"
	"os/exec"
	"strconv"
	"strings"
)

const chunkSize = 8192
const windowSize = 4096
const fuzzFactor = 2

var freqBins = [...]int16{40, 80, 120, 180, 300}

// Decode returns float32 slice of samples
func Decode(filename string) []float64 {
	decoder, err := mpg123.NewDecoder("")
	checkErr(err)

	err = decoder.Open(filename)
	checkErr(err)
	defer decoder.Close()

	rate, channels, _ := decoder.GetFormat()
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	var rawData []float64
	tmp := make([]int16, chunkSize)
	for {
		buf := make([]byte, 2*len(tmp))
		_, err := decoder.Read(buf)

		if err != nil {
			break
		}

		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
		if channels == 2 {
			for i := 2; i < len(tmp); i += 4 {
				rawData = append(rawData, (float64)(tmp[i-2]+tmp[i])/2/math.MaxInt32)
			}
		} else {
			for i := 0; i < len(tmp); i++ {
				rawData = append(rawData, (float64)(tmp[i])/math.MaxInt32)
			}
		}
	}

	decoder.Delete()
	return rawData
}

// Fingerprint returns a fingerprint of song
func Fingerprint(filename string) (hashArray []string) {
	rawData := Decode(filename)
	blockNum := len(rawData) / windowSize

	for i := 0; i < blockNum; i++ {
		complexArray := fft.FFTReal(rawData[i*windowSize : i*windowSize+windowSize])
		hashArray = append(hashArray, getKeyPoints(complexArray))
	}

	return hashArray
}

func getKeyPoints(compArr []complex128) string {
	highScores := make([]float64, len(freqBins))
	recordPoints := make([]uint, len(freqBins))

	for bin := freqBins[0]; bin < freqBins[len(freqBins)-1]; bin++ {
		magnitude := cmplx.Abs(compArr[bin])

		binIdx := 0
		for freqBins[binIdx] < bin {
			binIdx++
		}

		if magnitude > highScores[binIdx] {
			highScores[binIdx] = magnitude
			recordPoints[binIdx] = (uint)(bin)
		}
	}

	tmp := (recordPoints[3]-(recordPoints[3]%fuzzFactor))*1e8 +
		(recordPoints[2]-(recordPoints[2]%fuzzFactor))*1e5 +
		(recordPoints[1]-(recordPoints[1]%fuzzFactor))*1e2 +
		(recordPoints[0] - recordPoints[0]%fuzzFactor)

	// return hash(recordPoints)
	return strconv.Itoa((int)(tmp))
}

func hash(arr []uint) string {
	h := md5.New()
	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ""), "[]")
	io.WriteString(h, str)

	// return fmt.Sprintf("%x", h.Sum(nil))
	return str
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
