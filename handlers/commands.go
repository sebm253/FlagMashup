package handlers

import (
	"flag-mashup/data"

	"github.com/disgoorg/disgo/handler"
)

func NewHandler(codeData *data.CodeData) *Handler {
	h := &Handler{
		CodeData: codeData,
		Router:   handler.New(),
	}
	h.Command("/mashup", h.HandleMashup)
	h.Autocomplete("/mashup", h.HandleAutocomplete)
	return h
}

type Handler struct {
	CodeData *data.CodeData
	handler.Router
}
