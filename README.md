# ABíbliaDigital (Go + SQLite Embedded)

API RESTful moderna, estateless e de altíssima performance para a Bíblia Sagrada escrita em **Go** com banco de dados **SQLite embutido** em modo *read-only*.

O projeto conta com **66 livros** e **582.270 versículos** distribuídos em **19 versões** (incluindo NVI, ACF, AA, KJV, BBE, RVR, APEE e outras).

---

## 🙏 Agradecimentos e Contexto do Projeto

Este projeto nasceu como uma iniciativa de continuidade e gratidão ao projeto original [ABíbliaDigital](https://github.com/omarciovsena/abibliadigital), criado e mantido brilhantemente por [Márcio Sena](https://github.com/omarciovsena) e colaboradores ao longo de quase 8 anos.

Após o anúncio do encerramento das APIs do projeto original em 01/08/2026 motivado pelos altos custos de infraestrutura e tempo de manutenção, esta reescrita em **Go + SQLite Embedded** foi desenvolvida com o objetivo de manter esta valiosa ferramenta viva e acessível para todos, de forma **100% gratuita, stateless, de alta performance e pronta para deploy serverless com custo zero**.

Expressamos nossa profunda gratidão ao Márcio Sena e a todos os contribuidores do projeto original pela dedicação e impacto gerado na comunidade!

---

## ⚡ Principais Características

- **Ultra-Rápido**: Respostas em tempo sub-milissegundo (~700µs) graças ao SQLite embedded.
- **100% Stateless**: Sem necessidade de containers MongoDB, Redis ou servidores externos.
- **Busca de Alta Performance**: Suporte a busca textual com SQLite **FTS5** (insensível a acentos).
- **Pronto para Serverless**: Otimizado para deploy gratuito no **GCP Cloud Run** (cota 0 instâncias) e CDN **Cloudflare**.
- **Licença MIT**: 100% Open Source.

---

## 🚀 Como Executar Localmente

### Pré-requisitos
- Go 1.22 ou superior

### Comandos
```bash
# Clone o repositório
git clone https://github.com/omarcoscardoso/abibliadigital-api-br.git
cd abibliadigital-api-br

# Inicie o servidor
go run ./cmd/server/main.go
```

O servidor estará rodando em `http://localhost:8080` com a **Landing Page** na raiz `/` e os endpoints da API em `/api/...`.

---

## 🐳 Executando com Docker

```bash
# Subir ambiente local com Docker Compose
docker-compose up --build
```

---

## 📖 Endpoints Principais da API

| Método | Endpoint | Descrição |
|---|---|---|
| `GET` | `/api/check` | Status de saúde da API (`{"result":"success"}`) |
| `GET` | `/api/books` | Lista os 66 livros da Bíblia |
| `GET` | `/api/books/:abbrev` | Detalhes de um livro específico (`gn`, `ex`, `mt`, etc.) |
| `GET` | `/api/versions` | Lista as 19 versões disponíveis e contagem de versículos |
| `GET` | `/api/verses/:version/:abbrev/:chapter` | Retorna todos os versículos de um capítulo |
| `GET` | `/api/verses/:version/:abbrev/:chapter/:number` | Retorna um versículo específico |
| `GET` | `/api/verses/:version/random` | Retorna um versículo aleatório |
| `POST` | `/api/verses/search` | Busca por palavras-chave (`{"version":"nvi","search":"Deus"}`) |

---

## ☁️ Deploy no GCP Cloud Run

Para implantar no **GCP Cloud Run**:

```bash
./scripts/deploy_cloud_run.sh
```

Consulte o guia completo de deploy e integração com a Cloudflare em [`DEPLOYMENT.md`](/DEPLOYMENT.md).

---

## 📄 Licença

Este projeto está sob a licença **MIT** - veja o arquivo [`LICENSE.md`](/LICENSE.md) para mais detalhes.
