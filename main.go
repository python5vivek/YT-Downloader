package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kkdai/youtube/v2"
)

func sanitizeFileName(name string) string {

	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := re.ReplaceAllString(name, "_")
	return sanitized
}

func main() {
	a := app.New()
	w := a.NewWindow("YouTube Downloader")
	w.Resize(fyne.NewSize(500, 300))

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter YouTube URL")

	formatSelect := widget.NewRadioGroup([]string{"MP4", "Audio (M4A)"}, func(string) {})
	formatSelect.Horizontal = true
	formatSelect.SetSelected("Audio (M4A)")

	qualities := widget.NewSelect([]string{}, func(string) {})
	qualities.PlaceHolder = "Select quality"

	statusLabel := widget.NewLabel("Status: Idle")

	fetchBtn := widget.NewButton("Fetch Qualities", func() {
		url := urlEntry.Text
		if url == "" {
			statusLabel.SetText("Enter a URL first")
			return
		}

		go func() {
			statusLabel.SetText("Fetching video info...")
			client := youtube.Client{}
			video, err := client.GetVideo(url)
			if err != nil {
				statusLabel.SetText("Error: " + err.Error())
				return
			}

			var options []string
			if formatSelect.Selected == "MP4" {
				for _, f := range video.Formats.WithAudioChannels() {
					label := f.QualityLabel
					if label == "" {
						label = strings.Split(f.MimeType, ";")[0]
					}
					options = append(options, label)
				}
			} else {
				for _, f := range video.Formats.Type("audio") {
					label := strings.Split(f.MimeType, ";")[0]
					options = append(options, label)
				}
			}

			if len(options) == 0 {
				options = []string{"Default"}
			}

			qualities.Options = options
			qualities.SetSelected(options[0])
			statusLabel.SetText(fmt.Sprintf("Fetched %d options", len(options)))
		}()
	})

	downloadBtn := widget.NewButton("Download", func() {
		url := urlEntry.Text
		if url == "" {
			statusLabel.SetText("Enter a URL first")
			return
		}

		selectedQuality := qualities.Selected
		downloadFormat := formatSelect.Selected

		statusLabel.SetText("Downloading...")
		go func() {
			err := downloadYouTube(url, downloadFormat, selectedQuality)
			if err != nil {
				statusLabel.SetText("Error: " + err.Error())
			} else {
				statusLabel.SetText("Download complete!")
			}
		}()
	})

	content := container.NewVBox(
		urlEntry,
		formatSelect,
		fetchBtn,
		qualities,
		downloadBtn,
		statusLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func downloadYouTube(videoURL, format, quality string) error {
	client := youtube.Client{}
	video, err := client.GetVideo(videoURL)
	if err != nil {
		return err
	}

	var chosenFormat *youtube.Format

	if format == "MP4" {
		for _, f := range video.Formats.WithAudioChannels() {
			label := f.QualityLabel
			if label == "" {
				label = strings.Split(f.MimeType, ";")[0]
			}
			if label == quality || quality == "Default" {
				chosenFormat = &f
				break
			}
		}
		if chosenFormat == nil {
			chosenFormat = &video.Formats.WithAudioChannels()[0]
		}
	} else {
		for _, f := range video.Formats.Type("audio") {
			label := strings.Split(f.MimeType, ";")[0]
			if label == quality || quality == "Default" {
				chosenFormat = &f
				break
			}
		}
		if chosenFormat == nil {
			chosenFormat = &video.Formats.Type("audio")[0]
		}
	}

	stream, _, err := client.GetStream(video, chosenFormat)
	if err != nil {
		return err
	}

	var fileName string
	safeTitle := sanitizeFileName(video.Title)
	if format == "MP4" {
		fileName = safeTitle + ".mp4"
	} else {
		fileName = safeTitle + ".m4a"
	}

	file, err := os.Create(fileName)

	_, err = io.Copy(file, stream)
	file.Close()
	return err
}
