#!/usr/bin/env python3
# Test 10: Fibonacci sequence
import time

def fib(n):
    if n <= 1:
        return n
    return fib(n - 1) + fib(n - 2)

start = time.time()
result = 0
for i in range(25):
    result += fib(i)
end = time.time()

print(f"Result: {result}")
print(f"Time: {(end - start) * 1000:.2f} ms")
