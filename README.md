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
   -in /tmp/p0.wav \
   >/tmp/p0.pcm.dat

$ hh /tmp/p0.pcm.dat | head
00000000 00 00 00 00 00 00 00 00
00000008 00 00 00 00 00 00 00 00
00000016 00 00 00 00 00 00 00 00
00000024 00 00 00 00 00 00 00 00
00000032 00 00 00 00 00 00 00 00
00000040 00 00 00 00 00 00 00 00
00000048 00 00 00 00 00 00 00 00
00000056 00 00 00 00 00 00 00 00
00000064 00 00 00 00 00 00 00 00
00000072 00 00 00 00 00 00 00 00

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

$ hh /tmp/p0.te.dat | less
00000000 66 66 66 66 66 66 66 66
00000008 66 66 66 66 66 66 66 66
00000016 66 66 66 66 66 66 66 66
00000024 66 66 66 66 66 66 66 66
00000032 66 66 66 66 66 66 66 66
00000040 66 66 66 66 66 66 66 66
00000048 66 66 66 66 66 66 66 66
00000056 66 66 66 66 66 66 66 66
00000064 66 66 66 66 66 66 66 66
00000072 66 66 66 66 66 66 66 66
00000080 66 66 66 66 66 66 66 66
00000088 66 66 66 66 66 66 66 66
00000096 66 66 66 66 66 66 66 66
00000104 66 66 66 66 66 66 66 66
00000112 66 66 66 66 66 66 66 66
00000120 66 66 66 66 66 66 66 66
00000128 66 66 66 66 66 66 66 66
00000136 66 66 66 66 66 66 66 66
00000144 66 66 66 66 66 66 66 66
00000152 66 66 66 66 66 66 66 66
00000160 66 66 66 66 66 66 66 66
00000168 66 66 66 66 66 66 66 66
00000176 66 66 66 66 66 66 66 66
00000184 66 66 66 66 66 66 66 66
00000192 66 66 66 66 66 66 66 66
00000200 66 66 66 66 66 66 66 66
00000208 66 66 66 66 66 66 66 66
00000216 66 66 66 66 66 66 66 66
00000224 66 66 66 66 66 66 66 66
00000232 66 66 66 66 66 66 66 66
00000240 66 66 66 66 66 66 66 66
00000248 66 66 66 66 9a 45 84 08
00000256 57 7d 16 b3 a3 df 56 b7
00000264 7c 13 f9 f5 b1 f0 25 49
00000272 ab b6 73 d7 a6 dd 76 98
00000280 ef b8 0b 10 55 a6 3c f9
...
00000112 ea 64 4c 2b 57 29 19 e8
00000120 ff 1b f4 d6 fd 8b 32 3d
00000128 1a 76 b8 e7 2b a6 94 53
00000136 83 6e fc d9 6b 28 f0 9d
00000144 2e ab 6e a9 c4 a1 8c c1
00000152 c0 6f 42 66 af 79 b4 6f
00000160 8f 0e b4 42 3a c6 72 d1
00000168 13 19 6e 89 59 47 55 55
00000176 55
```

Second pattern active. [examples/p1.wav](examples/p1.wav)

Repeating the steps above for the second pattern, we have a basis for comparison in the encoded data:

```console
$ hh /tmp/p0.te.dat >/tmp/p0.te-hex.txt
$ hh /tmp/p1.te.dat >/tmp/p1.te-hex.txt
$ diff -u /tmp/p0.te-hex.txt /tmp/p1.te-hex.txt
--- /tmp/p0.te-hex.txt	2020-05-03 12:53:23.000000000 -0500
+++ /tmp/p1.te-hex.txt	2020-05-03 12:53:28.000000000 -0500
@@ -658,6 +658,6 @@
 00000136 83 6e fc d9 6b 28 f0 9d
 00000144 2e ab 6e a9 c4 a1 8c c1
 00000152 c0 6f 42 66 af 79 b4 6f
-00000160 8f 0e b4 42 3a c6 72 d1
-00000168 13 19 6e 89 59 47 55 55
-00000176 55
+00000160 8f 0e b4 42 3a c6 f2 ad
+00000168 69 4a 22 94 9b 92 aa aa
+00000176 aa
```

That is, eleven bytes differ in the backup stream when the active pattern is switched from the first pattern to the second.

Note that the last three bytes of each stream are repeated. Perhaps a checksum?

Repeating the steps above for the 3rd - 7th patterns yields similar results.

Eighth pattern active. [examples/p7.wav](examples/p7.wav)

```console
$ hh /tmp/p7.te.dat >/tmp/p7.te-hex.txt
$ diff -u /tmp/p0.te-hex.txt /tmp/p7.te-hex.txt
--- /tmp/p0.te-hex.txt	2020-05-03 12:53:23.000000000 -0500
+++ /tmp/p7.te-hex.txt	2020-05-03 12:53:51.000000000 -0500
@@ -654,10 +654,10 @@
 00000104 05 e3 40 06 53 8f c9 ae
 00000112 ea 64 4c 2b 57 29 19 e8
 00000120 ff 1b f4 d6 fd 8b 32 3d
-00000128 1a 76 b8 e7 2b a6 94 53
-00000136 83 6e fc d9 6b 28 f0 9d
-00000144 2e ab 6e a9 c4 a1 8c c1
-00000152 c0 6f 42 66 af 79 b4 6f
-00000160 8f 0e b4 42 3a c6 72 d1
-00000168 13 19 6e 89 59 47 55 55
+00000128 1a 76 b8 e7 69 f9 23 1d
+00000136 a1 f9 08 93 de f7 36 bf
+00000144 b1 1d ee d4 33 43 0e b8
+00000152 35 23 e9 35 a6 95 5b 16
+00000160 04 29 5d 48 f5 68 a5 4c
+00000168 53 90 31 33 16 8b 55 55
 00000176 55
```

Now a larger portion of the signal is altered. Repeating the steps above for the nineth - sixteenth patterns yields results closer to pattern eight than to pattern one. The bits preserving the active pattern ID are more complicated than they first appear.

See `rouge -help` for more options.

## Additional notes

Repeated export transmissions of the same active pattern yield identical hex decodings.

However, performing a factory reset (`Write + Pattern + insert batteries`) appears to alter

# REQUIREMENTS

* [Go](https://golang.org/) 1.12+
* [PO-32 Tonic](https://teenage.engineering/products/po-32)

## Recommended

* [Audacity](https://www.audacityteam.org/)
* [ffmpeg](https://www.ffmpeg.org/)
* [mplayer](http://www.mplayerhq.hu/)
* [hellcat](https://github.com/mcandre/hellcat)
* [objdump](https://linux.die.net/man/1/objdump)
* [diff](https://linux.die.net/man/1/diff)
* [git](https://git-scm.com/)

# BUILD + INSTALL

```console
$ git submodule update --init --recursive
$ go install ./...
```
