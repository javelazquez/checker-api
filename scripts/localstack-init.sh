#!/bin/sh

# Script de inicialización para LocalStack
# Este script se ejecuta automáticamente cuando LocalStack está listo
# No usar set -e para que no falle si hay errores menores

echo "🚀 Inicializando recursos de LocalStack..."

# Configurar variables
ENDPOINT="http://localhost:4566"
REGION="us-east-1"

# Esperar un poco más para asegurar que LocalStack esté completamente listo
sleep 3

# Verificar que LocalStack esté disponible (máximo 30 segundos)
for i in 1 2 3 4 5 6 7 8 9 10; do
  if curl -s "${ENDPOINT}/_localstack/health" >/dev/null 2>&1; then
    echo "✅ LocalStack está listo"
    break
  fi
  if [ $i -eq 10 ]; then
    echo "⚠️  LocalStack no está respondiendo, pero continuando..."
    exit 0
  fi
  sleep 3
done

# Crear tabla de DynamoDB usando curl (ignorar errores si ya existe)
echo "📊 Creando tabla DynamoDB: comparisons"
RESULT=$(curl -s -w "\n%{http_code}" -X POST "${ENDPOINT}/" \
  -H "Content-Type: application/x-amz-json-1.0" \
  -H "X-Amz-Target: DynamoDB_20120810.CreateTable" \
  -d "{
    \"TableName\": \"comparisons\",
    \"AttributeDefinitions\": [{\"AttributeName\": \"id\", \"AttributeType\": \"S\"}],
    \"KeySchema\": [{\"AttributeName\": \"id\", \"KeyType\": \"HASH\"}],
    \"BillingMode\": \"PAY_PER_REQUEST\"
  }" 2>/dev/null)

HTTP_CODE=$(echo "$RESULT" | tail -n1)
if [ "$HTTP_CODE" = "200" ]; then
  echo "✅ Tabla DynamoDB creada"
else
  echo "ℹ️  Tabla DynamoDB (puede que ya exista)"
fi

# Crear cola de SQS usando múltiples métodos
echo "📬 Creando cola SQS: comparison-queue"
QUEUE_CREATED=0

# Método 1: Verificar si ya existe
GET_RESPONSE=$(curl -s -X POST "${ENDPOINT}/" \
  -H "Content-Type: application/x-amz-json-1.0" \
  -H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
  -d "{\"QueueName\": \"comparison-queue\"}" 2>&1)

if echo "$GET_RESPONSE" | grep -q "QueueUrl"; then
  QUEUE_URL=$(echo "$GET_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4)
  echo "ℹ️  Cola SQS ya existe: $QUEUE_URL"
  QUEUE_CREATED=1
else
  # Método 2: Crear usando endpoint específico de SQS con account ID
  echo "Intentando crear con endpoint específico..."
  SQS_ENDPOINT="${ENDPOINT}/000000000000"
  CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${SQS_ENDPOINT}/" \
    -H "Content-Type: application/x-amz-json-1.0" \
    -H "X-Amz-Target: AWSSimpleQueueServiceV20121105.CreateQueue" \
    -d "{\"QueueName\": \"comparison-queue\"}" 2>&1)
  
  HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
  
  # Esperar y verificar
  sleep 5
  for i in 1 2 3; do
    GET_RESPONSE=$(curl -s -X POST "${ENDPOINT}/" \
      -H "Content-Type: application/x-amz-json-1.0" \
      -H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
      -d "{\"QueueName\": \"comparison-queue\"}" 2>&1)
    
    if echo "$GET_RESPONSE" | grep -q "QueueUrl"; then
      QUEUE_URL=$(echo "$GET_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4)
      echo "✅ Cola SQS creada: $QUEUE_URL"
      QUEUE_CREATED=1
      break
    fi
    sleep 3
  done
  
  # Método 3: Si aún no se creó, intentar con endpoint general
  if [ $QUEUE_CREATED -eq 0 ]; then
    echo "Intentando con endpoint general..."
    CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${ENDPOINT}/" \
      -H "Content-Type: application/x-amz-json-1.0" \
      -H "X-Amz-Target: AWSSimpleQueueServiceV20121105.CreateQueue" \
      -d "{\"QueueName\": \"comparison-queue\"}" 2>&1)
    
    HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
    
    # Esperar y verificar múltiples veces
    for i in 1 2 3 4 5; do
      sleep 5
      GET_RESPONSE=$(curl -s -X POST "${ENDPOINT}/" \
        -H "Content-Type: application/x-amz-json-1.0" \
        -H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
        -d "{\"QueueName\": \"comparison-queue\"}" 2>&1)
      
      if echo "$GET_RESPONSE" | grep -q "QueueUrl"; then
        QUEUE_URL=$(echo "$GET_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4)
        echo "✅ Cola SQS creada: $QUEUE_URL"
        QUEUE_CREATED=1
        break
      fi
    done
  fi
fi

if [ $QUEUE_CREATED -eq 0 ]; then
  echo "⚠️  No se pudo crear la cola SQS después de intentar múltiples métodos"
  echo "⚠️  La aplicación puede fallar al intentar consumir mensajes de SQS"
  echo "ℹ️  Puedes crear la cola manualmente con: make docker-init-resources"
fi

echo "✅ Inicialización de LocalStack completada"
exit 0
