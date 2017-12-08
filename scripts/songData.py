import sys
import os
from tinytag import TinyTag


class SongFile:
    def __init__(self, title, artist, source):
        if title is None:
            raise Exception("title is None")
        if source is None:
            raise Exception("source is None")

        self.title = title
        self.artist = artist or "Unknown"
        self.source = source

    def format(self, format_func):
        format_func(self.title, self.artist, self.source)


def write_csv(out, *args):
    # change double quotes to single quotes in fields
    # surround each field in double quotes
    # separate each field with commas
    print(','.join(map(lambda a: "\"{}\"".format(a.replace('"', '\'')), *args)), file=out)


# pulled from tinytag source (https://github.com/devsnd/tinytag)
fileExtensions = [".mp3", ".oga", ".ogg", ".opus", ".wav", ".flac", ".wma", ".m4a", ".mp4"]


def main(argv):
    search_path = argv[1]
    output_path = argv[2]
    songs = []

    for root, subDirs, files in os.walk(top=search_path, followlinks=True):
        for fileName in files:
            ext = os.path.splitext(fileName)[1]

            if ext not in fileExtensions:
                continue

            path = os.path.join(root, fileName)
            tag = TinyTag.get(path, duration=False)

            # make the path in the csv file relative to the expected file structure in a docker container
            # and make path separators Unix-style
            songs.append(
                SongFile(
                    tag.title,
                    tag.artist,
                    path.replace(search_path, "/songs/").replace('\\', '/')))

    with open(output_path, mode="w", encoding="utf_8") as f:
        for song in songs:
            song.format(lambda *args: write_csv(f, args))

    print("Created '{}' with {} songs.".format(output_path, len(songs)))


if __name__ == "__main__":
    main(sys.argv)
