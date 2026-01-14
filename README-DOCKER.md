# Dockerización de Checker API con LocalStack

Este proyecto está dockerizado y utiliza LocalStack para simular los servicios de AWS (SQS y DynamoDB) en un entorno local.

## Requisitos

- Docker
- Docker Compose
- Go 1.25+ (para ejecución local y debugging)

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

## Debugging Local (Recomendado para desarrollo)

Para desarrollo, puedes ejecutar la aplicación Go directamente en tu máquina (no en Docker) y conectarla a LocalStack que corre en Docker. Esto permite debugging completo con breakpoints en VS Code/Cursor.

### Requisitos para debugging local

- Go 1.25+ instalado
- VS Code o Cursor con la extensión de Go instalada
- LocalStack corriendo en Docker
- Archivo `.vscode/launch.json` (ya incluido en el proyecto)

### Pasos para debugging local

1. **Iniciar solo LocalStack** (sin la aplicación Go):

```bash
docker compose up -d localstack
```

2. **Verificar que LocalStack está corriendo**:

```bash
docker compose ps localstack
```

Debe mostrar `Status: Up (healthy)`

3. **Ejecutar la aplicación en modo debug desde VS Code/Cursor**:

   - Abre el proyecto en VS Code/Cursor
   - Ve a "Run and Debug" (Ctrl+Shift+D / Cmd+Shift+D)
   - Selecciona "Debug Local" del dropdown
   - Presiona F5 o haz clic en "Start Debugging"

4. **Colocar breakpoints y depurar**:

   - Coloca breakpoints en tu código
   - La aplicación se ejecutará localmente
   - Los breakpoints funcionarán correctamente
   - Puedes usar la Debug Console para inspeccionar variables

5. **Acceder a Swagger**:

Una vez que la aplicación esté corriendo:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API**: http://localhost:8080

### Ejecutar la aplicación localmente sin debugger

Si prefieres ejecutar sin el debugger de VS Code:

```bash
# Variables de entorno necesarias
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export KVS_REGION=us-east-1
export KVS_TABLE_NAME=comparisons
export KVS_KEY_ATTRIBUTE_NAME=id
export SQS_REGION=us-east-1
export SQS_QUEUE_URL=http://localhost:4566/000000000000/comparison-queue
export SERVER_PORT=8080
export ENV=local

# Ejecutar la aplicación
go run cmd/server/main.go
```

O crear un archivo `.env` en la raíz del proyecto con estas variables (el código las cargará automáticamente si `ENV=local`).

### Variables de entorno para debugging local

Las siguientes variables de entorno están configuradas en `.vscode/launch.json`:

- **AWS_ENDPOINT_URL**: `http://localhost:4566` - Endpoint de LocalStack
- **AWS_REGION**: `us-east-1` - Región de AWS
- **AWS_ACCESS_KEY_ID**: `test` - Credenciales de prueba
- **AWS_SECRET_ACCESS_KEY**: `test` - Credenciales de prueba
- **KVS_REGION**: `us-east-1` - Región para DynamoDB
- **KVS_TABLE_NAME**: `comparisons` - Nombre de la tabla DynamoDB
- **KVS_KEY_ATTRIBUTE_NAME**: `id` - Nombre del atributo clave
- **SQS_REGION**: `us-east-1` - Región para SQS
- **SQS_QUEUE_URL**: `http://localhost:4566/000000000000/comparison-queue` - URL de la cola SQS
- **SERVER_PORT**: `8080` - Puerto del servidor HTTP
- **ENV**: `local` - Entorno local

**Nota**: Asegúrate de que LocalStack esté corriendo antes de iniciar la aplicación localmente, ya que la aplicación necesita conectarse a DynamoDB y SQS.

### Ventajas del debugging local

- ✅ Breakpoints funcionan perfectamente
- ✅ Debug Console disponible
- ✅ Cambios en el código se reflejan inmediatamente (sin rebuild de Docker)
- ✅ Más rápido que debugging remoto
- ✅ Swagger disponible inmediatamente

## Configuración

### Variables de entorno

Las variables de entorno están configuradas en `docker-compose.yml`:

- **AWS_ENDPOINT_URL**: `http://localstack:4566` - Endpoint de LocalStack
- **KVS_TABLE_NAME**: `comparisons` - Nombre de la tabla DynamoDB
- **SQS_QUEUE_URL**: `http://localstack:4566/000000000000/comparison-queue` - URL de la cola SQS
- **SERVER_PORT**: `8080` - Puerto del servidor HTTP

### Recursos creados automáticamente

Los scripts de inicialización en `localstack-init/` crean automáticamente:

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
