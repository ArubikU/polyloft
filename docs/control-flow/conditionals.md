# Conditional Statements

Polyloft supports `if`, `else if`, and `else` for conditional branching.

## Syntax

### Basic If
```pf
if condition:
    // code if condition is true
end
```

### If-Else
```pf
if condition:
    // code if condition is true
else:
    // code if condition is false
end
```

### If-Else If-Else
```pf
if condition1:
    // code if condition1 is true
else:
    if condition2:
        // code if condition2 is true
    else:
        // code if all conditions are false
    end
end
```

## Examples

### Simple If
```pf
let age = 18
if age >= 18:
    println("You are an adult")
end
```

### If-Else
```pf
let age = 15
if age >= 18:
    println("Adult")
else:
    println("Minor")
end
```

### Multiple Conditions
```pf
let score = 85

if score >= 90:
    println("Grade: A")
else:
    if score >= 80:
        println("Grade: B")
    else:
        if score >= 70:
            println("Grade: C")
        else:
            if score >= 60:
                println("Grade: D")
            else:
                println("Grade: F")
            end
        end
    end
end
```

### Comparison Operators
```pf
let x = 10
let y = 20

if x == y:
    println("Equal")
end

if x != y:
    println("Not equal")
end

if x < y:
    println("x is less than y")
end

if x > y:
    println("x is greater than y")
end

if x <= y:
    println("x is less than or equal to y")
end

if x >= y:
    println("x is greater than or equal to y")
end
```

### Logical Operators
```pf
let age = 25
let hasLicense = true

// AND operator
if age >= 18 and hasLicense:
    println("Can drive")
end

// OR operator
if age < 18 or not hasLicense:
    println("Cannot drive")
end

// NOT operator
if not hasLicense:
    println("Need to get a license")
end
```

### Complex Conditions
```pf
let temperature = 25
let isRaining = false
let hasUmbrella = true

if temperature > 20 and not isRaining:
    println("Perfect day for a walk!")
else:
    if isRaining and hasUmbrella:
        println("Can still go out with umbrella")
    else:
        println("Better stay inside")
    end
end
```

### Nested Conditionals
```pf
let number = 15

if number > 0:
    if number % 2 == 0:
        println("Positive even number")
    else:
        println("Positive odd number")
    end
else:
    if number < 0:
        if number % 2 == 0:
            println("Negative even number")
        else:
            println("Negative odd number")
        end
    else:
        println("Zero")
    end
end
```

### Truthy and Falsy Values
```pf
// False values: false, nil, 0, empty string
// True values: everything else

if "":
    println("Empty string is truthy")  // Won't print
end

if "hello":
    println("Non-empty string is truthy")  // Will print
end

if 0:
    println("Zero is truthy")  // Won't print
end

if 1:
    println("Non-zero is truthy")  // Will print
end

if nil:
    println("nil is truthy")  // Won't print
end
```

### Ternary-like Pattern
```pf
// Polyloft doesn't have ternary operator, but you can use functions
def max(a, b):
    if a > b:
        return a
    else:
        return b
    end
end

let result = max(10, 20)
println(result)  // 20
```

## Best Practices

### ✅ DO - Use clear conditions
```pf
if user.isActive() and user.hasPermission("write"):
    allowEdit()
end
```

### ✅ DO - Extract complex conditions
```pf
def canAccessResource(user, resource):
    return user.isLoggedIn() and 
           user.hasPermission(resource) and 
           not resource.isLocked()
end

if canAccessResource(currentUser, document):
    showDocument()
end
```

### ✅ DO - Use guard clauses
```pf
def processUser(user):
    if user == nil:
        return
    end
    
    if not user.isValid():
        return
    end
    
    // Process valid user
    println("Processing #{user.name}")
end
```

### ❌ DON'T - Nest too deeply
```pf
// Bad
if condition1:
    if condition2:
        if condition3:
            if condition4:
                // Too deep!
            end
        end
    end
end

// Better
def checkConditions():
    if not condition1:
        return false
    end
    if not condition2:
        return false
    end
    if not condition3:
        return false
    end
    return condition4
end

if checkConditions():
    // Process
end
```

## See Also

- [Loops](loops.md)
- [Switch Statement](switch.md)
- [Exception Handling](exceptions.md)
