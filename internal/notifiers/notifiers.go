package notifiers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hibare/GoCommon/v2/pkg/notifiers/discord"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/constants"
	"github.com/hibare/go-yts/internal/db"
)

func Discord(ctx context.Context, movies []db.Movies) {
	if !config.Current.NotifierConfig.Discord.Enabled {
		slog.WarnContext(ctx, "Notifier is disabled")
		return
	}

	for _, v := range movies {
		slog.InfoContext(ctx, "Sending Discord notification", "movie", v.Title)

		message := discord.Message{
			Username:  constants.ProgramIdentifierFormatted,
			Content:   ":clapper: Movie Alert",
			AvatarURL: "https://i.imgur.com/4M34hi2.png",
			Embeds: []discord.Embed{
				{
					Title:       fmt.Sprintf("%s (%d)", v.Title, v.Year),
					Description: fmt.Sprintf("[View](%s)", v.Link),
					Color:       15258703,
					Image: discord.EmbedImage{
						URL: v.CoverImage,
					},
				},
			},
		}

		client, err := discord.NewClient(discord.Options{
			WebhookURL: config.Current.NotifierConfig.Discord.Webhook,
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create Discord client", "error", err)
			return
		}
		if err := client.Send(ctx, &message); err != nil {
			slog.ErrorContext(ctx, "Failed to send Discord notification", "error", err)
		} else {
			slog.InfoContext(ctx, "Discord notification sent", "movie", v.Title)
		}
	}
}

func Notify(ctx context.Context, movies []db.Movies) {
	if len(movies) == 0 {
		return
	}

	Discord(ctx, movies)
}
