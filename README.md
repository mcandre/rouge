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

## Additional notes

Repeated export transmissions of the same active pattern yield identical hex decodings.

The last three bytes of each stream are repeated. So far, they are restricted to either `0x555555` or `0xaaaaaa`.

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
