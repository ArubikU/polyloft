#!/usr/bin/env python3
# Test 9: Class instantiation and method calls
import time

class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def distance(self):
        return (self.x * self.x + self.y * self.y) ** 0.5

start = time.time()
total = 0
for i in range(10000):
    p = Point(i, i + 1)
    total += p.distance()
end = time.time()

print(f"Total: {total:.2f}")
print(f"Time: {(end - start) * 1000:.2f} ms")
