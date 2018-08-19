package model

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/IgaguriMK/bgslogviewer/api"
)

type Factions struct {
	ControllingFaction ControllingFaction `json:"controllingFaction"`
	Factions           []Faction          `json:"factions"`
	ID                 int64              `json:"id"`
	Name               string             `json:"name"`
	FetchedAt          time.Time          `json:"fetchedAt"`
}

type ControllingFaction struct {
	Allegiance string `json:"allegiance"`
	Government string `json:"government"`
	Name       string `json:"name"`
}

func FromApiResult(apiResult api.SystemFactions) Factions {
	var factions Factions

	factions.ID = apiResult.ID
	factions.Name = apiResult.Name

	factions.ControllingFaction = ControllingFaction{
		Allegiance: apiResult.ControllingFaction.Allegiance,
		Government: apiResult.ControllingFaction.Government,
		Name:       apiResult.ControllingFaction.Name,
	}

	for _, f := range apiResult.Factions {
		factions.Factions = append(
			factions.Factions,
			factionFromApi(f),
		)
	}

	return factions
}

func (f *Factions) GenStr(loc *time.Location, fmtStr string) {
	for i := 0; i < len(f.Factions); i++ {
		f.Factions[i].GenStr(loc, fmtStr)
	}
}

type Faction struct {
	Allegiance   string   `json:"allegiance"`
	Government   string   `json:"government"`
	IsPlayer     bool     `json:"isPlayer"`
	Name         string   `json:"name"`
	NewestStates States   `json:"newest"`
	History      []States `json:"history"`
}

func factionFromApi(apiFaction api.Faction) Faction {
	var faction Faction

	faction.Allegiance = apiFaction.Allegiance
	faction.Government = apiFaction.Government
	faction.IsPlayer = apiFaction.IsPlayer
	faction.Name = apiFaction.Name

	faction.NewestStates = States{
		Date:       apiFaction.LastUpdate,
		Influence:  apiFaction.Influence,
		Current:    apiFaction.State,
		Pending:    fromApiStateSlice(apiFaction.PendingStates),
		Recovering: fromApiStateSlice(apiFaction.RecoveringStates),
	}

	history := make(map[int64]*States)

	for ds, inf := range apiFaction.InfluenceHistory {
		d := roundedUnix(ds)

		h, ok := history[d]

		if !ok {
			history[d] = &States{
				Date:      d,
				Influence: inf,
			}
			continue
		}

		h.Influence = inf
	}

	for ds, st := range apiFaction.StateHistory {
		d := roundedUnix(ds)

		h, ok := history[d]
		if !ok {
			history[d] = &States{
				Date:    d,
				Current: st,
			}
			continue
		}

		h.Current = st
	}

	for ds, st := range apiFaction.PendingStatesHistory {
		d := roundedUnix(ds)

		h, ok := history[d]
		if !ok {
			history[d] = &States{
				Date:    d,
				Pending: fromApiStateSlice(st),
			}
			continue
		}

		h.Pending = fromApiStateSlice(st)
	}

	for ds, st := range apiFaction.RecoveringStatesHistory {
		d := roundedUnix(ds)

		h, ok := history[d]
		if !ok {
			history[d] = &States{
				Date:       d,
				Recovering: fromApiStateSlice(st),
			}
			continue
		}

		h.Recovering = fromApiStateSlice(st)
	}

	ds := make([]int64, 0, len(history))
	for d, _ := range history {
		ds = append(ds, d)
	}
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })

	for _, d := range ds {
		faction.History = append(
			faction.History,
			*history[d],
		)
	}

	return faction
}

func roundedUnix(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatalf("[ERROR] can't parse %q into int64", str)
	}

	return v - v%(1*60*60)
}

func (f *Faction) GenStr(loc *time.Location, fmtStr string) {
	f.NewestStates.GenStr(loc, fmtStr)

	for i := 0; i < len(f.History); i++ {
		f.History[i].GenStr(loc, fmtStr)
	}
}

type States struct {
	Date         int64   `json:"date"`
	DateStr      string  `json:"-"`
	Influence    float64 `json:"influence"`
	InfluenceStr string  `json:"-"`
	Current      string  `json:"current"`
	Pending      []State `json:"pending"`
	Recovering   []State `json:"recovering"`
}

func (s *States) GenStr(loc *time.Location, fmtStr string) {
	s.DateStr = time.Unix(s.Date, 0).In(loc).Format(fmtStr)

	if s.Influence > 0.0 {
		s.InfluenceStr = fmt.Sprintf("%.1f", s.Influence*100)
	} else {
		s.InfluenceStr = ""
	}

	genStateSliceStr(s.Pending)
	genStateSliceStr(s.Recovering)
}

type State struct {
	State    string `json:"state"`
	Trend    int64  `json:"trend"`
	TrendStr string `json:"-"`
}

func fromApiStateSlice(ss []api.State) []State {
	rss := make([]State, 0, len(ss))

	for _, s := range ss {
		rss = append(
			rss,
			State{
				State: s.State,
				Trend: s.Trend,
			},
		)
	}

	return rss
}

func genStateSliceStr(ss []State) {
	for i := 0; i < len(ss); i++ {
		ss[i].GenStr()
	}
}

func (s *State) GenStr() {
	switch s.Trend {
	case -1:
		s.TrendStr = "↓"
	case 0:
		s.TrendStr = "→"
	case 1:
		s.TrendStr = "↑"
	default:
		log.Printf("[WARNING] Unknown trend %d", s.Trend)
		s.TrendStr = "?"
	}
}
