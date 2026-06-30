#!/bin/bash

HOST="localhost:3000"
PASS=0
FAIL=0

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

check() {
  local desc=$1
  local result=$2
  local expect_error=$3

  if [ "$expect_error" = "true" ]; then
    if echo "$result" | grep -qi "code\|error"; then
      echo -e "  ${GREEN}PASS${RESET} $desc"
      PASS=$((PASS+1))
    else
      echo -e "  ${RED}FAIL${RESET} $desc"
      echo -e "       ${YELLOW}$result${RESET}"
      FAIL=$((FAIL+1))
    fi
  else
    if echo "$result" | grep -qi "orderId\|order_id"; then
      echo -e "  ${GREEN}PASS${RESET} $desc"
      PASS=$((PASS+1))
    else
      echo -e "  ${RED}FAIL${RESET} $desc"
      echo -e "       ${YELLOW}$result${RESET}"
      FAIL=$((FAIL+1))
    fi
  fi
}

echo -e "${BOLD}${CYAN}"
echo "  ╔══════════════════════════════════════════════╗"
echo "  ║       Testes - Microsservicos gRPC           ║"
echo "  ╚══════════════════════════════════════════════╝"
echo -e "${RESET}"

echo ""
echo -e "${BOLD}--- [1] Fluxo completo: Order -> Payment -> Shipping ---${RESET}"
echo -e "${CYAN}  Verifica se a cadeia completa entre os 3 servicos funciona.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 2, "unit_price": 100.0},
    {"product_code": "PROD-B", "quantity": 1, "unit_price": 80.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Fluxo completo (Order->Payment->Shipping)" "$RESULT" "false"

echo ""
echo -e "${BOLD}--- [2] Propagacao de erro: produto inexistente ---${RESET}"
echo -e "${CYAN}  Verifica se um erro gerado localmente no Order chega corretamente ao cliente.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-Z", "quantity": 1, "unit_price": 10.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Erro propagado: NOT_FOUND para produto inexistente" "$RESULT" "true"

echo ""
echo -e "${BOLD}--- [3] Propagacao de erro: estoque insuficiente ---${RESET}"
echo -e "${CYAN}  Verifica se a validacao de estoque funciona. PROD-D tem 5 unidades, pedindo 6.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-D", "quantity": 6, "unit_price": 900.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Erro propagado: INVALID_ARGUMENT para estoque insuficiente" "$RESULT" "true"

echo ""
echo -e "${BOLD}--- [4] Validacao local: limite de 50 itens ---${RESET}"
echo -e "${CYAN}  Verifica se o Order rejeita pedidos com mais de 50 itens sem nem consultar outros servicos.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-B", "quantity": 51, "unit_price": 80.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Validacao local: mais de 50 itens rejeitado" "$RESULT" "true"

echo ""
echo -e "${BOLD}--- [5] Teste de valor limite: exatamente 50 itens ---${RESET}"
echo -e "${CYAN}  Verifica o comportamento no valor exato do limite (boundary value analysis).${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-B", "quantity": 50, "unit_price": 1.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Valor limite: exatamente 50 itens aceito" "$RESULT" "false"

echo ""
echo -e "${BOLD}--- [6] Erro em servico downstream: Payment recusa valor alto ---${RESET}"
echo -e "${CYAN}  Verifica se uma falha no Payment (servico downstream) chega corretamente ao cliente.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 1, "unit_price": 1001.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Erro downstream: Payment rejeita valor acima de 1000" "$RESULT" "true"

echo ""
echo -e "${BOLD}--- [7] Teste de valor limite: exatamente R\$1000 no Payment ---${RESET}"
echo -e "${CYAN}  Verifica o valor exato no limite do servico downstream (boundary value analysis).${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": [
    {"product_code": "PROD-A", "quantity": 1, "unit_price": 1000.0}
  ]
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Valor limite: exatamente R\$1000 aceito pelo Payment" "$RESULT" "false"

echo ""
echo -e "${BOLD}--- [8] Concorrencia: 5 requisicoes simultaneas ---${RESET}"
echo -e "${CYAN}  Verifica se o sistema lida corretamente com multiplas requisicoes em paralelo.${RESET}"
echo "Disparando 5 requisicoes simultaneas..."
for i in 1 2 3 4 5; do
  grpcurl -d "{
    \"costumer_id\": $i,
    \"order_items\": [
      {\"product_code\": \"PROD-C\", \"quantity\": 1, \"unit_price\": 50.0}
    ]
  }" -plaintext $HOST Order/Create > /tmp/concurrent_$i.txt 2>&1 &
done
wait
CONCURRENT_OK=0
for i in 1 2 3 4 5; do
  if grep -qi "orderId\|order_id" /tmp/concurrent_$i.txt; then
    CONCURRENT_OK=$((CONCURRENT_OK+1))
  fi
done
echo "Requisicoes bem-sucedidas: $CONCURRENT_OK/5"
if [ "$CONCURRENT_OK" -eq 5 ]; then
  echo "[OK] Concorrencia: todas as requisicoes processadas"
  PASS=$((PASS+1))
else
  echo "[FAIL] Concorrencia: apenas $CONCURRENT_OK/5 bem-sucedidas"
  FAIL=$((FAIL+1))
fi

echo ""
echo -e "${BOLD}--- [9] Pedido vazio (zero itens) ---${RESET}"
echo -e "${CYAN}  Verifica se entrada invalida e rejeitada antes de chegar a qualquer servico externo.${RESET}"
RESULT=$(grpcurl -d '{
  "costumer_id": 1,
  "order_items": []
}' -plaintext $HOST Order/Create 2>&1)
echo "$RESULT"
check "Validacao: pedido sem itens rejeitado" "$RESULT" "true"

echo ""
echo ""
echo -e "${BOLD}${CYAN}  ╔══════════════════════════════════════════════╗${RESET}"
if [ "$FAIL" -eq 0 ]; then
  echo -e "${BOLD}${GREEN}  ║   Resultado: $PASS/$((PASS+FAIL)) testes passaram         ║${RESET}"
else
  echo -e "${BOLD}${RED}  ║   Resultado: $PASS passou(aram) | $FAIL falhou(aram)    ║${RESET}"
fi
echo -e "${BOLD}${CYAN}  ╚══════════════════════════════════════════════╝${RESET}"
