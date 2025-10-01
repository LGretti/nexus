# Nexus Service Desk API

## Sobre o Projeto

API em Go para o sistema Nexus, uma plataforma de service desk para controle de contratos, horas e chamados. Este é o serviço de back-end responsável por toda a lógica de negócio e persistência de dados.

**Versão Atual:** CRUD completo para Usuários, Empresas e Contratos, com uma arquitetura refatorada usando generics para melhor manutenibilidade.

## Tecnologias

*   **Linguagem:** Go (Golang)
*   **Roteador:** `net/http` (biblioteca padrão)
*   **Banco de Dados:** PostgreSQL

## Pré-requisitos

*   Go 1.21 ou superior
*   PostgreSQL

## Como Executar

1.  Clone o repositório.
2.  Certifique-se de que o PostgreSQL está em execução e acessível com a string de conexão padrão: `postgres://nexususer:postgres@localhost:5432/nexusdb?sslmode=disable`. O script `table.sql` na raiz pode ser usado para criar as tabelas.
3.  Navegue até a pasta da API: `cd api`
4.  Execute o servidor: `go run ./cmd/api/main.go`
5.  A API estará disponível em `http://localhost:8080`.

## Endpoints Disponíveis

**Nota:** Todos os endpoints que recebem um corpo de requisição esperam o `Content-Type: application/json`. As rotas com `.../` no final requerem a barra (`/`) para funcionar corretamente.

### Usuários

*   `POST /usuarios/`: Cria um novo usuário.
*   `GET /usuarios/`: Lista todos os usuários.
*   `GET /usuarios/{id}`: Busca um usuário específico pelo seu ID.
*   `PUT /usuarios/{id}`: Atualiza os dados de um usuário existente.
*   `DELETE /usuarios/{id}`: Deleta um usuário.

### Empresas

*   `POST /empresas/`: Cria uma nova empresa (com um único objeto JSON) ou múltiplas empresas (com um array de objetos JSON).
*   `GET /empresas/`: Lista todas as empresas.
*   `GET /empresas/{id}`: Busca uma empresa específica pelo seu ID.
*   `PUT /empresas/{id}`: Atualiza os dados de uma empresa existente.
*   `DELETE /empresas/{id}`: Deleta uma empresa.

### Contratos

*   `POST /contratos/`: Cria um novo contrato.
*   `GET /contratos/{id}`: Busca um contrato específico pelo seu ID.
*   `GET /empresas/{id}/contratos`: Lista todos os contratos de uma empresa específica.
*   `PUT /contratos/{id}`: Atualiza um contrato existente.
*   `DELETE /contratos/{id}`: Desativa um contrato (define `ativo` como `false`, não deleta o registro).

## Modelos de Dados (Exemplos)

### Usuario (Request/Response)

```json
{
    "id": 1,
    "nome": "Jules",
    "email": "jules@example.com",
    "perfil": "admin"
}
```

### Empresa (Request/Response)

```json
{
    "id": 1,
    "nome": "Acme Corp",
    "cnpj": "12.345.678/0001-99",
    "email_contato": "contact@acme.com"
}
```

### Contrato (Request/Response)

```json
{
    "id": 1,
    "empresa_id": 1,
    "tipo_contrato": "Retainer",
    "horas_contratadas": 100,
    "data_inicio": "2025-01-01T00:00:00Z",
    "data_fim": "2025-12-31T23:59:59Z",
    "ativo": true
}
```