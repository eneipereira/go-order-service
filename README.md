# Order Service API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Uma API de microsserviço robusta, escrita em Go, para gerenciar clientes, produtos e pedidos. Este projeto demonstra as melhores práticas de arquitetura de software, incluindo separação de camadas (controllers, services, repositories), tratamento de erros centralizado e documentação de API.

## ✨ Features

- **Gerenciamento de Clientes**: CRUD completo para clientes.
- **Gerenciamento de Produtos**: CRUD completo para produtos.
- **Gerenciamento de Pedidos**: Criação de pedidos com validação de estoque, consulta e gerenciamento de ciclo de vida (Pagar, Cancelar).
- **Arquitetura Limpa**: Código organizado em camadas para alta coesão e baixo acoplamento.
- **Roteamento com `chi`**: Roteador leve, rápido e idiomático.
- **Banco de Dados PostgreSQL**: Interação com o banco de dados usando o driver `pgx`.
- **Tratamento de Erros Centralizado**: Um middleware elegante que captura e formata erros de toda a aplicação.
- **Documentação com Swagger**: Documentação da API gerada automaticamente e disponível interativamente.
- **Configuração via Variáveis de Ambiente**: Suporte para arquivos `.env` para fácil configuração.

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: Go
- **Roteador HTTP**: go-chi/chi
- **Driver PostgreSQL**: jackc/pgx
- **Documentação da API**: swaggo/swag
- **Variáveis de Ambiente**: joho/godotenv

## 🚀 Getting Started

Siga os passos abaixo para configurar e executar o projeto localmente.

### 1. Pré-requisitos

- Go (versão 1.18 ou superior)
- Docker e Docker Compose
- `swag` CLI (para gerar a documentação)
  ```sh
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

### 2. Clonar o Repositório

```sh
git clone https://github.com/eneipereira/go-order-service.git
cd go-order-service
```

### 3. Configuração do Ambiente

Crie um arquivo `.env` na raiz do projeto, baseado no exemplo abaixo, para configurar a conexão com o banco de dados.

**`.env`**

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=orders_service
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable
PORT=8080
```

### 4. Configuração do Banco de Dados

Você pode iniciar uma instância do PostgreSQL usando o Docker Compose. Crie um arquivo `docker-compose.yml` com o seguinte conteúdo:

**`docker-compose.yml`**

```yml
version: "3.8"
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "${POSTGRES_PORT}:5432"
    volumes:
      - ./database/init.sql:/docker-entrypoint-initdb.d/init.sql
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Inicie o container do banco de dados:

```sh
docker-compose up -d
```

Após iniciar o banco, você precisará criar as tabelas. Conecte-se ao banco de dados (usando uma ferramenta como `psql`, DBeaver, etc.) e execute o script SQL em `database/init.sql`.

### 5. Instalar Dependências

```sh
go mod tidy
```

### 6. Executar a Aplicação

```sh
go run ./cmd/api/main.go
```

O servidor estará rodando em `http://localhost:8080`.

## 📖 Documentação da API (Swagger)

Após iniciar a aplicação, a documentação interativa da API estará disponível em:

**http://localhost:8080/swagger/index.html**

Para regenerar a documentação após fazer alterações nas anotações do código, execute o seguinte comando na raiz do projeto:

```sh
swag init -g cmd/api/main.go
```

## 🛣️ Endpoints da API

### Testes

O projeto possui uma suíte de testes unitários e de integração para garantir a qualidade e a estabilidade do código. Atualmente, o projeto possui 43.7% de cobertura global de testes, cobrindo as camadas de `model`, `service` e `controller`

Posteriormente, serão adicionados os testes para as demais camadas.

Para executar todos os testes, use o seguinte comando:

```sh
go test -v ./...
```

#### Cobertura de Testes

Para gerar um relatório de cobertura, primeiro execute o comando de teste com a flag `-coverprofile`:

```sh
go test -coverprofile=coverage.out ./...
```

Isso criará um arquivo `coverage.out`. Com ele, você pode:

1.  **Ver a porcentagem de cobertura total no terminal:**
    ```sh
    go tool cover -func=coverage.out
    ```
2.  **Gerar um relatório visual em HTML:**
    ```sh
    go tool cover -html=coverage.out -o coverage.html

```

### Health Check

- `GET /`: Verifica se a API está no ar.
- `GET /health`: Verifica se a API está no ar.

### Customers

- `POST /customers`: Cria um novo cliente.
- `GET /customers`: Lista todos os clientes (com paginação).
- `GET /customers/{id}`: Obtém um cliente por ID.

### Products

- `POST /products`: Cria um novo produto.
- `GET /products`: Lista todos os produtos (com paginação).
- `GET /products/{id}`: Obtém um produto por ID.

### Orders

- `POST /orders`: Cria um novo pedido.
- `GET /orders`: Lista todos os pedidos (com paginação).
- `GET /orders/{id}`: Obtém um pedido por ID.
- `POST /orders/{id}/pay`: Marca um pedido como pago.
- `POST /orders/{id}/cancel`: Cancela um pedido e reabastece o estoque.

---

**Autor**: Enei Pereira
```
