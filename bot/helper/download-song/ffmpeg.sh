#!/bin/sh
chmod +x script

#download song with python/yt-dlp
python bot/helper/download-song/main.py $1
sleep 0.1

# Convert mp3 to PCM
ffmpeg -i bot/helper/download-song/songs/current.mp3 -y -f s16le -ar 48000 -ac 2 bot/helper/download-song/songs/output
if [ $? -ne 0 ]; then
    echo "ffmpeg failed"
    exit 1
fi

#Encode PCM to opus
opusenc --bitrate 128 --raw --downmix-stereo bot/helper/download-song/songs/output bot/helper/download-song/songs/output.opus
if [ $? -ne 0 ]; then
    echo "opus-tools failed"
    exit 1
fi

rm bot/helper/download-song/songs/output
rm bot/helper/download-song/songs/current.mp3

# for testing return mp3 to check if encoded correctly
# ffmpeg -i bot/helper/download-song/songs/output.opus -y bot/helper/download-song/songs/output.mp3