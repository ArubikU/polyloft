# Stress Test Results

Comparison of Python 3 vs Polyloft performance

| Test | Description | Python (ms) | Polyloft (ms) | Ratio (Poly/Py) |
|------|-------------|-------------|---------------|-----------------|
| Test 1 | loop | ERROR | ERROR | - |
| Test 2 | array | 6.82 | 0 | 0x |
| Test 3 | string | ERROR | ERROR | - |
| Test 4 | nested | ERROR | ERROR | - |
| Test 5 | factorial | ERROR | ERROR | - |
| Test 6 | map | 2.44 | 0 | 0x |
| Test 7 | conditional | ERROR | ERROR | - |
| Test 8 | function | ERROR | ERROR | - |
| Test 9 | class | 4.35 | 0 | 0x |
| Test 10 | fibonacci | ERROR | ERROR | - |

## Analysis

- Ratio < 1.0: Polyloft is faster
- Ratio > 1.0: Python is faster
- The ratio shows how many times slower/faster Polyloft is compared to Python
