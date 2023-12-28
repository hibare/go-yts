package main

import (
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/history"
	"github.com/hibare/go-yts/internal/notifiers"
)

func ConstructURL(baseUrl *url.URL, refUrl string) (string, error) {
	ref, err := url.Parse(refUrl)
	if err != nil {
		return "", err
	}

	u := baseUrl.ResolveReference(ref)

	return u.String(), nil
}

func ticker() {
	slog.Info("[Start] Scraper task")

	movies := history.Movies{}
	urls := []string{"https://yts.mx/", "https://yts.autos/", "https://yts.rs/", "https://yts.lt/", "https://yts.do/"}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
	)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.Current.HTTPConfig.RequestTimeout,
			DualStack: true,
		}).DialContext,
	})

	c.SetRequestTimeout(config.Current.HTTPConfig.RequestTimeout)

	c.OnHTML("#popular-downloads", func(e *colly.HTMLElement) {
		temp := history.Movie{}
		e.ForEach("div .browse-movie-wrap", func(_ int, el *colly.HTMLElement) {
			var err error
			temp.Link = el.ChildAttr(".browse-movie-link", "href")
			temp.TimeStamp = time.Now()
			temp.Title = el.ChildText(".browse-movie-title")
			temp.Year = el.ChildText(".browse-movie-year")
			temp.CoverImage = el.ChildAttr("img", "src")

			temp.Link, err = ConstructURL(e.Request.URL, temp.Link)
			if err != nil {
				slog.Error("Failed to construct URL", "error", err)
				return
			}

			temp.CoverImage, err = ConstructURL(e.Request.URL, temp.CoverImage)
			if err != nil {
				slog.Error("Failed to construct URL", "error", err)
				return
			}

			movies[temp.Title] = temp
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referrer", "https://www.google.com/")
		slog.Info("Visiting URL", "url", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		slog.Info("Finished URL", "url", r.Request.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		slog.Info("Finished URL", "url", r.Request.URL.String(), "status_code", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.Error("Failed to load URL", "url", r.Request.URL.String(), "error", err)
	})

	q, _ := queue.New(
		2, &queue.InMemoryQueueStorage{MaxSize: 100},
	)

	for _, url := range urls {
		q.AddURL(url)
	}

	q.Run(c)

	slog.Info("Scraped movies", "total", len(movies))

	h := history.ReadHistory(config.Current.StorageConfig.DataDir, config.Current.StorageConfig.HistoryFile)
	diff := history.DiffHistory(movies, h)
	history.WriteHistory(diff, h, config.Current.StorageConfig.DataDir, config.Current.StorageConfig.HistoryFile)
	slog.Info("Found new movies", "total", len(diff))

	notifiers.Notify(diff)

	slog.Info("[End] Scraper task")
}

func main() {
	commonLogger.InitDefaultLogger()
	config.LoadConfig()
	slog.Info("Config", "cron", config.Current.Schedule, "request_timeout", config.Current.HTTPConfig.RequestTimeout, "data_dir", config.Current.StorageConfig.DataDir, "history_file", config.Current.StorageConfig.HistoryFile, "history_file", config.Current.StorageConfig.HistoryFile)

	slog.Info("Starting scheduler")

	s := gocron.NewScheduler(time.UTC)
	s.Cron(config.Current.Schedule).Do(ticker)
	s.StartBlocking()
}
