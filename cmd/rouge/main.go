package main

import (
	"github.com/mcandre/rouge"

	"flag"
	"fmt"
	"os"
)

var flagIn = flag.String("in", "", "WAV input file (default: raw stdin)")
var flagOut = flag.String("out", "", "WAV output file (default: raw stdout)")
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

		mod = rouge.NewWavModulator(file, dem.SampleRate(), dem.BitDepth(), dem.NumChannels(), dem.SampleRate(), dem.BitDepth(), dem.NumChannels(), dem.WavCategory())
	} else {
		mod = rouge.NewRawModulator(os.Stdout)
	}

	chIn := dem.Decoder()
	chOut, chErr := mod.Encoder()

	for m := range chIn {
		if m.Error != nil {
			panic(m.Error)
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
