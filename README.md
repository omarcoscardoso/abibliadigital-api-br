# ABíbliaDigital (Go + SQLite Embedded)

API RESTful moderna, estateless e de altíssima performance para a Bíblia Sagrada escrita em **Go** com banco de dados **SQLite embutido** em modo *read-only*.

O projeto conta com **66 livros** e **781.508 versículos** distribuídos em **26 versões** (incluindo NVI, NVT, NTLH, ACF, AA, BLIVRE, ALM1911, OL, MENS, AS21, KJV, BBE, RVR, APEE e outras em 14 idiomas).

---

## 🙏 Agradecimentos e Contexto do Projeto

Este projeto nasceu como uma iniciativa de continuidade e gratidão ao projeto original [ABíbliaDigital](https://github.com/omarciovsena/abibliadigital), criado e mantido brilhantemente por [Márcio Sena](https://github.com/omarciovsena) e colaboradores ao longo de quase 8 anos.

Após o anúncio do encerramento das APIs do projeto original em 01/08/2026 motivado pelos altos custos de infraestrutura e tempo de manutenção, esta reescrita em **Go + SQLite Embedded** foi desenvolvida com o objetivo de manter esta valiosa ferramenta viva e acessível para todos, de forma **100% gratuita, stateless, de alta performance e pronta para deploy serverless com custo zero**.

Expressamos nossa profunda gratidão ao Márcio Sena e a todos os contribuidores do projeto original pela dedicação e impacto gerado na comunidade!

Agradecemos também ao projeto open-source [damarals/biblias](https://github.com/damarals/biblias) de [Daniel Amaral](https://github.com/damarals), repositório fundamental que serviu de fonte de dados para diversas traduções em português incorporadas (como NVT, NTLH, Bíblia Livre, Almeida 1911, O Livro, A Mensagem, Almeida Século 21 e atualização da NVI).

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

## 🤖 Agent Skills para Desenvolvimento (Go)

Este repositório conta com um conjunto de **46 Agent Skills** em [`.agents/skills/`](/.agents/README.md) baseadas nas diretrizes de [`samber/cc-skills-golang`](https://github.com/samber/cc-skills-golang). Elas orientam assistentes de IA (como Antigravity, Cursor e Claude Code) a seguirem automaticamente as melhores práticas de Go (concorrência, segurança, testes, `slog`, etc.) durante o desenvolvimento.

---

## 📄 Licença

Este projeto está sob a licença **MIT** - veja o arquivo [`LICENSE.md`](/LICENSE.md) para mais detalhes.
