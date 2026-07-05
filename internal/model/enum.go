package model

import "fmt"

type SessionType int16

const (
	SessionTypeDoubles      SessionType = 1
	SessionTypeSingles      SessionType = 2
	SessionTypeTraining     SessionType = 3
	SessionTypeSinglesMatch SessionType = 4
	SessionTypeDoublesMatch SessionType = 5
)

type SessionCategory int16

const (
	SessionCategoryDaily    SessionCategory = 1
	SessionCategoryTraining SessionCategory = 2
	SessionCategoryMatch    SessionCategory = 3
)

type SessionSubCategory int16

const (
	SessionSubCategorySingles SessionSubCategory = 1
	SessionSubCategoryDoubles SessionSubCategory = 2
	SessionSubCategoryServe   SessionSubCategory = 3
	SessionSubCategoryOther   SessionSubCategory = 4
)

type MatchRank int16

const (
	MatchRankNone         MatchRank = 0
	MatchRankChampion     MatchRank = 1
	MatchRankRunnerUp     MatchRank = 2
	MatchRankThirdPlace   MatchRank = 3
	MatchRankSemiFinal    MatchRank = 4
	MatchRankQuarterFinal MatchRank = 5
	MatchRankRoundOf16    MatchRank = 6
	MatchRankGroupStage   MatchRank = 7
)

func (t SessionType) IsValid() bool {
	return t >= SessionTypeDoubles && t <= SessionTypeDoublesMatch
}

func (t SessionType) IsMatch() bool {
	return t == SessionTypeSinglesMatch || t == SessionTypeDoublesMatch
}

func (t SessionType) Label() string {
	switch t {
	case SessionTypeDoubles:
		return "双打"
	case SessionTypeSingles:
		return "单打"
	case SessionTypeTraining:
		return "训练"
	case SessionTypeSinglesMatch:
		return "单打比赛"
	case SessionTypeDoublesMatch:
		return "双打比赛"
	default:
		return "未知"
	}
}

func (t SessionType) ToCategoryPair() (SessionCategory, SessionSubCategory) {
	switch t {
	case SessionTypeDoubles:
		return SessionCategoryDaily, SessionSubCategoryDoubles
	case SessionTypeSingles:
		return SessionCategoryDaily, SessionSubCategorySingles
	case SessionTypeTraining:
		return SessionCategoryTraining, SessionSubCategoryOther
	case SessionTypeSinglesMatch:
		return SessionCategoryMatch, SessionSubCategorySingles
	case SessionTypeDoublesMatch:
		return SessionCategoryMatch, SessionSubCategoryDoubles
	default:
		return 0, 0
	}
}

func (c SessionCategory) IsValid() bool {
	return c == SessionCategoryDaily || c == SessionCategoryTraining || c == SessionCategoryMatch
}

func (c SessionCategory) Label() string {
	switch c {
	case SessionCategoryDaily:
		return "日常球局"
	case SessionCategoryTraining:
		return "训练"
	case SessionCategoryMatch:
		return "比赛"
	default:
		return "未知"
	}
}

func (s SessionSubCategory) IsValid() bool {
	return s == SessionSubCategorySingles || s == SessionSubCategoryDoubles || s == SessionSubCategoryServe || s == SessionSubCategoryOther
}

func (s SessionSubCategory) Label(category SessionCategory) string {
	switch category {
	case SessionCategoryDaily:
		switch s {
		case SessionSubCategorySingles:
			return "打单"
		case SessionSubCategoryDoubles:
			return "双打"
		}
	case SessionCategoryTraining:
		switch s {
		case SessionSubCategoryServe:
			return "发球"
		case SessionSubCategoryOther:
			return "其他"
		}
	case SessionCategoryMatch:
		switch s {
		case SessionSubCategorySingles:
			return "单打"
		case SessionSubCategoryDoubles:
			return "双打"
		}
	}
	return "未知"
}

func (c SessionCategory) IsValidSubCategory(subCategory SessionSubCategory) bool {
	switch c {
	case SessionCategoryDaily:
		return subCategory == SessionSubCategorySingles || subCategory == SessionSubCategoryDoubles
	case SessionCategoryTraining:
		return subCategory == SessionSubCategoryServe || subCategory == SessionSubCategoryOther
	case SessionCategoryMatch:
		return subCategory == SessionSubCategorySingles || subCategory == SessionSubCategoryDoubles
	default:
		return false
	}
}

func (c SessionCategory) ToSessionType(subCategory SessionSubCategory) (SessionType, bool) {
	if !c.IsValidSubCategory(subCategory) {
		return 0, false
	}
	switch c {
	case SessionCategoryDaily:
		if subCategory == SessionSubCategoryDoubles {
			return SessionTypeDoubles, true
		}
		return SessionTypeSingles, true
	case SessionCategoryTraining:
		return SessionTypeTraining, true
	case SessionCategoryMatch:
		if subCategory == SessionSubCategoryDoubles {
			return SessionTypeDoublesMatch, true
		}
		return SessionTypeSinglesMatch, true
	default:
		return 0, false
	}
}

func (c SessionCategory) TypeText(subCategory SessionSubCategory) string {
	return fmt.Sprintf("%s · %s", c.Label(), subCategory.Label(c))
}

func (r MatchRank) IsValid() bool {
	return r >= MatchRankNone && r <= MatchRankGroupStage
}

func (r MatchRank) Label() string {
	switch r {
	case MatchRankNone:
		return ""
	case MatchRankChampion:
		return "冠军"
	case MatchRankRunnerUp:
		return "亚军"
	case MatchRankThirdPlace:
		return "季军"
	case MatchRankSemiFinal:
		return "四强"
	case MatchRankQuarterFinal:
		return "八强"
	case MatchRankRoundOf16:
		return "16强"
	case MatchRankGroupStage:
		return "小组赛"
	default:
		return "未知"
	}
}
