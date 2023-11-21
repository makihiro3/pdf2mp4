#!/bin/bash
set -o pipefail
#ulimit -b 0         #socket buffer size
ulimit -c 0         #core file size                     (blocks, -c)
ulimit -d 524288    #data seg size                      (kbytes, -d)
ulimit -e 0         #scheduling priority                        (-e)
ulimit -f 131072    #file size                          (blocks, -f)
ulimit -i 8         #pending signals                            (-i)
#ulimit -k 0         #kqueues
ulimit -l 0         #max locked memory                  (kbytes, -l)
ulimit -m 524288    #max memory size                    (kbytes, -m)
ulimit -n 32        #open files                                 (-n)
#ulimit -p 8         #pipe size                       (512 bytes, -p)
ulimit -q 0         #POSIX message queues                (bytes, -q)
ulimit -r 0         #real-time non-blocking time  (microseconds, -R)
ulimit -s 8192      #stack size                         (kbytes, -s)
ulimit -t 300       #cpu time                          (seconds, -t)
ulimit -u 64        #max user processes                         (-u)
#ulimit -v 524288    #virtual memory                     (kbytes, -v)
ulimit -x 4         #file locks                                 (-x)
#ulimit -P 0         #pseudoterminals
ulimit -R 0         #real-time priority
#ulimit -T 64        #thread

DIR="$1"
SIZE="$2"
INTERVAL="$3"
PDFTOPPM_OPTS=()
if [[ $SIZE != 0 ]]; then
    PDFTOPPM_OPTS+=( -scale-to-x "-1" -scale-to-y "$SIZE" )
fi
RATE="1/2"
case "$INTERVAL" in
1) RATE=1 ;;
2) RATE="1/2" ;;
3) RATE="1/3" ;;
esac

set -xe
pdftoppm "${PDFTOPPM_OPTS[@]}" "$DIR/input.pdf" | ffmpeg -r "$RATE" -f ppm_pipe -i - -c:v libopenh264 -profile:v high -allow_skip_frames 1 -coder cabac -qmin 1 -qmax 50 -r 30 -y -f mp4 "$DIR/output.mp4"
