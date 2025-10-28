package ast

// Visitor is a function type for visiting AST nodes
type Visitor func(Node) bool // return false to stop traversal

// IterativeWalk traverses an AST iteratively using a manual stack
// This avoids deep recursion that can cause stack overflow on large ASTs
// The visitor function is called for each node. Return false to stop traversal.
func IterativeWalk(root Node, visitor Visitor) {
	if root == nil {
		return
	}

	// Manual stack to avoid recursion
	// Preallocate with reasonable initial capacity
	stack := make([]Node, 0, 64)
	stack = append(stack, root)

	for len(stack) > 0 {
		// Pop from stack
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n == nil {
			continue
		}

		// Visit node
		if !visitor(n) {
			return
		}

		// Push children onto stack
		// Process in reverse order so they're visited in the correct order
		switch node := n.(type) {
		case *Program:
			for i := len(node.Stmts) - 1; i >= 0; i-- {
				stack = append(stack, node.Stmts[i])
			}

		case *BinaryExpr:
			if node.Rhs != nil {
				stack = append(stack, node.Rhs)
			}
			if node.Lhs != nil {
				stack = append(stack, node.Lhs)
			}

		case *UnaryExpr:
			if node.X != nil {
				stack = append(stack, node.X)
			}

		case *CallExpr:
			for i := len(node.Args) - 1; i >= 0; i-- {
				stack = append(stack, node.Args[i])
			}
			if node.Callee != nil {
				stack = append(stack, node.Callee)
			}

		case *GenericCallExpr:
			for i := len(node.Args) - 1; i >= 0; i-- {
				stack = append(stack, node.Args[i])
			}

		case *SuperExpr:
			for i := len(node.Args) - 1; i >= 0; i-- {
				stack = append(stack, node.Args[i])
			}

		case *IndexExpr:
			if node.Index != nil {
				stack = append(stack, node.Index)
			}
			if node.X != nil {
				stack = append(stack, node.X)
			}

		case *FieldExpr:
			if node.X != nil {
				stack = append(stack, node.X)
			}

		case *InterpolatedStringLit:
			for i := len(node.Parts) - 1; i >= 0; i-- {
				stack = append(stack, node.Parts[i])
			}

		case *ArrayLit:
			for i := len(node.Elems) - 1; i >= 0; i-- {
				stack = append(stack, node.Elems[i])
			}

		case *MapLit:
			for i := len(node.Pairs) - 1; i >= 0; i-- {
				stack = append(stack, node.Pairs[i].Value)
			}

		case *LetStmt:
			if node.Value != nil {
				stack = append(stack, node.Value)
			}

		case *AssignStmt:
			if node.Value != nil {
				stack = append(stack, node.Value)
			}
			if node.Target != nil {
				stack = append(stack, node.Target)
			}

		case *ReturnStmt:
			if node.Value != nil {
				stack = append(stack, node.Value)
			}

		case *ExprStmt:
			if node.X != nil {
				stack = append(stack, node.X)
			}

		case *DefStmt:
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}

		case *IfStmt:
			for i := len(node.Else) - 1; i >= 0; i-- {
				stack = append(stack, node.Else[i])
			}
			for i := len(node.Clauses) - 1; i >= 0; i-- {
				clause := node.Clauses[i]
				for j := len(clause.Body) - 1; j >= 0; j-- {
					stack = append(stack, clause.Body[j])
				}
				if clause.Cond != nil {
					stack = append(stack, clause.Cond)
				}
			}

		case *ForInStmt:
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}
			if node.Where != nil {
				stack = append(stack, node.Where)
			}
			if node.Iterable != nil {
				stack = append(stack, node.Iterable)
			}

		case *LoopStmt:
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}
			if node.Condition != nil {
				stack = append(stack, node.Condition)
			}

		case *DoLoopStmt:
			if node.Condition != nil {
				stack = append(stack, node.Condition)
			}
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}

		case *TryStmt:
			for i := len(node.Finally) - 1; i >= 0; i-- {
				stack = append(stack, node.Finally[i])
			}
			for i := len(node.Catches) - 1; i >= 0; i-- {
				catch := node.Catches[i]
				for j := len(catch.Body) - 1; j >= 0; j-- {
					stack = append(stack, catch.Body[j])
				}
			}
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}

		case *ThrowStmt:
			if node.Value != nil {
				stack = append(stack, node.Value)
			}

		case *DeferStmt:
			if node.Call != nil {
				stack = append(stack, node.Call)
			}

		case *InterfaceDecl:
			for i := len(node.Methods) - 1; i >= 0; i-- {
				method := node.Methods[i]
				for j := len(method.DefaultBody) - 1; j >= 0; j-- {
					stack = append(stack, method.DefaultBody[j])
				}
			}
			for i := len(node.Fields) - 1; i >= 0; i-- {
				// Push FieldDecl as a node to visit
				stack = append(stack, &node.Fields[i])
			}

		case *ClassDecl:
			if node.Constructor != nil {
				for i := len(node.Constructor.Body) - 1; i >= 0; i-- {
					stack = append(stack, node.Constructor.Body[i])
				}
			}
			for i := len(node.Methods) - 1; i >= 0; i-- {
				method := node.Methods[i]
				for j := len(method.Body) - 1; j >= 0; j-- {
					stack = append(stack, method.Body[j])
				}
			}
			for i := len(node.Fields) - 1; i >= 0; i-- {
				// Push FieldDecl as a node to visit
				stack = append(stack, &node.Fields[i])
			}

		case *EnumDecl:
			if node.Constructor != nil {
				for i := len(node.Constructor.Body) - 1; i >= 0; i-- {
					stack = append(stack, node.Constructor.Body[i])
				}
			}
			for i := len(node.Methods) - 1; i >= 0; i-- {
				method := node.Methods[i]
				for j := len(method.Body) - 1; j >= 0; j-- {
					stack = append(stack, method.Body[j])
				}
			}
			for i := len(node.Fields) - 1; i >= 0; i-- {
				// Push FieldDecl as a node to visit
				stack = append(stack, &node.Fields[i])
			}
			for i := len(node.Values) - 1; i >= 0; i-- {
				value := node.Values[i]
				for j := len(value.Args) - 1; j >= 0; j-- {
					stack = append(stack, value.Args[j])
				}
			}

		case *RecordDecl:
			for i := len(node.Methods) - 1; i >= 0; i-- {
				method := node.Methods[i]
				for j := len(method.Body) - 1; j >= 0; j-- {
					stack = append(stack, method.Body[j])
				}
			}

		case *FieldDecl:
			if node.InitValue != nil {
				stack = append(stack, node.InitValue)
			}

		case *InstanceOfExpr:
			if node.Expr != nil {
				stack = append(stack, node.Expr)
			}

		case *TypeExpr:
			if node.Expr != nil {
				stack = append(stack, node.Expr)
			}

		case *LambdaExpr:
			for i := len(node.BlockBody) - 1; i >= 0; i-- {
				stack = append(stack, node.BlockBody[i])
			}
			if node.Body != nil {
				stack = append(stack, node.Body)
			}

		case *ThreadSpawnExpr:
			for i := len(node.Body) - 1; i >= 0; i-- {
				stack = append(stack, node.Body[i])
			}

		case *ThreadJoinExpr:
			if node.Thread != nil {
				stack = append(stack, node.Thread)
			}

		case *SelectStmt:
			for i := len(node.Cases) - 1; i >= 0; i-- {
				c := node.Cases[i]
				for j := len(c.Body) - 1; j >= 0; j-- {
					stack = append(stack, c.Body[j])
				}
				if c.Channel != nil {
					stack = append(stack, c.Channel)
				}
			}

		case *SwitchStmt:
			for i := len(node.Default) - 1; i >= 0; i-- {
				stack = append(stack, node.Default[i])
			}
			for i := len(node.Cases) - 1; i >= 0; i-- {
				c := node.Cases[i]
				for j := len(c.Body) - 1; j >= 0; j-- {
					stack = append(stack, c.Body[j])
				}
				for j := len(c.Values) - 1; j >= 0; j-- {
					stack = append(stack, c.Values[j])
				}
			}
			if node.Expr != nil {
				stack = append(stack, node.Expr)
			}

		case *TernaryExpr:
			if node.FalseBranch != nil {
				stack = append(stack, node.FalseBranch)
			}
			if node.TrueBranch != nil {
				stack = append(stack, node.TrueBranch)
			}
			if node.Condition != nil {
				stack = append(stack, node.Condition)
			}

		case *RangeExpr:
			if node.End != nil {
				stack = append(stack, node.End)
			}
			if node.Start != nil {
				stack = append(stack, node.Start)
			}

		// Leaf nodes - no children to add
		case *Ident, *NumberLit, *StringLit, *BytesLit, *BoolLit, *NilLit,
			*BreakStmt, *ContinueStmt, *ImportStmt, *TypeAliasStmt,
			*ChannelExpr:
			// No children
		}
	}
}

// CountNodes counts the number of nodes in an AST using iterative traversal
func CountNodes(root Node) int {
	count := 0
	IterativeWalk(root, func(n Node) bool {
		count++
		return true
	})
	return count
}

// FindNodes finds all nodes matching a predicate using iterative traversal
func FindNodes(root Node, predicate func(Node) bool) []Node {
	var result []Node
	IterativeWalk(root, func(n Node) bool {
		if predicate(n) {
			result = append(result, n)
		}
		return true
	})
	return result
}

// FindFirstNode finds the first node matching a predicate
func FindFirstNode(root Node, predicate func(Node) bool) Node {
	var result Node
	IterativeWalk(root, func(n Node) bool {
		if predicate(n) {
			result = n
			return false // Stop traversal
		}
		return true
	})
	return result
}
