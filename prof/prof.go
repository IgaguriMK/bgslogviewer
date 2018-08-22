package prof

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IgaguriMK/bgslogviewer/config"
)

var saveCh chan *simpleProf

func init() {
	if config.EnableProf {
		saveCh = make(chan *simpleProf, 2)

		go func() {
			for p := range saveCh {
				bs, err := json.Marshal(p)
				if err != nil {
					log.Println("[ERROR] Profiler error: ", err)
					return
				}

				log.Println("[PROF]", string(bs))
			}
		}()
	}
}

func NewProfiler() Prof {
	if config.EnableProf {
		return &simpleProf{
			Params: make(map[string]interface{}),
			Times:  make([]lap, 0, 8),
		}
	} else {
		return new(nopProf)
	}
}

type Prof interface {
	AddParam(key string, val interface{})
	Start(name string)
	Mark(name string)
	End(code int)
}

type nopProf struct{}

func (_ *nopProf) AddParam(key string, val interface{}) {}
func (_ *nopProf) Start(name string)                    {}
func (_ *nopProf) Mark(name string)                     {}
func (_ *nopProf) End(code int)                         {}

type simpleProf struct {
	Params   map[string]interface{} `json:"p"`
	Times    []lap                  `json:"ts"`
	Code     int                    `json:"c"`
	lastName string
	lastTime time.Time
}

type lap struct {
	Name         string  `json:"n"`
	DurationMSec float64 `json:"ms"`
}

func (sp *simpleProf) AddParam(key string, val interface{}) {
	sp.Params[key] = val
}

func (sp *simpleProf) Start(name string) {
	sp.lastName = name
	sp.lastTime = time.Now()
}

func (sp *simpleProf) Mark(name string) {
	now := time.Now()

	sp.Times = append(
		sp.Times,
		lap{
			Name:         sp.lastName,
			DurationMSec: now.Sub(sp.lastTime).Seconds() * 1000,
		},
	)

	sp.lastName = name
	sp.lastTime = now
}

func (sp *simpleProf) End(code int) {
	now := time.Now()

	sp.Times = append(
		sp.Times,
		lap{
			Name:         sp.lastName,
			DurationMSec: now.Sub(sp.lastTime).Seconds() * 1000,
		},
	)

	sp.Code = code

	saveCh <- sp
}
