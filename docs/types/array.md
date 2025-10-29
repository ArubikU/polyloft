# Array Type

Arrays in Polyloft are dynamic, mutable collections of elements.

## Creating Arrays

### Array Literals
```pf
let empty = []
let numbers = [1, 2, 3, 4, 5]
let mixed = [1, "text", true, 3.14]
let nested = [[1, 2], [3, 4], [5, 6]]
```

### Array Constructor
```pf
let arr = Array()
arr.add(1)
arr.add(2)
arr.add(3)
```

## Properties

### `length()`
Returns the number of elements in the array.

**Returns:** Int

```pf
let arr = [1, 2, 3, 4, 5]
println(arr.length())  // 5
```

### `isEmpty()`
Checks if the array is empty.

**Returns:** Bool

```pf
println([].isEmpty())        // true
println([1, 2].isEmpty())    // false
```

## Modification Methods

### `add(element)` / `push(element)`
Adds an element to the end of the array.

**Parameters:**
- `element` (Any): Element to add

**Returns:** void

```pf
let arr = [1, 2, 3]
arr.add(4)
arr.push(5)
println(arr)  // [1, 2, 3, 4, 5]
```

### `pop()`
Removes and returns the last element.

**Returns:** Any

```pf
let arr = [1, 2, 3]
let last = arr.pop()
println(last)  // 3
println(arr)   // [1, 2]
```

### `shift()`
Removes and returns the first element.

**Returns:** Any

```pf
let arr = [1, 2, 3]
let first = arr.shift()
println(first)  // 1
println(arr)    // [2, 3]
```

### `unshift(element)`
Adds an element to the beginning of the array.

**Parameters:**
- `element` (Any): Element to add

**Returns:** void

```pf
let arr = [2, 3]
arr.unshift(1)
println(arr)  // [1, 2, 3]
```

### `set(index, value)`
Sets the value at the specified index.

**Parameters:**
- `index` (Int): Position
- `value` (Any): New value

**Returns:** void

```pf
let arr = [1, 2, 3]
arr.set(1, 10)
println(arr)  // [1, 10, 3]
```

### `reverse()`
Reverses the array in place.

**Returns:** void

```pf
let arr = [1, 2, 3, 4, 5]
arr.reverse()
println(arr)  // [5, 4, 3, 2, 1]
```

### `sort()`
Sorts the array in place.

**Returns:** void

```pf
let arr = [3, 1, 4, 1, 5, 9]
arr.sort()
println(arr)  // [1, 1, 3, 4, 5, 9]
```

### `clear()`
Removes all elements from the array.

**Returns:** void

```pf
let arr = [1, 2, 3]
arr.clear()
println(arr)  // []
```

## Access Methods

### `get(index)`
Returns the element at the specified index.

**Parameters:**
- `index` (Int): Position

**Returns:** Any

```pf
let arr = [10, 20, 30]
println(arr.get(1))  // 20
```

### `indexOf(element)`
Returns the index of the first occurrence of element, or -1 if not found.

**Parameters:**
- `element` (Any): Element to find

**Returns:** Int

```pf
let arr = [10, 20, 30, 20]
println(arr.indexOf(20))   // 1
println(arr.indexOf(999))  // -1
```

### `contains(element)`
Checks if the array contains the element.

**Parameters:**
- `element` (Any): Element to check

**Returns:** Bool

```pf
let arr = [1, 2, 3]
println(arr.contains(2))  // true
println(arr.contains(5))  // false
```

### `slice(start, end?)`
Returns a new array containing elements from start to end (exclusive).

**Parameters:**
- `start` (Int): Starting index
- `end` (Int, optional): Ending index

**Returns:** Array

```pf
let arr = [1, 2, 3, 4, 5]
println(arr.slice(1, 3))  // [2, 3]
println(arr.slice(2))     // [3, 4, 5]
```

## Transformation Methods

### `map(function)`
Creates a new array with the results of calling the function on every element.

**Parameters:**
- `function` (Function): Transformation function

**Returns:** Array

```pf
let numbers = [1, 2, 3, 4, 5]
let doubled = numbers.map((x) => x * 2)
println(doubled)  // [2, 4, 6, 8, 10]
```

### `filter(predicate)`
Creates a new array with elements that pass the test.

**Parameters:**
- `predicate` (Function): Test function

**Returns:** Array

```pf
let numbers = [1, 2, 3, 4, 5, 6]
let evens = numbers.filter((x) => x % 2 == 0)
println(evens)  // [2, 4, 6]
```

### `reduce(function, initial)`
Reduces the array to a single value.

**Parameters:**
- `function` (Function): Reducer function
- `initial` (Any): Initial value

**Returns:** Any

