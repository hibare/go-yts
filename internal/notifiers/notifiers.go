package notifiers

import (
	"fmt"

	"github.com/hibare/GoCommon/v2/pkg/notifiers/discord"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/constants"
	"github.com/hibare/go-yts/internal/history"
	log "github.com/sirupsen/logrus"
)

func Discord(movies history.Movies) {
	if !config.Current.NotifierConfig.Discord.Enabled {
		log.Error("Notifier is disabled")
		return
	}

	for k, v := range movies {
		log.Infof("Sending Discord notification for %s\n", k)

		message := discord.Message{
			Username:  constants.ProgramIdentifierFormatted,
			Content:   ":clapper: Movie Alert",
			AvatarURL: "https://i.imgur.com/4M34hi2.png",
			Embeds: []discord.Embed{
				{
					Title:       fmt.Sprintf("%s (%s)", v.Title, v.Year),
					Description: fmt.Sprintf("[View](%s)", v.Link),
					Color:       15258703,
					Image: discord.EmbedImage{
						URL: v.CoverImage,
					},
				},
			},
		}

		if err := message.Send(config.Current.NotifierConfig.Discord.Webhook); err != nil {
			log.Error(err)
		}
	}
}

func Notify(movies history.Movies) {
	if len(movies) == 0 {
		return
	}

	Discord(movies)
}
