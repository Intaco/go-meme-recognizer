package search

//import "github.com/abogovski/Go-TelegramBotAPI/tgbot"

// Query search query
type Query struct {
	ClientID int64 //tgbot.Integer1
	Query    string
	IsURL    bool // text otherwise
}

// Entry single result entry
type Entry struct {
	Title string
	URL   string
}

// ToMarkdown entry as markdown url
func (e Entry) ToMarkdown() string {
	return "[" + e.Title + "](" + e.URL + ")"
}

// ProcessedQuery processed query to be sent back to frontend
type ProcessedQuery struct {
	Query   Query
	Results []Entry
}

// ToMarkdown represent Query as markdown data
func (pq ProcessedQuery) ToMarkdown() string {
	md := "*Query:* " + pq.Query.Query + "\n\n*Results:*"
	for _, entry := range pq.Results {
		md = md + "\n" + entry.ToMarkdown()
	}
	return md
}
