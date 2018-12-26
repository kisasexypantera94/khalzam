package fingerprint

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/jfreymuth/oggvorbis"
	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
	"io"
	"log"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const chunkSize = 1024
const fftWindowSize = 4096
const fuzzFactor = 2

var freqBins = [...]int16{40, 80, 120, 180, 300}

// DecodeMp3Int16 returns float32 slice of samples
func DecodeMp3Int16(filename string) []float64 {
	decoder, err := mpg123.NewDecoder("")
	checkErr(err)

	err = decoder.Open(filename)
	checkErr(err)
	defer decoder.Close()

	rate, channels, _ := decoder.GetFormat()
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	var pcmLeft []float32
	var pcmRight []float32
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
				pcmLeft = append(pcmLeft, (float32)(left)/(float32)(math.MaxInt16))
				pcmRight = append(pcmRight, (float32)(right)/(float32)(math.MaxInt16))
			}
		} else {
			for i := 0; i < len(tmp); i++ {
				mono := tmp[i]
				pcmLeft = append(pcmLeft, (float32)(mono)/(float32)(math.MaxInt16))
			}
		}
	}

	pcm64 := make([]float64, len(pcmLeft)+len(pcmRight))
	for i := range pcmLeft {
		pcm64[i] = (float64)(pcmLeft[i])
	}
	for i := range pcmRight {
		pcm64[i+len(pcmLeft)] = (float64)(pcmRight[i])
	}

	decoder.Delete()
	return pcm64
}

// DecodeOggFloat64 decodes ogg files
func DecodeOggFloat64(filename string) []float64 {
	f, _ := os.Open(filename)
	defer f.Close()
	var r io.Reader
	r = f
	pcm32, _, _ := oggvorbis.ReadAll(r)
	pcm64 := make([]float64, len(pcm32))
	for i := 0; i < len(pcm32); i++ {
		pcm64[i] = (float64)(pcm32[i])
	}
	fn, _ := os.Create(filename + ".raw")
	defer fn.Close()
	binary.Write(fn, binary.LittleEndian, pcm64)

	return pcm64
}

// DecodeWavFloat64 decodes wav file to slice of float64 values
func DecodeWavFloat64(filename string) []float64 {
	file, _ := os.Open(filename)
	reader := wav.NewReader(file)

	defer file.Close()

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

	f, _ := os.Create(filename + ".raw")
	binary.Write(f, binary.LittleEndian, pcm)

	return pcm
}

// Fingerprint returns a fingerprint of song
func Fingerprint(filename string) (hashArray []string) {
	var pcm64 []float64
	switch filepath.Ext(filename) {
	case ".mp3":
		pcm64 = DecodeMp3Int16(filename)

	case ".wav":
		pcm64 = DecodeWavFloat64(filename)

	case ".ogg":
		pcm64 = DecodeOggFloat64(filename)
	}

	blockNum := len(pcm64) / fftWindowSize

	for i := 0; i < blockNum; i++ {
		complexArray := fft.FFTReal(pcm64[i*fftWindowSize : i*fftWindowSize+fftWindowSize])
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

// stereoToMonoFFMPEG converts file to mono using ffmpeg
func stereoToMonoFFMPEG(filename string) string {
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
