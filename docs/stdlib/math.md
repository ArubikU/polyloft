# Math Module

The Math module provides mathematical functions and constants for numerical computations.

## Constants

### Math.PI
```polyloft
Math.PI  // 3.141592653589793
```
The mathematical constant π (pi).

**Type**: `Float`

**Example**:
```polyloft
let circumference = 2 * Math.PI * radius
```

### Math.E
```polyloft
Math.E  // 2.718281828459045
```
The mathematical constant e (Euler's number).

**Type**: `Float`

**Example**:
```polyloft
let growth = Math.E * rate
```

## Functions

### Math.abs(x)
Returns the absolute value of a number.

**Parameters**:
- `x: Number` - The input number

**Returns**: `Float` - The absolute value

**Example**:
```polyloft
Math.abs(-5.5)    // 5.5
Math.abs(3.2)     // 3.2
Math.abs(0)       // 0
```

### Math.sqrt(x)
Returns the square root of a number.

**Parameters**:
- `x: Number` - The input number (must be non-negative)

**Returns**: `Float` - The square root

**Example**:
```polyloft
Math.sqrt(9)      // 3.0
Math.sqrt(16)     // 4.0
Math.sqrt(2)      // 1.4142135623730951
```

### Math.pow(base, exponent)
Returns base raised to the power of exponent.

**Parameters**:
- `base: Number` - The base number
- `exponent: Number` - The exponent

**Returns**: `Float` - The result of base^exponent

**Example**:
```polyloft
Math.pow(2, 3)    // 8.0
Math.pow(10, 2)   // 100.0
Math.pow(5, 0)    // 1.0
Math.pow(2, -1)   // 0.5
```

### Math.floor(x)
Returns the largest integer less than or equal to x.

**Parameters**:
- `x: Number` - The input number

**Returns**: `Float` - The floor value

**Example**:
```polyloft
Math.floor(3.7)   // 3.0
Math.floor(3.2)   // 3.0
Math.floor(-3.7)  // -4.0
```

### Math.ceil(x)
Returns the smallest integer greater than or equal to x.

**Parameters**:
- `x: Number` - The input number

**Returns**: `Float` - The ceiling value

**Example**:
```polyloft
Math.ceil(3.2)    // 4.0
Math.ceil(3.7)    // 4.0
Math.ceil(-3.2)   // -3.0
```

### Math.round(x)
Returns the nearest integer to x, rounding half away from zero.

**Parameters**:
- `x: Number` - The input number

**Returns**: `Float` - The rounded value

**Example**:
```polyloft
Math.round(3.5)   // 4.0
Math.round(3.4)   // 3.0
Math.round(-3.5)  // -4.0
```

### Math.min(a, b)
Returns the smaller of two numbers.

**Parameters**:
- `a: Number` - First number
- `b: Number` - Second number

**Returns**: `Float` - The minimum value

**Example**:
```polyloft
Math.min(5, 3)    // 3.0
Math.min(-1, 10)  // -1.0
```

### Math.max(a, b)
Returns the larger of two numbers.

**Parameters**:
- `a: Number` - First number
- `b: Number` - Second number

**Returns**: `Float` - The maximum value

**Example**:
```polyloft
Math.max(5, 3)    // 5.0
Math.max(-1, 10)  // 10.0
```

### Math.clamp(value, min, max)
Clamps a value between a minimum and maximum.

**Parameters**:
- `value: Number` - The value to clamp
- `min: Number` - The minimum value
- `max: Number` - The maximum value

**Returns**: `Float` - The clamped value

**Example**:
```polyloft
Math.clamp(5, 0, 10)   // 5.0
Math.clamp(-5, 0, 10)  // 0.0
Math.clamp(15, 0, 10)  // 10.0
```

## Trigonometric Functions

### Math.sin(x)
Returns the sine of x (x in radians).

**Parameters**:
- `x: Number` - Angle in radians

**Returns**: `Float` - The sine value

**Example**:
```polyloft
Math.sin(0)           // 0.0
Math.sin(Math.PI / 2) // 1.0
Math.sin(Math.PI)     // ~0.0
```

### Math.cos(x)
Returns the cosine of x (x in radians).

**Parameters**:
- `x: Number` - Angle in radians

**Returns**: `Float` - The cosine value

**Example**:
```polyloft
Math.cos(0)           // 1.0
Math.cos(Math.PI / 2) // ~0.0
Math.cos(Math.PI)     // -1.0
```

### Math.tan(x)
Returns the tangent of x (x in radians).

**Parameters**:
- `x: Number` - Angle in radians

**Returns**: `Float` - The tangent value

**Example**:
```polyloft
Math.tan(0)           // 0.0
Math.tan(Math.PI / 4) // 1.0
```

## Random Number Generation

### Math.random()
Returns a random float between 0.0 (inclusive) and 1.0 (exclusive).

**Returns**: `Float` - A random number

**Example**:
```polyloft
let r = Math.random()  // e.g., 0.6342341
let dice = Math.floor(Math.random() * 6) + 1  // Random number 1-6
```

## Practical Examples

### Distance Between Points
```polyloft
def distance(x1: Float, y1: Float, x2: Float, y2: Float) -> Float:
    let dx = x2 - x1
    let dy = y2 - y1
    return Math.sqrt(dx * dx + dy * dy)
end

let d = distance(0.0, 0.0, 3.0, 4.0)  // 5.0
```

### Circle Area
```polyloft
def circleArea(radius: Float) -> Float:
    return Math.PI * radius * radius
end

println(circleArea(5.0))  // 78.53981633974483
```

### Degrees to Radians
```polyloft
def toRadians(degrees: Float) -> Float:
    return degrees * Math.PI / 180.0
end

def toDegrees(radians: Float) -> Float:
    return radians * 180.0 / Math.PI
end

let rad = toRadians(90.0)   // 1.5707963267948966
let deg = toDegrees(rad)    // 90.0
```

### Random Range
```polyloft
def randomRange(min: Int, max: Int) -> Int:
    return Math.floor(Math.random() * (max - min + 1)) + min
end

let randomDice = randomRange(1, 6)
```

### Sigmoid Function
```polyloft
def sigmoid(x: Float) -> Float:
    return 1.0 / (1.0 + Math.pow(Math.E, -x))
end

println(sigmoid(0))   // 0.5
println(sigmoid(10))  // ~1.0
```

### Pythagorean Theorem
```polyloft
def hypotenuse(a: Float, b: Float) -> Float:
    return Math.sqrt(a * a + b * b)
end

let c = hypotenuse(3.0, 4.0)  // 5.0
```

## Notes

- All trigonometric functions use radians, not degrees
- `Math.random()` is seeded automatically with the current time
- For seeded random numbers, use `Sys.seed()` first
- Mathematical operations follow IEEE 754 floating-point standard
- Division by zero returns infinity or NaN as appropriate

## See Also

- [Sys Module](sys.md) - System-level random with seed control
- [Number Methods](number.md) - Built-in number methods
- [Variables and Types](../language/variables-and-types.md) - Numeric types
