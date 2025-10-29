# Math Module

The `Math` module provides mathematical functions and constants.

## Constants

### `Math.PI`
The mathematical constant π (pi) ≈ 3.14159265359

**Example:**
```pf
let circumference = 2 * Math.PI * radius
let area = Math.PI * radius * radius
```

### `Math.E`
Euler's number ≈ 2.71828182846

**Example:**
```pf
let result = Math.pow(Math.E, 2)  // e²
```

## Functions

### `Math.abs(x)`
Returns the absolute value of a number.

**Parameters:**
- `x` (Number): Input value

**Returns:** Float

**Examples:**
```pf
println(Math.abs(-5))      // 5.0
println(Math.abs(3.14))    // 3.14
println(Math.abs(-2.5))    // 2.5
```

### `Math.floor(x)`
Rounds down to nearest integer.

**Parameters:**
- `x` (Number): Input value

**Returns:** Float

**Examples:**
```pf
println(Math.floor(3.7))   // 3.0
println(Math.floor(-2.3))  // -3.0
```

### `Math.ceil(x)`
Rounds up to nearest integer.

**Parameters:**
- `x` (Number): Input value

**Returns:** Float

**Examples:**
```pf
println(Math.ceil(3.2))    // 4.0
println(Math.ceil(-2.7))   // -2.0
```

### `Math.round(x)`
Rounds to nearest integer.

**Parameters:**
- `x` (Number): Input value

**Returns:** Float

**Examples:**
```pf
println(Math.round(3.5))   // 4.0
println(Math.round(3.4))   // 3.0
println(Math.round(-2.5))  // -2.0
```

### `Math.sqrt(x)`
Returns square root.

**Parameters:**
- `x` (Number): Input value (must be ≥ 0)

**Returns:** Float

**Examples:**
```pf
println(Math.sqrt(16))     // 4.0
println(Math.sqrt(2))      // 1.414...
println(Math.sqrt(0))      // 0.0
```

### `Math.pow(base, exponent)`
Returns base raised to exponent power.

**Parameters:**
- `base` (Number): Base value
- `exponent` (Number): Exponent value

**Returns:** Float

**Examples:**
```pf
println(Math.pow(2, 3))    // 8.0 (2³)
println(Math.pow(5, 2))    // 25.0 (5²)
println(Math.pow(2, -1))   // 0.5 (2⁻¹)
println(Math.pow(4, 0.5))  // 2.0 (√4)
```

### Trigonometric Functions

### `Math.sin(x)`
Returns sine of x (radians).

**Examples:**
```pf
println(Math.sin(0))          // 0.0
println(Math.sin(Math.PI / 2)) // 1.0
```

### `Math.cos(x)`
Returns cosine of x (radians).

**Examples:**
```pf
println(Math.cos(0))       // 1.0
println(Math.cos(Math.PI)) // -1.0
```

### `Math.tan(x)`
Returns tangent of x (radians).

**Examples:**
```pf
println(Math.tan(0))       // 0.0
println(Math.tan(Math.PI / 4))  // ~1.0
```

### `Math.min(a, b)`
Returns the smaller of two numbers.

**Parameters:**
- `a` (Number): First value
- `b` (Number): Second value

**Returns:** Float

**Examples:**
```pf
println(Math.min(5, 3))    // 3.0
println(Math.min(-2, -5))  // -5.0
```

### `Math.max(a, b)`
Returns the larger of two numbers.

**Parameters:**
- `a` (Number): First value
- `b` (Number): Second value

**Returns:** Float

**Examples:**
```pf
println(Math.max(5, 3))    // 5.0
println(Math.max(-2, -5))  // -2.0
```

### `Math.clamp(value, min, max)`
Constrains a value between min and max.

**Parameters:**
- `value` (Number): Value to clamp
- `min` (Number): Minimum bound
- `max` (Number): Maximum bound

**Returns:** Float

**Examples:**
```pf
println(Math.clamp(5, 0, 10))   // 5.0
println(Math.clamp(-5, 0, 10))  // 0.0
println(Math.clamp(15, 0, 10))  // 10.0
```

### `Math.random()`
Returns a random float between 0.0 and 1.0.

**Returns:** Float

**Examples:**
```pf
let rand = Math.random()
println(rand)  // e.g., 0.42387...

// Random integer between 1-100
let randInt = Math.floor(Math.random() * 100) + 1
```

## Common Patterns

### Distance Between Points
```pf
def distance(x1, y1, x2, y2):
    let dx = x2 - x1
    let dy = y2 - y1
    return Math.sqrt(dx * dx + dy * dy)
end

let dist = distance(0, 0, 3, 4)
println(dist)  // 5.0
```

### Angle Conversion
```pf
def degreesToRadians(degrees):
    return degrees * Math.PI / 180
end

def radiansToDegrees(radians):
    return radians * 180 / Math.PI
end

let angle = 90
let rad = degreesToRadians(angle)
println(Math.sin(rad))  // 1.0
```

### Circle Calculations
```pf
def circleArea(radius):
    return Math.PI * Math.pow(radius, 2)
end

def circleCircumference(radius):
    return 2 * Math.PI * radius
end

println(circleArea(5))           // 78.53...
println(circleCircumference(5))  // 31.41...
```

### Random Range
```pf
def randomRange(min, max):
    return min + Math.random() * (max - min)
end

def randomInt(min, max):
    return Math.floor(randomRange(min, max + 1))
end

println(randomRange(10, 20))  // e.g., 15.7
println(randomInt(1, 6))      // e.g., 4 (dice roll)
```

## See Also

- [Sys Module](sys.md) - System utilities
- [Number Types](../types/numbers.md) - Int and Float types
