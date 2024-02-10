package handlers

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func (h *Handler) HandleAutocomplete(event *handler.AutocompleteEvent) error {
	data := event.Data
	focused := data.Focused()
	var choices []discord.AutocompleteChoice

	input := data.String(focused.Name)
	codes := h.CodeData.Map()
	if len(input) == 0 {
		keys := maps.Keys(codes)
		for i := 0; i < 25; i++ {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  codes[keys[i]],
				Value: keys[i],
			})
		}
	} else {
		for code, name := range codes {
			if strings.Contains(strings.ToLower(name), strings.ToLower(input)) {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  name,
					Value: code,
				})
			}
		}
		choices = choices[:min(len(choices), 25)] // me when the
	}
	sortChoices(choices)
	return event.AutocompleteResult(choices)
}

func sortChoices(choices []discord.AutocompleteChoice) {
	slices.SortFunc(choices, func(a, b discord.AutocompleteChoice) int {
		return strings.Compare(a.(discord.AutocompleteChoiceString).Name, b.(discord.AutocompleteChoiceString).Name)
	})
}
