# pdf2mp4

server-side converter from pdf to mp4.

# Using components

- [FFmpeg](https://ffmpeg.org/)
- [OpenH264](http://www.openh264.org/)
- pdftoppm of [Poppler](https://poppler.freedesktop.org/)

debian package(libopenh264-cisco6) use pre-build library by Cisco Systems, Inc.

```
+ pdftoppm -scale-to-x -1 -scale-to-y 1080 /tmp/example-1155511624/input.pdf
+ ffmpeg -r 1 -f ppm_pipe -i - -c:v libopenh264 -profile:v main -allow_skip_frames 1 -r 30 -y -f mp4 /tmp/example-1155511624/output.mp4
ffmpeg version 6.0 Copyright (c) 2000-2023 the FFmpeg developers
  built with gcc 11 (Ubuntu 11.4.0-1ubuntu1~22.04)
  configuration: --disable-all --disable-everything --enable-libopenh264 --enable-demuxer=image_ppm_pipe --enable-parser=pnm --enable-decoder=ppm --enable-encoder=libopenh264 --enable-protocol='fd,file' --enable-muxer=mp4 --enable-filter='scale,format' --enable-ffmpeg --enable-swscale --enable-avcodec --enable-avutil --enable-avfilter --enable-avformat
  libavutil      58.  2.100 / 58.  2.100
  libavcodec     60.  3.100 / 60.  3.100
  libavformat    60.  3.100 / 60.  3.100
  libavfilter     9.  3.100 /  9.  3.100
  libswscale      7.  1.100 /  7.  1.100
Input #0, ppm_pipe, from 'fd:':
  Duration: N/A, bitrate: N/A
  Stream #0:0: Video: ppm, rgb24, 1920x1080, 1 fps, 1 tbr, 1 tbn
Stream mapping:
  Stream #0:0 -> #0:0 (ppm (native) -> h264 (libopenh264))
Output #0, mp4, to '/tmp/example-1155511624/output.mp4':
  Metadata:
    encoder         : Lavf60.3.100
  Stream #0:0: Video: h264 (avc1 / 0x31637661), yuv420p(tv, progressive), 1920x1080, q=2-31, 30 fps, 15360 tbn
    Metadata:
      encoder         : Lavc60.3.100 libopenh264
    Side data:
      cpb: bitrate max/min/avg: 0/0/2000000 buffer size: 0 vbv_delay: N/A
frame=  303 fps=0.0 q=-0.0 Lsize=    1001kB time=00:00:10.96 bitrate= 747.8kbits/s dup=319 drop=0 speed=11.3x
video:999kB audio:0kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: 0.213499%
[libopenh264 @ 0x56003ab80d40] 27 frames skipped
```
