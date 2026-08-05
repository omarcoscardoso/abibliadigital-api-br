package models

// Abbrev represents language abbreviation mapping (pt, en)
type Abbrev struct {
	Pt string `json:"pt"`
	En string `json:"en"`
}

// Book represents a Bible book entity
type Book struct {
	ID        int    `json:"-"`
	Abbrev    Abbrev `json:"abbrev"`
	Author    string `json:"author"`
	Chapters  int    `json:"chapters"`
	Comment   string `json:"comment,omitempty"`
	Group     string `json:"group"`
	Name      string `json:"name"`
	Order     int    `json:"-"`
	Testament string `json:"testament"`
}

// VerseHeader represents book info nested inside verse responses
type VerseHeader struct {
	Abbrev  Abbrev `json:"abbrev"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

// SingleVerseResponse represents response for GET /api/verses/:version/:abbrev/:chapter/:number and random
type SingleVerseResponse struct {
	Book    VerseHeader `json:"book"`
	Chapter int         `json:"chapter"`
	Number  int         `json:"number"`
	Text    string      `json:"text"`
}

// ChapterMeta represents the metadata for a chapter response
type ChapterMeta struct {
	Number int `json:"number"`
	Verses int `json:"verses"`
}

// VerseItem represents verse item within chapter list
type VerseItem struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

// ChapterResponse represents response for GET /api/verses/:version/:abbrev/:chapter
type ChapterResponse struct {
	Book    VerseHeader `json:"book"`
	Chapter ChapterMeta `json:"chapter"`
	Verses  []VerseItem `json:"verses"`
}

// SearchRequest represents body for POST /api/verses/search
type SearchRequest struct {
	Version string `json:"version"`
	Search  string `json:"search"`
}

// SearchVerseItem represents verse item returned in search response
type SearchVerseItem struct {
	Book    Book   `json:"book"`
	Chapter int    `json:"chapter"`
	Number  int    `json:"number"`
	Text    string `json:"text"`
}

// SearchResponse represents response for POST /api/verses/search
type SearchResponse struct {
	Occurrence int               `json:"occurrence"`
	Version    string            `json:"version"`
	Verses     []SearchVerseItem `json:"verses"`
}

// VersionResponse represents response for GET /api/versions
type VersionResponse struct {
	Version string `json:"version"`
	Verses  int    `json:"verses"`
}

// ErrorResponse represents standard error JSON response
type ErrorResponse struct {
	Msg string `json:"msg"`
}

// CheckResponse represents response for GET /api/check and /health endpoints
type CheckResponse struct {
	Result    string `json:"result"`
	Status    string `json:"status"`
	Database  string `json:"database"`
	Uptime    string `json:"uptime,omitempty"`
	Timestamp string `json:"timestamp"`
}

