package utils

import "time"

type Movie struct {
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	CoverImage string    `json:"cover_image"`
	Year       string    `json:"year"`
	TimeStamp  time.Time `json:"timestamp"`
}
