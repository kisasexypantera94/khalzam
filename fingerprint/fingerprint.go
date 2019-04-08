package fingerprint

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/jfreymuth/oggvorbis"
	"github.com/kisasexypantera94/go-mpg123/mpg123"
	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
	"io"
	"log"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
#include <mpg123.h>
#cgo LDFLAGS: -lmpg123
*/
import "C"

const chunkSize = 2048
const fftWindowSize = 4096
const fuzzFactor = 2

var freqBins = [...]int16{40, 80, 120, 180, 300}

// DecodeMp3 decodes mp3 files using `libmpg123`
func DecodeMp3(filename string) []float64 {
	decoder, err := mpg123.NewDecoder("", C.MPG123_MONO_MIX|C.MPG123_FORCE_FLOAT)
	checkErr(err)

	err = decoder.Open(filename)
	checkErr(err)
	defer decoder.Close()

	decoder.GetFormat()

	var pcm64 []float64
	tmp := make([]float32, chunkSize/4)
	for {
		buf := make([]byte, chunkSize)
		_, err := decoder.Read(buf)

		if err != nil {
			break
		}

		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, tmp)
		for i := 0; i < len(tmp); i++ {
			mono := tmp[i]
			pcm64 = append(pcm64, (float64)(mono))
		}
	}

	decoder.Delete()
	return pcm64
}

// DecodeOgg decodes ogg files
func DecodeOgg(filename string) []float64 {
	f, _ := os.Open(filename)
	defer f.Close()
	var r io.Reader
	r = f
	pcm32, _, _ := oggvorbis.ReadAll(r)
	pcm64 := make([]float64, len(pcm32))
	for i := 0; i < len(pcm32); i++ {
		pcm64[i] = (float64)(pcm32[i])
	}

	return pcm64
}

// DecodeWav decodes wav files
func DecodeWav(filename string) []float64 {
	file, _ := os.Open(filename)
	defer file.Close()
	reader := wav.NewReader(file)

	var pcm []float64
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			pcm = append(pcm, (reader.FloatValue(sample, 0)+reader.FloatValue(sample, 1))/2)
		}
	}

	return pcm
}

// Fingerprint constructs fingerprint for song
func Fingerprint(filename string) (hashArray []int, err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("Fingerprint: file not found")
	}

	var pcm64 []float64
	switch filepath.Ext(filename) {
	case ".mp3":
		pcm64 = DecodeMp3(filename)
	case ".wav":
		pcm64 = DecodeWav(filename)
	case ".ogg":
		pcm64 = DecodeOgg(filename)
	default:
		return nil, fmt.Errorf("Fingerprint: invalid file")
	}

	blockNum := len(pcm64) / fftWindowSize
	for i := 0; i < blockNum; i++ {
		complexArray := fft.FFTReal(pcm64[i*fftWindowSize : i*fftWindowSize+fftWindowSize])
		hashArray = append(hashArray, getKeyPoints(complexArray))
	}
	return hashArray, nil
}

type output struct {
	idx int
	val int
}

// ParallelFingerprint constructs fingerprint for song.
// It perform FFT on chunks in goroutines.
func ParallelFingerprint(filename string) (hashArray []int, err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("Fingerprint: file not found")
	}

	var pcm64 []float64
	switch filepath.Ext(filename) {
	case ".mp3":
		pcm64 = DecodeMp3(filename)
	case ".wav":
		pcm64 = DecodeWav(filename)
	case ".ogg":
		pcm64 = DecodeOgg(filename)
	default:
		return nil, fmt.Errorf("Fingerprint: invalid file")
	}

	blockNum := len(pcm64) / fftWindowSize
	wg := new(sync.WaitGroup)
	ch := make(chan output, 100)
	hashArray = make([]int, blockNum, blockNum)

	for i := 0; i < blockNum; i++ {
		wg.Add(1)
		go func(idx int, ch chan output, wg *sync.WaitGroup) {
			defer wg.Done()
			complexArray := fft.FFTReal(pcm64[idx*fftWindowSize : idx*fftWindowSize+fftWindowSize])
			out := output{idx, getKeyPoints(complexArray)}
			ch <- out
		}(i, ch, wg)
	}

	quit := make(chan bool)
	go func() {
		wg.Wait()
		quit <- true
	}()

	for {
		select {
		case output := <-ch:
			hashArray[output.idx] = output.val
		case <-quit:
			return hashArray, nil
		}
	}
}

func getKeyPoints(compArr []complex128) int {
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

func hash(arr []uint) int {
	tmp := (arr[3]-(arr[3]%fuzzFactor))*1e8 +
		(arr[2]-(arr[2]%fuzzFactor))*1e5 +
		(arr[1]-(arr[1]%fuzzFactor))*1e2 +
		(arr[0] - (arr[0] % fuzzFactor))

	return int(tmp)
}

func hashMd5(arr []uint) string {
	h := md5.New()
	str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ""), "[]")
	io.WriteString(h, str)

	return str
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
