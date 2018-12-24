package fingerprint

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mjibson/go-dsp/fft"
	"io"
	"log"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const chunkSize = 1024
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

	var rawData1 []float32
	var rawData2 []float32
	tmp := make([]int16, chunkSize/2)
	for {
		buf := make([]byte, chunkSize)
		_, err := decoder.Read(buf)

		if err != nil {
			break
		}

		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
		if channels == 2 {
			for i := 0; i < len(tmp); i += 2 {
				left := (tmp[i])
				right := (tmp[i+1])
				rawData1 = append(rawData1, (float32)(left)/(float32)(math.MaxInt16))
				rawData2 = append(rawData2, (float32)(right)/(float32)(math.MaxInt16))
			}
		} else {
			for i := 0; i < len(tmp); i++ {
				x := tmp[i]
				rawData1 = append(rawData1, (float32)(x)/(float32)(math.MaxInt16))
			}
		}
	}

	rawData64 := make([]float64, len(rawData1) + len(rawData2))
	for i := range rawData1 {
		rawData64[i] = (float64)(rawData1[i])
	}
	for i := range rawData2 {
		rawData64[i + len(rawData1)] = (float64)(rawData2[i])
	}

	decoder.Delete()
	return rawData64
}

func Decode2(filename string) []float64 {
	file, _ := os.Open(filename)
	defer file.Close()
	d, _ := mp3.NewDecoder(file)
	defer d.Close()

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	written, _ := io.Copy(w, d)
	w.Flush()
	r := bufio.NewReader(&buf)

	tmp := make([]int16, written/2)
	var rawData []float32
	binary.Read(r, binary.LittleEndian, tmp)
	for i := 0; i < len(tmp); i += 2 {
		x := (tmp[i] + tmp[i+1]) / 2
		rawData = append(rawData, (float32)(x)/(float32)(math.MaxInt16))
	}

	rawData64 := make([]float64, len(rawData))
	for i := range rawData {
		fmt.Println(rawData[i])
		rawData64[i] = (float64)(rawData[i])
	}
	f, _ := os.Create(filename + "raw.raw")
	binary.Write(f, binary.LittleEndian, rawData64)
	f.Close()

	return rawData64
}

// Fingerprint returns a fingerprint of song
func Fingerprint(filename string) (hashArray []string) {
	rawData := Decode(filename)
	f, _ := os.Create(filename + "raw.raw")
	binary.Write(f, binary.LittleEndian, rawData)
	f.Close()
	blockNum := len(rawData) / windowSize

	for i := 0; i < blockNum; i++ {
		complexArray := fft.FFTReal(rawData[i*windowSize : i*windowSize+windowSize])
		// complexArray := make([]complex128, windowSize)
		// for j := 0; j < windowSize; j++ {
		// 	complexArray[j] = complex(rawData[i*windowSize+j], 0)
		// }
		// plan := fftw.NewPlan1d(complexArray, false, false)

		// // perform Fourier transform
		// plan.Execute()
		// plan.Free()
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
