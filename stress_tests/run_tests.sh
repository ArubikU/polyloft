#!/bin/bash

# Stress Test Runner - Compares Python vs Polyloft performance

echo "==================================="
echo "STRESS TEST COMPARISON: Python vs Polyloft"
echo "==================================="
echo ""

cd "$(dirname "$0")"

# Check if polyloft binary exists
if [ ! -f "../polyloft" ]; then
    echo "Building polyloft..."
    cd ..
    go build -o polyloft cmd/polyloft/main.go
    cd stress_tests
fi

results_file="results.md"
echo "# Stress Test Results" > $results_file
echo "" >> $results_file
echo "Comparison of Python 3 vs Polyloft performance" >> $results_file
echo "" >> $results_file
echo "| Test | Description | Python (ms) | Polyloft (ms) | Ratio (Poly/Py) |" >> $results_file
echo "|------|-------------|-------------|---------------|-----------------|" >> $results_file

for i in {1..10}; do
    py_file="test${i}_*.py"
    pf_file="test${i}_*.pf"
    
    py_actual=$(ls test${i}_*.py 2>/dev/null | head -1)
    pf_actual=$(ls test${i}_*.pf 2>/dev/null | head -1)
    
    if [ -z "$py_actual" ] || [ -z "$pf_actual" ]; then
        continue
    fi
    
    # Get test name
    test_name=$(basename $py_actual .py | sed 's/test[0-9]*_//')
    
    echo "Running Test $i: $test_name"
    echo "  Python..."
    
    # Run Python test 3 times and take average
    py_sum=0
    py_runs=3
    for run in {1..3}; do
        py_output=$(python3 "$py_actual" 2>&1)
        py_time=$(echo "$py_output" | grep "Time:" | awk '{print $2}' | sed 's/ms//')
        if [ ! -z "$py_time" ]; then
            py_sum=$(echo "$py_sum + $py_time" | bc)
        fi
    done
    py_avg=$(echo "scale=2; $py_sum / $py_runs" | bc)
    
    echo "  Polyloft..."
    
    # Run Polyloft test 3 times and take average
    pf_sum=0
    pf_runs=3
    for run in {1..3}; do
        pf_output=$(../polyloft run "$pf_actual" 2>&1)
        pf_time=$(echo "$pf_output" | grep "Time:" | awk '{print $2}' | sed 's/ms//')
        if [ ! -z "$pf_time" ]; then
            pf_sum=$(echo "$pf_sum + $pf_time" | bc)
        fi
    done
    pf_avg=$(echo "scale=2; $pf_sum / $pf_runs" | bc)
    
    # Calculate ratio
    if [ ! -z "$py_avg" ] && [ ! -z "$pf_avg" ] && [ $(echo "$py_avg > 0" | bc) -eq 1 ]; then
        ratio=$(echo "scale=2; $pf_avg / $py_avg" | bc)
        
        echo "  Python avg: ${py_avg} ms"
        echo "  Polyloft avg: ${pf_avg} ms"
        echo "  Ratio: ${ratio}x"
        echo ""
        
        # Add to results
        echo "| Test $i | $test_name | $py_avg | $pf_avg | ${ratio}x |" >> $results_file
    else
        echo "  Error: Could not parse times"
        echo "| Test $i | $test_name | ERROR | ERROR | - |" >> $results_file
    fi
done

echo "" >> $results_file
echo "## Analysis" >> $results_file
echo "" >> $results_file
echo "- Ratio < 1.0: Polyloft is faster" >> $results_file
echo "- Ratio > 1.0: Python is faster" >> $results_file
echo "- The ratio shows how many times slower/faster Polyloft is compared to Python" >> $results_file

echo "==================================="
echo "Results saved to: $results_file"
echo "==================================="
cat $results_file
