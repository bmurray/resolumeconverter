package resolume

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type Resolume struct {
	baseUrl *url.URL

	log *slog.Logger
}

type ResolumeOption func(*Resolume)

func NewResolume(baseUrl *url.URL, opts ...ResolumeOption) *Resolume {
	r := &Resolume{
		baseUrl: baseUrl,
		log:     slog.Default().With("pkg", "resolume"),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r Resolume) GetComposition(ctx context.Context) (Composition, error) {
	var v Composition
	u, err := r.baseUrl.Parse("composition")
	if err != nil {
		return v, err
	}
	r.log.Info("Getting composition", "url", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return v, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// return v, enc.Encode(v)
}

func (r Resolume) GetSelectedClip(ctx context.Context) (Clip, error) {
	var v Clip
	u, err := r.baseUrl.Parse("composition/clips/selected")
	if err != nil {
		return v, err
	}
	r.log.Info("Getting selected clip", "url", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return v, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
func (r Resolume) GetLayers(ctx context.Context) ([]Layer, error) {
	comp, err := r.GetComposition(ctx)
	if err != nil {
		return nil, err
	}
	return comp.Layers, nil
}
