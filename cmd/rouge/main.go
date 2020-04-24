package main

import (
	"github.com/mcandre/rouge"

	"flag"
	"fmt"
	"os"
)

var flagIn = flag.String("in", "", "WAV input file (default: raw stdin)")
var flagOut = flag.String("out", "", "WAV output file (default: raw stdout)")
var flagSampleRateIn = flag.Int("sampleRateIn", -1, "Specify sample rate in, e.g. with raw stdin data")
var flagBitDepthIn = flag.Int("bitDepthIn", -1, "Specify bit depth in, e.g. with raw stdin data")
var flagChannelsIn = flag.Int("channelsIn", -1, "Specify number of audio channels in, e.g. with raw stdin data")
var flagCategoryIn = flag.Int("categoryIn", -1, "Specify WAV category in, e.g. with raw stdin data. PCM is category 1.")
var flagVersion = flag.Bool("version", false, "Show version information")

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(rouge.Version)
		os.Exit(0)
	}

	var dem rouge.Demodulator

	if *flagIn != "" {
		file, err := os.Open(*flagIn)

		if err != nil {
			panic(err)
		}

		dm, err := rouge.NewWavDemodulator(file)

		if err != nil {
			panic(err)
		}

		dem = dm
	} else {
		dem = rouge.NewRawDemodulator(os.Stdin)
	}

	var mod rouge.Modulator

	if *flagOut != "" {
		file, err := os.Create(*flagOut)

		if err != nil {
			panic(err)
		}

		sampleRate := dem.SampleRate()

		if *flagSampleRateIn != -1 {
			sampleRate = uint32(*flagSampleRateIn)
		}

		bitDepth := dem.BitDepth()

		if *flagBitDepthIn != -1 {
			bitDepth = uint16(*flagBitDepthIn)
		}

		channelsIn := dem.NumChannels()

		if *flagChannelsIn != -1 {
			channelsIn = uint16(*flagChannelsIn)
		}

		categoryIn := dem.WavCategory()

		if *flagCategoryIn != -1 {
			categoryIn = uint16(*flagCategoryIn)
		}

		mod = rouge.NewWavModulator(file, sampleRate, bitDepth, channelsIn, sampleRate, bitDepth, channelsIn, categoryIn)
	} else {
		mod = rouge.NewRawModulator(os.Stdout)
	}

	chIn := dem.Decoder()
	chOut, chErr := mod.Encoder()

	for m := range chIn {
		if m.Error != nil {
			panic(*m.Error)
		}

		select {
		case err := <-chErr:
			panic(err)
		case chOut<-m:
		}

		if m.Done {
			break
		}
	}
}
