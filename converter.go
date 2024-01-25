package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bmurray/resolumeconverter/encoder"
	"github.com/bmurray/resolumeconverter/resolume"
	"github.com/google/uuid"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	baseUrlString := flag.String("base-url", "http://127.0.0.1:8089/api/v1/", "Base URL of Resolume")
	flag.Parse()

	baseUrl, err := url.Parse(*baseUrlString)
	if err != nil {
		slog.Error("Error parsing base URL", "error", err)
	}
	r := resolume.NewResolume(baseUrl)

	args := flag.Args()
	if len(args) == 0 {
		slog.Error("No command specified")
		os.Exit(1)
	}

	switch args[0] {
	case "layers":
		layers(ctx, r, args[1:])
	case "composition":
		composition(ctx, r, args[1:])
	case "convert":
		convert(ctx, r, args[1:])
	default:
		slog.Error("Unknown command", "command", args[0])
	}

}
func layers(ctx context.Context, res *resolume.Resolume, args []string) {

	if len(args) == 0 {
		listLayers(ctx, res)
		return
	}
	switch args[0] {
	case "list":
		listLayers(ctx, res)
	case "get":
		getLayers(ctx, res)
	default:
		slog.Error("Unknown command", "command", fmt.Sprintf("layer %s", args[0]))
	}
}

func listLayers(ctx context.Context, r *resolume.Resolume) {
	layers, err := r.GetLayers(ctx)
	if err != nil {
		slog.Error("Error getting layers", "error", err)
	}
	for idx, layer := range layers {
		fmt.Printf("%d: %s (%d)\n", idx, layer.Name.Value, layer.Id)
	}
}

func getLayers(ctx context.Context, r *resolume.Resolume) {
	layers, err := r.GetLayers(ctx)
	if err != nil {
		slog.Error("Error getting layers", "error", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(layers)
}

func composition(ctx context.Context, res *resolume.Resolume, args []string) {

	if len(args) == 0 {
		getComposition(ctx, res)
		return
	}
	switch args[0] {
	case "get":
		getComposition(ctx, res)
	default:
		slog.Error("Unknown command", "command", fmt.Sprintf("composition %s", args[0]))
	}
}

func getComposition(ctx context.Context, r *resolume.Resolume) {
	composition, err := r.GetComposition(ctx)
	if err != nil {
		slog.Error("Error getting composition", "error", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(composition)
}

func convert(ctx context.Context, r *resolume.Resolume, args []string) {
	if len(args) < 3 {
		slog.Error("No file specified")
		return
	}
	inDir := args[0]
	audioOutDir := args[1]
	videoOutDir := args[2]
	err := convertAudioFiles(ctx, inDir, audioOutDir)
	if err != nil {
		slog.Error("Error converting files", "error", err)
		return
	}

	matched, err := matchAudioVideo(ctx, audioOutDir, videoOutDir)
	if err != nil {
		slog.Error("Error matching audio and video files", "error", err)
		return
	}

	for _, m := range matched {
		slog.Info("Matched", "file", m)
	}

}

func matchAudioVideo(ctx context.Context, audioOutDir, videoOutDir string) ([]string, error) {

	audiofiles, err := filepath.Glob(filepath.Join(audioOutDir, "*.m4a"))
	if err != nil {
		slog.Error("Error globbing audio files", "error", err)
		return nil, err
	}

	videoFiles, err := filepath.Glob(filepath.Join(videoOutDir, "*.mp4"))
	if err != nil {
		slog.Error("Error globbing video files", "error", err)
		return nil, err
	}

	files := make(map[string]bool)
	for _, af := range audiofiles {
		afbase := filepath.Base(af)
		afbase = afbase[:len(afbase)-len(".m4a")]
		files[afbase] = false
	}
	for _, vf := range videoFiles {
		vfbase := filepath.Base(vf)
		vfbase = vfbase[:len(vfbase)-len(".mp4")]
		if _, ok := files[vfbase]; ok {
			files[vfbase] = true
		}
	}
	matches := make([]string, 0, len(files))
	for k, v := range files {
		if !v {
			slog.Warn("No video match for audio file", "audio", k)
			continue
		}
		matches = append(matches, k)
	}

	return matches, nil
}

func convertAudioFiles(ctx context.Context, inDir, outDir string) error {

	runId, err := uuid.NewRandom()
	if err != nil {
		slog.Error("Error generating run ID", "error", err)
		return err
	}
	logfile := filepath.Join("logs", runId.String(), "resolumeconverter.log")
	err = os.MkdirAll(filepath.Dir(logfile), 0755)
	if err != nil {
		slog.Error("Error creating log directory", "error", err)
		return err
	}
	flog, err := os.Create(logfile)
	if err != nil {
		slog.Error("Error creating log file", "error", err)
		return err
	}
	defer flog.Close()

	enc := encoder.NewEncoder(
		encoder.WithStdout(flog),
		encoder.WithStderr(flog),
	)
	_ = enc
	_ = outDir

	files, err := os.ReadDir(inDir)
	if err != nil {
		slog.Error("Error reading directory", "error", err)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[0] == '.' {
			continue
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		fname := filepath.Join(inDir, file.Name())

		slog.Info("Converting", "file", fname)
		err := enc.Encode(ctx, fname, outDir)
		if err != nil {
			slog.Error("Error converting", "error", err)
			return err
		}
	}

	return nil
}
