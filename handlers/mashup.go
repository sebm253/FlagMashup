package handlers

import (
	"bytes"
	"flag-mashup/utils"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

const (
	defaultMaxColors = 3
)

func (h *Handler) HandleMashup(event *handler.CommandEvent) error {
	data := event.SlashCommandInteractionData()
	codes := h.CodeData.Map()
	messageBuilder := discord.NewMessageCreateBuilder()
	src := data.String("source")
	if _, ok := codes[src]; !ok {
		return event.CreateMessage(messageBuilder.
			SetContentf("Code for source (**%s**) not found.", src).
			SetEphemeral(true).
			Build())
	}
	dst := data.String("destination")
	if _, ok := codes[dst]; !ok {
		return event.CreateMessage(messageBuilder.
			SetContentf("Code for destination (**%s**) not found.", dst).
			SetEphemeral(true).
			Build())
	}
	if src == dst {
		return event.CreateMessage(messageBuilder.
			SetContentf("What's there to mashup?").
			SetEphemeral(true).
			Build())
	}

	hide, ok := data.OptBool("hide")
	if !ok {
		hide = true
	}
	if err := event.DeferCreateMessage(hide); err != nil {
		return err
	}

	maxColors, ok := data.OptInt("maximum")
	if !ok {
		maxColors = defaultMaxColors
	}

	buf := new(bytes.Buffer)
	if err := utils.MashupFlags(src, dst, maxColors, h.CodeData, buf); err != nil {
		_, err := event.CreateFollowupMessage(messageBuilder.
			SetContentf("Could not mashup flags: %s", err.Error()).
			Build())
		return err
	}
	_, err := event.CreateFollowupMessage(messageBuilder.
		AddFile("mashup.png", fmt.Sprintf("A flag mashup of %s and %s.", codes[src], codes[dst]), buf).
		Build())
	return err
}
