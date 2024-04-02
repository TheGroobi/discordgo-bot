import yt_dlp
import json
import sys

def downloadSong(url:str):
    print("Starting song download...")

    ydl_opts = {
        'format': 'bestaudio/best',
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'opus',
            'preferredquality': '192',
        }],
        'audio_only': True,
        'ffmpeg_location': 'C:/ffmpeg/bin',
        'noplaylist': True,
        'max_filesize': 10000000,
        'outtmpl': './songs/currentSong.dca',
    }
    
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        try:
            ydl.download(url)
            sys.stderr.write("Normal Output")
        except Exception as e:
            sys.stderr.write("Exception: {}\n".format(str(e)))
            sys.exit(1)
    return None

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <URL to a song>")
        sys.exit(1)
    argument = sys.argv[1]
    downloadSong(argument)