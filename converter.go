package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bmurray/resolumeconverter/encoder"
	"github.com/bmurray/resolumeconverter/resolume"
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
	case "clips":
		clips(ctx, r, args[1:])
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

func clips(ctx context.Context, res *resolume.Resolume, args []string) {

	if len(args) == 0 {
		// listClips(ctx, res)
		return
	}

	switch args[0] {
	case "list":
		// listClips(ctx, res)
	case "get":
		// getClips(ctx, res)
	case "thumbnail":
		getThumbnail(ctx, res, args[1:])
	case "selected":
		getSelectedClip(ctx, res)
	default:
		slog.Error("Unknown command", "command", fmt.Sprintf("clip %s", args[0]))
	}
}

func getThumbnail(ctx context.Context, r *resolume.Resolume, args []string) {
	if len(args) == 0 {
		slog.Error("No clip ID specified")
		return
	}
	clipId, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error("Error parsing clip ID", "error", err)
		return
	}

	thumbnail, err := r.GetThumbnail(ctx, clipId)
	if err != nil {
		slog.Error("Error getting thumbnail", "error", err)
		return
	}
	defer thumbnail.Close()

	fname := fmt.Sprintf("%s.png", clipId)
	f, err := os.Create(fname)
	if err != nil {
		slog.Error("Error creating file", "error", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, thumbnail)
	if err != nil {
		slog.Error("Error copying thumbnail", "error", err)
		return
	}

}

