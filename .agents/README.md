# Agent Skills (.agents)

Este diretório contém o conjunto de **Agent Skills** (habilidades e instruções estruturadas para assistentes de IA como Antigravity, Cursor, Claude Code, Copilot, OpenCode e Gemini CLI) utilizadas no desenvolvimento da API **ABíbliaDigital**.

---

## 📌 Origem das Habilidades

- **Repositório Origem:** [`samber/cc-skills-golang`](https://github.com/samber/cc-skills-golang)
- **Autor Original:** Samuel Berthe ([@samber](https://github.com/samber))
- **Licença:** MIT

As habilidades fornecem diretrizes idiomáticas, padrões de concorrência, boas práticas de segurança, observabilidade (`slog`), tratamento de erros e testes automatizados otimizados para projetos modernos em **Go (1.22+)**.

---

## 📁 Estrutura das Skills

Para manter o repositório leve, limpo e livre de artefatos de avaliação ou sub-repositórios `.git`, extraímos apenas as especificações das habilidades para:

```text
.agents/skills/
├── golang-code-style/
├── golang-concurrency/
├── golang-context/
├── golang-database/
├── golang-error-handling/
├── golang-observability/
├── golang-performance/
├── golang-security/
├── golang-testing/
└── ... (46 habilidades disponíveis)
```

Cada pasta contém um arquivo `SKILL.md` principal e, opcionalmente, uma subpasta `references/` com exemplos e diretrizes aprofundadas.

---

## 🔄 Como Atualizar ou Sincronizar

Caso queira atualizar as habilidades com a versão mais recente do repositório de origem:

```bash
# Atualizar extraindo as pastas de skills do repositório original
mkdir -p /tmp/cc-skills-golang
git clone https://github.com/samber/cc-skills-golang.git /tmp/cc-skills-golang
cp -r /tmp/cc-skills-golang/skills/* .agents/skills/
rm -rf /tmp/cc-skills-golang
```
