package handlers

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"golang.org/x/exp/maps"
)

func (h *Handler) HandleAutocomplete(event *handler.AutocompleteEvent) error {
	data := event.Data
	focused, found := data.Find(func(option discord.AutocompleteOption) bool {
		return option.Focused
	})
	if !found {
		return nil
	}
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
		return event.AutocompleteResult(choices)
	}
	for code, name := range codes {
		if strings.Contains(strings.ToLower(name), strings.ToLower(input)) {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  name,
				Value: code,
			})
		}
	}
	return event.AutocompleteResult(choices[:min(len(choices), 25)]) // me when the
}
