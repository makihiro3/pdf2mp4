FROM golang:1.24.5 AS build-server

WORKDIR /work
COPY ./ ./
RUN make

FROM buildpack-deps:trixie AS builder
RUN apt-get update
RUN apt-get install -y --no-install-recommends libopenh264-dev nasm
WORKDIR /usr/src

ARG VERSION=8.0
RUN curl -LO https://ffmpeg.org/releases/ffmpeg-$VERSION.tar.xz
RUN curl -LO https://ffmpeg.org/releases/ffmpeg-$VERSION.tar.xz.asc
COPY ffmpeg-devel.asc ./
RUN gpg --import ffmpeg-devel.asc
RUN gpg --verify ffmpeg-$VERSION.tar.xz.asc ffmpeg-$VERSION.tar.xz
RUN tar xaf ffmpeg-$VERSION.tar.xz
WORKDIR /usr/src/ffmpeg-$VERSION
ENV CFLAGS="-Wall -Wformat -Wformat=2 -Wconversion -Wimplicit-fallthrough -Werror=format-security -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=3 -D_GLIBCXX_ASSERTIONS -fstrict-flex-arrays=3 -fstack-clash-protection -fstack-protector-strong -Wl,-z,nodlopen -Wl,-z,noexecstack -Wl,-z,relro -Wl,-z,now -Wl,--as-needed -Wl,--no-copy-dt-needed-entries -Wtrampolines -Wbidi-chars=any -fPIE -pie -fcf-protection=full -fno-delete-null-pointer-checks -fno-strict-overflow -fno-strict-aliasing -ftrivial-auto-var-init=zero -Werror=implicit -Werror=incompatible-pointer-types -Werror=int-conversion -fexceptions"
RUN  ./configure --disable-all --disable-everything --enable-libopenh264 --enable-demuxer=image_ppm_pipe --enable-parser=pnm --enable-decoder=ppm --enable-encoder=libopenh264 --enable-protocol='fd,file' --enable-muxer=mp4 --enable-filter='scale,format' --enable-ffmpeg --enable-swscale --enable-avcodec --enable-avutil --enable-avfilter --enable-avformat
RUN make -j $(nproc)
RUN make install

FROM debian:trixie
RUN  sed -i -e 's/ main/ main contrib/g' /etc/apt/sources.list.d/debian.sources && apt-get update && apt-get install -y --no-install-recommends tini libopenh264-cisco8 poppler-utils poppler-data fonts-noto-core fonts-noto-cjk fonts-noto-cjk-extra && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN install -d -o root -g root -m 0755 /sock
COPY --from=builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build-server /work/backend /pdf2mp4
COPY run.sh /run.sh
USER root
ENV LISTEN=/sock/listen.sock
ENV CACHE=/data
CMD ["/usr/bin/tini", "--", "/pdf2mp4"]
STOPSIGNAL SIGINT
