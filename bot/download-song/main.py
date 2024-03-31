import yt_dlp
import json
import sys

print("Starting main.py script...")


def downloadSong(url:str):
    ydl_opts = {
        'format': 'bestaudio/best',
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'mp3',
            'preferredquality': '192',
        }],
        'audio_only': True,
        'ffmpeg_location': 'C:/ffmpeg/bin',
        'noplaylist': True,
        'max_filesize': 10000000,
        'outtmpl': './songs/currentSong',
    }
    
    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        try:
            ydl.download(url)
        except ydl.utils.DownloadError:
            return 'Video not found or inaccessible.'
    return None

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <URL to a song>")
        sys.exit(1)
    argument = sys.argv[1]
    downloadSong(argument)