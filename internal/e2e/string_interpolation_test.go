package e2e

import (
"bytes"
"testing"

"github.com/ArubikU/polyloft/internal/engine"
"github.com/ArubikU/polyloft/internal/lexer"
"github.com/ArubikU/polyloft/internal/parser"
)

func TestStringInterpolationEnhancements(t *testing.T) {
tests := []struct {
name     string
code     string
expected string
}{
{
name: "Simple variable interpolation",
code: `
let x = 5
println("x = #{x}")
`,
expected: "x = 5\n",
},
{
name: "Function call in interpolation",
code: `
def double(n: Int) -> Int:
    return n * 2
end
let x = 5
println("double(x) = #{double(x)}")
`,
expected: "double(x) = 10\n",
},
{
name: "Ternary operator in interpolation",
code: `
let x = 5
let y = 10
println("max = #{x > y ? x : y}")
`,
expected: "max = 10\n",
},
{
name: "Binary operations in interpolation",
code: `
let x = 5
let y = 10
println("sum = #{x + y}")
println("product = #{x * y}")
`,
expected: "sum = 15\nproduct = 50\n",
},
{
name: "Complex expression in interpolation",
code: `
def double(n: Int) -> Int:
    return n * 2
end
let x = 5
let y = 10
println("result = #{double(x) + double(y)}")
`,
expected: "result = 30\n",
},
{
name: "Nested function calls in interpolation",
code: `
def double(n: Int) -> Int:
    return n * 2
end
def triple(n: Int) -> Int:
    return n * 3
end
let x = 5
println("nested = #{triple(double(x))}")
`,
expected: "nested = 30\n",
},
{
name: "Ternary with operations in interpolation",
code: `
let x = 5
let y = 10
println("result = #{(x + y) > 10 ? x + y : 0}")
`,
expected: "result = 15\n",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
var out bytes.Buffer
lx := &lexer.Lexer{}
tokens := lx.Scan([]byte(tt.code))
p := parser.New(tokens)
prog, err := p.Parse()
if err != nil {
t.Fatalf("parse error: %v", err)
}

opts := engine.Options{Stdout: &out}
_, err = engine.Eval(prog, opts)
if err != nil {
t.Fatalf("eval error: %v", err)
}

got := out.String()
if got != tt.expected {
t.Errorf("expected %q, got %q", tt.expected, got)
}
})
}
}
