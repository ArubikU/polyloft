# Guía Rápida de Importaciones - Polyloft

## Respuesta a tu Pregunta

Tienes la siguiente estructura:
```
src/
├── miarchivo.pf
├── carpeta1/
│   └── archivo2.pf
└── carpeta2/
    └── archivo3.pf
```

### ¿Cómo importar en archivo2.pf desde miarchivo.pf?

En `carpeta1/archivo2.pf`:
```polyloft
import miarchivo { funcionDeseada, ClaseDeseada }
```

### ¿Cómo importar en archivo3.pf desde archivo2.pf?

En `carpeta2/archivo3.pf`:
```polyloft
import carpeta1.archivo2 { funcionDeseada, ClaseDeseada }
```

## Ejemplo Completo Funcional

Este proyecto incluye un ejemplo completo que demuestra exactamente esto.

### Paso 1: Ejecutar el ejemplo

Desde este directorio (`example_project/`):
```bash
polyloft run src/carpeta2/archivo3.pf
```

### Paso 2: Ver el código

Mira los siguientes archivos para entender cómo funciona:
- `src/miarchivo.pf` - Define funciones y clases
- `src/carpeta1/archivo2.pf` - Importa desde miarchivo.pf
- `src/carpeta2/archivo3.pf` - Importa desde archivo2.pf y miarchivo.pf

## Sintaxis General

```polyloft
// Importar símbolos específicos
import ruta.al.modulo { Simbolo1, Simbolo2 }

// Importar todo el módulo en un namespace
import ruta.al.modulo

// Usar símbolo importado con namespace
modulo.Simbolo1()
```

## Puntos Clave

1. ✅ Usa puntos (`.`) para separar directorios: `carpeta1.archivo2`
2. ✅ NO incluyas la extensión `.pf` en las importaciones
3. ✅ Ejecuta desde el directorio que contiene `src/`
4. ✅ Las rutas son relativas al directorio `src/`

## Salida del Ejemplo

```
=== Ejemplo de importación entre carpetas ===
Hola, Mundo!
Valor de MiClase: 42
Resultado: Hola, Mundo! - procesado
Clase desde carpeta1: Mi objeto
Hola, Usuario!
Constante: 42
```

Para más detalles, consulta el archivo `README.md` completo.
