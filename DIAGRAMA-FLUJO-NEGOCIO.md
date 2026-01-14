# Diagrama de Flujo - Checker API

## Flujo Principal: Comparación de Respuestas de API

Este diagrama describe el flujo de negocio de la aplicación desde la perspectiva del usuario.

```mermaid
flowchart TD
    A[Usuario: Envía dos respuestas de API para comparar] --> B{Validar que las respuestas sean válidas}
    B -->|Respuestas inválidas| C[Retornar error al usuario]
    B -->|Respuestas válidas| D[Comparar las respuestas]
    
    D --> E{¿Las respuestas son iguales?}
    
    E -->|Sí, son iguales| F[Generar resultado: Son iguales]
    E -->|No, son diferentes| G[Identificar todas las diferencias]
    
    G --> H[Analizar diferencias en:<br/>- Código de estado HTTP<br/>- Headers HTTP<br/>- Cuerpo JSON]
    
    H --> I[Generar lista detallada de diferencias]
    I --> J[Guardar resultado de la comparación]
    
    F --> J
    J --> K[Asignar ID único a la comparación]
    K --> L[Devolver resultado al usuario:<br/>- ID de la comparación<br/>- Si son iguales o no<br/>- Lista de diferencias encontradas]
    
    L --> M[Fin]
    C --> M
```

## Flujo Secundario: Consultar Comparación Guardada

```mermaid
flowchart TD
    A[Usuario: Solicita una comparación por ID] --> B{¿El ID existe?}
    B -->|No existe| C[Retornar error: Comparación no encontrada]
    B -->|Existe| D[Buscar comparación en el almacenamiento]
    D --> E[Devolver resultado al usuario:<br/>- ID de la comparación<br/>- Si las respuestas son iguales<br/>- Lista de diferencias encontradas]
    E --> F[Fin]
    C --> F
```

## Descripción del Proceso de Negocio

### ¿Qué hace la aplicación?

La aplicación **Checker API** permite comparar dos respuestas de API (HTTP) para identificar si son iguales o qué diferencias existen entre ellas.

### ¿Qué compara?

1. **Código de Estado HTTP**: Verifica si ambas respuestas tienen el mismo código (200, 404, 500, etc.)
2. **Headers HTTP**: Compara los encabezados de ambas respuestas
3. **Cuerpo JSON**: Compara el contenido JSON del cuerpo de las respuestas

### ¿Qué información devuelve?

- Un **ID único** que identifica la comparación realizada
- Un indicador que dice si las respuestas son **iguales o diferentes**
- Una **lista detallada de todas las diferencias** encontradas, indicando:
  - El tipo de diferencia (código de estado, header, campo del body, etc.)
  - La ubicación de la diferencia (por ejemplo: "body.user.name")
  - El valor en la respuesta fuente
  - El valor en la respuesta destino
  - Una descripción legible de la diferencia

### Casos de Uso

1. **Comparar respuestas de diferentes versiones de una API** para verificar que producen los mismos resultados
2. **Validar migraciones** de sistemas comparando respuestas antiguas vs nuevas
3. **Testing de integración** para asegurar que los cambios no rompen la compatibilidad
4. **Auditoría y seguimiento** de cambios en respuestas de API guardando los resultados de las comparaciones
