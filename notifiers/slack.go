package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hibare/GoYTS/utils"
)

type slackPayload struct {
	Text string `json:"text"`
}

func Slack(webhook string, movies map[string]utils.Movie) {

	for k, v := range movies {
		log.Printf("Sending Slack notification for %s\n", k)
		data := slackPayload{}
		data.Text = fmt.Sprintf("Hey found new movie on YTS\n\n*%s* `<%s|view>`\n\n%s", v.Title, v.Link, v.CoverImage)
		jsonReq, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonReq))
		if err != nil {
			log.Println(err)
			continue
		}
		res, _ := http.DefaultClient.Do(req)
		log.Printf("Slack notification status %v", res.StatusCode)
	}
}
