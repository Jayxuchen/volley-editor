package video

import (
	"fmt"
	"strings"

	"github.com/mowshon/moviego"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	labelFontSize = 96
	scoreFontSize = 224
	labelY        = 300
	scoreY        = 100
	homeX         = 100
	awayX         = 600
	bgBoxW        = 1000 // make wider
	bgBoxH        = 320  // tall enough for labels + scores
	bgBoxX        = 50   // left offset
	bgBoxBottom   = 85   // distance from bottom of screen
	boxColor      = "black@0.5"
	homeColor     = "deepskyblue"
	awayColor     = "orangered"
)

func buildTextFilters(overlays []ScoreOverlay) string {
	var filters []string
	for _, o := range overlays {
		backgroundBox := fmt.Sprintf(
			"drawbox=enable='between(t,%.2f,%.2f)':x=%d:y=ih-%d-%d:w=%d:h=%d:color=%s:t=fill",
			o.Start, o.End,
			bgBoxX, bgBoxH, bgBoxBottom, // position
			bgBoxW, bgBoxH,
			boxColor,
		)

		homeLabel := fmt.Sprintf(
			"drawtext=enable='between(t,%.2f,%.2f)':text='Home':x=%d:y=h-th-20-%d:fontsize=%d:fontcolor=%s:box=0",
			o.Start, o.End, homeX, labelY, labelFontSize, homeColor,
		)
		awayLabel := fmt.Sprintf(
			"drawtext=enable='between(t,%.2f,%.2f)':text='Away':x=%d:y=h-th-%d:fontsize=%d:fontcolor=%s:box=0",
			o.Start, o.End, awayX, labelY, labelFontSize, awayColor,
		)

		homeScore := fmt.Sprintf(
			"drawtext=enable='between(t,%.2f,%.2f)':text='%d':x=%d:y=h-th-%d:fontsize=%d:fontcolor=%s:box=0",
			o.Start, o.End, o.HomeScore, homeX, scoreY, scoreFontSize, homeColor,
		)
		awayScore := fmt.Sprintf(
			"drawtext=enable='between(t,%.2f,%.2f)':text='%d':x=%d:y=h-th-%d:fontsize=%d:fontcolor=%s:box=0",
			o.Start, o.End, o.AwayScore, awayX, scoreY, scoreFontSize, awayColor,
		)

		filters = append(filters, backgroundBox, homeLabel, awayLabel, homeScore, awayScore)
	}
	return strings.Join(filters, ",")
}

func drawTextWithShadow(start, end float64, text string, x, y, fontsize int, color string) (string, string) {
	shadow := fmt.Sprintf(
		"drawtext=enable='between(t,%.2f,%.2f)':text='%s':x=%d:y=h-th-%d:fontsize=%d:fontcolor=black@0.6:box=0",
		start, end, text, x+2, y+2, fontsize,
	)
	main := fmt.Sprintf(
		"drawtext=enable='between(t,%.2f,%.2f)':text='%s':x=%d:y=h-th-%d:fontsize=%d:fontcolor=%s:box=0",
		start, end, text, x, y, fontsize, color,
	)
	return shadow, main
}

func ApplyScoreOverlay(inputPath, outputPath string, overlays []ScoreOverlay) error {
	fmt.Println(inputPath)
	_, err := moviego.Load(inputPath) // still useful to validate and probe
	if err != nil {
		return err
	}

	textFilter := buildTextFilters(overlays)

	return ffmpeg.Input(inputPath).
		Output(outputPath,
			ffmpeg.KwArgs{
				"vf":      textFilter,
				"c:v":     "h264_videotoolbox", // h.264 + videotoolbox hardware acceleration
				"preset":  "slow",              // slower = better compression
				"b:v":     "20M",               // 20Mbps bitrate for 4k
				"c:a":     "aac",               //  set audio codec
				"pix_fmt": "yuv420p",           // good compatibility
			}).
		OverWriteOutput().
		Run()
}
