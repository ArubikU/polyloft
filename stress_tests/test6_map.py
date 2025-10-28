#!/usr/bin/env python3
# Test 6: Map/dictionary operations
import time

start = time.time()
data = {}
for i in range(10000):
    data[str(i)] = i * 2

total = sum(data.values())
end = time.time()

print(f"Map size: {len(data)}")
print(f"Sum: {total}")
print(f"Time: {(end - start) * 1000:.2f} ms")
