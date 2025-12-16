
# Nexus Service Desk

Sistema completo de Service Desk para gestão de consultoria, controle de contratos e apontamento de horas. O projeto consiste em uma API robusta em Go e um Frontend moderno em Next.js.

## Tecnologias

### Backend (API)
*  **Linguagem:** Go (Golang)
*  **Router:** chi (go-chi)
*  **Banco de Dados:** PostgreSQL
*  **Arquitetura:** Camadas (Handlers, Repositories, Models) com injeção de dependência.

### Frontend (Web)
*  **Framework:** Next.js 14 (App Router)
*  **Linguagem:** TypeScript
*  **Estado/Cache:** TanStack Query (React Query)
*  **Estilização:** Tailwind CSS
  

## Como Executar

### Pré-requisitos

* Go 1.21+
* Node.js 18+
* PostgreSQL rodando (Database: `nexusdb`)  

### Modo Automático
Na raiz do projeto, apenas execute o script:
``` bash
./dev.bat
```

Isso abrirá duas janelas de terminal: uma para o Backend (porta 8080) e outra para o Frontend (porta 3000).

### Modo Manual
1.  **Backend:**
``` bash
cd api
go run cmd/api/main.go
```

2. **Frontend:**
``` bash
cd  frontend
npm  install
npm  run  dev
```

Acesse o sistema em: http://localhost:3000


## API Endpoints

A API roda em `http://localhost:8080` e todas as rotas são prefixadas com `/api`.

### Empresas (Companies)
| **Método** | **Rota** | **Descrição** |
|--|--|--|
| `GET` | `/api/companies` | Lista todas as empresas |
| `POST` | `/api/companies` | Cadastra nova empresa |
| `GET` | `/api/companies/{id}` | Detalhes da empresa |
| `PUT` | `/api/companies/{id}` | Atualiza empresa |
| `DELETE` | `/api/companies/{id}` | Remove empresa |

### Contratos (Contracts)
| **Método** | **Rota** | **Descrição** |
|--|--|--|
| `GET` | `/api/contracts` | Lista contratos (inclui nome da empresa) |
| `POST` | `/api/contracts` | Cria contrato vinculado a uma empresa |
| `GET` | `/api/contracts/{id}` | Detalhes do contrato |
| `GET` | `/api/contracts/{id}/appointments` | **Relatório:** Atendimentos deste contrato |

### Usuários (Users)
| **Método** | **Rota** | **Descrição** |
|--|--|--|
| `GET` | `/api/users` | Lista consultores e admins |
| `POST` | `/api/users` | Cadastra usuário |
| `GET` | `/api/users/{id}/appointments` | **Produtividade:** Horas deste consultor |

### Apontamentos (Appointments)
| **Método** | **Rota** | **Descrição** |
|--|--|--|
| `POST` | `/api/appointments` | Lança horas (Start/End Time) |
| `GET` | `/api/appointments` | Visão Geral (Admin) |
| `DELETE` | `/api/appointments/{id}` | Remove lançamento |



## Modelos de Dados (JSON)

A API utiliza `camelCase` para compatibilidade com o Frontend.

**Empresa**
``` json
{
    "id": 1,
    "name": "Ademicon",
    "cnpj": "12.345.678/0001-90",
    "email": "contato@ademicon.com.br"
}
```

**Contrato**
``` json
{
    "id": 10,
    "companyId": 1,
    "companyName": "Ademicon",
    "title": "Ademicon - Suporte",
    "contractType": "Mensal",
    "totalHours": 100,
    "isActive": true
}
```


**Atendimento**
``` json
{
    "contractId": 10,
    "userId": 5,
    "startTime": "2025-12-16T08:00:00Z",
    "endTime": "2025-12-16T12:00:00Z",
    "description": "Correção de bug crítico",
    "totalHours": 4.0
}
```