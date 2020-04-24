package main

import (
	"github.com/mcandre/rouge"

	"flag"
	"fmt"
	"os"
)

var flagIn = flag.String("in", "", "WAV input file (default: raw stdin)")
var flagOut = flag.String("out", "", "WAV output file (default: raw stdout)")
var flagSampleRateIn = flag.Int("sampleRateIn", 0, "Specify sample rate in, e.g. with raw stdin data")
var flagBitDepthIn = flag.Int("bitDepthIn", 0, "Specify bit depth in, e.g. with raw stdin data")
var flagChannelsIn = flag.Int("channelsIn", 0, "Specify number of audio channels in, e.g. with raw stdin data")
var flagCategoryIn = flag.Int("categoryIn", 0, "Specify WAV category in, e.g. with raw stdin data. PCM is category 1.")
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

		if *flagSampleRateIn != 0 {
			sampleRate = *flagSampleRateIn
		}

		if sampleRate == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -sampleRateIn\n")
			os.Exit(1)
		}

		bitDepth := dem.BitDepth()

		if *flagBitDepthIn != 0 {
			bitDepth = *flagBitDepthIn
		}

		if bitDepth == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -bitDepthIn\n")
			os.Exit(1)
		}

		numChannels := dem.NumChannels()

		if *flagChannelsIn != 0 {
			numChannels = *flagChannelsIn
		}

		if numChannels == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -channelsIn\n")
			os.Exit(1)
		}

		wavCategory := dem.WavCategory()

		if *flagCategoryIn != 0 {
			wavCategory = *flagCategoryIn
		}

		if wavCategory == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -categoryIn\n")
			os.Exit(1)
		}

		mod = rouge.NewWavModulator(file, uint32(sampleRate), uint16(bitDepth), uint16(numChannels), uint32(sampleRate), uint16(bitDepth), uint16(numChannels), uint16(wavCategory))
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
