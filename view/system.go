package view

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/IgaguriMK/bgslogviewer/lib"
	"github.com/IgaguriMK/bgslogviewer/model"
)

const (
	timeFormat     = "2006-01-02 (15)"
	timeLongFormat = "2006-01-02 15:04:05"
)

const (
	imgFederation = "/static/img/federation.png"
	imgEmpire     = "/static/img/empire.png"
	imgAlliance   = "/static/img/alliance.png"
)

var statTemplate *template.Template

func init() {
	var err error

	statTemplate, err = template.New("systemstats.html.tpl").ParseFiles("template/systemstats.html.tpl")
	if err != nil {
		log.Fatal("[FATAL] Failed parse template: ", err)
	}
}

func System(res *bytes.Buffer, factions model.Factions) error {
	factions.GenStr(time.UTC, timeFormat)

	values := buildSystemMain(factions)

	err := statTemplate.Execute(res, values)
	if err != nil {
		return err
	}

	return nil
}

type SystemMain struct {
	SystemName      string
	CachedAt        string
	Overview        []FactionOverview
	RetreatedExists bool
	Retreated       []FactionOverview
	History         []History
}

func buildSystemMain(factions model.Factions) SystemMain {
	var s SystemMain

	s.SystemName = factions.Name
	s.CachedAt = factions.FetchedAt.In(time.UTC).Format(timeLongFormat) + " UTC"

	n := len(factions.Factions)
	controlName := factions.ControllingFaction.Name

	s.Overview = make([]FactionOverview, 0, n)
	s.Retreated = make([]FactionOverview, 0, n)
	s.History = make([]History, 0, n)

	for _, f := range factions.Factions {
		o, exist := buildFactionOverview(f, controlName)
		if exist {
			s.Overview = append(s.Overview, o)
		} else {
			s.Retreated = append(s.Retreated, o)
			s.RetreatedExists = true
		}

		s.History = append(s.History, buildHistory(f))
	}

	return s
}

type FactionOverview struct {
	IsControl     bool
	IsPF          bool
	HasAllegiance bool
	AllegianceImg string
	Name          string
	NameHash      string
	Government    string
	Influence     string
	State         string
	Recovering    string
	Pending       string
	LastUpdate    string
}

func buildFactionOverview(f model.Faction, control string) (overview FactionOverview, isExist bool) {
	isExist = f.NewestStates.Influence > 0.0

	overview = FactionOverview{
		IsControl:  f.Name == control,
		IsPF:       f.IsPlayer,
		Name:       f.Name,
		NameHash:   lib.NameHash("faction history", f.Name),
		Government: f.Government,
		Influence:  formatInfluence(f.NewestStates.Influence),
		State:      f.NewestStates.Current,
		Recovering: formatStates(f.NewestStates.Recovering),
		Pending:    formatStates(f.NewestStates.Pending),
		LastUpdate: time.Unix(f.NewestStates.Date, 0).Format(timeFormat),
	}

	overview.HasAllegiance = true
	switch f.Allegiance {
	case "Federation":
		overview.AllegianceImg = imgFederation
	case "Empire":
		overview.AllegianceImg = imgEmpire
	case "Alliance":
		overview.AllegianceImg = imgAlliance
	default:
		overview.HasAllegiance = false
	}

	return
}

type History struct {
	Name     string
	NameHash string
	Records  []HistoryRecord
}

func buildHistory(f model.Faction) History {
	var h History

	h.Name = f.Name
	h.NameHash = lib.NameHash("faction history", f.Name)
	h.Records = make([]HistoryRecord, 0, len(f.History))

	for _, hh := range f.History {
		h.Records = append(
			h.Records,
			buildHistoryRecord(hh),
		)
	}

	return h
}

type HistoryRecord struct {
	Date       string
	Influence  RecordField
	State      RecordField
	Recovering RecordField
	Pending    RecordField
}

func buildHistoryRecord(h model.States) HistoryRecord {
	return HistoryRecord{
		Date:       time.Unix(h.Date, 0).Format(timeFormat),
		Influence:  RecordField{h.ValidInfluence, formatInfluence(h.Influence)},
		State:      RecordField{h.ValidCurrent, h.Current},
		Recovering: RecordField{h.ValidRecovering, formatStates(h.Recovering)},
		Pending:    RecordField{h.ValidPending, formatStates(h.Pending)},
	}
}

type RecordField struct {
	IsValid bool
	Value   string
}

func formatStates(states []model.State) string {
	strs := make([]string, 0)

	for _, s := range states {
		var trend string
		switch s.Trend {
		case -1:
			trend = "↓"
		case 0:
			trend = "→"
		case 1:
			trend = "↑"
		default:
			log.Printf("error: Unknown trend %d", s.Trend)
			trend = "(?)"
		}

		strs = append(strs, s.State+trend)
	}

	return strings.Join(strs, " ")
}

func formatInfluence(inf float64) string {
	return fmt.Sprintf("%.1f", 100*inf)
}
