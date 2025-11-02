# Quick Import Guide - Polyloft

## Answer to Your Question

You have the following structure:
```
src/
├── miarchivo.pf
├── carpeta1/
│   └── archivo2.pf
└── carpeta2/
    └── archivo3.pf
```

### How to import in archivo2.pf from miarchivo.pf?

In `carpeta1/archivo2.pf`:
```polyloft
import miarchivo { desiredFunction, DesiredClass }
```

### How to import in archivo3.pf from archivo2.pf?

In `carpeta2/archivo3.pf`:
```polyloft
import carpeta1.archivo2 { desiredFunction, DesiredClass }
```

## Complete Working Example

This project includes a complete working example that demonstrates exactly this.

### Step 1: Run the example

From this directory (`example_project/`):
```bash
polyloft run src/carpeta2/archivo3.pf
```

### Step 2: Look at the code

Check the following files to understand how it works:
- `src/miarchivo.pf` - Defines functions and classes
- `src/carpeta1/archivo2.pf` - Imports from miarchivo.pf
- `src/carpeta2/archivo3.pf` - Imports from archivo2.pf and miarchivo.pf

## General Syntax

```polyloft
// Import specific symbols
import path.to.module { Symbol1, Symbol2 }

// Import entire module in a namespace
import path.to.module

// Use imported symbol with namespace
module.Symbol1()
```

## Key Points

1. ✅ Use dots (`.`) to separate directories: `carpeta1.archivo2`
2. ✅ DO NOT include the `.pf` extension in imports
3. ✅ Run from the directory that contains `src/`
4. ✅ Paths are relative to the `src/` directory

## Example Output

```
=== Ejemplo de importación entre carpetas ===
Hola, Mundo!
Valor de MiClase: 42
Resultado: Hola, Mundo! - procesado
Clase desde carpeta1: Mi objeto
Hola, Usuario!
Constante: 42
```

For more details, see the complete `README.md` file.
