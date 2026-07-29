# Guia de Deploy e Infraestrutura: ABíbliaDigital em Go

Este documento descreve como rodar localmente, empacotar via Docker e realizar o deploy serverless da API **ABíbliaDigital** no **GCP Cloud Run** integrado com a **Cloudflare**.

---

## 🚀 1. Executando Localmente (sem Docker)

### Pré-requisitos:
- Go 1.22+
- Python 3.x (apenas se for reconstruir o banco `biblia.db` manualmente)

### Passos:
1. Execute o servidor Go:
   ```bash
   go run ./cmd/server/main.go
   ```
   *Nota: O servidor Go verifica a presença de `biblia.db` na inicialização. Caso o arquivo não exista, ele executa a compilação automática das 19 versões a partir de `data/`.*

2. (Opcional) Para reconstruir manualmente a base de dados SQLite com as 19 versões:
   ```bash
   python3 scripts/convert_all_versions.py
   ```

3. Teste a API:
   ```bash
   curl http://localhost:8080/api/check
   ```

---

## 🐳 2. Executando via Docker & Docker Compose

### Usando Docker Compose:
```bash
docker-compose up --build -d
```

### Testando a imagem Docker compilada:
O `Dockerfile` compila o banco `biblia.db` automaticamente durante a construção da imagem Docker (multi-stage build):
```bash
docker build -t abibliadigital:latest .
docker run -p 8080:8080 abibliadigital:latest
```

---

## ☁️ 3. Deploy no GCP Cloud Run (100% Gratuito)

O Cloud Run oferece um **Free Tier generoso** (2 milhões de requisições/mês, 360.000 GB-segundos de memória) e reduz a `0` instâncias quando ociosas.

### Pré-requisitos:
- Conta no Google Cloud Platform (GCP)
- Google Cloud SDK (`gcloud`) instalado e autenticado (`gcloud auth login`)

### Comando Automático:
```bash
./scripts/deploy_cloud_run.sh
```

### Comando Manual:
```bash
# Setar o seu ID do projeto no GCP
gcloud config set project SEU_PROJECT_ID

# Submeter a imagem via Cloud Build
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/abibliadigital:latest .

# Fazer o deploy no Cloud Run
gcloud run deploy abibliadigital \
  --image gcr.io/SEU_PROJECT_ID/abibliadigital:latest \
  --platform managed \
  --region us-central1 \
  --memory 128Mi \
  --min-instances 0 \
  --max-instances 10 \
  --allow-unauthenticated
```

---

## 🌐 4. Configuração na Cloudflare (Domínio Personalizado & CDN)

Para mapear seu domínio (ex: `abibliadigital.api.br`) e obter cache gratuito nas bordas globalmente:

1. **Domínio Customizado no Cloud Run:**
   - Acesse o console do **GCP Cloud Run** -> **Manage Custom Domains**.
   - Adicione o seu domínio (ex: `abibliadigital.api.br`).
   - Copie o registro CNAME fornecido pelo GCP (ex: `ghs.googlehosted.com`).

2. **Configuração de DNS na Cloudflare:**
   - No painel da **Cloudflare**, acesse **DNS** -> **Add record**:
     - **Type:** `CNAME`
     - **Name:** `api` (ou `@`)
     - **Target:** `ghs.googlehosted.com`
     - **Proxy status:** `Proxied` (Ícone de nuvem Laranja 🟠 ativado).

3. **Configuração de SSL/TLS na Cloudflare:**
   - Acesse **SSL/TLS** -> Defina o modo como **Full (strict)**.

4. **Regras de Caching e Cabeçalhos HTTP:**
   - A API Go já emite o cabeçalho `Cache-Control: public, max-age=31536000` em requisições GET com sucesso.
   - Com o proxy da Cloudflare ativado (Nuvem Laranja), as requisições repetidas de leitura de versículos e livros serão servidas diretamente da memória cache da borda da Cloudflare em **< 10ms**, sem consumir cota do Cloud Run!
