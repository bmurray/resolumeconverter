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

type Todo map[string]interface{}

type Clip struct {
	Audio               Todo      `json:"audio"`
	BeatSnap            Todo      `json:"beatsnap"`
	Connected           Connected `json:"connected"`
	Dashboard           Todo      `json:"dashboard"`
	FaderStart          Todo      `json:"faderstart"`
	Id                  int       `json:"id"`
	IgnoreColumnTrigger Todo      `json:"ignorecolumntrigger"`
	Name                Name      `json:"name"`
	Selected            Todo      `json:"selected"`
	Target              Todo      `json:"target"`
	Thumbnail           Todo      `json:"thumbnail"`
	TransportType       Todo      `json:"transporttype"`
	TriggerStyle        Todo      `json:"triggerstyle"`
	Video               ClipVideo `json:"video"`
	Transport           Todo      `json:"transport"`
}

type Name struct {
	Id        int       `json:"id"`
	Value     string    `json:"value"`
	ValueType ValueType `json:"valuetype"`
}

type Connected struct {
	Value     string    `json:"value"`
	ValueType ValueType `json:"valuetype"`
	Options   []string  `json:"options"`
	Index     int       `json:"index"`
	Id        int       `json:"id,omitempty"`
}

type ClipVideo struct {
	A            Todo     `json:"a"`
	B            Todo     `json:"b"`
	Description  string   `json:"description"`
	Effects      []Todo   `json:"effects"`
	FileInfo     FileInfo `json:"fileinfo"`
	G            Todo     `json:"g"`
	Height       int      `json:"height"`
	Mixer        Todo     `json:"mixer"`
	Opacity      Todo     `json:"opacity"`
	R            Todo     `json:"r"`
	Resize       Todo     `json:"resize"`
	SourceParams Todo     `json:"sourceparams"`
	Width        int      `json:"width"`
}

type FileInfo struct {
	Path       string     `json:"path"`
	Exists     bool       `json:"exists"`
	Duration   string     `json:"duration,omitempty"`
	DurationMS float32    `json:"duration_ms,omitempty"`
	Framerate  *Framerate `json:"framerate,omoitempty"`
	Width      int        `json:"width,omitempty"`
	Height     int        `json:"height,omitempty"`
}

type Framerate struct {
	Num int `json:"num"`
	Den int `json:"den"`
}
type Video struct {
	Autosize json.RawMessage `json:"autosize"`
	Effects  []Effect        `json:"effects"`
	Height   int             `json:"height"`
	Width    int             `json:"width"`
	Mixer    json.RawMessage `json:"mixer"`
	Opacity  json.RawMessage `json:"opacity"`
}

type Effect Todo

type Parameter struct {
	ValueType string `json:"valuetype"`
	Id        int    `json:"id"`
	Value     any    `json:"value"`
}

type ValueType string
