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

## Debugging

El proyecto incluye configuración para debugging remoto usando Delve en un contenedor Docker.

### Requisitos para debugging

- Docker y Docker Compose (ya instalados)
- VS Code con la extensión de Go instalada
- Archivo `.vscode/launch.json` (ya incluido en el proyecto)

### Pasos para ejecutar en modo debug

1. **Construir la imagen de debug** (solo la primera vez o si cambias código):

```bash
docker compose build checker-api-debug
```

O usando Make:

```bash
make docker-debug-build
```

2. **Levantar el servicio de debug**:

```bash
docker compose up -d checker-api-debug
```

O usando Make:

```bash
make docker-debug
```

Esto iniciará:
- **LocalStack**: Servicio que simula SQS y DynamoDB
- **Checker API (Debug)**: La aplicación Go con Delve debugger en el puerto 2345

3. **Verificar que todo esté corriendo**:

```bash
docker compose ps
```

Debes ver:
- `checker-api-localstack` - Status: Up (healthy)
- `checker-api-debug` - Status: Up

4. **Ver los logs para confirmar que la aplicación inició**:

```bash
docker compose logs checker-api-debug --tail=20
```

O usando Make:

```bash
make docker-logs-debug
```

Deberías ver mensajes como:
- "HTTP server starting on port 8080"
- "Application started successfully"
- "API server listening at: [::]:2345"

### Conectarse desde VS Code

1. **Coloca un breakpoint** en tu código (por ejemplo, en el handler que quieres depurar)

2. **Abre VS Code** en el directorio del proyecto

3. **Conéctate al debugger**:
   - Abre la pestaña "Run and Debug" (Ctrl+Shift+D / Cmd+Shift+D)
   - Selecciona "Debug Docker (Remote)" del dropdown
   - Presiona F5 o haz clic en "Start Debugging"

4. **Verifica la conexión**:
   - Deberías ver en la barra inferior "Debugging" con el estado "connected"
   - La pestaña "DEBUG CONSOLE" mostrará información de la conexión

5. **Ejecuta tu código**:
   - Accede a Swagger UI: http://localhost:8080/swagger/index.html
   - Ejecuta un request desde Swagger
   - El código se detendrá en tu breakpoint

### Acceso a los servicios en modo debug

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Debugger (Delve)**: localhost:2345 (para VS Code)
- **LocalStack Dashboard**: http://localhost:4566/_localstack/health

### Notas importantes sobre debugging

- El programa inicia automáticamente con `--continue`, por lo que Swagger está disponible inmediatamente
- Los breakpoints funcionan mejor si te conectas desde VS Code antes de ejecutar el código que quieres depurar
- Si los breakpoints no se activan, intenta:
  1. Reiniciar el contenedor: `docker compose restart checker-api-debug`
  2. Conectarte desde VS Code inmediatamente después del reinicio
  3. Luego ejecutar el request desde Swagger

### Detener el modo debug

```bash
docker compose stop checker-api-debug
```

O para detener todo:

```bash
docker compose down
```

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
