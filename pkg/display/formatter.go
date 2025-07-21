// Pretty much this whole file is a copy of Ashish Kumar's Mufetch formatter.go file
package display

import (
	"fmt"
	"image"
	"github.com/yubaix/mofetch/pkg/omdb"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

// ANSI color codes for terminal output
const (
	ColorReset   = "\033[0m"
	ColorRedBold = "\033[1;31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorPurple  = "\033[35m"
	ColorWhite   = "\033[37m"
	ColorCyan    = "\033[36m"
	ColorBold    = "\033[1m"
)

// BOLD, RED
// yellow
// green
// blue
// purple

// ImageRenderer handles terminal image rendering using Unicode blocks
type ImageRenderer struct {
	width  int
	height int
}

// NewImageRenderer creates an image renderer with specified size
func NewImageRenderer(size int) *ImageRenderer {
	return &ImageRenderer{
		width:  size,
		height: size,
	}
}

// RenderImageLines converts image URL to terminal-displayable lines
func (r *ImageRenderer) RenderImageLines(imageURL string) []string {
	if imageURL == "" {
		return r.getPlaceholderLines()
	}

	img, err := r.downloadImage(imageURL)
	if err != nil {
		return r.getPlaceholderLines()
	}

	return r.getBlockArtLines(img)
}

// downloadImage fetches and decodes image from URL
func (r *ImageRenderer) downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	return img, err
}

func (r *ImageRenderer) getBlockArtLines(img image.Image) []string {
	// 3:2 aspect ratio for movie posters
	const aspectWidth = 2
	const aspectHeight = 3

	dynamicHeight := (r.width * aspectHeight) / aspectWidth

	resized := imaging.Resize(img, r.width, dynamicHeight, imaging.Lanczos)
	bounds := resized.Bounds()

	var lines []string

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var line strings.Builder
		line.WriteString(" ") // Left padding

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := resized.At(x, y)
			r8, g8, b8, _ := pixelColor.RGBA()
			r := uint8(r8 >> 8)
			g := uint8(g8 >> 8)
			b := uint8(b8 >> 8)

			// ANSI true color block
			line.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm  \033[0m", r, g, b))
		}
		lines = append(lines, line.String())
	}

	return lines
}

// getPlaceholderLines creates a placeholder when no image is available
func (r *ImageRenderer) getPlaceholderLines() []string {
	lines := []string{
		fmt.Sprintf("%sNo Poster Available%s", ColorWhite, ColorReset),
	}

	return lines
}

// Runtime conversion stuff
func getFormattedRuntime(runtime string) string {
	parts := strings.Fields(runtime)
	minStr := parts[0]
	minInt, err := strconv.Atoi(minStr)
	if err != nil {
		panic(err)
	}
	dur := time.Duration(minInt) * time.Minute
	hours := int(dur.Hours())
	mins := int(dur.Minutes()) % 60
	fRuntime := fmt.Sprintf("%dh %02dm", hours, mins)
	return fRuntime
}

// Plot truncation (cuts into the poster otherwise)
func truncatePlot(text string) string {
	words := strings.Fields(text)
	if len(words) <= 10 {
		return text
	}
	return strings.Join(words[:10], " ") + "..."
}

// DisplayFilm renders film information and poster in terminal
func DisplayFilm(film omdb.Film, apiKey omdb.Api, imageSize int) {
	renderer := NewImageRenderer(imageSize)

	imageLines := renderer.RenderImageLines(film.Poster)

	infoLines := []string{
		formatInfoLine(film.Title, ColorRedBold),
		" ",
		formatInfoLine("- "+film.Year, ColorYellow),
		formatInfoLine("- "+film.Director, ColorGreen),
		formatInfoLine("- "+getFormattedRuntime(film.Runtime), ColorBlue),
		formatInfoLine("- "+film.ImdbRating+"/10", ColorPurple),
	}

	displaySidebySide(imageLines, infoLines)
}

// DisplayFilmVerbose renders film information and poster in terminal with additional details
func DisplayFilmVerbose(film omdb.Film, apiKey omdb.Api, imageSize int) {
	renderer := NewImageRenderer(imageSize)

	imageLines := renderer.RenderImageLines(film.Poster)
	infoLines := []string{
		formatInfoLine(film.Title, ColorRedBold),
		" ",
		formatInfoLineVerbose("Plot", " - "+truncatePlot(film.Plot), ColorYellow),
		formatInfoLineVerbose("Actors", " - "+film.Actors, ColorGreen),
		formatInfoLineVerbose("Director", " - "+film.Director, ColorBlue),
		formatInfoLineVerbose("Year", " - "+film.Year, ColorCyan),
		formatInfoLineVerbose("Runtime", " - "+getFormattedRuntime(film.Runtime), ColorPurple),
		formatInfoLineVerbose("Language(s)", " - "+film.Language, ColorCyan),
		formatInfoLineVerbose("Box Office", " - "+film.Boxoffice, ColorBlue),
		formatInfoLineVerbose("IMDb Rating", " - "+film.ImdbRating+"/10", ColorGreen),
		formatInfoLineVerbose("MPAA Rating", " - "+film.Rating, ColorYellow),
	}

	displaySidebySide(imageLines, infoLines)
}

// displaySideBySideWithLinks renders image and info side-by-side
func displaySidebySide(imageLines, infoLines []string) {
	maxLines := len(imageLines)

	for i := 0; i < maxLines; i++ {
		var img, info string

		if i < len(imageLines) {
			img = imageLines[i]
		} else {
			img = strings.Repeat(" ", 46) // blank space to align image column
		}

		if i < len(infoLines) {
			info = infoLines[i]
		} else {
			info = ""
		}

		// Left-align image column to 46 chars, then print info
		fmt.Printf("%-46s   %s\n", img, info)
	}
}

// formatInfoLine creates consistently formatted text
func formatInfoLine(value, color string) string {
	return fmt.Sprintf("%s%s%s", color, value, ColorReset)
}

// Sanme as the above, but for the verbose version
// Adds padding for the text, and adds labels
func formatInfoLineVerbose(label, value, color string) string {
	const minPadding = 2
	const maxLabelWidth = 12

	labelWidth := len(label)
	padding := maxLabelWidth - labelWidth + minPadding
	if padding < minPadding {
		padding = minPadding
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s",
		ColorBold, label, ColorReset,
		strings.Repeat(" ", padding),
		color, value, ColorReset)
}
