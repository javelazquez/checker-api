# Dockerización de Checker API con LocalStack

Este proyecto está dockerizado y utiliza LocalStack para simular los servicios de AWS (SQS y DynamoDB) en un entorno local.

## Requisitos

- Docker
- Docker Compose

## Uso

### Iniciar los servicios

Para iniciar la aplicación junto con LocalStack:

```bash
docker-compose up --build
```

Esto iniciará:
- **LocalStack**: Servicio que simula SQS y DynamoDB en el puerto 4566
- **Checker API**: La aplicación Go en el puerto 8080

### Iniciar en segundo plano

```bash
docker-compose up -d --build
```

### Ver logs

```bash
# Ver todos los logs
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f checker-api
docker-compose logs -f localstack
```

### Detener los servicios

```bash
docker-compose down
```

Para eliminar también los volúmenes (datos de LocalStack):

```bash
docker-compose down -v
```

## Configuración

### Variables de entorno

Las variables de entorno están configuradas en `docker-compose.yml`:

- **AWS_ENDPOINT_URL**: `http://localstack:4566` - Endpoint de LocalStack
- **KVS_TABLE_NAME**: `comparisons` - Nombre de la tabla DynamoDB
- **SQS_QUEUE_URL**: `http://localstack:4566/000000000000/comparison-queue` - URL de la cola SQS
- **SERVER_PORT**: `8080` - Puerto del servidor HTTP

### Recursos creados automáticamente

El script de inicialización (`scripts/localstack-init.sh`) crea automáticamente:

1. **Tabla DynamoDB**: `comparisons`
   - Clave primaria: `id` (String)
   - Modo de facturación: PAY_PER_REQUEST

2. **Cola SQS**: `comparison-queue`

## Acceso a los servicios

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **LocalStack Dashboard**: http://localhost:4566/_localstack/health

## Solución de problemas

### La aplicación no puede conectarse a LocalStack

Asegúrate de que:
1. LocalStack esté completamente iniciado (verifica los logs)
2. La variable `AWS_ENDPOINT_URL` esté configurada correctamente
3. Ambos servicios estén en la misma red Docker

### Los recursos no se crean automáticamente

El script de inicialización se ejecuta cuando LocalStack está listo. Si hay problemas:
1. Verifica los logs de LocalStack: `docker-compose logs localstack`
2. Verifica que el script tenga permisos de ejecución
3. Puedes crear los recursos manualmente usando la API de LocalStack

### Reconstruir la imagen

Si haces cambios en el código:

```bash
docker-compose build --no-cache checker-api
docker-compose up
```
