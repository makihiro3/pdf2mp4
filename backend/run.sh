#!/bin/sh
set -xe
pdftoppm "$1" | ffmpeg -r 1/2 -f ppm_pipe -i - -c:v libopenh264 -profile:v main -allow_skip_frames 1 -r 30 -y -f mp4 "$2"
