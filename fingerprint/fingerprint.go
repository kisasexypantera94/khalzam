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
	"strings"
)

const windowSize = 8192

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
	tmp := make([]int16, windowSize)
	for {
		buf := make([]byte, 2*len(tmp))
		_, err := decoder.Read(buf)

		if err != nil {
			break
		}

		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
		if channels == 2 {
			for i := 1; i < len(tmp); i += 2 {
				rawData = append(rawData, (float64)(tmp[i-1]+tmp[i])/2/math.MaxInt16)
			}
		} else {
			for i := 0; i < len(tmp); i++ {
				rawData = append(rawData, (float64)(tmp[i])/math.MaxInt16)
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

	return hash(recordPoints)
}

func hash(arr []uint) string {
	h := md5.New()
	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ""), "[]")
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))
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
