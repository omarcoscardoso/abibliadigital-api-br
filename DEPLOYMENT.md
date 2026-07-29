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

### Deploy Manual:
```bash
# Setar o seu ID do projeto no GCP
gcloud config set project SEU_PROJECT_ID

# Submeter a imagem via Cloud Build
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/abibliadigital:latest .

# Fazer o deploy no Cloud Run
gcloud run deploy abibliadigital-api-br \
  --image gcr.io/SEU_PROJECT_ID/abibliadigital:latest \
  --platform managed \
  --region us-central1 \
  --memory 128Mi \
  --min-instances 0 \
  --max-instances 10 \
  --allow-unauthenticated
```

---

## ⚡ 4. Deploy Automático via Cloud Build (CI/CD)

O arquivo [`cloudbuild.yaml`](file:///home/cardoso/projetos/abibliadigital/cloudbuild.yaml) orquestra a compilação, o envio da imagem para o **Container Registry (`gcr.io`)** e o deploy no **Cloud Run** (`abibliadigital-api-br`).

### Pré-requisitos e Checklist do GCP:

1. **Ativar as APIs necessárias no GCP:**
   ```bash
   gcloud services enable \
     cloudbuild.googleapis.com \
     run.googleapis.com \
     containerregistry.googleapis.com
   ```

2. **Conceder permissões para a Service Account do Cloud Build:**
   A conta de serviço do Cloud Build (`[PROJECT_NUMBER]@cloudbuild.gserviceaccount.com`) precisa das seguintes permissões no IAM:
   - **Cloud Run Admin** (`roles/run.admin`)
   - **Service Account User** (`roles/iam.serviceAccountUser`)

   *Comandos para conceder as permissões via `gcloud`:*
   ```bash
   PROJECT_NUMBER=$(gcloud projects describe $(gcloud config get-value project) --format='value(projectNumber)')

   # Conceder Cloud Run Admin
   gcloud projects add-iam-policy-binding $(gcloud config get-value project) \
     --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
     --role="roles/run.admin"

   # Conceder Service Account User
   gcloud projects add-iam-policy-binding $(gcloud config get-value project) \
     --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
     --role="roles/iam.serviceAccountUser"

   # Conceder Artifact Registry Writer
   gcloud projects add-iam-policy-binding $(gcloud config get-value project) \
     --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
     --role="roles/artifactregistry.writer"
   ```

4. **Criar o Trigger no Cloud Build Console:**
   - Acesse **GCP Console** -> **Cloud Build** -> **Triggers**.
   - Clique em **Create Trigger**.
   - **Name:** `deploy-abibliadigital-api-br`
   - **Event:** `Push to a branch` (branch `^main$`).
   - **Source:** Seu repositório GitHub (`omarcoscardoso/abibliadigital-api-br`).
   - **Configuration:** `Cloud Build configuration file (yaml or json)` -> `cloudbuild.yaml`.
   - Salve a Trigger. Qualquer `git push` para a branch `main` disparará o deploy automático!

---

## 🌐 5. Configuração Completa na Cloudflare (Edge Cache & WAF)

Para integrar a **Cloudflare** como camada de Edge Cache e WAF Gratuito na frente do **GCP Cloud Run**:

### 1. Configuração de DNS e Proxying (Nuvem Laranja 🟠)

#### Opção A: Domínio Personalizado Mapeado no Cloud Run (Recomendado)
1. No console do **GCP Cloud Run**, vá em **Custom Domains** -> **Add Mapping** e registre seu domínio/subdomínio (ex: `api.abibliadigital.com.br`).
2. O GCP gerará um destino CNAME (ex: `ghs.googlehosted.com`).
3. No painel da **Cloudflare** -> **DNS** -> **Records**, crie o registro:
   - **Type:** `CNAME`
   - **Name:** `api` (ou `@` para o domínio raiz)
   - **Target:** `ghs.googlehosted.com`
   - **Proxy status:** `Proxied` (Nuvem Laranja 🟠 **ativada**)

#### Opção B: CNAME Direto para a URL Padrão do Cloud Run
1. No painel da **Cloudflare** -> **DNS** -> **Records**:
   - **Type:** `CNAME`
   - **Name:** `api`
   - **Target:** `abibliadigital-api-br-xxxxx-uc.a.run.app` (URL gerada pelo Cloud Run)
   - **Proxy status:** `Proxied` (Nuvem Laranja 🟠 **ativada**)

---

### 2. Configuração de SSL/TLS (Criptografia de Ponta a Ponta)
- No painel da Cloudflare, acesse **SSL/TLS** -> **Overview**.
- Selecione a opção **Full (strict)**.
- Isso garante conexão segura (HTTPS) criptografada entre o Usuário -> Cloudflare -> GCP Cloud Run.

---

### 3. Criar Regra de Cache (Cache Rules) — CRUCIAL ⚡
> [!IMPORTANT]
> Por padrão, a Cloudflare **NÃO** faz cache de respostas JSON/API apenas com o cabeçalho `Cache-Control`. É **obrigatório** criar uma Cache Rule para ativar o Edge Cache em rotas da API.

1. No painel da Cloudflare, acesse **Caching** -> **Cache Rules** -> **Create Rule**.
2. **Rule name:** `Cache API GET Requests`
3. **When incoming requests match:**
   - Campo: `URI Path` -> Operador: `starts with` -> Valor: `/api/`
   - **AND**
   - Campo: `Request Method` -> Operador: `equals` -> Valor: `GET`
4. **Cache eligibility:** Selecione `Eligible for cache` (Cache Everything).
5. **Edge TTL:** Selecione `Respect origin` (a API Go já envia `Cache-Control: public, max-age=31536000`).
6. **Browser TTL:** Selecione `Respect origin`.
7. Clique em **Deploy**.

Com esta regra ativa:
- **1ª requisição (Cache Miss):** Requisita ao GCP Cloud Run, lê o banco SQLite e guarda o resultado na borda da Cloudflare.
- **Requisições subsequentes (Cache Hit):** Respondidas em **< 15ms** direto da rede global da Cloudflare, sem gastar cota do GCP nem acordar o container!

---

### 4. WAF Gratuito e Segurança na Borda
1. **Bot Fight Mode:**
   - Acesse **Security** -> **Bots** -> Ative o **Bot Fight Mode** (bloqueia bots maliciosos conhecidos e requisições automatizadas abusivas).
2. **Rate Limiting (Plano Gratuito):**
   - Acesse **Security** -> **WAF** -> **Rate limiting rules** -> **Create rule**.
   - **Name:** `API Rate Limit`
   - **Match:** `URI Path starts with /api/`
   - **Rate limit:** `100 requisições` por `1 minuto` por IP.
   - **Action:** `Block` ou `Managed Challenge`.
3. **Otimizações Globais:**
   - Acesse **Speed** -> **Optimization**:
     - Ative **Brotli** (compressão ultraeficiente para respostas JSON).
     - Ative **HTTP/3 (with QUIC)** para menor latência em apps móveis.

