# rouge: a work in progress BPSK modem

```text
 ^   A
>θ.θ>/
   ♯
```

# DISCLAIMER

This project is intended purely for educational research purposes. Make sure to support your local synth shop! The findings here are provisional and incomplete. No guarantees. Your mileage may vary, etc. etc.

# EXAMPLES

PO-32 factory default backup (first pattern active). [examples/p0.wav](examples/p0.wav)

## Read/Write audio file samples

Rouge can manipulate WAVE PCM mono samples as 32-bit 2's complement integers.

```console
$ mplayer examples/p0.wav

$ rouge \
   -in examples/p0.wav \
   >/tmp/p0.pcm.dat

$ rouge \
   -out /tmp/p0.copy.wav \
   -sampleRate 44100 \
   -bitDepth 16 \
   -channels 1 \
   -category 1 \
   </tmp/p0.pcm.dat

$ mplayer /tmp/p0.copy.wav
```

`diff` can also confirm PCM data equivalence, as long as the metadata is identical.

## Read BPSK samples

The next step is to apply the right phase shift keying algorithm to the signal. A simple BPSK attempt over a six-peak window yields some early results.

```console
$ rouge \
   -bpskIn /tmp/p0.pcm.dat \
   -innerThreshold 1000 \
   -outerThreshold 1200 \
   -bitWindow 6 \
   >/tmp/p0.te.dat

$ ls -Ahl /tmp/p0.te.dat
-rw-r--r--  1 andrew  staff   5.2K May  3 12:53 /Users/andrew/Downloads/p0.te.dat

That's in the ball park for a file size to completely model all the PO-32 Tonic parameters.

$ hexdump /tmp/p0.te.dat | head
0000000 66 66 66 66 66 66 66 66 66 66 66 66 66 66 66 66
*
00000f0 66 66 66 66 66 66 66 66 66 66 66 66 9a 45 84 08
0000100 57 7d 16 b3 a3 df 56 b7 7c 13 f9 f5 b1 f0 25 49
0000110 ab b6 73 d7 a6 dd 76 98 ef b8 0b 10 55 a6 3c f9
0000120 f9 81 3b 68 dd 03 69 a9 f7 fe c8 28 d8 35 33 c6
0000130 f4 66 86 6c 51 65 1b e2 b2 01 c9 3f 96 3a 91 3a
0000140 de ac 98 a4 b8 3e 13 d1 4a 3f d5 fb 8c c3 f6 40
0000150 d6 57 e9 ee 52 66 a3 7d 7b 5b ea fa 4f 41 28 5b
0000160 20 64 1d fe ed de 17 54 76 dc c7 90 08 e9 0e 8c

$ hexdump /tmp/p0.te.dat >/tmp/p0.te-hex.txt
```

Second pattern active. [examples/p1.wav](examples/p1.wav)

Repeating the steps above for the second pattern, we have a basis for comparison in the encoded data:

```console
$ diff -u /tmp/p0.te-hex.txt /tmp/p1.te-hex.txt
--- /tmp/p0.te-hex.txt	2020-05-03 14:03:49.000000000 -0500
+++ /tmp/p1.te-hex.txt	2020-05-03 14:02:23.000000000 -0500
@@ -315,6 +315,6 @@
 0001470 ea 64 4c 2b 57 29 19 e8 ff 1b f4 d6 fd 8b 32 3d
 0001480 1a 76 b8 e7 2b a6 94 53 83 6e fc d9 6b 28 f0 9d
 0001490 2e ab 6e a9 c4 a1 8c c1 c0 6f 42 66 af 79 b4 6f
-00014a0 8f 0e b4 42 3a c6 72 d1 13 19 6e 89 59 47 55 55
-00014b0 55
+00014a0 8f 0e b4 42 3a c6 f2 ad 69 4a 22 94 9b 92 aa aa
+00014b0 aa
 00014b1
```

That is, eleven bytes differ in the backup stream when the active pattern is switched from the first pattern to the second.

Repeating the steps above for the 3rd - 16th patterns yields similar results, though the size of the difference varies from somewhat larger to somewhat smaller. The bits preserving the active pattern ID are more complicated than they first appear.

