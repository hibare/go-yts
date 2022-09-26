package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hibare/GoYTS/utils"
)

type embedImage struct {
	Url string `json:"url"`
}

type embed struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Color       int        `json:"color"`
	Image       embedImage `json:"image"`
}

type discordPayload struct {
	Username  string  `json:"username"`
	AvatarUrl string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []embed `json:"embeds"`
}

func Discord(webhook string, movies map[string]utils.Movie) {
	for k, v := range movies {
		log.Printf("Sending Discord notification for %s\n", k)

		embed := embed{}
		embed.Title = v.Title
		embed.Description = fmt.Sprintf("[View](%s)", v.Link)
		embed.Color = 15258703
		embed.Image.Url = v.CoverImage

		data := discordPayload{}
		data.Username = "GoYTS"
		data.AvatarUrl = "https://i.imgur.com/4M34hi2.png"
		data.Content = ":clapper: Movie Alert"
		data.Embeds = append(data.Embeds, embed)

		jsonReq, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonReq))
		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		res, _ := http.DefaultClient.Do(req)

		log.Printf("Discord notification status %v", res.StatusCode)
	}
}
