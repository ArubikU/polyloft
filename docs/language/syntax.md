# Basic Syntax

This guide covers the fundamental syntax rules of the Polyloft programming language.

## Comments

### Single-Line Comments
```polyloft
// This is a single-line comment
let x = 42  // Comments can appear after code
```

### Multi-Line Comments
```polyloft
/*
 * This is a multi-line comment
 * It can span multiple lines
 */
let y = 10
```

## Variables

### Declaration Keywords

#### `let` - Mutable Variable
```polyloft
let x = 42
x = 50  // OK: can be reassigned

let name: String = "Alice"
name = "Bob"  // OK
```

#### `const` - Constant (Immutable)
```polyloft
const PI = 3.14159
// PI = 3.14  // Error: cannot reassign constant

const greeting = "Hello"
```

#### `var` - Uninitialized Variable
```polyloft
var count: Int
count = 10  // Must assign before use

var name: String
// println(name)  // Error: used before assignment
```

### Type Annotations

Optional but recommended for clarity:

```polyloft
// With type annotation
let age: Int = 30
let price: Float = 19.99
let name: String = "Alice"

// Type inference
let x = 42           // Inferred as Int
let y = 3.14        // Inferred as Float
let z = "Hello"     // Inferred as String
```

## Literals

### Numbers

```polyloft
// Integers
let decimal = 42
let negative = -100

// Floats
let pi = 3.14159
let scientific = 1.5e10
let negative_float = -2.5
```

### Strings

```polyloft
// Single quotes
let single = 'Hello'

// Double quotes
let double = "World"

// String interpolation
let name = "Alice"
let greeting = "Hello, #{name}!"  // "Hello, Alice!"

// Multi-line strings
let multiline = "Line 1
Line 2
Line 3"

// Escape sequences
let escaped = "Quote: \" Newline: \n Tab: \t"
```

### Booleans

```polyloft
let isTrue = true
let isFalse = false
```

### Arrays

```polyloft
// Array literal
let numbers = [1, 2, 3, 4, 5]
let mixed = [1, "two", 3.0, true]

// Empty array
let empty = []

// Access elements
let first = numbers[0]  // 1
let last = numbers[4]   // 5

// Modify elements
numbers[0] = 10
```

### Maps

```polyloft
// Map literal (object notation)
let person = {
    "name": "Alice",
    "age": 30,
    "city": "New York"
}

// Access properties
let name = person["name"]
let age = person.age

// Modify properties
person["age"] = 31
person.city = "Boston"
```

## Operators

### Arithmetic Operators

```polyloft
let a = 10
let b = 3

let sum = a + b        // 13
let diff = a - b       // 7
let product = a * b    // 30
let quotient = a / b   // 3 (integer division)
let remainder = a % b  // 1
```

### Comparison Operators

```polyloft
let x = 10
let y = 20

x == y   // false (equal)
x != y   // true (not equal)
x < y    // true (less than)
x <= y   // true (less than or equal)
x > y    // false (greater than)
x >= y   // false (greater than or equal)
```

### Logical Operators

```polyloft
let a = true
let b = false

a and b   // false (logical AND)
a or b    // true (logical OR)
not a     // false (logical NOT)

// Short-circuit evaluation
true or expensiveFunction()   // expensiveFunction not called
false and expensiveFunction() // expensiveFunction not called
```

### String Operators

```polyloft
let first = "Hello"
let second = "World"

let combined = first + " " + second  // "Hello World"

// String interpolation (preferred)
let name = "Alice"
let greeting = "Hello, #{name}!"
```

## Control Flow

### If-Else Statements

```polyloft
let age = 18

if age >= 18:
    println("Adult")
elif age >= 13:
    println("Teenager")
else:
    println("Child")
end
```

### While Loops

```polyloft
let count = 0

let count = 0
loop
    println(count)
    count = count + 1
end
```

### For Loops

```polyloft
// Range-based for loop
for i in range(5):
    println(i)  // 0, 1, 2, 3, 4
end

// Range with start and end
for i in range(1, 6):
    println(i)  // 1, 2, 3, 4, 5
end

// Range with step
for i in range(0, 10, 2):
    println(i)  // 0, 2, 4, 6, 8
end

// Iterate over array
let fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    println(fruit)
end
```

### Break and Continue

```polyloft
// Break: exit loop
for i in range(10):
    if i == 5:
        break
    end
    println(i)  // 0, 1, 2, 3, 4
end

// Continue: skip to next iteration
for i in range(10):
    if i % 2 == 0:
        continue
    end
    println(i)  // 1, 3, 5, 7, 9
end
```

## Functions

### Function Definition

```polyloft
// Basic function
def greet():
    println("Hello!")
end

// Function with parameters
def greet(name: String):
    println("Hello, #{name}!")
end

// Function with return type
def add(a: Int, b: Int) -> Int:
    return a + b
end

// Function with default parameters
def greet(name: String, greeting: String):
    if not greeting:
        greeting = "Hello"
    end
    println("#{greeting}, #{name}!")
end
```

