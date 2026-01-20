# YT-Downloader

An application to download YouTube videos.

## Features
- Download YouTube videos (video/audio)
- Save files in a target directory
- Simple command-line interface

## Requirements
- Go 1.20+ (this repository is implemented in Go)
- Optional: ffmpeg (for audio extraction or format conversion)

## Build
1. Install Go: https://go.dev/dl/
2. From the repository root:
   - go build -o yt-downloader

## Usage
Run the built binary or use `go run`:

- Build and run:
  - ./yt-downloader --url "https://www.youtube.com/watch?v=VIDEO_ID" --output ./downloads

- Or with go run (from project root):
  - go run ./... --url "https://www.youtube.com/watch?v=VIDEO_ID" --output ./downloads

Common flags (adjust if your binary uses different flags):
- --url     The YouTube video URL (required)
- --output  Output directory (optional; default: current directory)
- --format  Desired output format: mp4 / mp3 / best (optional)

Check the source files for the exact flag names if the above differ.

## Contributing
Contributions are welcome. Please open an issue or submit a pull request with a clear description of your changes and any testing instructions.

## License
This repository does not include a license file. If you want to use a permissive license, consider adding an MIT license. Contact the repository owner if you want me to add a LICENSE file.

## Contact
Created by @python5vivek — open issues on this repo for questions or feature requests.