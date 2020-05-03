package main

import (
	"github.com/mcandre/rouge"

	"flag"
	"fmt"
	"os"
)

var flagIn = flag.String("in", "", "WAV input file (default: raw stdin)")
var flagOut = flag.String("out", "", "WAV output file (default: raw stdout)")
var flagBPSKIn = flag.String("bpskIn", "", "BPSK input file")
var flagBitWindow = flag.Int("bitWindow", 0, "Specify how many peaks/vallyes constitute a BPSK bit")
var flagInnerThreshold = flag.Float64("innerThreshold", 0.0, "Specify maximum amplitude for BPSK lulls")
var flagOuterThreshold = flag.Float64("outerThreshold", 0.0, "Specify minimum amplitude for BPSK peaks")
var flagSampleRate = flag.Int("sampleRate", 0, "Specify sample rate in, e.g. with raw stdin data")
var flagBitDepth = flag.Int("bitDepth", 0, "Specify bit depth in, e.g. with raw stdin data")
var flagChannels = flag.Int("channels", 0, "Specify number of audio channels in, e.g. with raw stdin data")
var flagCategory = flag.Int("category", 0, "Specify WAV category in, e.g. with raw stdin data. PCM is category 1.")
var flagVersion = flag.Bool("version", false, "Show version information")

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(rouge.Version)
		os.Exit(0)
	}

	var dem rouge.Demodulator

	if *flagBPSKIn != "" {
		if *flagBitWindow == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -bitWindow\n")
			os.Exit(1)
		}

		if *flagInnerThreshold == 0.0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -innerThreshold\n")
			os.Exit(1)
		}

		if *flagOuterThreshold == 0.0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -outerThreshold\n")
			os.Exit(1)
		}

		file, err := os.Open(*flagBPSKIn)

		if err != nil {
			panic(err)
		}

		dem = rouge.NewBPSKDemodulator(file, *flagBitWindow, *flagInnerThreshold, *flagOuterThreshold)
	} else if *flagIn != "" {
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

		if *flagSampleRate != 0 {
			sampleRate = *flagSampleRate
		}

		if sampleRate == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -sampleRate\n")
			os.Exit(1)
		}

		bitDepth := dem.BitDepth()

		if *flagBitDepth != 0 {
			bitDepth = *flagBitDepth
		}

		if bitDepth == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -bitDepth\n")
			os.Exit(1)
		}

		numChannels := dem.NumChannels()

		if *flagChannels != 0 {
			numChannels = *flagChannels
		}

		if numChannels == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -channels\n")
			os.Exit(1)
		}

		wavCategory := dem.WavCategory()

		if *flagCategory != 0 {
			wavCategory = *flagCategory
		}

		if wavCategory == 0 {
			fmt.Fprintf(os.Stderr, "Missing -in or -category\n")
			os.Exit(1)
		}

		mod = rouge.NewWavModulator(rouge.WavModulatorConfig{
			File: file,
			SampleRate: uint32(sampleRate),
			BitDepth: uint16(bitDepth),
			NumChannels: uint16(numChannels),
			WavCategory: uint16(wavCategory),
		})
	} else {
		mod = rouge.NewRawModulator(os.Stdout)
	}

	chIn := dem.Decoder()
	chDone, chOut, chErr := mod.Encoder()

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
			<-chDone
			break
		}
	}
}
