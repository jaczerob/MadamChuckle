package toontown

import (
	"encoding/json"
	"time"
)

// obtained from https://github.com/ToontownRewritten/api-doc/blob/master/field-offices.md#zone-id-lookup
var zoneIDLookup = map[string]string{
	"3100": "Walrus Way",
	"3200": "Sleet Street",
	"3300": "Polar Place",
	"4100": "Alto Avenue",
	"4200": "Baritone Boulevard",
	"4300": "Tenor Terrace",
	"5100": "Elm Street",
	"5200": "Maple Street",
	"5300": "Oak Street",
	"9100": "Lullaby Lane",
	"9200": "Pajama Place",
}

// allows for easy unmarshal of Zones from JSON
type Zone struct {
	Name string
}

func (z *Zone) UnmarshalJSON(bytes []byte) (err error) {
	var raw string

	if err = json.Unmarshal(bytes, &raw); err != nil {
		return
	}

	z.Name = zoneIDLookup[raw]
	return
}

// allows for easy unmarshal of UNIX timestamps from JSON
type Timestamp struct{ time.Time }

func (t *Timestamp) UnmarshalJSON(bytes []byte) (err error) {
	var raw int64

	if err = json.Unmarshal(bytes, &raw); err != nil {
		return
	}

	t.Time = time.Unix(raw, 0)
	return
}

type ToontownObject interface{}

// base type for all Toontown Rewritten API responses
type Base struct {
	ToontownObject
	LastUpdated Timestamp `json:"lastUpdated,omitempty"` // when the server last updated invasion data
}

type Invasion struct {
	AsOf     Timestamp `json:"asOf"`     // when the invasion started
	Type     string    `json:"type"`     // the cog invading
	Progress string    `json:"progress"` // how many cogs have been beated (xxxx/yyyy)
}

// holds Toontown Rewritten Invasion API response data
type InvasionData struct {
	Base
	Invasions map[string]Invasion `json:"invasions"` // a mapping of the district to invasion
}

// holds Toontown Rewritten Population API response data
type PopulationData struct {
	Base
	TotalPopulation      int            `json:"totalPopulation"`
	PopulationByDistrict map[string]int `json:"populationByDistrict"`
}

type FieldOffice struct {
	Department rune      `json:"department"` // the type of cog ('s')
	Difficulty int       `json:"difficulty"` // how many stars the field office is (1, 2, 3)
	Annexes    int       `json:"annexes"`    // how many annexes are left
	Open       bool      `json:"open"`       // whether or not the field office is open
	Expiring   Timestamp `json:"expiring"`   // when the field office will expire, usually nil
}

// holds Toontown Rewritten Field Office API response data
type FieldOfficeData struct {
	Base
	FieldOffices map[Zone]FieldOffice `json:"fieldOffices"` // a mapping showing field office data in each zone
}

// holds Toontown Rewritten Status API response data
type Status struct {
	Base
	Open   bool   `json:"open"`             // tells whether or not the server is open
	Banner string `json:"banner,omitempty"` // if the server is not open, this shows why
}

type Reward struct {
	Title       string
	Description string
	Points      int
}

// holds Toontown Rewritten Status API response data
type SillyMeter struct {
	Base
	State              string    `json:"state"`               // Active, Reward, Inactive
	Health             int       `json:"health"`              // how much HP is left
	Winner             string    `json:"winner"`              // whomever won, if reward state is Reward, otherwise nil
	NextUpdate         Timestamp `json:"nextUpdateTimestamp"` // when the silly meter will next be updated by the server
	AsOf               Timestamp `json:"asOf"`                // when the silly meter cycle started
	RewardTitles       []string  `json:"rewards"`
	RewardDescriptions []string  `json:"rewardDescriptions"`
	RewardPoints       []int     `json:"rewardPoints"`
}

func (s *SillyMeter) Rewards() (rewards []Reward) {
	rewards = make([]Reward, 0)

	for i := 0; i < len(s.RewardTitles); i++ {
		reward := Reward{
			Title:       s.RewardTitles[i],
			Description: s.RewardDescriptions[i],
			Points:      s.RewardPoints[i],
		}

		rewards = append(rewards, reward)
	}

	return
}
