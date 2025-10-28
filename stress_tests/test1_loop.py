#!/usr/bin/env python3
# Test 1: Simple loop with arithmetic
import time

start = time.time()
result = 0
for i in range(1, 1000001):
    result += i * i
end = time.time()

print(f"Result: {result}")
print(f"Time: {(end - start) * 1000:.2f} ms")
