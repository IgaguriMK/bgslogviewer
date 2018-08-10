package api

type SystemFactions struct {
	ControllingFaction ControllingFaction `json:"controllingFaction"`
	Factions           []Faction          `json:"factions"`
	ID                 int64              `json:"id"`
	ID64               int64              `json:"id64"`
	Name               string             `json:"name"`
	URL                string             `json:"url"`
}

type ControllingFaction struct {
	Allegiance string `json:"allegiance"`
	Government string `json:"government"`
	ID         int64  `json:"id"`
	Name       string `json:"name"`
}

type Faction struct {
	Allegiance              string             `json:"allegiance"`
	Government              string             `json:"government"`
	ID                      int64              `json:"id"`
	Influence               float64            `json:"influence"`
	IsPlayer                bool               `json:"isPlayer"`
	LastUpdate              int64              `json:"lastUpdate"`
	Name                    string             `json:"name"`
	State                   string             `json:"state"`
	PendingStates           []State            `json:"pendingStates"`
	RecoveringStates        []State            `json:"recoveringStates"`
	InfluenceHistory        map[string]float64 `json:"influenceHistory"`
	PendingStatesHistory    map[string][]State `json:"pendingStatesHistory"`
	RecoveringStatesHistory map[string][]State `json:"recoveringStatesHistory"`
	StateHistory            map[string]string  `json:"stateHistory"`
}

type State struct {
	State string `json:"state"`
	Trend int64  `json:"trend"`
}
