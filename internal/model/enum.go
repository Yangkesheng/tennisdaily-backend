package model

type SessionType int16

const (
	SessionTypeDoubles      SessionType = 1
	SessionTypeSingles      SessionType = 2
	SessionTypeTraining     SessionType = 3
	SessionTypeSinglesMatch SessionType = 4
	SessionTypeDoublesMatch SessionType = 5
)

type MatchRank int16

const (
	MatchRankNone         MatchRank = 0
	MatchRankChampion     MatchRank = 1
	MatchRankRunnerUp     MatchRank = 2
	MatchRankSemiFinal    MatchRank = 3
	MatchRankQuarterFinal MatchRank = 4
	MatchRankGroupStage   MatchRank = 5
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
	case MatchRankSemiFinal:
		return "四强"
	case MatchRankQuarterFinal:
		return "八强"
	case MatchRankGroupStage:
		return "小组赛"
	default:
		return "未知"
	}
}
