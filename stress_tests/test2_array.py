#!/usr/bin/env python3
# Test 2: Array operations
import time

start = time.time()
arr = []
for i in range(100000):
    arr.append(i)

total = sum(arr)
end = time.time()

print(f"Array length: {len(arr)}")
print(f"Sum: {total}")
print(f"Time: {(end - start) * 1000:.2f} ms")
