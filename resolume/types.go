package resolume

import "encoding/json"

type Composition struct {
	Audio            json.RawMessage `json:"audio"`
	Bypassed         json.RawMessage `json:"bypassed"`
	ClipBeatSnap     json.RawMessage `json:"clipbeatsnap"`
	ClipTriggerStyle json.RawMessage `json:"cliptriggerstyle"`
	Columns          []Column        `json:"columns"`
	Crossfader       json.RawMessage `json:"crossfader"`
	Dashboard        json.RawMessage `json:"dashboard"`
	Decks            []Deck          `json:"decks"`
	Layergroups      []LayerGroup    `json:"layergroups"`
	Layers           []Layer         `json:"layers"`
	Master           json.RawMessage `json:"master"`
	Name             json.RawMessage `json:"name"`
	Selected         json.RawMessage `json:"selected"`
	Speed            json.RawMessage `json:"speed"`
	Tempcontroller   json.RawMessage `json:"tempcontroller"`
	Video            json.RawMessage `json:"video"`
}
type Deck any

type LayerGroup any

type Column any

type Layer struct {
	Id                  int             `json:"id"`
	Audio               json.RawMessage `json:"audio"`
	Autopilot           json.RawMessage `json:"autopilot"`
	Bypassed            json.RawMessage `json:"bypassed"`
	Clips               []Clip          `json:"clips"`
	Colorid             json.RawMessage `json:"colorid"`
	CrossFaderGroup     json.RawMessage `json:"crossfadegroup"`
	Dashboard           json.RawMessage `json:"dashboard"`
	FaderStart          json.RawMessage `json:"faderstart"`
	IgnoreColumnTrigger json.RawMessage `json:"ignorecolumntrigger"`
	MaskMode            json.RawMessage `json:"maskmode"`
	Master              json.RawMessage `json:"master"`
	Name                Parameter       `json:"name"`
	Selected            json.RawMessage `json:"selected"`
	Solo                json.RawMessage `json:"solo"`
	Transition          json.RawMessage `json:"transition"`
	Video               Video           `json:"video"`
}

type Clip struct {
	Audio               json.RawMessage `json:"audio"`
	BeatSnap            json.RawMessage `json:"beatsnap"`
	Connected           json.RawMessage `json:"connected"`
	Dashboard           json.RawMessage `json:"dashboard"`
	FaderStart          json.RawMessage `json:"faderstart"`
	Id                  int             `json:"id"`
	IgnoreColumnTrigger json.RawMessage `json:"ignorecolumntrigger"`
	Name                json.RawMessage `json:"name"`
	Selected            json.RawMessage `json:"selected"`
	Target              json.RawMessage `json:"target"`
	Thumbnail           json.RawMessage `json:"thumbnail"`
	TransportType       json.RawMessage `json:"transporttype"`
	TriggerStyle        json.RawMessage `json:"triggerstyle"`
	Video               ClipVideo       `json:"video"`
}

type ClipVideo struct {
	A            json.RawMessage `json:"a"`
	B            json.RawMessage `json:"b"`
	Description  string          `json:"description"`
	Effects      []Effect        `json:"effects"`
	FileInfo     json.RawMessage `json:"fileinfo"`
	G            json.RawMessage `json:"g"`
	Height       int             `json:"height"`
	Mixer        json.RawMessage `json:"mixer"`
	Opacity      json.RawMessage `json:"opacity"`
	R            json.RawMessage `json:"r"`
	Resize       json.RawMessage `json:"resize"`
	SourceParams json.RawMessage `json:"sourceparams"`
	Width        int             `json:"width"`
}

type Video struct {
	Autosize json.RawMessage `json:"autosize"`
	Effects  []Effect        `json:"effects"`
	Height   int             `json:"height"`
	Width    int             `json:"width"`
	Mixer    json.RawMessage `json:"mixer"`
	Opacity  json.RawMessage `json:"opacity"`
}

type Effect any

type Parameter struct {
	ValueType string `json:"valuetype"`
	Id        int    `json:"id"`
	Value     any    `json:"value"`
}
