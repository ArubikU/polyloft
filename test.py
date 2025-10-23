import time

# time.time() returns seconds as a float
start = time.time()
for i in range(1000000):
    pass
elapsed = time.time() - start

# time.time() is in seconds, so no conversion needed
print(f"Elapsed time: {elapsed:.6f} seconds")
# 0.310948 seconds