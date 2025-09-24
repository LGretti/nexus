# Nexus Service Desk API

## Sobre o Projeto

API em Go para o sistema Nexus, uma plataforma de service desk para controle de contratos, horas e chamados. Este é o serviço de back-end responsável por toda a lógica de negócio e persistência de dados.

**Versão Inicial (MVP):** Foco no CRUD de Empresas e seus Contratos de horas.

## Tecnologias

*   **Linguagem:** Go (Golang)
*   **Roteador:** `net/http` (biblioteca padrão)
*   **Banco de Dados:** PostgreSQL (a ser confirmado)

## Pré-requisitos

*   Go 1.21 ou superior
*   PostgreSQL (ou o banco de dados de sua escolha)

## Como Executar

1.  Clone o repositório.
2.  Configure as variáveis de ambiente em um arquivo `.env` (ex: string de conexão com o banco).
3.  Navegue até a pasta da API: `cd api`
4.  Execute o servidor: `go run ./cmd/api`
5.  A API estará disponível em `http://localhost:8080`.

## Endpoints Disponíveis

### Empresas

*   `POST /empresas`: Cria uma nova empresa.
*   `POST /empresas/lote`: Cria múltiplas empresas em uma única requisição.
*   `GET /empresas`: Lista todas as empresas.
*   `GET /empresas/{id}`: Busca uma empresa específica pelo seu ID.
*   `PUT /empresas/{id}`: Atualiza os dados de uma empresa existente.
*   `DELETE /empresas/{id}`: Deleta uma empresa.

### Contratos

*   `POST /empresas/{id}/contratos`: Cria um novo contrato para uma empresa.
*   `GET /empresas/{id}/contratos`: Lista todos os contratos de uma empresa.
*   `PUT /contratos/{id}`: Atualiza um contrato existente.
*   `DELETE /contratos/{id}`: Deleta um contrato.

## Modelos de Dados (Exemplos)

### Empresa

```json
{
    "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "nome_fantasia": "Empresa Exemplo",
    "razao_social": "Empresa Exemplo LTDA",
    "cnpj": "12.345.678/0001-99"
}
```

### Contrato

```json
{
    "id": "f0e9d8c7-b6a5-4321-fedc-ba9876543210",
    "empresa_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "horas_contratadas": 50,
    "valor_hora": 150.50,
    "data_inicio": "2025-01-01T00:00:00Z",
    "data_fim": "2025-12-31T23:59:59Z",
    "ativo": true
}
```