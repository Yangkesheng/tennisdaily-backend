package model

import "fmt"

type StringingAxis string

const (
	StringingMain  StringingAxis = "main"
	StringingCross StringingAxis = "cross"
)

type StringingChoice struct {
	CustomDiameterMm *float64 `json:"customDiameterMm"`
	StringID         *int64   `json:"stringId"`
	VariantID        *int64   `json:"variantId"`
	CustomName       string   `json:"customName"`
	Color            string   `json:"color"`
}
type StringingLinesRequest struct {
	Main  StringingChoice `json:"main"`
	Cross StringingChoice `json:"cross"`
}
type StringingLine struct {
	IsSelectable      *bool         `json:"isSelectable,omitempty" gorm:"-"`
	ID                int64         `json:"-" gorm:"primaryKey"`
	StringingRecordID int64         `json:"-"`
	Axis              StringingAxis `json:"-"`
	StringID          *int64        `json:"stringId"`
	VariantID         *int64        `json:"variantId"`
	Name              string        `json:"name" gorm:"column:name_snapshot"`
	DiameterMm        *float64      `json:"diameterMm" gorm:"column:diameter_mm_snapshot"`
	Color             string        `json:"color" gorm:"column:color_snapshot"`
}

func (StringingLine) TableName() string { return "racket_stringing_lines" }

type StringingLinesResponse struct {
	Main  StringingLine `json:"main"`
	Cross StringingLine `json:"cross"`
}

func StringingLinesFromRecord(r RacketStringingRecord) *StringingLinesResponse {
	if len(r.Lines) != 2 {
		return nil
	}
	var result StringingLinesResponse
	for _, line := range r.Lines {
		if line.Axis == StringingMain {
			result.Main = line
		}
		if line.Axis == StringingCross {
			result.Cross = line
		}
	}
	if result.Main.Name == "" || result.Cross.Name == "" {
		return nil
	}
	return &result
}
func StringingLineLabel(line StringingLine) string {
	if line.DiameterMm != nil {
		return fmt.Sprintf("%s · %.2f mm", line.Name, *line.DiameterMm)
	}
	return line.Name
}
func StringingDisplayName(r RacketStringingRecord) string {
	lines := StringingLinesFromRecord(r)
	if lines == nil {
		return r.StringName
	}
	main, cross := StringingLineLabel(lines.Main), StringingLineLabel(lines.Cross)
	if main == cross {
		return main
	}
	return "竖 " + main + " / 横 " + cross
}
