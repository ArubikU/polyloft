# Ejemplo de Importaciones en Polyloft

Este proyecto demuestra cómo usar el sistema de importación en Polyloft para importar código entre diferentes archivos y carpetas.

## Estructura del Proyecto

```
example_project/
├── src/
│   ├── miarchivo.pf              # Archivo en la raíz de src/
│   ├── carpeta1/
│   │   └── archivo2.pf           # Archivo en subcarpeta
│   └── carpeta2/
│       └── archivo3.pf           # Archivo en otra subcarpeta
```

## Cómo Funcionan las Importaciones

### 1. Importar desde un Archivo Padre (Parent Directory)

**Archivo: `carpeta1/archivo2.pf`**

Para importar desde `miarchivo.pf` (que está en el directorio padre):

```polyloft
import miarchivo { saludar, MiClase, CONSTANTE }
```

**Explicación:**
- El sistema de importación busca archivos relativos al directorio `src/` del proyecto
- Especifica el nombre del archivo sin la extensión `.pf`
- Entre llaves `{ }` lista los símbolos que quieres importar (funciones, clases, constantes)
- Como `carpeta1/archivo2.pf` y `miarchivo.pf` están ambos bajo `src/`, puedes usar el nombre directamente

### 2. Importar desde una Carpeta Hermana (Sibling Directory)

**Archivo: `carpeta2/archivo3.pf`**

Para importar desde `carpeta1/archivo2.pf`:

```polyloft
import carpeta1.archivo2 { procesarDatos, ClaseCarpeta1 }
```

**Explicación:**
- Usa la ruta relativa desde el directorio `src/`
- Navega por las carpetas usando puntos `.`
- Especifica el archivo sin extensión
- Lista los símbolos a importar entre llaves

### 3. Importaciones Múltiples

Puedes importar de múltiples archivos en el mismo programa:

```polyloft
import carpeta1.archivo2 { procesarDatos, ClaseCarpeta1 }
import miarchivo { saludar, CONSTANTE }
```

### 4. Sintaxis de Importación General

```polyloft
import ruta.del.modulo { Simbolo1, Simbolo2, Simbolo3 }
```

o para importar todo en un namespace:

```polyloft
import ruta.del.modulo
```

Y luego usar como `modulo.Simbolo1`, `modulo.Simbolo2`, etc.

## Ejecutar el Ejemplo

**IMPORTANTE**: Debes ejecutar el comando desde la raíz del proyecto (el directorio `example_project/`), ya que el sistema de importación busca archivos relativos al directorio `src/`.

Desde el directorio del proyecto:

```bash
cd example_project
polyloft run src/carpeta2/archivo3.pf
```

O si estás en la raíz del repositorio:

```bash
cd /ruta/a/polyloft
polyloft run example_project/src/carpeta2/archivo3.pf
```

## Reglas Importantes

1. **Ejecutar desde la Raíz del Proyecto**: Los comandos deben ejecutarse desde el directorio que contiene `src/`
2. **Rutas Relativas a src/**: Las importaciones son relativas al directorio `src/` del proyecto
3. **Sin Extensión**: No incluyas `.pf` en las rutas de importación
4. **Puntos como Separadores**: Usa `.` para separar directorios (ej: `carpeta1.archivo2`)
5. **Especificar Símbolos**: Lista explícitamente qué quieres importar entre `{ }`
6. **Orden de Búsqueda**: 
   - Primero busca relativamente desde el archivo actual
   - Luego busca en `src/` (relativo a donde ejecutas el comando)
   - Luego busca en `libs/`

## Salida Esperada

Al ejecutar `archivo3.pf`:

```
=== Ejemplo de importación entre carpetas ===
Hola, Mundo!
Valor de MiClase: 42
Resultado: Hola, Mundo! - procesado
Clase desde carpeta1: Mi objeto
Hola, Usuario!
Constante: 42
```

## Notas Adicionales

- Las importaciones son procesadas en tiempo de ejecución
- Los módulos se cachean después de la primera carga
- Puedes importar clases, funciones y constantes
- El sistema soporta importaciones cíclicas con cuidado
