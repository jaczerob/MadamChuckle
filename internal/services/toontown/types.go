package toontown

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ============================ INVASIONS ============================= //

type Invasion struct {
	AsOf     int    `json:"asOf"`
	Type     string `json:"type"`
	Progress string `json:"progress"`
}

type InvasionData struct {
	LastUpdated int                 `json:"lastUpdated"`
	Invasions   map[string]Invasion `json:"invasions"`
}

func (i *InvasionData) ToEmbed() (e *discordgo.MessageEmbed, err error) {
	districtField := &discordgo.MessageEmbedField{
		Name:   "District",
		Value:  "",
		Inline: true,
	}

	invasionField := &discordgo.MessageEmbedField{
		Name:   "Invasion",
		Value:  "",
		Inline: true,
	}

	for district, invasion := range i.Invasions {
		districtField.Value += fmt.Sprintf("**%s**\n", district)
		invasionField.Value += fmt.Sprintf("%s: %s\n", invasion.Type, invasion.Progress)
	}

	e = &discordgo.MessageEmbed{
		Title: "Invasions",
		Fields: []*discordgo.MessageEmbedField{
			districtField,
			invasionField,
		},
	}

	return
}

/////////////////////////////////////////////////////////////////////////

// ============================ POPULATION ============================ //

type PopulationData struct {
	LastUpdated          int            `json:"lastUpdated"`
	TotalPopulation      int            `json:"totalPopulation"`
	PopulationByDistrict map[string]int `json:"populationByDistrict"`
}

func (p *PopulationData) ToEmbed() (e *discordgo.MessageEmbed, err error) {
	districtField := &discordgo.MessageEmbedField{
		Name:   "District",
		Value:  "",
		Inline: true,
	}

	populationField := &discordgo.MessageEmbedField{
		Name:   "Population",
		Value:  "",
		Inline: true,
	}

	for district, population := range p.PopulationByDistrict {
		districtField.Value += fmt.Sprintf("**%s**\n", district)

		statusEmoji := "🔵"
		if population >= 300 {
			statusEmoji = "🟢"
		} else if population >= 500 {
			statusEmoji = "🔴"
		}

		populationField.Value += fmt.Sprintf("%s %d\n", statusEmoji, population)
	}

	e = &discordgo.MessageEmbed{
		Title: "Invasions",
		Fields: []*discordgo.MessageEmbedField{
			districtField,
			populationField,
		},
	}

	return
}

/////////////////////////////////////////////////////////////////////////

// ========================== FIELD OFFICES ========================== //

type FieldOffice struct {
	Department string `json:"department"`
	Difficulty int    `json:"difficulty"`
	Annexes    int    `json:"annexes"`
	Open       bool   `json:"open"`
	Expiring   int    `json:"expiring"`
}

type FieldOfficeData struct {
	LastUpdated  int                    `json:"lastUpdated"`
	FieldOffices map[string]FieldOffice `json:"fieldOffices"`
}

func (f *FieldOfficeData) ToEmbed() (e *discordgo.MessageEmbed, err error) {
	zoneField := &discordgo.MessageEmbedField{
		Name:   "Zone",
		Value:  "",
		Inline: true,
	}

	openField := &discordgo.MessageEmbedField{
		Name:   "Open?",
		Value:  "",
		Inline: true,
	}

	for zone, fieldOffice := range f.FieldOffices {
		zoneField.Value += fmt.Sprintf("**%s** %s\n", zone, strings.Repeat("⭐️", fieldOffice.Difficulty))

		if fieldOffice.Open {
			annexStr := "annexes"
			if fieldOffice.Annexes == 1 {
				annexStr = "annex"
			}

			openField.Value += fmt.Sprintf("Yes (%d %s)\n", fieldOffice.Annexes, annexStr)
		} else {
			openField.Value += "No\n"
		}
	}

	e = &discordgo.MessageEmbed{
		Title: "Field Offices",
		Fields: []*discordgo.MessageEmbedField{
			zoneField,
			openField,
		},
	}

	return
}

/////////////////////////////////////////////////////////////////////////
