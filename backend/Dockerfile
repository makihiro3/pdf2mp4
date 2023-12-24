FROM golang:1.21.4 AS build-server

WORKDIR /work
COPY ./ ./
ARG CGO_ENABLED=0
RUN go build -o pdf2mp4

FROM buildpack-deps:jammy AS builder
RUN apt-get update
RUN apt-get install -y --no-install-recommends libopenh264-dev nasm
WORKDIR /usr/src

RUN curl -LO https://ffmpeg.org/releases/ffmpeg-6.1.tar.xz
RUN curl -LO https://ffmpeg.org/releases/ffmpeg-6.1.tar.xz.asc
COPY ffmpeg-devel.asc ./
RUN gpg --import ffmpeg-devel.asc
RUN gpg --verify ffmpeg-6.1.tar.xz.asc ffmpeg-6.1.tar.xz
RUN tar xaf ffmpeg-6.1.tar.xz
WORKDIR /usr/src/ffmpeg-6.1
ENV CFLAGS="-fanalyzer -Wall -Wextra -Wpedantic -Wformat=2 -Wformat-overflow=2 -Wformat-truncation=2 -Wformat-security -Wnull-dereference -Wstack-protector -Wtrampolines -Walloca -Wvla -Warray-bounds=2 -Wimplicit-fallthrough=3 -Wtraditional-conversion -Wshift-overflow=2 -Wcast-qual -Wstringop-overflow=4 -Wconversion -Warith-conversion -Wlogical-op -Wduplicated-cond -Wduplicated-branches -Wformat-signedness -Wshadow -Wstrict-overflow=4 -Wundef -Wstrict-prototypes -Wswitch-default -Wswitch-enum -Wstack-usage=1000000 -Wcast-align=strict -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=3 -D_GLIBCXX_ASSERTIONS -fcf-protection=full -fstack-protector-strong -fstack-clash-protection -fPIE -pie -fsanitize=bounds -fsanitize-undefined-trap-on-error -Wl,-z,nodlopen -Wl,-z,relro -Wl,-z,now -Wl,-z,noexecstack -Wl,-z,separate-code"
RUN  ./configure --disable-all --disable-everything --enable-libopenh264 --enable-demuxer=image_ppm_pipe --enable-parser=pnm --enable-decoder=ppm --enable-encoder=libopenh264 --enable-protocol='fd,file' --enable-muxer=mp4 --enable-filter='scale,format' --enable-ffmpeg --enable-swscale --enable-avcodec --enable-avutil --enable-avfilter --enable-avformat
RUN make -j $(nproc)
RUN make install

FROM ubuntu:jammy
RUN apt-get update && apt-get install -y --no-install-recommends tini libopenh264-cisco6 poppler-utils && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN install -d -o root -g root -m 0755 /sock
COPY --from=builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build-server /work/pdf2mp4 /pdf2mp4
COPY run.sh /run.sh
USER root
ENV LISTEN=/sock/listen.sock
ENV CACHE=/data
ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/pdf2mp4"]
STOPSIGNAL SIGINT
