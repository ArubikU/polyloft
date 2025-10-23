# Control-flow Documentation

This document covers control flow statements in Polyloft, including conditional statements, loops, and pattern matching.

## Table of Contents

- [If-Elif-Else](#if-elif-else)
- [For Loops](#for-loops)
- [While Loops](#while-loops)
- [Loop Statement](#loop-statement)
- [Switch Statement](#switch-statement)
- [Select Statement](#select-statement)
- [Break and Continue](#break-and-continue)

## If-Elif-Else

Conditional execution based on boolean expressions.

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

## For Loops

Iterate over ranges or collections.

### Range-based For Loop

```polyloft
// Iterate over a range
for i in range(5):
    println(i)  // Prints 0, 1, 2, 3, 4
end

// Iterate over an array
let fruits = ["apple", "banana", "orange"]
for fruit in fruits:
    println(fruit)
end
```

## While Loops

Execute a block while a condition is true.

```polyloft
let count = 0
while count < 5:
    println("Count: " + count.toString())
    count = count + 1
end
```

## Loop Statement

Infinite loop that can be exited with `break`.

```polyloft
let i = 0
loop:
    if i >= 5:
        break
    end
    println(i)
    i = i + 1
end
```

## Switch Statement

Pattern matching for values, types, and enums. The switch statement provides a clean way to handle multiple conditions.

### Value Matching

Match against specific values. Multiple values can be specified in a single case using commas.

```polyloft
let day = 3

switch day:
    case 1:
        println("Monday")
    case 2:
        println("Tuesday")
    case 3:
        println("Wednesday")
    case 4:
        println("Thursday")
    case 5:
        println("Friday")
    case 6, 7:  // Multiple values in one case
        println("Weekend")
    default:
        println("Invalid day")
end
```

### String Matching

```polyloft
let fruit = "apple"

switch fruit:
    case "apple":
        println("It's an apple")
    case "banana":
        println("It's a banana")
    case "orange", "tangerine":
        println("It's a citrus fruit")
    default:
        println("Unknown fruit")
end
```

### Type Matching

Match based on the runtime type of a value and bind it to a variable.

```polyloft
let value = 42

// Direct type matching with variable binding
switch value:
    case (x: Int):
        println("Integer: " + x.toString())
    case (s: String):
        println("String: " + s)
    case (f: Float):
        println("Float: " + f.toString())
    case (b: Bool):
        println("Boolean: " + b.toString())
    default:
        println("Unknown type")
end

// Type matching using Sys.type()
switch Sys.type(value):
    case "int", "Int":
        println("It's an integer")
    case "string", "String":
        println("It's a string")
    default:
        println("Other type")
end
```

### Enum Matching

Match against enum values.

```polyloft
enum Color:
    RED
    GREEN
    BLUE
end

let currentColor = Color.RED

switch currentColor:
    case Color.RED:
        println("The color is RED")
    case Color.GREEN:
        println("The color is GREEN")
    case Color.BLUE:
        println("The color is BLUE")
    default:
        println("Unknown color")
end
```

### Key Features

- **No Fall-through**: Unlike some languages (like C/Java), Polyloft's switch doesn't fall through to the next case. Each case executes only if matched, then exits the switch.
- **Multiple Values**: Use commas to match multiple values in a single case: `case 1, 2, 3:`
- **Type Safety**: Type matching with variable binding ensures type safety: `case (x: Int):`
- **Default Case**: The `default` case is optional and executes if no other case matches.

## Select Statement

Handle multiple channel operations. The select statement is used for concurrent programming with channels.

```polyloft
let ch1 = Channel<Int>()
let ch2 = Channel<String>()

select:
    case ch1.recv() as value:
        println("Got from ch1: #{value}")
    case ch2.recv() as msg:
        println("Got from ch2: #{msg}")
    case closed ch1:
        println("Channel ch1 is closed")
end
```

See [Concurrency Examples](../examples/concurrency.md) for more details on channels and select.

## Break and Continue

Control loop execution.

### Break

Exit the current loop immediately.

```polyloft
for i in range(10):
    if i == 5:
        break
    end
    println(i)  // Prints 0, 1, 2, 3, 4
end
```

### Continue

Skip to the next iteration of the loop.

```polyloft
for i in range(10):
    if i % 2 == 0:
        continue  // Skip even numbers
    end
    println(i)  // Prints 1, 3, 5, 7, 9
end
```

## See Also

- [Language Overview](overview.md)
- [Basic Syntax](syntax.md)
- [Functions](functions.md)
- [Enums](enums.md)
