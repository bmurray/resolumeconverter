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
	case "layer":
		layers(ctx, r, args[1:])
	case "composition":
		composition(ctx, r, args[1:])
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
