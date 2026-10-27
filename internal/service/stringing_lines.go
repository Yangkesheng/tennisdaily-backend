package service

import (
	"fmt"
	"math"
	"strings"
	"tennisdaily-backend/internal/model"
	"unicode/utf8"
)

type stringCatalogReader interface {
	CatalogString(int64) (*model.StringLibrary, error)
}

func sameOptionalID(a, b *int64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
func resolveStringingChoice(repo stringCatalogReader, choice model.StringingChoice, axis model.StringingAxis, previous []model.StringingLine) (model.StringingLine, error) {
	line := model.StringingLine{Axis: axis, Color: strings.TrimSpace(choice.Color)}
	if utf8.RuneCountInString(line.Color) > 50 {
		return line, NewInvalidRequestError("颜色过长")
	}
	if choice.StringID == nil {
		name := strings.TrimSpace(choice.CustomName)
		if choice.VariantID != nil || name == "" || utf8.RuneCountInString(name) > 100 {
			return line, NewInvalidRequestError("请填写球线名称")
		}
		if choice.CustomDiameterMm != nil {
			d := *choice.CustomDiameterMm
			if math.IsNaN(d) || math.IsInf(d, 0) || d < 0.5 || d > 2 || math.Abs(d*100-math.Round(d*100)) > 1e-8 {
				return line, NewInvalidRequestError("线径须为0.50至2.00毫米，最多两位小数")
			}
			line.DiameterMm = choice.CustomDiameterMm
		}
		line.Name = name
		return line, nil
	}
	if choice.CustomDiameterMm != nil || *choice.StringID <= 0 || strings.TrimSpace(choice.CustomName) != "" || (choice.VariantID != nil && *choice.VariantID <= 0) {
		return line, ErrInvalidRequest
	}
	// An unchanged side is historical data, even when its catalog variant was removed.
	for _, old := range previous {
		if old.Axis == axis && sameOptionalID(old.StringID, choice.StringID) && sameOptionalID(old.VariantID, choice.VariantID) && old.Color == line.Color {
			old.ID = 0
			return old, nil
		}
	}
	item, err := repo.CatalogString(*choice.StringID)
	if err != nil {
		return line, err
	}
	// Unavailable references may only be retained unchanged on the same side of an existing record.
	if item == nil || !item.Published {
		for _, old := range previous {
			if old.Axis == axis && sameOptionalID(old.StringID, choice.StringID) && sameOptionalID(old.VariantID, choice.VariantID) && old.Color == line.Color {
				old.ID = 0
				return old, nil
			}
		}
		return line, NewInvalidRequestError("该球线已下架，请重新选择或手动填写")
	}
	line.StringID = choice.StringID
	line.VariantID = choice.VariantID
	line.Name = item.Brand + " " + item.Model
	if choice.VariantID != nil {
		found := false
		for _, v := range item.Variants {
			if v.ID == *choice.VariantID {
				diameter := v.DiameterMm
				line.DiameterMm = &diameter
				found = true
				break
			}
		}
		if !found {
			return line, NewInvalidRequestError("线径不属于该球线，请重新选择")
		}
	}
	// Preserve snapshots when the user did not change this side, even if the catalog was renamed.
	for _, old := range previous {
		if old.Axis == axis && sameOptionalID(old.StringID, choice.StringID) && sameOptionalID(old.VariantID, choice.VariantID) && old.Color == line.Color {
			old.ID = 0
			return old, nil
		}
	}
	return line, nil
}
func (s *RacketService) resolveStringingLines(req model.StringingLinesRequest, previous []model.StringingLine) ([]model.StringingLine, error) {
	main, err := resolveStringingChoice(s.repo, req.Main, model.StringingMain, previous)
	if err != nil {
		return nil, err
	}
	cross, err := resolveStringingChoice(s.repo, req.Cross, model.StringingCross, previous)
	if err != nil {
		return nil, err
	}
	return []model.StringingLine{main, cross}, nil
}
func (s *RacketService) LatestStringing(userID, racketID int64) (*model.StringingRecordResponse, error) {
	if _, err := s.requireRacket(userID, racketID); err != nil {
		return nil, err
	}
	record, err := s.repo.LatestStringing(userID, racketID)
	if err != nil || record == nil {
		return nil, err
	}
	result := model.NewStringingRecordResponse(*record)
	if result.Lines != nil {
		for _, line := range []*model.StringingLine{&result.Lines.Main, &result.Lines.Cross} {
			available := true
			if line.StringID != nil {
				item, lookupErr := s.repo.CatalogString(*line.StringID)
				if lookupErr != nil {
					return nil, lookupErr
				}
				available = item != nil && item.Published
				if available && line.VariantID != nil {
					available = false
					for _, v := range item.Variants {
						if v.ID == *line.VariantID {
							available = true
							break
						}
					}
				}
			}
			line.IsSelectable = &available
		}
	}
	return &result, nil
}

type RecentStringChoice struct {
	model.StringingLine
	IsSelectable bool `json:"isSelectable"`
}

func (s *RacketService) RecentStrings(userID int64) ([]RecentStringChoice, error) {
	result := []RecentStringChoice{}
	seen := map[string]bool{}
	availability := map[int64]*model.StringLibrary{}
	for offset := 0; ; offset += 100 {
		records, err := s.repo.StringingHistoryPage(userID, offset)
		if err != nil {
			return nil, err
		}
		for _, record := range records {
			lines := record.Lines
			// Historical names are offered as explicit manual choices, never auto-matched to the library.
			if len(lines) == 0 && strings.TrimSpace(record.StringName) != "" {
				lines = []model.StringingLine{{Name: record.StringName}}
			}
			for _, line := range lines {
				key := "manual:" + strings.ToLower(strings.Join(strings.Fields(line.Name), " ")) + ":" + line.Color
				if line.DiameterMm != nil {
					key += fmt.Sprintf(":%.2f", *line.DiameterMm)
				}
				selectable := true
				if line.StringID != nil {
					key = fmt.Sprintf("catalog:%d:", *line.StringID) + line.Color
					if line.VariantID != nil {
						key += fmt.Sprintf(":%d", *line.VariantID)
					}
					item, ok := availability[*line.StringID]
					if !ok {
						item, err = s.repo.CatalogString(*line.StringID)
						if err != nil {
							return nil, err
						}
						availability[*line.StringID] = item
					}
					selectable = item != nil && item.Published
					if selectable && line.VariantID != nil {
						selectable = false
						for _, v := range item.Variants {
							if v.ID == *line.VariantID {
								selectable = true
								break
							}
						}
					}
				}
				if seen[key] {
					continue
				}
				seen[key] = true
				result = append(result, RecentStringChoice{StringingLine: line, IsSelectable: selectable})
				if len(result) == 6 {
					return result, nil
				}
			}
		}
		if len(records) < 100 {
			break
		}
	}
	return result, nil
}
