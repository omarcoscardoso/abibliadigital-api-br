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
# Executar a suíte de testes em Go
go test -v ./...

# Verificar formatação padrão do código Go
go fmt ./...
```

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
