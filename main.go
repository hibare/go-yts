package main

import (
	"log"
	"net/url"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/hibare/GoYTS/notifiers"
	"github.com/hibare/GoYTS/utils"
)

var (
	schedule       string
	dataDir        string
	historyFile    string
	slackWebhook   string
	discordWebhook string
)

func init() {
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	dataDir = config.DataDir
	historyFile = config.HistoryFile
	schedule = config.Schedule
	slackWebhook = config.SlackWebhook
	discordWebhook = config.DiscordWebhook
}

func ticker() {
	log.Println("[Start] Scraper task")

	movies := map[string]utils.Movie{}
	urls := []string{"https://yts.mx/", "https://wvw.yts.vc/yify/", "https://yts.lt/"}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
	)

	c.OnHTML("#popular-downloads", func(e *colly.HTMLElement) {
		temp := utils.Movie{}
		e.ForEach("div .browse-movie-wrap", func(_ int, el *colly.HTMLElement) {
			temp.Link = el.ChildAttr(".browse-movie-link", "href")
			temp.TimeStamp = time.Now()
			temp.Title = el.ChildText(".browse-movie-title")
			temp.Year = el.ChildText(".browse-movie-year")
			temp.CoverImage = el.ChildAttr("img", "src")

			base, err := url.Parse(temp.Link)
			if err != nil {
				log.Fatal(err)
			}

			ref, err := url.Parse(temp.CoverImage)
			if err != nil {
				log.Fatal(err)
			}

			u := base.ResolveReference(ref)
			temp.CoverImage = u.String()

			movies[temp.Title] = temp
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referrer", "https://www.google.com/")
		log.Println("Visiting URL:", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Finished URL:", r.Request.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("visited URL:", r.Request.URL, "Status Code: ", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Failed to load URL:", r.Request.URL, "Error:", err)
	})

	q, _ := queue.New(
		2, &queue.InMemoryQueueStorage{MaxSize: 100},
	)

	for _, url := range urls {
		q.AddURL(url)
	}
	q.Run(c)

	log.Printf("Scraped %d movies\n", len(movies))

	history := utils.ReadHistory(dataDir, historyFile)
	diff := utils.DiffHistory(movies, history)
	log.Printf("Found %d new movies", len(diff))
	utils.WriteHistory(diff, history, dataDir, historyFile)

	if len(slackWebhook) > 0 {
		notifiers.Slack(slackWebhook, diff)
	}

	if len(discordWebhook) > 0 {
		notifiers.Discord(discordWebhook, diff)
	}

	log.Println("[End] Scraper task")
}

func main() {
	s := gocron.NewScheduler(time.UTC)
	s.Cron(schedule).Do(ticker)
	log.Println("Starting scheduler")
	log.Printf("Cron %s\n", schedule)
	log.Printf("Data directory %s\n", dataDir)
	log.Printf("History file %s\n", historyFile)
	s.StartBlocking()
}
