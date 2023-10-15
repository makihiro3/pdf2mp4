FROM golang:1.21.3 AS build-server

WORKDIR /work
COPY ./ ./
ARG CGO_ENABLED=0
RUN go build -o pdf2mp4

FROM buildpack-deps:jammy AS builder
RUN apt-get update
RUN apt-get install -y --no-install-recommends libopenh264-dev nasm
WORKDIR /usr/src

RUN curl -LO https://ffmpeg.org/releases/ffmpeg-6.0.tar.xz
RUN curl -LO https://ffmpeg.org/releases/ffmpeg-6.0.tar.xz.asc
COPY ffmpeg-devel.asc ./
RUN gpg --import ffmpeg-devel.asc
RUN gpg --verify ffmpeg-6.0.tar.xz.asc ffmpeg-6.0.tar.xz
RUN tar xaf ffmpeg-6.0.tar.xz
WORKDIR ffmpeg-6.0
RUN  ./configure --disable-all --disable-everything --enable-libopenh264 --enable-demuxer=image_ppm_pipe --enable-parser=pnm --enable-decoder=ppm --enable-encoder=libopenh264 --enable-protocol='fd,file' --enable-muxer=mp4 --enable-filter='scale,format' --enable-ffmpeg --enable-swscale --enable-avcodec --enable-avutil --enable-avfilter --enable-avformat
RUN make -j $(nproc)
RUN make install

FROM ubuntu:jammy
RUN apt-get update && apt-get install -y --no-install-recommends libopenh264-cisco6 poppler-utils && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN install -d -o root -g daemon -m 0775 /sock
COPY --from=builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build-server /work/pdf2mp4 /pdf2mp4
COPY run.sh /run.sh
USER daemon
ENV LISTEN=/sock/listen.sock
CMD ["/pdf2mp4"]
STOPSIGNAL SIGINT
