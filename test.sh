#!/bin/bash

HOST="localhost:3000"

echo "========================================"
echo " Testes do microsservico Order"
echo "========================================"

echo ""
echo "--- Teste 1: Pedido valido ---"
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 2, "unit_price": 1500.0},
    {"product_code": "PROD-B", "quantity": 1, "unit_price": 80.0}
  ]
}' -plaintext $HOST Order/Create

echo ""
echo "--- Teste 2: Produto inexistente ---"
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-Z", "quantity": 1, "unit_price": 10.0}
  ]
}' -plaintext $HOST Order/Create

echo ""
echo "--- Teste 3: Quantidade insuficiente em estoque ---"
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-D", "quantity": 99, "unit_price": 900.0}
  ]
}' -plaintext $HOST Order/Create

echo ""
echo "--- Teste 4: Pedido com mais de 50 itens ---"
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-B", "quantity": 51, "unit_price": 80.0}
  ]
}' -plaintext $HOST Order/Create

echo ""
echo "--- Teste 5: Pagamento acima de 1000 ---"
grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 1, "unit_price": 1500.0}
  ]
}' -plaintext $HOST Order/Create

echo ""
echo "========================================"
echo " Testes concluidos"
echo "========================================"
