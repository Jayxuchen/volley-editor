package video

import (
	"fmt"
	"strings"

	"github.com/mowshon/moviego"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	fontSize   = 197
	boxBorderW = 20
	xPos       = 10
	yOffset    = 10
	fontColor  = "white"
	boxColor   = "black@0.5"
)

func buildTextFilters(overlays []ScoreOverlay) string {

	var filters []string
	for _, o := range overlays {
		filter := fmt.Sprintf(
			"drawtext=enable='between(t,%.2f,%.2f)':text='%s':x=%d:y=h-th-%d:fontsize=%d:fontcolor=%s:box=1:boxcolor=%s:boxborderw=%d",
			o.Start, o.End,
			strings.ReplaceAll(o.Text, ":", "\\:"),
			xPos, yOffset,
			fontSize, fontColor, boxColor, boxBorderW,
		)
		filters = append(filters, filter)
	}
	return strings.Join(filters, ",")
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
				"c:v":     "libx265", // keep HEVC encoding
				"crf":     "18",      // default 23
				"preset":  "slow",    // slower = better compression
				"pix_fmt": "yuv420p", // good compatibility
			}).
		OverWriteOutput().
		Run()
}