### Lambda Expressions

```polyloft
// Lambda syntax
let add = (a: Int, b: Int) -> Int => a + b

// Lambda with single parameter
let square = (x: Int) -> Int => x * x

// Lambda with multiple statements
let process = (x: Int) -> Int => {
    let doubled = x * 2
    let squared = doubled * doubled
    return squared
}

// Lambda without parameters
let greet = () -> Void => println("Hello!")
```

### Function Calls

```polyloft
greet()
greet("Alice")
let sum = add(5, 3)
```

## Classes

### Basic Class

```polyloft
class Person:
    Person(name: String, age: Int):
        this.name = name
        this.age = age
    end
    
    def greet():
        println("Hello, I'm #{this.name}")
    end
end

// Create instance
let person = Person("Alice", 30)
person.greet()
```

### Class with Fields

```polyloft
class Rectangle:
    private let width: Float
    private let height: Float
    
    Rectangle(width: Float, height: Float):
        this.width = width
        this.height = height
    end
    
    def area() -> Float:
        return this.width * this.height
    end
end
```

## Interfaces

```polyloft
interface Drawable:
    def draw(): Void
end

class Circle implements Drawable:
    def draw():
        println("Drawing circle")
    end
end
```

## Blocks and Scope

### Block Structure

Polyloft uses `end` keywords to denote block endings:

```polyloft
if condition:
    // block
end

loop
    // block
end

for i in range(10):
    // block
end

def function():
    // block
end

class MyClass:
    // block
end
```

### Scope Rules

```polyloft
let x = 10  // Global scope

def function():
    let y = 20  // Function scope
    println(x)  // Can access global
    
    if true:
        let z = 30  // Block scope
        println(y)  // Can access function scope
    end
    
    // println(z)  // Error: z not in scope
end

// println(y)  // Error: y not in scope
```

## Indentation

Polyloft is indentation-aware but not required:

```polyloft
// Preferred: indented
if condition:
    doSomething()
    doSomethingElse()
end

// Also valid: not indented
if condition:
doSomething()
doSomethingElse()
end
```

## Statement Separators

Newlines separate statements:

```polyloft
let x = 10
let y = 20
let z = x + y
```

Multiple statements on one line (not recommended):

```polyloft
let x = 10; let y = 20; let z = x + y
```

## Naming Conventions

### Variables and Functions
- Use camelCase: `myVariable`, `calculateTotal`
- Descriptive names: `userCount` not `uc`

### Classes and Interfaces
- Use PascalCase: `MyClass`, `Drawable`
- Noun-based names: `User`, `Animal`

### Constants
- Use UPPER_SNAKE_CASE: `MAX_SIZE`, `DEFAULT_TIMEOUT`

## Reserved Keywords

```
and        as         async      await      break
case       catch      chan       class      const
continue   def        default    defer      do
elif       else       end        enum       extends
false      finally    for        if         implements
import     in         interface  is         let
new        not        null       or         private
protected  public     range      return     sealed
select     spawn      static     super      this
thread     throw      true       try        type
var        while      with
```

## Imports

Polyloft supports importing modules and packages to organize code and reuse functionality.

### Basic Import

Import a module by its dotted path:

```polyloft
// Import an entire module
import math

// Use imported functions
let result = math.sqrt(16)
```

### Selective Import

Import specific items from a module:

```polyloft
// Import specific functions or classes
import math { sqrt, pow, PI }

// Use imported items directly
let result = sqrt(16)
let squared = pow(2, 3)
let circle = PI * 2
```

### Import Syntax

The general syntax for imports is:

```polyloft
// Full module import
import module.name

// Selective import
import module.name { Item1, Item2, Item3 }

// Nested module import
import package.subpackage.module

// Selective import from nested module
import package.subpackage.module { ClassName, functionName }
```

### Using Imports

After importing, you can use the imported items in your code:

```polyloft
// Import from a custom module
import utils.string { capitalize, reverse }

let text = "hello world"
let capitalized = capitalize(text)  // "Hello World"
let reversed = reverse(text)         // "dlrow olleh"
```

### Standard Library Imports

Polyloft provides a standard library with various utility modules:

```polyloft
// Math utilities
import math { sqrt, pow, sin, cos }

// System utilities (Sys is available by default)
// No need to import Sys - it's a built-in global
let time = Sys.time()
let type = Sys.type(42)

// File I/O
import io { readFile, writeFile }

let content = readFile("data.txt")
writeFile("output.txt", "Hello, World!")
```

### Package Management

For information about installing and managing packages from the registry, see the main [README](../../README.md#package-registry).

## See Also

- [Variables and Types](variables-and-types.md) - Type system details
- [Control Flow](control-flow.md) - Control structures
- [Functions](functions.md) - Function features
- [Classes and Objects](classes-and-objects.md) - OOP features

