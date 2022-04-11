package ttr

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jaczerob/madamchuckle/pkg/toontown"
)

func InvasionDataToEmbed(i *toontown.InvasionData) (e *discordgo.MessageEmbed, err error) {
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

func PopulationDataToEmbed(p *toontown.PopulationData) (e *discordgo.MessageEmbed, err error) {
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

func FieldOfficeDataToEmbed(f *toontown.FieldOfficeData) (e *discordgo.MessageEmbed, err error) {
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
		zoneField.Value += fmt.Sprintf("**%s** %s\n", zone.Name, strings.Repeat("⭐️", fieldOffice.Difficulty))

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

func StatusToEmbed(s *toontown.Status) (e *discordgo.MessageEmbed, err error) {
	e = &discordgo.MessageEmbed{
		Title: "Toontown Rewritten Server Status",
	}

	if s.Open {
		e.Description = "Toontown Rewritten is currently open!"
	} else {
		e.Description = fmt.Sprintf("⚠️Toontown Rewritten is currently closed!⚠️\n%s", s.Banner)
	}

	return
}

func SillyMeterToEmbed(s *toontown.SillyMeter) (e *discordgo.MessageEmbed, err error)