func getSelectedClip(ctx context.Context, r *resolume.Resolume) {
	clips, err := r.GetSelectedClip(ctx)
	if err != nil {
		slog.Error("Error getting clips", "error", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(clips)
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

	if len(args) == 0 {
		slog.Error("No command specified")
		os.Exit(1)
	}

	switch args[0] {
	case "input":
		if len(args) < 2 {
			slog.Error("No input dir specified specified")
			return
		}
		convertInputs(ctx, args[1])
	case "audio":
		// Only do audio conversion
		if len(args) < 3 {
			slog.Error("No input dir specified specified")
			return
		}
		inDir := args[1]
		audioOutDir := args[2]
		convertAudioFiles(ctx, encoder.NewEncoder(), inDir, audioOutDir)
	case "input-audio":
		// Convert input and audio
		if len(args) < 3 {
			slog.Error("No input dir specified specified")
			return
		}
		inDir := args[1]
		audioOutDir := args[2]

		convertInputs(ctx, inDir)
		convertAudioFiles(ctx, encoder.NewEncoder(), inDir, audioOutDir)

	case "import":
		convertImport(ctx, r, args[1:])
	default:
		slog.Error("Unknown command", "command", args[0])
	}
}
func convertImport(ctx context.Context, r *resolume.Resolume, args []string) {

	if len(args) < 2 {
		slog.Error("No input dir specified specified")
		return
	}
	indir := args[0]
	layer, err := strconv.Atoi(args[1])
	if err != nil {
		slog.Error("Error parsing layer", "error", err)
		return
	}

	files, err := filepath.Glob(filepath.Join(indir, "*.mov"))
	if err != nil {
		slog.Error("Error globbing files", "error", err)
		return
	}

	for _, file := range files {
		err := convertAddToResolume(ctx, r, file, layer)
		if err != nil {
			slog.Error("Error converting file", "error", err)
			return
		}
	}
}

func convertInputs(ctx context.Context, inDir string) {
	// if len(args) < 1 {
	// 	slog.Error("No input dir specified specified")
	// 	return
	// }
	// inDir := args[0]
	files, err := filepath.Glob(filepath.Join(inDir, "*.mp4"))
	if err != nil {
		slog.Error("Error globbing files", "error", err)
		return
	}
	enc := encoder.NewEncoder()
	for _, file := range files {
		slog.Info("Converting", "file", file)

		audioTitle, err := getAudioTitle(ctx, enc, file)
		if err != nil {
			slog.Error("Error getting audio title", "error", err)
			return
		}
		slog.Info("Audio title", "title", audioTitle)
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		base = base[:len(base)-len(ext)]
		outFile := filepath.Join(inDir, audioTitle+".mp4")
		if st, err := os.Stat(outFile); err == nil && st.Size() > 0 {
			slog.Info("Skipping", "file", file)
			continue
		}
		err = os.Rename(file, outFile)
		if err != nil {
			slog.Error("Error renaming file", "error", err)
			return
		}
	}
}

// func correctAudioVideos(ctx context.Context, matched []string, audioOutDir, videoOutDir string) ([]string, error) {
// 	correct := make([]string, 0, len(matched))
// 	for _, m := range matched {

//			newFile, err := correctAudioVideo(ctx, audioOutDir, videoOutDir, m)
//			if err != nil {
//				slog.Error("Error correcting audio video", "error", err)
//				return nil, err
//			}
//			correct = append(correct, newFile)
//		}
//		return correct, nil
//	}

func convertAddToResolume(ctx context.Context, r *resolume.Resolume, file string, layer int) error {

	exists, err := clipExists(ctx, r, file)
	if err != nil {
		slog.Error("Error checking if clip exists", "error", err)
		return err
	}
	if exists {
		slog.Info("Clip already exists", "clip", file)
		return nil
	}

	_, clip, err := r.FindEmptyClip(ctx, layer, layer)
	if err != nil {
		slog.Error("Error finding empty clip", "error", err)
		return err
	}
	err = r.OpenClip(ctx, clip.Id, file)
	if err != nil {
		slog.Error("Error opening clip", "error", err)
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}
	val := make(map[string]any)
	// val["name"] = struct {
	// 	Valuetype string `json:"valuetype"`
	// 	Id        int    `json:"id"`
	// 	Value     string `json:"value"`
	// }{
	// 	Valuetype: "ParamString",
	// 	Id:        clip.Name.Id,
	// 	Value:     audioTitle,
	// }
	val["target"] = struct {
		// Id        any    `json:"id"`
		ValueType string `json:"valuetype"`
		Value     string `json:"value"`
		Index     int    `json:"index"`
		// Options   any    `json:"options"`
	}{
		// Id:        clip.Target["id"],
		ValueType: "ParamChoice",
		Value:     "Denon Player Determined",
		Index:     4,
		// Options:   clip.Target["options"],
	}
	val["transporttype"] = struct {
		// 	Id        int    `json:"id"`
		ValueType string `json:"valuetype"`
		Value     string `json:"value"`
		Index     int    `json:"index"`
	}{
		// 	Id:        rand.Intn(1000000),
		ValueType: "ParamChoice",
		Value:     "Denon DJ",
		Index:     4,
	}
	err = r.SetClipRaw(ctx, clip.Id, val)
	if err != nil {
		slog.Error("Error setting clip", "error", err)
		return err
	}
	return nil

}

func clipExists(ctx context.Context, r *resolume.Resolume, videoFile string) (bool, error) {
	comp, err := r.GetComposition(ctx)
	if err != nil {
		slog.Error("Error getting composition", "error", err)
		return false, err
	}

	for _, layer := range comp.Layers {
		for _, clip := range layer.Clips {
			if clip.Video.FileInfo.Path == videoFile {
				return true, nil
			}
		}
	}
	return false, nil
}

func getAudioTitle(ctx context.Context, enc *encoder.Encoder, audioFile string) (string, error) {

	ff, err := enc.GetAudioTitle(ctx, audioFile)
	if err != nil {
		slog.Error("Error getting audio title", "error", err)
		return "", err
	}

	return ff, nil
}

func convertAudioFiles(ctx context.Context, enc *encoder.Encoder, inDir, outDir string) error {

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
		basename := filepath.Base(fname)
		ext := filepath.Ext(basename)
		basename = basename[:len(basename)-len(ext)]
		outFile := filepath.Join(outDir, basename+".m4a")

		err := enc.Encode(ctx, fname, outFile)
		if err != nil {
			slog.Error("Error converting", "error", err)
			return err
		}
	}

	return nil
}
