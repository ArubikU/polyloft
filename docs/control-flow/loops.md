# Loops

Polyloft supports `for` and `loop` statements for iteration.

## For Loops

### Range-based For Loop
```pf
for i in range(10):
    println(i)  // 0, 1, 2, ..., 9
end
```

### Array/List Iteration
```pf
let fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    println(fruit)
end
```

### Range with Start and End
```pf
for i in range(1, 6):
    println(i)  // 1, 2, 3, 4, 5
end
```

### Range with Step
```pf
for i in range(0, 10, 2):
    println(i)  // 0, 2, 4, 6, 8
end
```

### Inclusive Range (...)
```pf
for i in 1...5:
    println(i)  // 1, 2, 3, 4, 5
end
```

### Map Iteration
```pf
let person = {name: "Alice", age: 25, city: "NYC"}
for key in person:
    println("#{key}: #{person[key]}")
end
```

### Map Iteration with Destructuring
You can destructure key-value pairs while iterating over maps:
```pf
let person = {name: "Alice", age: 25, city: "NYC"}
for key, value in person:
    println("#{key}: #{value}")
end
```

### Array Destructuring
When iterating over arrays of arrays, you can destructure each element:
```pf
let pairs = [[1, 2], [3, 4], [5, 6]]
for a, b in pairs:
    println("#{a} + #{b} = #{a + b}")
end
// Outputs:
// 1 + 2 = 3
// 3 + 4 = 7
// 5 + 6 = 11
```

### For Loop with Where Clause
The `where` clause allows you to filter elements during iteration:
```pf
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let sum = 0

for n in numbers where n > 5:
    sum = sum + n
end

println("Sum: #{sum}")  // Sum: 40 (6+7+8+9+10)
```

### Where Clause with Modulo
```pf
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evenSum = 0

for n in numbers where n % 2 == 0:
    evenSum = evenSum + n
end

println("Even sum: #{evenSum}")  // Even sum: 30
```

### Combining Destructuring and Where Clause
You can use both destructuring and filtering together:
```pf
let data = {a: 10, b: 5, c: 20, d: 3}
let sum = 0

for key, value in data where value > 5:
    sum = sum + value
end

println("Filtered sum: #{sum}")  // Filtered sum: 30 (10+20)
```

### Where Clause with Complex Conditions
```pf
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let count = 0

for n in numbers where (n > 3) && (n < 8):
    count = count + 1
end

println("Count: #{count}")  // Count: 4 (numbers 4, 5, 6, 7)
```

### Break Statement
```pf
for i in range(10):
    if i == 5:
        break  // Exit loop when i is 5
    end
    println(i)
end
// Outputs: 0, 1, 2, 3, 4
```

### Continue Statement
```pf
for i in range(10):
    if i % 2 == 0:
        continue  // Skip even numbers
    end
    println(i)
end
// Outputs: 1, 3, 5, 7, 9
```

## Loop Statement

The `loop` statement provides conditional looping (similar to while in other languages).

### Infinite Loop
```pf
// Old style (backward compatible)
loop
    println("Forever")
    if shouldStop:
        break
    end
end

// New style with colon
loop:
    println("Forever")
    if shouldStop:
        break
    end
end
```

### Loop with Condition
```pf
let count = 0
loop count < 5:
    println(count)
    count = count + 1
end
// Outputs: 0, 1, 2, 3, 4
```

### Loop with Break
```pf
let i = 0
loop:
    if i >= 5:
        break
    end
    println(i)
    i = i + 1
end
```

### Loop with Continue
```pf
let i = 0
loop i < 10:
    i = i + 1
    if i % 2 == 0:
        continue
    end
    println(i)
end
// Outputs: 1, 3, 5, 7, 9
```

## Do-Loop Statement

The `do-loop` statement executes the block at least once before checking the condition.

### Basic Do-Loop
```pf
let count = 0
do:
    println(count)
    count = count + 1
loop count < 5
// Outputs: 0, 1, 2, 3, 4
```

### Do-Loop Always Executes Once
```pf
let x = 10
do:
    println("This runs once")
    x = x + 1
loop x < 5
// Outputs: "This runs once" (even though condition is false)
```

## Examples

### Sum of Array Elements
```pf
let numbers = [1, 2, 3, 4, 5]
let sum = 0

for num in numbers:
    sum = sum + num
end

println("Sum: #{sum}")  // Sum: 15
```

### Find Maximum Value
```pf
let numbers = [23, 45, 12, 67, 34]
let max = numbers[0]

for num in numbers:
    if num > max:
        max = num
    end
end

println("Maximum: #{max}")  // Maximum: 67
```

### Count Occurrences
```pf
let text = "hello world"
let target = "l"
let count = 0

for char in text:
    if char == target:
        count = count + 1
    end
end

println("'#{target}' appears #{count} times")
```

### Nested Loops
```pf
// Multiplication table
for i in range(1, 6):
    for j in range(1, 6):
        println("#{i} x #{j} = #{i * j}")
    end
    println("---")
end
```

