# Bytes

La clase `Bytes` proporciona funcionalidades para trabajar con arrays de bytes.

## Constructores

```polyloft
// Crear array de bytes vacío
let b1 = Bytes()

// Crear array de bytes con tamaño específico
let b2 = Bytes(100)

// Crear desde array de enteros (0-255)
let b3 = Bytes([72, 101, 108, 108, 111])
```

## Métodos Estáticos

### fromString(str: String) -> Bytes

Crea un `Bytes` desde un string.

```polyloft
let bytes = Bytes.fromString("Hello, World!")
```

### fromHex(hex: String) -> Bytes

Crea un `Bytes` desde una cadena hexadecimal.

```polyloft
let bytes = Bytes.fromHex("48656c6c6f")
```

## Métodos de Instancia

### size() -> Int

Retorna el tamaño del array de bytes.

```polyloft
let bytes = Bytes.fromString("Hello")
Sys.println(bytes.size())  // 5
```

### get(index: Int) -> Int

Obtiene el byte en la posición especificada (0-255).

```polyloft
let bytes = Bytes.fromString("ABC")
Sys.println(bytes.get(0))  // 65 (ASCII 'A')
```

### set(index: Int, value: Int) -> Void

Establece el byte en la posición especificada.

```polyloft
let bytes = Bytes(3)
bytes.set(0, 65)  // 'A'
bytes.set(1, 66)  // 'B'
bytes.set(2, 67)  // 'C'
```

### toString() -> String

Convierte los bytes a string.

```polyloft
let bytes = Bytes.fromString("Hello")
Sys.println(bytes.toString())  // "Hello"
```

### toHex() -> String

Convierte los bytes a representación hexadecimal.

```polyloft
let bytes = Bytes.fromString("Hello")
Sys.println(bytes.toHex())  // "48656c6c6f"
```

### toArray() -> Array

Convierte los bytes a array de enteros.

```polyloft
let bytes = Bytes.fromString("ABC")
let arr = bytes.toArray()
// arr = [65, 66, 67]
```

### slice(start: Int, end: Int) -> Bytes

Crea un nuevo `Bytes` con una porción del array original.

```polyloft
let bytes = Bytes.fromString("Hello World")
let slice = bytes.slice(0, 5)
Sys.println(slice.toString())  // "Hello"
```

### equals(other: Bytes) -> Bool

Compara dos arrays de bytes.

```polyloft
let b1 = Bytes.fromString("Hello")
let b2 = Bytes.fromString("Hello")
let b3 = Bytes.fromString("World")

Sys.println(b1.equals(b2))  // true
Sys.println(b1.equals(b3))  // false
```

## Ejemplo Completo

```polyloft
// Crear y manipular bytes
let bytes = Bytes.fromString("Hello")

// Ver tamaño
Sys.println("Size: " + bytes.size())

// Modificar bytes
bytes.set(0, 74)  // Cambiar 'H' por 'J'
Sys.println(bytes.toString())  // "Jello"

// Convertir a hex
Sys.println("Hex: " + bytes.toHex())

// Crear slice
let slice = bytes.slice(1, 4)
Sys.println(slice.toString())  // "ell"

// Comparar
let other = Bytes.fromString("Jello")
Sys.println(bytes.equals(other))  // true
```
