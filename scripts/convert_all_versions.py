#!/usr/bin/env python3
import json
import sqlite3
import os
import sys
import glob

DB_PATH = "biblia.db"

VERSION_MAPPING = {
    "pt_nvi.json": "nvi",
    "pt_acf.json": "acf",
    "pt_aa.json": "aa",
    "en_kjv.json": "kjv",
    "en_bbe.json": "bbe",
    "es_rvr.json": "rvr",
    "fr_apee.json": "apee",
    "de_schlachter.json": "schlachter",
    "el_greek.json": "greek",
    "eo_esperanto.json": "esperanto",
    "fi_finnish.json": "finnish",
    "fi_pr.json": "pr",
    "ko_ko.json": "ko",
    "ro_cornilescu.json": "cornilescu",
    "ru_synodal.json": "synodal",
    "vi_vietnamese.json": "vietnamese",
    "zh_cuv.json": "cuv",
    "zh_ncv.json": "ncv",
    "ar_svd.json": "svd",
}

def init_db(conn):
    cursor = conn.cursor()
    
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

    cursor.execute("CREATE INDEX IF NOT EXISTS idx_books_abbrev_pt ON books(abbrev_pt);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_books_abbrev_en ON books(abbrev_en);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_lookup_pt ON verses(version, book_abbrev_pt, chapter, number);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_lookup_en ON verses(version, book_abbrev_en, chapter, number);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_versao_livro_capitulo_pt ON verses(version, book_abbrev_pt, chapter);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_versao_livro_capitulo_en ON verses(version, book_abbrev_en, chapter);")
    cursor.execute("CREATE INDEX IF NOT EXISTS idx_verses_version ON verses(version);")

    conn.commit()

def populate_database(conn):
    cursor = conn.cursor()

    books_file = "data/books.json"
    if not os.path.exists(books_file):
        print(f"Error: {books_file} not found.")
        sys.exit(1)

    print(f"Loading books metadata from {books_file}...")
    with open(books_file, "r", encoding="utf-8") as f:
        books_data = json.load(f)

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
        book_id_map[abbrev_pt.lower()] = (book_id, abbrev_pt, abbrev_en)
        book_id_map[abbrev_en.lower()] = (book_id, abbrev_pt, abbrev_en)

    json_dir = "data/json"
    json_files = sorted(glob.glob(os.path.join(json_dir, "*.json")))
    
    total_verses_inserted = 0

    for filepath in json_files:
        filename = os.path.basename(filepath)
        version_code = VERSION_MAPPING.get(filename, filename.replace(".json", ""))

        print(f"Processing version: '{version_code}' from {filename}...")
        with open(filepath, "r", encoding="utf-8-sig") as f:
            version_data = json.load(f)

        verses_to_insert = []
        for book_obj in version_data:
            abbrev = book_obj.get("abbrev", "").lower()
            book_info = book_id_map.get(abbrev)

            if not book_info:
                # Try fallback matching
                for key, val in book_id_map.items():
                    if key in abbrev or abbrev in key:
                        book_info = val
                        break

            if not book_info:
                book_id, b_pt, b_en = 0, abbrev, abbrev
            else:
                book_id, b_pt, b_en = book_info

            chapters = book_obj.get("chapters", [])
            for c_idx, chapter_verses in enumerate(chapters, start=1):
                for v_idx, text in enumerate(chapter_verses, start=1):
                    if isinstance(text, str):
                        verses_to_insert.append((book_id, b_pt, b_en, c_idx, v_idx, version_code, text, ""))

        cursor.executemany("""
            INSERT INTO verses (book_id, book_abbrev_pt, book_abbrev_en, chapter, number, version, text, comment)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """, verses_to_insert)

        print(f"  -> Inserted {len(verses_to_insert)} verses for version '{version_code}'")
        total_verses_inserted += len(verses_to_insert)

    print("Building FTS5 index...")
    cursor.execute("""
        INSERT INTO verses_fts (verse_id, version, book_abbrev_pt, book_abbrev_en, chapter, number, text)
        SELECT id, version, book_abbrev_pt, book_abbrev_en, chapter, number, text FROM verses
    """)

    conn.commit()
    print(f"Conversion complete! Total: {len(books_data)} books, {total_verses_inserted} verses across {len(json_files)} versions.")

def main():
    if os.path.exists(DB_PATH):
        os.remove(DB_PATH)

    conn = sqlite3.connect(DB_PATH)
    init_db(conn)
    populate_database(conn)
    conn.close()

if __name__ == "__main__":
    main()