```pf
let numbers = [1, 2, 3, 4, 5]
let sum = numbers.reduce((acc, x) => acc + x, 0)
println(sum)  // 15
```

### `concat(array)`
Returns a new array by concatenating arrays.

**Parameters:**
- `array` (Array): Array to concatenate

**Returns:** Array

```pf
let arr1 = [1, 2, 3]
let arr2 = [4, 5, 6]
let combined = arr1.concat(arr2)
println(combined)  // [1, 2, 3, 4, 5, 6]
```

### `join(separator)`
Joins all elements into a string with the separator.

**Parameters:**
- `separator` (String): Separator string

**Returns:** String

```pf
let arr = ["Hello", "beautiful", "world"]
println(arr.join(" "))  // "Hello beautiful world"
println(arr.join(", "))  // "Hello, beautiful, world"
```

## Iteration Methods

### `forEach(function)`
Executes a function for each element.

**Parameters:**
- `function` (Function): Function to execute

**Returns:** void

```pf
let numbers = [1, 2, 3, 4, 5]
numbers.forEach((x) => println(x * x))
// Outputs: 1, 4, 9, 16, 25
```

### `find(predicate)`
Returns the first element that satisfies the predicate.

**Parameters:**
- `predicate` (Function): Test function

**Returns:** Any (or nil if not found)

```pf
let numbers = [1, 2, 3, 4, 5]
let found = numbers.find((x) => x > 3)
println(found)  // 4
```

### `findIndex(predicate)`
Returns the index of the first element that satisfies the predicate.

**Parameters:**
- `predicate` (Function): Test function

**Returns:** Int (or -1 if not found)

```pf
let numbers = [1, 2, 3, 4, 5]
let index = numbers.findIndex((x) => x > 3)
println(index)  // 3
```

### `every(predicate)`
Tests whether all elements pass the test.

**Parameters:**
- `predicate` (Function): Test function

**Returns:** Bool

```pf
let numbers = [2, 4, 6, 8]
println(numbers.every((x) => x % 2 == 0))  // true
```

### `some(predicate)`
Tests whether at least one element passes the test.

**Parameters:**
- `predicate` (Function): Test function

**Returns:** Bool

```pf
let numbers = [1, 3, 5, 6]
println(numbers.some((x) => x % 2 == 0))  // true
```

## Examples

### Array Indexing
```pf
let arr = [10, 20, 30, 40, 50]
println(arr[0])   // 10
println(arr[2])   // 30
println(arr[-1])  // 50 (last element)
```

### Array Slicing
```pf
let arr = [1, 2, 3, 4, 5]
println(arr[1:3])  // [2, 3]
println(arr[:3])   // [1, 2, 3]
println(arr[2:])   // [3, 4, 5]
```

### Iterating Arrays
```pf
let fruits = ["apple", "banana", "cherry"]

// With for loop
for fruit in fruits:
    println(fruit)
end

// With index
for i in range(fruits.length()):
    println("#{i}: #{fruits[i]}")
end
```

### Building Arrays
```pf
let squares = []
for i in range(1, 6):
    squares = squares.concat([i * i])
end
println(squares)  // [1, 4, 9, 16, 25]
```

### Filtering and Mapping
```pf
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let result = numbers
    .filter((x) => x % 2 == 0)
    .map((x) => x * x)
println(result)  // [4, 16, 36, 64, 100]
```

### Finding Elements
```pf
let users = [
    {name: "Alice", age: 25},
    {name: "Bob", age: 30},
    {name: "Charlie", age: 35}
]

let bob = users.find((u) => u.name == "Bob")
println(bob.age)  // 30
```

### Reducing to Sum
```pf
let numbers = [1, 2, 3, 4, 5]
let sum = numbers.reduce((total, n) => total + n, 0)
println(sum)  // 15
```

### Flattening Arrays
```pf
let nested = [[1, 2], [3, 4], [5, 6]]
let flat = nested.reduce((acc, arr) => acc.concat(arr), [])
println(flat)  // [1, 2, 3, 4, 5, 6]
```

## Best Practices

### ✅ DO - Use array methods for transformations
```pf
// Good: Functional approach
let doubled = numbers.map((x) => x * 2)
```

### ✅ DO - Check bounds before accessing
```pf
if index >= 0 and index < arr.length():
    let value = arr[index]
end
```

### ❌ DON'T - Modify array while iterating
```pf
// Bad: Can cause issues
for item in arr:
    arr.pop()  // Modifying while iterating
end

// Better: Create new array
let filtered = arr.filter((item) => shouldKeep(item))
```

## See Also

- [List Type](list.md)
- [String Type](string.md)
- [Map Type](map.md)
- [Loops](../control-flow/loops.md)
