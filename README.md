# pdf2mp4

server-side converter from pdf to mp4.

# Using components

- [FFmpeg](https://ffmpeg.org/)
- [OpenH264](http://www.openh264.org/): OpenH264 Video Codec provided by Cisco Systems, Inc.
- pdftoppm of [Poppler](https://poppler.freedesktop.org/)

debian package(libopenh264-cisco8) use pre-build library by Cisco Systems, Inc.

```
+ setpriv --reuid=daemon --regid=daemon --init-groups --no-new-privs pdftoppm -scale-to-x -1 -scale-to-y 1080 /tmp/example-3515530111/input.pdf
+ setpriv --reuid=daemon --regid=daemon --init-groups --no-new-privs ffmpeg -r 1 -f ppm_pipe -i - -c:v libopenh264 -profile:v high -allow_skip_frames 1 -coder cabac -qmin 1 -qmax 50 -r 30 -y -f mp4 /tmp/example-3515530111/output.mp4
ffmpeg version 8.0 Copyright (c) 2000-2025 the FFmpeg developers
  built with gcc 14 (Debian 14.2.0-19)
  configuration: --disable-all --disable-everything --enable-libopenh264 --enable-demuxer=image_ppm_pipe --enable-parser=pnm --enable-decoder=ppm --enable-encoder=libopenh264 --enable-protocol='fd,file' --enable-muxer=mp4 --enable-filter='scale,format' --enable-ffmpeg --enable-swscale --enable-avcodec --enable-avutil --enable-avfilter --enable-avformat
  libavutil      60.  8.100 / 60.  8.100
  libavcodec     62. 11.100 / 62. 11.100
  libavformat    62.  3.100 / 62.  3.100
  libavfilter    11.  4.100 / 11.  4.100
  libswscale      9.  1.100 /  9.  1.100
Input #0, ppm_pipe, from 'fd:':
  Duration: N/A, bitrate: N/A
  Stream #0:0: Video: ppm, rgb24, 1920x1080, 1 fps, 1 tbr, 1 tbn
Stream mapping:
  Stream #0:0 -> #0:0 (ppm (native) -> h264 (libopenh264))
Output #0, mp4, to '/tmp/example-3515530111/output.mp4':
  Metadata:
    encoder         : Lavf62.3.100
  Stream #0:0: Video: h264 (avc1 / 0x31637661), yuv420p(tv, progressive), 1920x1080, q=2-31, 30 fps, 15360 tbn
    Metadata:
      encoder         : Lavc62.11.100 libopenh264
    Side data:
      cpb: bitrate max/min/avg: 0/0/2000000 buffer size: 0 vbv_delay: N/A
[out#0/mp4 @ 0x5c64b9a78980] video:998KiB audio:0KiB subtitle:0KiB other streams:0KiB global headers:0KiB muxing overhead: 0.214058%
frame=  303 fps= 29 q=-0.0 Lsize=    1000KiB time=00:00:11.00 bitrate= 744.6kbits/s dup=319 drop=0 speed=1.04x elapsed=0:00:10.52
[libopenh264 @ 0x5c64b9a74f80] 27 frames skipped
```
