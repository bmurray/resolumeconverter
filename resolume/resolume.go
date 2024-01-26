package resolume

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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
	// r.log.Info("Getting composition", "url", u.String())
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
func (r Resolume) FindEmptyClip(ctx context.Context, startLayer, endLayer int) (layer_id int, clip_id Clip, err error) {
	if startLayer > endLayer {
		return 0, Clip{}, fmt.Errorf("startLayer must be less than endLayer")
	}
	comp, err := r.GetComposition(ctx)
	if err != nil {
		return 0, Clip{}, err
	}
	if len(comp.Layers) < endLayer {
		return 0, Clip{}, fmt.Errorf("endLayer is greater than number of layers")
	}

	for i := startLayer; i <= endLayer; i++ {
		layer := comp.Layers[i]
		for _, clip := range layer.Clips {
			if clip.Connected.Value == "Empty" {
				return layer.Id, clip, nil
			}
		}
	}
	return 0, Clip{}, fmt.Errorf("no empty clips found")
}
func (r Resolume) OpenClip(ctx context.Context, clipId int, filePath string) error {

	fmpath := fmt.Sprintf("file://%s", filePath)
	furl, err := url.Parse(fmpath)
	if err != nil {
		return err
	}

	// slog.Info("Opening clip", "path", furl.String())
	b := strings.NewReader(furl.String())
	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/clips/by-id/%d/open", clipId))
	if err != nil {
		return err
	}
	// r.log.Info("Opening clip", "url", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil

}
func (r Resolume) GetClip(ctx context.Context, clipId int) (Clip, error) {
	var v Clip
	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/clips/by-id/%d", clipId))
	if err != nil {
		return v, err
	}
	// r.log.Info("Getting clip", "url", u.String())
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
func (r Resolume) SetClip(ctx context.Context, clipId int, clip Clip) error {
	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(clip)
	if err != nil {
		return err
	}
	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/clips/by-id/%d", clipId))
	if err != nil {
		return err
	}
	// r.log.Info("Setting clip", "url", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
func (r Resolume) SetClipRaw(ctx context.Context, clipId int, val map[string]any) error {
	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(val)
	if err != nil {
		return err
	}
	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/clips/by-id/%d", clipId))
	if err != nil {
		return err
	}
	// r.log.Info("Setting clip", "url", u.String(), "val", b.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), &b)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
func (r Resolume) SetClipByLayerClipRaw(ctx context.Context, layerId, clipId int, val map[string]any) error {
	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(val)
	if err != nil {
		return err
	}
	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/layers/%d/clips/%d", layerId, clipId))
	if err != nil {
		return err
	}
	// r.log.Info("Setting clip", "url", u.String(), "val", b.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), &b)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (r Resolume) GetSelectedClip(ctx context.Context) (Clip, error) {
	var v Clip
	u, err := r.baseUrl.Parse("composition/clips/selected")
	if err != nil {
		return v, err
	}
	// r.log.Info("Getting selected clip", "url", u.String())
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
func (r Resolume) GetThumbnail(ctx context.Context, clipId int) (io.ReadCloser, error) {

	u, err := r.baseUrl.Parse(fmt.Sprintf("composition/clips/by-id/%s/thumbnail", clipId))
	if err != nil {
		return nil, err
	}

	// r.log.Info("Getting thumbnail", "url", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
