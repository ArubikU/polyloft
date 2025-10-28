#!/usr/bin/env python3
# Test 3: String concatenation
import time

start = time.time()
result = ""
for i in range(10000):
    result += str(i)
end = time.time()

print(f"String length: {len(result)}")
print(f"Time: {(end - start) * 1000:.2f} ms")
