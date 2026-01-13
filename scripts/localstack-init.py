#!/usr/bin/env python3
"""
Script de inicialización para LocalStack
Crea los recursos necesarios (DynamoDB table y SQS queue)
"""

import sys
import time

# Intentar importar boto3, si no está disponible usar método alternativo
try:
    import boto3
    HAS_BOTO3 = True
except ImportError:
    HAS_BOTO3 = False
    print("⚠️  boto3 no está disponible, usando método alternativo con curl")

# Configuración
ENDPOINT_URL = "http://localhost:4566"
REGION = "us-east-1"

if HAS_BOTO3:
    # Esperar a que LocalStack esté listo
    print("⏳ Esperando a que LocalStack esté listo...")
    max_retries = 30
    retry_count = 0

    while retry_count < max_retries:
        try:
            # Intentar crear un cliente de DynamoDB
            dynamodb = boto3.client(
                'dynamodb',
                endpoint_url=ENDPOINT_URL,
                region_name=REGION,
                aws_access_key_id='test',
                aws_secret_access_key='test'
            )
            # Verificar que el servicio esté disponible
            dynamodb.list_tables()
            print("✅ LocalStack está listo")
            break
        except Exception as e:
            retry_count += 1
            if retry_count >= max_retries:
                print(f"❌ Error: No se pudo conectar a LocalStack después de {max_retries} intentos")
                sys.exit(0)  # No fallar, solo continuar
            time.sleep(2)
    else:
        print("⚠️  Timeout esperando a LocalStack, continuando...")
        sys.exit(0)

    # Crear cliente de DynamoDB
    dynamodb = boto3.client(
        'dynamodb',
        endpoint_url=ENDPOINT_URL,
        region_name=REGION,
        aws_access_key_id='test',
        aws_secret_access_key='test'
    )

    # Crear cliente de SQS
    sqs = boto3.client(
        'sqs',
        endpoint_url=ENDPOINT_URL,
        region_name=REGION,
        aws_access_key_id='test',
        aws_secret_access_key='test'
    )

    # Crear tabla de DynamoDB
    print("📊 Creando tabla DynamoDB: comparisons")
    try:
        dynamodb.create_table(
            TableName='comparisons',
            AttributeDefinitions=[
                {
                    'AttributeName': 'id',
                    'AttributeType': 'S'
                }
            ],
            KeySchema=[
                {
                    'AttributeName': 'id',
                    'KeyType': 'HASH'
                }
            ],
            BillingMode='PAY_PER_REQUEST'
        )
        print("✅ Tabla DynamoDB 'comparisons' creada exitosamente")
    except dynamodb.exceptions.ResourceInUseException:
        print("ℹ️  La tabla DynamoDB 'comparisons' ya existe")
    except Exception as e:
        print(f"⚠️  Error al crear la tabla DynamoDB: {e}")

    # Crear cola de SQS
    print("📬 Creando cola SQS: comparison-queue")
    try:
        # Primero verificar si ya existe
        try:
            response = sqs.get_queue_url(QueueName='comparison-queue')
            queue_url = response['QueueUrl']
            print(f"ℹ️  La cola SQS 'comparison-queue' ya existe: {queue_url}")
        except Exception as get_error:
            # La cola no existe, intentar crearla
            print(f"La cola no existe, creándola... (error al verificar: {get_error})")
            try:
                response = sqs.create_queue(QueueName='comparison-queue')
                queue_url = response['QueueUrl']
                print(f"✅ Cola SQS creada exitosamente: {queue_url}")
                
                # Verificar que se creó correctamente
                time.sleep(2)
                verify_response = sqs.get_queue_url(QueueName='comparison-queue')
                verify_url = verify_response['QueueUrl']
                print(f"✅ Cola SQS verificada: {verify_url}")
            except Exception as create_error:
                print(f"⚠️  Error al crear la cola SQS: {create_error}")
                print(f"   Tipo de error: {type(create_error).__name__}")
                # Intentar listar todas las colas para debug
                try:
                    queues = sqs.list_queues()
                    print(f"Colas disponibles: {queues.get('QueueUrls', [])}")
                except Exception as list_error:
                    print(f"Error al listar colas: {list_error}")
                import traceback
                traceback.print_exc()
    except Exception as e:
        print(f"⚠️  Error general al crear/obtener la cola SQS: {e}")
        print(f"   Tipo de error: {type(e).__name__}")
        import traceback
        traceback.print_exc()
else:
    # Fallback: usar método con subprocess y curl si boto3 no está disponible
    import subprocess
    import json
    
    print("📊 Creando tabla DynamoDB: comparisons")
    result = subprocess.run([
        'curl', '-s', '-X', 'POST', f'{ENDPOINT_URL}/',
        '-H', 'Content-Type: application/x-amz-json-1.0',
        '-H', 'X-Amz-Target: DynamoDB_20120810.CreateTable',
        '-d', json.dumps({
            'TableName': 'comparisons',
            'AttributeDefinitions': [{'AttributeName': 'id', 'AttributeType': 'S'}],
            'KeySchema': [{'AttributeName': 'id', 'KeyType': 'HASH'}],
            'BillingMode': 'PAY_PER_REQUEST'
        })
    ], capture_output=True, text=True)
    
    if result.returncode == 0:
        print("✅ Tabla DynamoDB creada o ya existe")
    else:
        print(f"⚠️  Error al crear tabla DynamoDB: {result.stderr}")
    
    print("📬 Creando cola SQS: comparison-queue")
    result = subprocess.run([
        'curl', '-s', '-X', 'POST', f'{ENDPOINT_URL}/',
        '-H', 'Content-Type: application/x-amz-json-1.0',
        '-H', 'X-Amz-Target: AWSSimpleQueueServiceV20121105.CreateQueue',
        '-d', json.dumps({'QueueName': 'comparison-queue'})
    ], capture_output=True, text=True)
    
    if result.returncode == 0 and 'QueueUrl' in result.stdout:
        queue_url = json.loads(result.stdout).get('QueueUrl', '')
        print(f"✅ Cola SQS creada: {queue_url}")
    else:
        print(f"⚠️  Error al crear cola SQS: {result.stderr}")

print("✅ Inicialización de LocalStack completada")
