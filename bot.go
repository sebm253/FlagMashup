package main

import (
	"context"
	"flag-mashup/data"
	"flag-mashup/handlers"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/lmittmann/tint"
)

func main() {
	codeData := &data.CodeData{}
	codeData.Populate()

	logger := tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logger))

	slog.Info("starting the bot...", slog.String("disgo.version", disgo.Version))

	client, err := disgo.New(os.Getenv("FLAG_MASHUP_TOKEN"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNone),
			gateway.WithPresenceOpts(gateway.WithCompetingActivity("vexillology"))),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsNone)),
		bot.WithEventListeners(handlers.NewHandler(codeData)))
	if err != nil {
		panic(err)
	}

	defer client.Close(context.TODO())

	if err := client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	slog.Info("flag mashup bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
