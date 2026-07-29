# Guia de Contribuição: ABíbliaDigital (Go)

Agradecemos o seu interesse em contribuir com a **ABíbliaDigital**! Este documento orienta como colaborar com o projeto seguindo as boas práticas do ecossistema Go.

---

## 📋 Diretrizes Gerais

- **Reportando Issues**: A lista de issues destina-se a relatos de bugs e sugestões de melhorias. Descreva de forma clara os passos para reproduzir qualquer problema.
- **Pull Requests**: Todas as contribuições devem ser enviadas via Pull Request direcionadas à branch `main`.

---

## 🛠️ Ambiente de Desenvolvimento e Comandos

### Pré-requisitos
- Go 1.22 ou superior

### Validação Local de Código
Antes de abrir um Pull Request, certifique-se de que o código passa nos testes automatizados e segue a formatação padrão em Go:

```bash
# Executar a suíte de testes em Go (necessita do banco biblia.db gerado na raiz)
go test -v ./...

# Verificar formatação padrão do código Go
go fmt ./...
```

### 🤖 Agent Skills & Assistentes de IA
O repositório inclui um pacote de **Agent Skills** em [`.agents/skills/`](/.agents/README.md) derivadas de [`samber/cc-skills-golang`](https://github.com/samber/cc-skills-golang). Caso utilize editores e ferramentas com suporte a agentes de IA (como Antigravity, Cursor ou Claude Code), as diretrizes de concorrência, tratamento de erros e segurança em Go serão lidas e aplicadas automaticamente durante a edição de código.

### Validação Automática (CI/CD)
Após abrir ou atualizar um Pull Request para a branch `main`, o GitHub Actions executará automaticamente o workflow **CI - Run Go Tests**. 

Este workflow irá:
1. Baixar as dependências do projeto.
2. Gerar o banco de dados temporário `biblia.db` a partir das bases JSON utilizando `scripts/convert_all_versions.py`.
3. Executar o comando `go test -v ./...`.

**Nota:** Os testes devem obrigatoriamente passar com sucesso no CI antes que o Pull Request possa ser aprovado e integrado à branch `main`.

---

## 📝 Convenção de Commits (Conventional Commits)

Todas as mensagens de commit devem seguir rigorosamente o padrão **Conventional Commits** em **Português (Brasil)**:

`<tipo>(<escopo>): <descrição em minúsculas sem ponto final>`

### Tipos Permitidos:
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Alterações apenas em documentação
- `style`: Formatação, espaços em branco ou estilo sem alterar lógica
- `refactor`: Refatoração de código
- `perf`: Melhoria de desempenho
- `test`: Adição ou correção de testes
- `chore`: Tarefas de manutenção e dependências

### Exemplos:
- `feat(api): adicionar suporte ao formato XML`
- `fix(db): corrigir ordenacao dos versiculos no capitulo`
- `docs(readme): atualizar instruções de execução local`
- `chore(deps): atualizar versão do driver sqlite`

---

## 🌿 Fluxo de Branches

1. Crie uma branch a partir da `main`:
   - `feature/nome-da-funcionalidade`
   - `bugfix/descricao-do-bug`
   - `chore/manutencao`
2. Garanta que todos os testes passem localmente.
3. Submeta o Pull Request descrevendo claramente as alterações realizadas.
