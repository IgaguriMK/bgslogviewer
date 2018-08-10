package api

type FactionsStatus struct {
	ControllingFaction ControllingFaction `json:"controllingFaction"`
	Factions           []Faction          `json:"factions"`
	ID                 int64              `json:"id"`
	ID64               int64              `json:"id64"`
	Name               string             `json:"name"`
	URL                string             `json:"url"`
}

type Faction struct {
	Allegiance              string             `json:"allegiance"`
	Government              string             `json:"government"`
	ID                      int64              `json:"id"`
	Influence               float64            `json:"influence"`
	InfluenceHistory        map[string]float64 `json:"influenceHistory"`
	IsPlayer                bool               `json:"isPlayer"`
	LastUpdate              int64              `json:"lastUpdate"`
	Name                    string             `json:"name"`
	PendingStates           []State            `json:"pendingStates"`
	PendingStatesHistory    map[string][]State `json:"pendingStatesHistory"`
	RecoveringStates        []State            `json:"recoveringStates"`
	RecoveringStatesHistory map[string][]State `json:"recoveringStatesHistory"`
	State                   string             `json:"state"`
	StateHistory            map[string]string  `json:"stateHistory"`
}

type ControllingFaction struct {
	Allegiance string `json:"allegiance"`
	Government string `json:"government"`
	ID         int64  `json:"id"`
	Name       string `json:"name"`
}

type State struct {
	State string `json:"state"`
	Trend int64  `json:"trend"`
}
