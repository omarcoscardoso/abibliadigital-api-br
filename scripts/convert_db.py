#!/usr/bin/env python3
import json
import sqlite3
import os
import sys

DB_PATH = "biblia.db"

def init_db(conn):
    cursor = conn.cursor()
    
    # Table books
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        abbrev_pt TEXT NOT NULL UNIQUE,
        abbrev_en TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        author TEXT NOT NULL DEFAULT '',
        chapters INTEGER NOT NULL DEFAULT 0,
        group_name TEXT NOT NULL DEFAULT '',
        testament TEXT NOT NULL DEFAULT '',
        comment TEXT NOT NULL DEFAULT '',
        book_order INTEGER NOT NULL DEFAULT 0
    );
    """)
    
    # Table verses
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS verses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        book_id INTEGER NOT NULL REFERENCES books(id),
        book_abbrev_pt TEXT NOT NULL,
        book_abbrev_en TEXT NOT NULL,
        chapter INTEGER NOT NULL,
        number INTEGER NOT NULL,
        version TEXT NOT NULL,
        text TEXT NOT NULL,
        comment TEXT NOT NULL DEFAULT ''
    );
    """)

    # FTS5 for full-text search
    cursor.execute("""
    CREATE VIRTUAL TABLE IF NOT EXISTS verses_fts USING fts5(
        verse_id UNINDEXED,
        version UNINDEXED,
        book_abbrev_pt UNINDEXED,
        book_abbrev_en UNINDEXED,
        chapter UNINDEXED,
        number UNINDEXED,
        text,
        tokenize='unicode61 remove_diacritics 2'
    );
    """)

    # Indexes
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_books_abbrev_pt ON books(abbrev_pt);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_books_abbrev_en ON books(abbrev_en);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_lookup_pt ON verses(version, book_abbrev_pt, chapter, number);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_lookup_en ON verses(version, book_abbrev_en, chapter, number);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_versao_livro_capitulo_pt ON verses(version, book_abbrev_pt, chapter);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_versao_livro_capitulo_en ON verses(version, book_abbrev_en, chapter);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_version ON verses(version);")

    conn.commit()

def load_data(conn, books_file, verses_file):
    cursor = conn.cursor()

    if not os.path.exists(books_file):
        books_file = "__test__/mock/books.json"
    if not os.path.exists(verses_file):
        verses_file = "__test__/mock/verses.json"

    print(f"Reading books from: {books_file}")
    print(f"Reading verses from: {verses_file}")

    with open(books_file, "r", encoding="utf-8") as f:
        books_data = json.load(f)

    with open(verses_file, "r", encoding="utf-8") as f:
        verses_data = json.load(f)

    print(f"Loading {len(books_data)} books...")
    book_id_map = {}
    for idx, b in enumerate(books_data, start=1):
        abbrev_dict = b.get("abbrev", {})
        if isinstance(abbrev_dict, str):
            abbrev_pt = abbrev_dict
            abbrev_en = abbrev_dict
        else:
            abbrev_pt = abbrev_dict.get("pt", "")
            abbrev_en = abbrev_dict.get("en", abbrev_pt)

        name = b.get("name", "")
        author = b.get("author", "")
        chapters = b.get("chapters", 0)
        group_name = b.get("group", "")
        testament = b.get("testament", "")
        comment = b.get("comment", "")
        order = b.get("order", idx)

        cursor.execute("""
            INSERT OR REPLACE INTO books 
            (abbrev_pt, abbrev_en, name, author, chapters, group_name, testament, comment, book_order)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (abbrev_pt, abbrev_en, name, author, chapters, group_name, testament, comment, order))

        book_id = cursor.lastrowid
        book_id_map[book_id] = (book_id, abbrev_pt, abbrev_en)
        book_id_map[abbrev_pt] = (book_id, abbrev_pt, abbrev_en)
        book_id_map[abbrev_en] = (book_id, abbrev_pt, abbrev_en)

    print(f"Loading {len(verses_data)} verses...")
    verses_to_insert = []

    for v in verses_data:
        # Check if verse has book_id or nested book object
        b_id = v.get("book_id")
        abbrev_pt = ""
        abbrev_en = ""

        if b_id and b_id in book_id_map:
            book_id, abbrev_pt, abbrev_en = book_id_map[b_id]
        else:
            book_obj = v.get("book", {})
            abbrev_dict = book_obj.get("abbrev", {})
            if isinstance(abbrev_dict, str):
                abbrev_pt = abbrev_dict
                abbrev_en = abbrev_dict
            else:
                abbrev_pt = abbrev_dict.get("pt", "")
                abbrev_en = abbrev_dict.get("en", "")

            book_info = book_id_map.get(abbrev_pt) or book_id_map.get(abbrev_en)
            if book_info:
                book_id, abbrev_pt, abbrev_en = book_info
            else:
                book_id = 0

        chapter = v.get("chapter", 0)
        number = v.get("number", 0)
        version = v.get("version", "nvi")
        text = v.get("text", "")
        comment = v.get("comment", "")

        verses_to_insert.append((book_id, abbrev_pt, abbrev_en, chapter, number, version, text, comment))

    cursor.executemany("""
        INSERT INTO verses (book_id, book_abbrev_pt, book_abbrev_en, chapter, number, version, text, comment)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    """, verses_to_insert)

    # Populate FTS5
    cursor.execute("""
        INSERT INTO verses_fts (verse_id, version, book_abbrev_pt, book_abbrev_en, chapter, number, text)
        SELECT id, version, book_abbrev_pt, book_abbrev_en, chapter, number, text FROM verses
    """)

    conn.commit()
    print(f"Database conversion successfully completed! Inserted {len(books_data)} books and {len(verses_data)} verses.")

def main():
    books_file = sys.argv[1] if len(sys.argv) > 1 else "raw_data/books_with_en.json"
    verses_file = sys.argv[2] if len(sys.argv) > 2 else "raw_data/verses.json"

    if os.path.exists(DB_PATH):
        os.remove(DB_PATH)

    conn = sqlite3.connect(DB_PATH)
    init_db(conn)
    load_data(conn, books_file, verses_file)
    conn.close()

if __name__ == "__main__":
    main()
