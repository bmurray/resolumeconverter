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

	"math/rand"

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
	if len(args) < 5 {
		slog.Error("No file specified")
		return
	}
	inDir := args[0]
	audioOutDir := args[1]
	videoOutDir := args[2]

	startLayer, err := strconv.Atoi(args[3])
	if err != nil {
		slog.Error("Error parsing start layer", "error", err)
		return
	}

	if startLayer < 0 {
		startLayer = 0
	}
	endLayer, err := strconv.Atoi(args[4])

	if err != nil {
		slog.Error("Error parsing end layer", "error", err)
		return
	}
	if endLayer < 0 {
		endLayer = 0
	}

	runId, err := uuid.NewRandom()
	if err != nil {
		slog.Error("Error generating run ID", "error", err)
		return
	}
	logfile := filepath.Join("logs", runId.String(), "resolumeconverter.log")
	err = os.MkdirAll(filepath.Dir(logfile), 0755)
	if err != nil {
		slog.Error("Error creating log directory", "error", err)
		return
	}
	flog, err := os.Create(logfile)
	if err != nil {
		slog.Error("Error creating log file", "error", err)
		return
	}
	defer flog.Close()

	enc := encoder.NewEncoder(
		encoder.WithStdout(flog),
		encoder.WithStderr(flog),
	)

	err = convertAudioFiles(ctx, enc, inDir, audioOutDir)
	if err != nil {
		slog.Error("Error converting files", "error", err)
		return
	}

	matched, err := matchAudioVideo(ctx, audioOutDir, videoOutDir)
	if err != nil {
		slog.Error("Error matching audio and video files", "error", err)
		return
	}

	template, err := r.GetSelectedClip(ctx)
	if err != nil {
		slog.Error("Error getting selected clip", "error", err)
		return
	}

	template = wipeIdenfiers(template)

	for _, m := range matched {
		// slog.Info("Matched", "file", m)
		if err := addToResolume(ctx, r, enc, template, audioOutDir, videoOutDir, m, startLayer, endLayer); err != nil {
			slog.Error("Error adding to resolume", "error", err)
			return
		}

		slog.Info("Added", "file", m)
		return
	}
}

