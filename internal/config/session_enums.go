package config

import (
	"fmt"

	"tennisdaily-backend/internal/model"
)

type SessionEnumsConfig struct {
	Categories []SessionCategoryConfig `yaml:"categories"`
}

type SessionCategoryConfig struct {
	Value         int16                      `yaml:"value"`
	Label         string                     `yaml:"label"`
	SubCategories []SessionSubCategoryConfig `yaml:"subCategories"`
}

type SessionSubCategoryConfig struct {
	Value int16  `yaml:"value"`
	Label string `yaml:"label"`
}

type SessionCategoryOption struct {
	Category    model.SessionCategory
	SubCategory model.SessionSubCategory
	LegacyType  model.SessionType
}

type SessionCategoryResolver struct {
	categories []SessionCategoryConfig
	labels     map[model.SessionCategory]string
	subLabels  map[model.SessionCategory]map[model.SessionSubCategory]string
	options    map[model.SessionCategory]map[model.SessionSubCategory]SessionCategoryOption
}

func NewSessionCategoryResolver(enums SessionEnumsConfig) (*SessionCategoryResolver, error) {
	if len(enums.Categories) == 0 {
		return nil, fmt.Errorf("sessionEnums.categories is required")
	}

	resolver := &SessionCategoryResolver{
		categories: make([]SessionCategoryConfig, 0, len(enums.Categories)),
		labels:     make(map[model.SessionCategory]string, len(enums.Categories)),
		subLabels:  make(map[model.SessionCategory]map[model.SessionSubCategory]string, len(enums.Categories)),
		options:    make(map[model.SessionCategory]map[model.SessionSubCategory]SessionCategoryOption, len(enums.Categories)),
	}

	for _, category := range enums.Categories {
		categoryValue := model.SessionCategory(category.Value)
		if category.Value <= 0 || category.Label == "" {
			return nil, fmt.Errorf("invalid session category config: value=%d label=%q", category.Value, category.Label)
		}
		if _, exists := resolver.labels[categoryValue]; exists {
			return nil, fmt.Errorf("duplicate session category config: value=%d", category.Value)
		}
		if len(category.SubCategories) == 0 {
			return nil, fmt.Errorf("session category has no sub categories: value=%d", category.Value)
		}

		resolver.labels[categoryValue] = category.Label
		resolver.subLabels[categoryValue] = make(map[model.SessionSubCategory]string, len(category.SubCategories))
		resolver.options[categoryValue] = make(map[model.SessionSubCategory]SessionCategoryOption, len(category.SubCategories))

		subCategories := make([]SessionSubCategoryConfig, 0, len(category.SubCategories))
		for _, subCategory := range category.SubCategories {
			subCategoryValue := model.SessionSubCategory(subCategory.Value)
			if subCategory.Value <= 0 || subCategory.Label == "" {
				return nil, fmt.Errorf("invalid session sub category config: category=%d value=%d label=%q", category.Value, subCategory.Value, subCategory.Label)
			}
			if _, exists := resolver.subLabels[categoryValue][subCategoryValue]; exists {
				return nil, fmt.Errorf("duplicate session sub category config: category=%d value=%d", category.Value, subCategory.Value)
			}

			legacyType, ok := legacyTypeForCategoryOption(categoryValue, subCategoryValue)
			if !ok {
				return nil, fmt.Errorf("session category option cannot map legacy type: category=%d subCategory=%d", category.Value, subCategory.Value)
			}

			resolver.subLabels[categoryValue][subCategoryValue] = subCategory.Label
			resolver.options[categoryValue][subCategoryValue] = SessionCategoryOption{
				Category:    categoryValue,
				SubCategory: subCategoryValue,
				LegacyType:  legacyType,
			}
			subCategories = append(subCategories, subCategory)
		}

		resolver.categories = append(resolver.categories, SessionCategoryConfig{
			Value:         category.Value,
			Label:         category.Label,
			SubCategories: subCategories,
		})
	}

	return resolver, nil
}

func (r *SessionCategoryResolver) Categories() []SessionCategoryConfig {
	categories := make([]SessionCategoryConfig, 0, len(r.categories))
	for _, category := range r.categories {
		subCategories := make([]SessionSubCategoryConfig, len(category.SubCategories))
		copy(subCategories, category.SubCategories)
		categories = append(categories, SessionCategoryConfig{
			Value:         category.Value,
			Label:         category.Label,
			SubCategories: subCategories,
		})
	}
	return categories
}

func (r *SessionCategoryResolver) Label(category model.SessionCategory) (string, bool) {
	label, ok := r.labels[category]
	return label, ok
}

func (r *SessionCategoryResolver) SubCategoryLabel(category model.SessionCategory, subCategory model.SessionSubCategory) (string, bool) {
	labels, ok := r.subLabels[category]
	if !ok {
		return "", false
	}
	label, ok := labels[subCategory]
	return label, ok
}

func (r *SessionCategoryResolver) Option(category model.SessionCategory, subCategory model.SessionSubCategory) (SessionCategoryOption, bool) {
	options, ok := r.options[category]
	if !ok {
		return SessionCategoryOption{}, false
	}
	option, ok := options[subCategory]
	return option, ok
}

func (r *SessionCategoryResolver) OptionByLegacyType(sessionType model.SessionType) (SessionCategoryOption, bool) {
	category, subCategory := sessionType.ToCategoryPair()
	if category == 0 || subCategory == 0 {
		return SessionCategoryOption{}, false
	}
	return r.Option(category, subCategory)
}

func (r *SessionCategoryResolver) TypeText(category model.SessionCategory, subCategory model.SessionSubCategory) (string, bool) {
	categoryLabel, ok := r.Label(category)
	if !ok {
		return "", false
	}
	subCategoryLabel, ok := r.SubCategoryLabel(category, subCategory)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%s · %s", categoryLabel, subCategoryLabel), true
}

func legacyTypeForCategoryOption(category model.SessionCategory, subCategory model.SessionSubCategory) (model.SessionType, bool) {
	if category == model.SessionCategoryTraining {
		return model.SessionTypeTraining, true
	}
	return category.ToSessionType(subCategory)
}
