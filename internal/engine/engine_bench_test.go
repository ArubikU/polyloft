package engine

import (
"testing"

"github.com/ArubikU/polyloft/internal/ast"
"github.com/ArubikU/polyloft/internal/common"
)

// Benchmark variable access
func BenchmarkVariableAccess(b *testing.B) {
env := common.NewEnv()
env.Set("x", 42)
env.Set("y", "hello")
env.Set("z", true)

b.Run("IntLookup", func(b *testing.B) {
expr := &ast.Ident{Name: "x"}
for i := 0; i < b.N; i++ {
_, _ = evalExpr(env, expr)
}
})

b.Run("StringLookup", func(b *testing.B) {
expr := &ast.Ident{Name: "y"}
for i := 0; i < b.N; i++ {
_, _ = evalExpr(env, expr)
}
})

b.Run("BoolLookup", func(b *testing.B) {
expr := &ast.Ident{Name: "z"}
for i := 0; i < b.N; i++ {
_, _ = evalExpr(env, expr)
}
})
}