func addToResolume(ctx context.Context, r *resolume.Resolume, enc *encoder.Encoder, template resolume.Clip, audioOutDir, videoOutDir, name string, startLayer, endLayer int) error {

	audioFile := filepath.Join(audioOutDir, name+".m4a")
	videoFile := filepath.Join(videoOutDir, name+".mov")
	_ = audioFile
	_ = videoFile

	exists, err := clipExists(ctx, r, videoFile)
	if err != nil {
		slog.Error("Error checking if clip exists", "error", err)
		return err
	}
	if exists {
		slog.Info("Clip already exists", "clip", name)
		return nil
	}

	layer, clip, err := r.FindEmptyClip(ctx, startLayer, endLayer)
	if err != nil {
		slog.Error("Error finding empty clip", "error", err)
		return err
	}
	slog.Info("Found empty clip", "layer", layer, "clip", clip)

	audioTitle, err := getAudioTitle(ctx, enc, audioFile)
	if err != nil {
		slog.Error("Error getting audio title", "error", err)
		return err
	}
	slog.Info("Audio title", "title", audioTitle)

	template.Name.Value = name
	template.Video.Description = name
	template.Video.FileInfo.Path = videoFile
	template.Video.FileInfo.Duration = ""
	template.Video.FileInfo.DurationMS = 0
	template.Video.FileInfo.Framerate = nil
	template.Video.FileInfo.Width = 0
	template.Video.FileInfo.Height = 0

	err = r.OpenClip(ctx, clip, videoFile)
	if err != nil {
		slog.Error("Error opening clip", "error", err)
		return err
	}

	val := make(map[string]any)
	val["id"] = clip
	val["target"] = struct {
		Id        int    `json:"id"`
		ValueType string `json:"valuetype"`
		Value     string `json:"value"`
		Index     int    `json:"index"`
	}{
		Id:        rand.Intn(1000000),
		ValueType: "ParamChoice",
		Value:     "Denon Player Determined",
		Index:     4,
	}
	val["transporttype"] = struct {
		Id        int    `json:"id"`
		ValueType string `json:"valuetype"`
		Value     string `json:"value"`
		Index     int    `json:"index"`
	}{
		Id:        rand.Intn(1000000),
		ValueType: "ParamChoice",
		Value:     "Denon DJ",
		Index:     4,
	}

	err = r.SetClipRaw(ctx, clip, val)
	if err != nil {
		slog.Error("Error setting clip", "error", err)
		return err
	}
	// err = r.SetClip(ctx, clip, template)
	// if err != nil {
	// 	slog.Error("Error setting clip", "error", err)
	// 	return err
	// }

	// thumbnail, err := enc.GetThumbnail(ctx, videoFile)
	// if err != nil {
	// 	slog.Error("Error getting thumbnail", "error", err)
	// 	return err
	// }
	// defer thumbnail.Close()
	// of, err := os.Create("test.png")
	// if err != nil {
	// 	slog.Error("Error creating thumbnail file", "error", err)
	// 	return err
	// }
	// defer of.Close()
	// _, err = io.Copy(of, thumbnail)
	// if err != nil {
	// 	slog.Error("Error copying thumbnail", "error", err)
	// 	return err
	// }

	// clip, err := r.AddClip(ctx, name, audioFile, videoFile)
	// if err != nil {
	// 	slog.Error("Error adding clip", "error", err)
	// 	return err
	// }
	// slog.Info("Added clip", "clip", clip.Id)

	return nil
}
func wipeIdenfiers(template resolume.Clip) resolume.Clip {
	template.Id = 0
	template.Name.Id = 0
	template.Audio = wipeTodo(template.Audio)
	template.BeatSnap = wipeTodo(template.BeatSnap)
	template.Dashboard = wipeTodo(template.Dashboard)
	template.FaderStart = wipeTodo(template.FaderStart)
	template.IgnoreColumnTrigger = wipeTodo(template.IgnoreColumnTrigger)
	template.Selected = wipeTodo(template.Selected)
	template.Target = wipeTodo(template.Target)
	template.Thumbnail = wipeTodo(template.Thumbnail)
	template.TransportType = wipeTodo(template.TransportType)
	template.TriggerStyle = wipeTodo(template.TriggerStyle)
	template.Connected.Id = 0
	template.Video.A = wipeTodo(template.Video.A)
	template.Video.B = wipeTodo(template.Video.B)
	template.Video.Mixer = wipeTodo(template.Video.Mixer)
	template.Video.Opacity = wipeTodo(template.Video.Opacity)
	template.Video.R = wipeTodo(template.Video.R)
	template.Video.Resize = wipeTodo(template.Video.Resize)
	template.Video.SourceParams = wipeTodo(template.Video.SourceParams)

	for idx, effect := range template.Video.Effects {
		template.Video.Effects[idx] = wipeTodo(effect)
	}
	return template
}
func wipeTodo(todo resolume.Todo) resolume.Todo {
	if todo == nil {
		return nil
	}
	delete(todo, "id")
	return todo
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

func matchAudioVideo(ctx context.Context, audioOutDir, videoOutDir string) ([]string, error) {

	audiofiles, err := filepath.Glob(filepath.Join(audioOutDir, "*.m4a"))
	if err != nil {
		slog.Error("Error globbing audio files", "error", err)
		return nil, err
	}

	videoFiles, err := filepath.Glob(filepath.Join(videoOutDir, "*.mov"))
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
		vfbase = vfbase[:len(vfbase)-len(".mov")]
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
		err := enc.Encode(ctx, fname, outDir)
		if err != nil {
			slog.Error("Error converting", "error", err)
			return err
		}
	}

	return nil
}
