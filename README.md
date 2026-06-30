# Microsserviços com gRPC

Projeto de microsserviços desenvolvido para a disciplina de Programação Distribuída do IFPB, aplicando arquitetura hexagonal e comunicação via gRPC.

## Arquitetura

```
Cliente → Order (porta 3000)
              → Payment (porta 3001)
              → Shipping (porta 3002)
```

Cada microsserviço possui seu próprio banco de dados MySQL.

## Microsserviços

| Serviço | Porta | Descrição |
|---|---|---|
| Order | 3000 | Recebe pedidos, valida estoque, coordena pagamento e entrega |
| Payment | 3001 | Processa cobranças (limite de R$1000 por pedido) |
| Shipping | 3002 | Calcula e registra prazo de entrega (1 dia + 1 dia a cada 5 unidades) |

## Pré-requisitos

- Docker
- Docker Compose
- grpcurl (para testes)

## Executar com Docker Compose

Na pasta `microservices/`, execute:

```bash
docker-compose up --build
```

Isso sobe automaticamente:
- MySQL com os bancos `order`, `payment` e `shipping`
- Os três microsserviços
- Serviço `seed` que insere 5 produtos no estoque automaticamente

## Produtos cadastrados no estoque

O serviço `seed` insere automaticamente os seguintes produtos ao subir:

| Código | Nome | Quantidade |
|---|---|---|
| PROD-A | Notebook | 10 |
| PROD-B | Mouse | 50 |
| PROD-C | Teclado | 30 |
| PROD-D | Monitor | 5 |
| PROD-E | Headset | 20 |

Para adicionar mais produtos manualmente, conecte ao banco `order` e execute:

```sql
INSERT INTO products (product_code, name, quantity, created_at, updated_at)
VALUES ('PROD-F', 'Webcam', 15, NOW(), NOW());
```

## Testes

### Script de testes automático

Com todos os serviços rodando, execute:

```bash
sh test.sh
```

O script testa os seguintes cenários:

| Teste | Esperado |
|---|---|
| Pedido válido | Sucesso com order_id |
| Produto inexistente | Erro NOT_FOUND |
| Quantidade insuficiente em estoque | Erro INVALID_ARGUMENT |
| Mais de 50 itens | Erro INVALID_ARGUMENT |
| Pagamento acima de R$1000 | Erro INVALID_ARGUMENT |

### Teste manual com grpcurl

```bash
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 2, "unit_price": 10.0}
  ]
}' -plaintext localhost:3000 Order/Create
```

### Com o cliente Python

```bash
cd client
python client.py
```

## Executar localmente (sem Docker)

### 1. Subir o MySQL

```bash
docker run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=minhasenha -v "$(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql" mysql
```

### 2. Subir o Payment

```bash
cd payment
DB_DRIVER=mysql DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/payment APPLICATION_PORT=3001 ENV=development go run cmd/main.go
```

### 3. Subir o Shipping

```bash
cd shipping
DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/shipping APPLICATION_PORT=3002 ENV=development go run cmd/main.go
```

### 4. Subir o Order

```bash
cd order
DATA_SOURCE_URL=root:minhasenha@tcp(127.0.0.1:3306)/order APPLICATION_PORT=3000 ENV=development PAYMENT_SERVICE_URL=localhost:3001 SHIPPING_SERVICE_URL=localhost:3002 go run cmd/main.go
```

## Variáveis de Ambiente

### Order

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATA_SOURCE_URL` | URL do banco MySQL | `root:senha@tcp(127.0.0.1:3306)/order` |
| `APPLICATION_PORT` | Porta do serviço | `3000` |
| `ENV` | Ambiente | `development` |
| `PAYMENT_SERVICE_URL` | URL do serviço Payment | `localhost:3001` |
| `SHIPPING_SERVICE_URL` | URL do serviço Shipping | `localhost:3002` |

### Payment

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATA_SOURCE_URL` | URL do banco MySQL | `root:senha@tcp(127.0.0.1:3306)/payment` |
| `APPLICATION_PORT` | Porta do serviço | `3001` |
| `ENV` | Ambiente | `development` |

### Shipping

| Variável | Descrição | Exemplo |
|---|---|---|
| `DATA_SOURCE_URL` | URL do banco MySQL | `root:senha@tcp(127.0.0.1:3306)/shipping` |
| `APPLICATION_PORT` | Porta do serviço | `3002` |
| `ENV` | Ambiente | `development` |
