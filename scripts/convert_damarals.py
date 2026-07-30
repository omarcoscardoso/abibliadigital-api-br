#!/usr/bin/env python3
"""
convert_damarals.py

Script genérico para converter qualquer versão do repositório 'damarals/biblias'
(localizada em data/canonical/<VERSAO>/) para o formato JSON do projeto 'abibliadigital'
(salva em data/json/<idioma>_<versao>.json).

Uso:
    python3 scripts/convert_damarals.py <diretorio_origem> <arquivo_saida>

Exemplo:
    python3 scripts/convert_damarals.py tmp_biblias/data/canonical/NVT data/json/pt_nvt.json
"""

import json
import os
import sys
import glob

BOOKS_FILE = "data/books.json"

def load_books_metadata(books_path=BOOKS_FILE):
    if not os.path.exists(books_path):
        print(f"Erro: Arquivo de metadados dos livros não encontrado em '{books_path}'", file=sys.stderr)
        sys.exit(1)

    with open(books_path, "r", encoding="utf-8") as f:
        books_data = json.load(f)

    # Mapeamento de id para abbrev em português
    id_to_abbrev = {}
    for b in books_data:
        b_id = b["id"]
        abbrev_dict = b.get("abbrev", {})
        if isinstance(abbrev_dict, str):
            abbrev_pt = abbrev_dict
        else:
            abbrev_pt = abbrev_dict.get("pt", "")
        id_to_abbrev[b_id] = abbrev_pt.lower()

    return books_data, id_to_abbrev

def convert_version(src_dir, output_file):
    if not os.path.isdir(src_dir):
        print(f"Erro: Diretório de origem '{src_dir}' não existe.", file=sys.stderr)
        sys.exit(1)

    _, id_to_abbrev = load_books_metadata()

    # Buscar todos os arquivos json do livro (desconsiderando meta.json)
    json_files = sorted(glob.glob(os.path.join(src_dir, "*.json")))
    book_files = [f for f in json_files if not os.path.basename(f).lower() == "meta.json"]

    if not book_files:
        print(f"Erro: Nenhum arquivo JSON de livro encontrado em '{src_dir}'.", file=sys.stderr)
        sys.exit(1)

    # Carregar livros e ordenar pelo ID do livro
    books_by_id = {}
    for fpath in book_files:
        with open(fpath, "r", encoding="utf-8") as f:
            book_obj = json.load(f)

        b_id = book_obj.get("id")
        if b_id is None:
            print(f"Aviso: Ignorando arquivo '{fpath}' pois não possui a chave 'id'.", file=sys.stderr)
            continue

        books_by_id[b_id] = book_obj

    output_data = []

    # Processar na ordem dos 66 livros (IDs 1 a 66)
    for b_id in range(1, 67):
        if b_id not in books_by_id:
            print(f"Aviso: Livro com ID {b_id} não encontrado na versão de origem.", file=sys.stderr)
            continue

        book_src = books_by_id[b_id]
        abbrev_pt = id_to_abbrev.get(b_id, book_src.get("abbrev", "").lower())

        chapters_list = []
        # Garantir ordenação numérica dos capítulos
        raw_chapters = book_src.get("chapters", [])
        sorted_chapters = sorted(raw_chapters, key=lambda c: c.get("number", 0))

        for chap in sorted_chapters:
            raw_verses = chap.get("verses", [])
            sorted_verses = sorted(raw_verses, key=lambda v: v.get("number", 0))
            chapter_verses = [v.get("text", "").strip() for v in sorted_verses]
            chapters_list.append(chapter_verses)

        output_data.append({
            "abbrev": abbrev_pt,
            "chapters": chapters_list
        })

    # Garantir criação do diretório de destino
    os.makedirs(os.path.dirname(os.path.abspath(output_file)), exist_ok=True)

    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(output_data, f, ensure_ascii=False)

    print(f"Conversão concluída com sucesso! Gerado '{output_file}' com {len(output_data)} livros.")

def main():
    if len(sys.argv) < 3:
        print("Uso: python3 scripts/convert_damarals.py <diretorio_origem> <arquivo_saida>")
        print("Exemplo: python3 scripts/convert_damarals.py tmp_biblias/data/canonical/NVT data/json/pt_nvt.json")
        sys.exit(1)

    src_dir = sys.argv[1]
    output_file = sys.argv[2]

    convert_version(src_dir, output_file)

if __name__ == "__main__":
    main()
