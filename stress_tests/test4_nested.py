#!/usr/bin/env python3
# Test 4: Nested loops
import time

start = time.time()
total = 0
for i in range(1, 501):
    for j in range(1, 501):
        total += i * j
end = time.time()

print(f"Total: {total}")
print(f"Time: {(end - start) * 1000:.2f} ms")
