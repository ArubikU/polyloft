#!/usr/bin/env python3
# Test 5: Factorial computation
import time

def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)

start = time.time()
result = 0
for i in range(1, 101):
    result += factorial(i % 20)  # Keep recursion depth reasonable
end = time.time()

print(f"Result: {result}")
print(f"Time: {(end - start) * 1000:.2f} ms")
