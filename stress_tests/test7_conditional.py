#!/usr/bin/env python3
# Test 7: Conditional logic
import time

start = time.time()
count = 0
for i in range(100000):
    if i % 2 == 0:
        count += 1
    elif i % 3 == 0:
        count += 2
    else:
        count += 3
end = time.time()

print(f"Count: {count}")
print(f"Time: {(end - start) * 1000:.2f} ms")
