package encoder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"log/slog"
)

type Encoder struct {
	stdout io.Writer
	stderr io.Writer
}

type EncoderOption func(*Encoder)

func WithStdout(w io.Writer) EncoderOption {
	return func(e *Encoder) {
		e.stdout = w
	}
}

func WithStderr(w io.Writer) EncoderOption {
	return func(e *Encoder) {
		e.stderr = w
	}
}

func NewEncoder(opts ...EncoderOption) *Encoder {
	e := &Encoder{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e Encoder) Encode(ctx context.Context, inFile, outDir string) error {

	basename := filepath.Base(inFile)
	ext := filepath.Ext(basename)
	basename = basename[:len(basename)-len(ext)]
	outFile := filepath.Join(outDir, basename+".m4a")
	if st, err := os.Stat(outFile); err == nil && st.Size() > 0 {
		slog.Info("Skipping", "file", inFile)
		return nil
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inFile, "-vn", "-acodec", "copy", outFile)
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	return cmd.Run()
}

func (e Encoder) GetThumbnail(ctx context.Context, inFile string) (io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inFile, "-s", "320x240", "-vframes", "1", "-c:v", "png", "-f", "image2pipe", "-")
	cmd.Stderr = e.stderr
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return pipe, nil

}

func (e Encoder) GetAudioTitle(ctx context.Context, inFile string) (string, error) {
	data, err := e.GetMetadata(ctx, inFile)
	if err != nil {
		return "", err
	}
	if data.Format.Tags == nil {
		return "", fmt.Errorf("no tags")
	}
	return data.Format.Tags["title"], nil
}

type ffmetadata struct {
	Streams []ffstream `json:"streams"`
	Format  ffformat   `json:"format"`
}

type ffstream struct {
}
type ffformat struct {
	Filename string            `json:"filename"`
	Tags     map[string]string `json:"tags"`
}

func (e Encoder) GetMetadata(ctx context.Context, inFile string) (ffmetadata, error) {

	// ffprobe -show_format -show_streams -output_format json -i input.mp4

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ffprobe", "-show_format", "-show_streams", "-output_format", "json", "-i", inFile)
	cmd.Stderr = e.stderr
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return ffmetadata{}, err
	}
	data := ffmetadata{}
	dec := json.NewDecoder(pipe)
	err = cmd.Start()
	if err != nil {
		return ffmetadata{}, err
	}

	err = dec.Decode(&data)
	if err != nil {
		return ffmetadata{}, err
	}
	return data, nil
}

// encode audio
//         ffmpeg -i "$i" -vn -acodec copy "$ofname"

//encode video (stripping)
//     #     ffmpeg -i "$i" -an -vcodec copy "$ofvideo"

// encode video (DXV3)