Based on the complexity of individual sound configurations and the size of the relative diff between pattern backup signals, we can conclude that only a handful of pattern settings are present in the diffs, enough to signify which sounds trigger on which 16 sequencer steps. If any effect or motion effect automation is applied, that remains default and would appear to be absent from the signal.

A sound pattern preset could be directly represented with two bytes for trigger sequence on/off state across sixteen steps, with two additional bytes for accents, and two more bytes for fills, amounting to six bytes.

```
1: {
	Triggers: #-----#---#--#--
	Accents:  ----------------
	Fills:    ----------------
}
```

With eight drum sounds, the pattern configuration is 48 bytes, or 24 hex pairs. In the neighborhood of how many hexpairs differ between the first factory pattern configuration and the selection of the second factory pattern configuration.

## Additional notes

Repeated export transmissions of the same active pattern yield identical hex decodings.

The last three bytes of each stream are repeated. So far, they are restricted to either `0x555555` or `0xaaaaaa`.

There are some string values in presets:

```
MicrotonicPresetV3: {
	Tempo: 121.00000000 bpm
	Pattern: a
	StepRate: 1/16
	Swing: 0.00000000%
	FillRate: 2.00000000x
	MastVol: 0.00000000 dB
	Mutes: { Off, Off, Off, Off, Off, Off, Off, Off }
	DrumPatches: {
		1: {
			Name: "SC BD Power"
			Modified: true
			Path: "./By Category/Bass Drum Patches/SC BD Power.mtdrum"
```

Unknown whether the strings are preserved in PO-32 signals, they may just be for VST lookup convenience.

## End Goal

Find an analog + digital decoding sufficient to customize sounds:

```
MicrotonicDrumPatchV3: {
	OscWave: Sine
	OscFreq: 616.14874268 Hz
	OscAtk: 0.00000000 ms
	OscDcy: 71.25314331 ms
	ModMode: Sine
	ModRate: 279.64397346 Hz
	ModAmt: +32.85273205 sm
	NFilMod: BP
	NFilFrq: 1005.14013672 Hz
	NFilQ: 0.70710677
	NStereo: Off
	NEnvMod: Exp
	NEnvAtk: 0.00000000 ms
	NEnvDcy: 60.78934776 ms
	Mix: 50.00000000 / 50.00000000
	DistAmt: 0.00000000
	EQFreq: 1111.70104980 Hz
	EQGain: +34.53779602 dB
	Level: 0.00000000 dB
	Pan: 0.00000000
	OscVel: 28.42212677%
	NVel: 34.16033936%
	ModVel: 0.00000000%
}

MicroTonicDrumPatchV1={
	Name="NE RND Zeroto9"
	Modified=true
	Path="/Library/Audio/Presets/Sonic Charge/MicroTonic Drum Patches/Effect Patches/NE RND Zeroto9.mtdp"
	OscWave=Triangle
	OscFreq=742.02076800Hz
	OscDcy=58.21029589ms
	ModMode=Noise
	ModRate=0 Hz
	ModAmt=+22.19031754sm
	NFilMod=BP
	NFilFrq=149.74901463Hz
	NFilQ=1008.40029486
	NStereo=Off
	NEnvMod=Exp
	NEnvAtk=0.00000000ms
	NEnvDcy=216.86110703ms
	Mix=66.51028097/33.48971903
	DistAmt=17.84428060
	EQFreq=5377.53820357Hz
	EQGain=+28.41660023dB
	Level=+10.00000000dB
	Pan=0.00000000
	Output=A
	OscVel=100.00000000%
	NVel=100.00000000%
	ModVel=100.00000000%
}
```

# REQUIREMENTS

* [Go](https://golang.org/) 1.12+
* [PO-32 Tonic](https://teenage.engineering/products/po-32)

## Recommended

* [Audacity](https://www.audacityteam.org/)
* [ffmpeg](https://www.ffmpeg.org/)
* [mplayer](http://www.mplayerhq.hu/)
* [hexdump](http://man7.org/linux/man-pages/man1/hexdump.1.html)
* [diff](https://linux.die.net/man/1/diff)
* [git](https://git-scm.com/)

# BUILD + INSTALL

```console
$ git submodule update --init --recursive
$ go install ./...
```
