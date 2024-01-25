package encoder

import (
	"context"
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

func (e *Encoder) Encode(ctx context.Context, inFile, outDir string) error {

	basename := filepath.Base(inFile)
	ext := filepath.Ext(basename)
	basename = basename[:len(basename)-len(ext)]
	outFile := filepath.Join(outDir, basename+".m4a")
	if st, err := os.Stat(outFile); err == nil && st.Size() > 0 {
		slog.Info("Skipping", "file", inFile)
		return nil
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inFile, "-an", "-acodec", "copy", outFile)
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	return cmd.Run()

}

// encode audio
//         ffmpeg -i "$i" -vn -acodec copy "$ofname"

//encode video (stripping)
//     #     ffmpeg -i "$i" -an -vcodec copy "$ofvideo"

// encode video (DXV3)
