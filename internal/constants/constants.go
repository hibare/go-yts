package constants

import "time"

const (
	ProgramIdentifier          = "go-yts"
	ProgramIdentifierFormatted = "GoYTS"
	DefaultRequestTimeout      = 60 * time.Second
	DefaultHistoryFilename     = "history.json"
	DefaultDataDir             = "/data"
	DefaultSchedule            = "0 */4 * * *"
	DefaultSQLiteDB            = "movies.db"
)
