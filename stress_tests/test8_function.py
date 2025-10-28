#!/usr/bin/env python3
# Test 8: Function calls
import time

def calculate(x, y):
    return x * x + y * y

start = time.time()
result = 0
for i in range(50000):
    result += calculate(i, i + 1)
end = time.time()

print(f"Result: {result}")
print(f"Time: {(end - start) * 1000:.2f} ms")
