package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"abibliadigital/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	connStr := fmt.Sprintf("file:%s?mode=ro", dbPath)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	rand.Seed(time.Now().UnixNano())

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// GetBooks returns all books ordered by book_order
func (s *Store) GetBooks() ([]models.Book, error) {
	query := `
		SELECT abbrev_pt, abbrev_en, name, author, chapters, group_name, testament, comment, book_order 
		FROM books 
		ORDER BY book_order ASC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(
			&b.Abbrev.Pt, &b.Abbrev.En, &b.Name, &b.Author,
			&b.Chapters, &b.Group, &b.Testament, &b.Comment, &b.Order,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

// GetBook returns a book by Portuguese or English abbreviation
func (s *Store) GetBook(abbrev string) (*models.Book, error) {
	query := `
		SELECT abbrev_pt, abbrev_en, name, author, chapters, group_name, testament, comment, book_order 
		FROM books 
		WHERE LOWER(abbrev_pt) = LOWER(?) OR LOWER(abbrev_en) = LOWER(?)
		LIMIT 1
	`
	var b models.Book
	err := s.db.QueryRow(query, abbrev, abbrev).Scan(
		&b.Abbrev.Pt, &b.Abbrev.En, &b.Name, &b.Author,
		&b.Chapters, &b.Group, &b.Testament, &b.Comment, &b.Order,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetChapter returns all verses for a given version, book abbreviation and chapter
func (s *Store) GetChapter(version, abbrev string, chapter int) (*models.ChapterResponse, error) {
	book, err := s.GetBook(abbrev)
	if err != nil || book == nil {
		return nil, err
	}

	query := `
		SELECT number, text 
		FROM verses 
		WHERE version = ? AND (LOWER(book_abbrev_pt) = LOWER(?) OR LOWER(book_abbrev_en) = LOWER(?)) AND chapter = ?
		ORDER BY number ASC
	`
	rows, err := s.db.Query(query, version, abbrev, abbrev, chapter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []models.VerseItem
	for rows.Next() {
		var v models.VerseItem
		if err := rows.Scan(&v.Number, &v.Text); err != nil {
			return nil, err
		}
		verses = append(verses, v)
	}

	if len(verses) == 0 {
		return nil, nil
	}

	return &models.ChapterResponse{
		Book: models.VerseHeader{
			Abbrev:  book.Abbrev,
			Name:    book.Name,
			Author:  book.Author,
			Group:   book.Group,
			Version: version,
		},
		Chapter: models.ChapterMeta{
			Number: chapter,
			Verses: len(verses),
		},
		Verses: verses,
	}, nil
}

// GetVerse returns a specific verse
func (s *Store) GetVerse(version, abbrev string, chapter, number int) (*models.SingleVerseResponse, error) {
	book, err := s.GetBook(abbrev)
	if err != nil || book == nil {
		return nil, err
	}

	query := `
		SELECT text 
		FROM verses 
		WHERE version = ? AND (LOWER(book_abbrev_pt) = LOWER(?) OR LOWER(book_abbrev_en) = LOWER(?)) AND chapter = ? AND number = ?
		LIMIT 1
	`
	var text string
	err = s.db.QueryRow(query, version, abbrev, abbrev, chapter, number).Scan(&text)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &models.SingleVerseResponse{
		Book: models.VerseHeader{
			Abbrev:  book.Abbrev,
			Name:    book.Name,
			Author:  book.Author,
			Group:   book.Group,
			Version: version,
		},
		Chapter: chapter,
		Number:  number,
		Text:    text,
	}, nil
}

// GetRandomVerse returns a random verse for a version and optional book abbreviation
func (s *Store) GetRandomVerse(version, abbrev string) (*models.SingleVerseResponse, error) {
	var book *models.Book
	var err error

	if abbrev != "" {
		book, err = s.GetBook(abbrev)
	}

	if book == nil {
		books, err := s.GetBooks()
		if err != nil || len(books) == 0 {
			return nil, err
		}
		book = &books[rand.Intn(len(books))]
	}

	if book.Chapters <= 0 {
		book.Chapters = 1
	}

	randChapter := rand.Intn(book.Chapters) + 1
	chapterResp, err := s.GetChapter(version, book.Abbrev.Pt, randChapter)

	// Fallback retry if chapter has no verses
	if err != nil || chapterResp == nil || len(chapterResp.Verses) == 0 {
		query := `
			SELECT chapter, number, text 
			FROM verses 
			WHERE version = ? AND (LOWER(book_abbrev_pt) = LOWER(?) OR LOWER(book_abbrev_en) = LOWER(?))
			ORDER BY RANDOM() 
			LIMIT 1
		`
		var chap, num int
		var text string
		err := s.db.QueryRow(query, version, book.Abbrev.Pt, book.Abbrev.Pt).Scan(&chap, &num, &text)
		if err != nil {
			return nil, err
		}
		return &models.SingleVerseResponse{
			Book: models.VerseHeader{
				Abbrev:  book.Abbrev,
				Name:    book.Name,
				Author:  book.Author,
				Group:   book.Group,
				Version: version,
			},
			Chapter: chap,
			Number:  num,
			Text:    text,
		}, nil
	}

	randIdx := rand.Intn(len(chapterResp.Verses))
	selectedVerse := chapterResp.Verses[randIdx]

	return &models.SingleVerseResponse{
		Book:    chapterResp.Book,
		Chapter: randChapter,
		Number:  selectedVerse.Number,
		Text:    selectedVerse.Text,
	}, nil
}

// Search returns verses matching the search string for a version
func (s *Store) Search(version, searchText string) (*models.SearchResponse, error) {
	books, err := s.GetBooks()
	if err != nil {
		return nil, err
	}

	bookMapEn := make(map[string]models.Book)
	for _, b := range books {
		bookMapEn[b.Abbrev.En] = b
		bookMapEn[b.Abbrev.Pt] = b
	}

	words := strings.Fields(searchText)
	if len(words) == 0 {
		return &models.SearchResponse{
			Occurrence: 0,
			Version:    version,
			Verses:     []models.SearchVerseItem{},
		}, nil
	}

	// Build WHERE clause with LIKE for each word
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "version = ?")
	args = append(args, version)

	for _, word := range words {
		conditions = append(conditions, "LOWER(text) LIKE ?")
		args = append(args, "%"+strings.ToLower(word)+"%")
	}

	query := fmt.Sprintf(`
		SELECT book_abbrev_pt, book_abbrev_en, chapter, number, text 
		FROM verses 
		WHERE %s 
		ORDER BY id ASC
	`, strings.Join(conditions, " AND "))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SearchVerseItem
	for rows.Next() {
		var abbrevPt, abbrevEn string
		var chap, num int
		var text string

		if err := rows.Scan(&abbrevPt, &abbrevEn, &chap, &num, &text); err != nil {
			return nil, err
		}

		b, ok := bookMapEn[abbrevEn]
		if !ok {
			b = bookMapEn[abbrevPt]
		}

		items = append(items, models.SearchVerseItem{
			Book:    b,
			Chapter: chap,
			Number:  num,
			Text:    text,
		})
	}

	if items == nil {
		items = []models.SearchVerseItem{}
	}

	return &models.SearchResponse{
		Occurrence: len(items),
		Version:    version,
		Verses:     items,
	}, nil
}

// GetVersions returns all versions and total verse counts
func (s *Store) GetVersions() ([]models.VersionResponse, error) {
	query := `
		SELECT version, COUNT(*) as count 
		FROM verses 
		GROUP BY version 
		ORDER BY version ASC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.VersionResponse
	for rows.Next() {
		var v models.VersionResponse
		if err := rows.Scan(&v.Version, &v.Verses); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}
