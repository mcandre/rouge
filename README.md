# rouge: a work in progress BPSK modem

```text
 ^   A
>θ.θ>/
   ♯
```

# EXAMPLES

PO-32 factory default backup. [examples/pattern-0.wav](examples/pattern-0.wav)

## Read/Write audio file samples

Rouge can manipulate WAVE PCM mono samples as 32-bit 2's complement integers.

```console
$ mplayer examples/pattern-0.wav

$ ffmpeg \
   -i examples/pattern-0.wav \
   -ac 1 \
   /tmp/pattern-0-mono.wav

$ rouge \
   -in /tmp/pattern-0-mono.wav \
   >/tmp/pattern-0-mono.dat

$ hh /tmp/pattern-0-mono.dat | head
00000000 00 00 00 0c ff ff ff ef
00000008 ff ff fe b5 ff ff f8 75
00000016 ff ff f3 73 ff ff f3 ce
00000024 ff ff f4 ba ff ff fa 01
00000032 00 00 00 61 00 00 07 ea
00000040 00 00 0b 6e 00 00 0c df
00000048 00 00 0b 27 00 00 07 0c
00000056 00 00 00 34 ff ff f9 2a
00000064 ff ff f4 ae ff ff f3 5a
00000072 ff ff f4 7a ff ff f8 77

$ rouge \
   -out /tmp/pattern-0-mono.copy.wav \
   -sampleRate 44100 \
   -bitDepth 16 \
   -channels 1 \
   -category 1 \
   </tmp/pattern-0-mono.dat

$ mplayer /tmp/pattern-0-mono.copy.wav

$ diff \
   /tmp/pattern-0-mono.wav \
   /tmp/pattern-0-mono.copy.wav;
   echo "$?"
0
```

## Read BPSK samples

```console
$ rouge \
   -bpskIn /tmp/pattern-0-mono.dat \
   -innerThreshold 1000 \
   -outerThreshold 1200 \
   -bitWindow 6 \
   >/tmp/pattern-0.dat

$ hh /tmp/pattern-0.dat | less
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
00000288 f9 81 3b 68 dd 03 69 a9
00000296 f7 fe c8 28 d8 35 33 c6
00000304 f4 66 86 6c 51 65 1b e2
00000312 b2 01 c9 3f 96 3a 91 3a
...
0000112 ea 64 4c 2b 57 29 19 e8
00000120 ff 1b f4 d6 fd 8b 32 3d
00000128 1a 76 b8 e7 2b a6 94 53
00000136 83 6e fc d9 6b 28 64 42
00000144 34 8b fd 5f 20 3d 79 fb
00000152 42 b4 de 64 19 88 21 a9
00000160 40 54 22 54 99 e4 83 4c
00000168 e5 ee fe 5d 7c f9 55 55
00000176 55
```

A backup with the second pattern as the active pattern. [examples/pattern-1.wav](examples/pattern-1.wav)

Repeating the steps above for the second pattern, we have a basis for comparison in the encoded data:

```console
$ hh /tmp/pattern-0.dat >/tmp/pattern-0-hex.txt
$ hh /tmp/pattern-1.dat >/tmp/pattern-1-hex.txt
$ diff -u /tmp/pattern-0-hex.txt /tmp/pattern-1-hex.txt
--- /tmp/pattern-0-hex.txt      2020-05-02 22:23:54.000000000 -0500
+++ /tmp/pattern-1-hex.txt      2020-05-02 22:23:57.000000000 -0500
@@ -658,6 +658,6 @@
 00000136 83 6e fc d9 6b 28 64 42
 00000144 34 8b fd 5f 20 3d 79 fb
 00000152 42 b4 de 64 19 88 21 a9
-00000160 40 54 22 54 99 e4 83 4c
-00000168 e5 ee fe 5d 7c f9 55 55
-00000176 55
+00000160 40 54 22 54 99 e4 7c cf
+00000168 42 62 73 bf 41 ac aa aa
+00000176 aa
```

That is, eleven bytes differ in the backup stream when the active pattern is switched from the first pattern to the second. Note that the last three bytes of each stream are repeated. Perhaps a checksum?

See `rouge -help` for more options.

# REQUIREMENTS

* [Go](https://golang.org/) 1.12+
* [PO-32 Tonic](https://teenage.engineering/products/po-32)

## Recommended

* [Audacity](https://www.audacityteam.org/)
* [ffmpeg](https://www.ffmpeg.org/)
* [mplayer](http://www.mplayerhq.hu/)
* [hellcat](https://github.com/mcandre/hellcat)
* [objdump](https://linux.die.net/man/1/objdump)

# BUILD + INSTALL

```console
$ go install ./...
```