### Matrix Traversal
```pf
let matrix = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
]

for row in matrix:
    for value in row:
        Sys.print("#{value} ")
    end
    println("")
end
```

### Filtering
```pf
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evens = []

for num in numbers:
    if num % 2 == 0:
        evens = evens.concat([num])
    end
end

println(evens)  // [2, 4, 6, 8, 10]
```

### Accumulation
```pf
let numbers = [1, 2, 3, 4, 5]
let product = 1

for num in numbers:
    product = product * num
end

println("Product: #{product}")  // Product: 120
```

### Countdown
```pf
let countdown = 5
loop countdown > 0:
    println(countdown)
    countdown = countdown - 1
    Sys.sleep(1000)
end
println("Blast off!")
```

### Input Loop
```pf
loop:
    let input = Sys.input("Enter command (quit to exit): ")
    
    if input == "quit":
        break
    end
    
    println("You entered: #{input}")
end
```

### Search and Stop
```pf
let numbers = [10, 23, 45, 67, 89, 12, 34]
let target = 67
let found = false
let index = -1

for i in range(numbers.length()):
    if numbers[i] == target:
        found = true
        index = i
        break
    end
end

if found:
    println("Found #{target} at index #{index}")
else:
    println("#{target} not found")
end
```

### Building Strings
```pf
let words = ["Hello", "beautiful", "world"]
let sentence = ""

for i in range(words.length()):
    sentence = sentence + words[i]
    if i < words.length() - 1:
        sentence = sentence + " "
    end
end

println(sentence)  // Hello beautiful world
```

### Fibonacci Sequence
```pf
let n = 10
let a = 0
let b = 1

println("Fibonacci sequence:")
for i in range(n):
    println(a)
    let temp = a + b
    a = b
    b = temp
end
```

### User Input Validation
```pf
let valid = false
let attempts = 0
let maxAttempts = 3

loop not valid and attempts < maxAttempts:
    let input = Sys.input("Enter a number between 1-10: ")
    attempts = attempts + 1
    
    // Parse input as int (simplified)
    if input >= "1" and input <= "9":
        valid = true
        println("Valid input!")
    else:
        println("Invalid! Try again.")
    end
end

if not valid:
    println("Too many failed attempts")
end
```

### Menu System
```pf
loop:
    println("\nMenu:")
    println("1. Option A")
    println("2. Option B")
    println("3. Exit")
    
    let choice = Sys.input("Select: ")
    
    if choice == "1":
        println("You selected Option A")
    else:
        if choice == "2":
            println("You selected Option B")
        else:
            if choice == "3":
                println("Goodbye!")
                break
            else:
                println("Invalid choice")
            end
        end
    end
end
```

## Best Practices

### ✅ DO - Use for loops for known iterations
```pf
// Good: Known number of iterations
for i in range(10):
    process(i)
end
```

### ✅ DO - Use loop for unknown iterations
```pf
// Good: Continue until condition met
loop not finished:
    processData()
    finished = checkIfDone()
end
```

### ✅ DO - Use break to exit early
```pf
for item in items:
    if item == target:
        found = true
        break  // Stop searching
    end
end
```

### ✅ DO - Use do-loop when you need at least one iteration
```pf
do:
    let input = Sys.input("Enter value: ")
    processInput(input)
loop askAgain()
```

### ❌ DON'T - Modify collection while iterating
```pf
// Bad: Modifying during iteration
let numbers = [1, 2, 3, 4, 5]
for num in numbers:
    numbers.remove(num)  // Can cause issues
end

// Better: Create new collection
let numbers = [1, 2, 3, 4, 5]
let filtered = []
for num in numbers:
    if num % 2 != 0:
        filtered = filtered.concat([num])
    end
end
```

### ❌ DON'T - Create infinite loops accidentally
```pf
// Bad: Forgot to increment
let i = 0
loop i < 10:
    println(i)
    // Missing: i = i + 1
end

// Good: Ensure loop progresses
let i = 0
loop i < 10:
    println(i)
    i = i + 1
end
```

## Common Patterns

### Loop with Index
```pf
let fruits = ["apple", "banana", "cherry"]
for i in range(fruits.length()):
    println("#{i}: #{fruits[i]}")
end
```

### Reverse Iteration
```pf
let numbers = [1, 2, 3, 4, 5]
let i = numbers.length() - 1
loop i >= 0:
    println(numbers[i])
    i = i - 1
end
```

### Loop Until Success
```pf
let success = false
let attempts = 0
let maxAttempts = 3

loop not success and attempts < maxAttempts:
    success = tryOperation()
    attempts = attempts + 1
    
    if not success:
        println("Attempt #{attempts} failed")
        Sys.sleep(1000)
    end
end
```

### Polling Pattern
```pf
loop:
    let status = checkStatus()
    
    if status == "ready":
        break
    end
    
    println("Waiting...")
    Sys.sleep(1000)
end
println("Ready!")
```

## See Also

- [Conditionals](conditionals.md)
- [Arrays](../types/array.md)
- [Range Type](../types/range.md)
- [Builtin Interfaces](../language/builtin-interfaces.md)
