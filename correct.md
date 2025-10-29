evalGenericCallExpr
	var gtypes []GenericType
	for _, tp := range expr.TypeParams {
		if tp.IsWildcard {
			var boundTypeName string
			if len(tp.Bounds) > 0 {
				boundTypeName = tp.Bounds[0]
			}

			bound := common.GenericBound{
				Name:       ast.Type{Name: boundTypeName},
				Variance:   tp.Variance, // "extends", "super", or "unbounded"
				IsVariadic: tp.IsVariadic,
			}

			gtypes = append(gtypes, GenericType{
				Bounds: []common.GenericBound{bound},
			})
		} else {
			// Regular type parameter without variance
			// Don't add to allArgs - type parameters are NOT constructor arguments

			// Create a GenericType for regular type parameters
			bound := common.GenericBound{
				Name:       ast.Type{Name: tp.Name},
				Variance:   tp.Variance, // This will be "" for non-wildcard types
				IsVariadic: tp.IsVariadic,
			}
			gtypes = append(gtypes, GenericType{
				Bounds: []common.GenericBound{bound},
			})
		}
	}


NO esta guardando d ela forma correcta los tipos con wildcard ni los reconoce bien para despues, cuando typecheck.go los valide, ejemplo hacer ast.Type{Name: "? extends integer"} es generico e incorrecto solo se deberia hacer en el caso que sea "?" y luego en Bound.Extends defina el ClassDefinition correcto de el nombre
type GenericBound struct {
	Name       ast.Type
	Variance   string // "in" (contravariance), "out" (covariance), or "" (invariant)
	IsVariadic bool
	Extends    *ClassDefinition
	Implements *InterfaceDefinition
}
como puedes ver la estructura tiene Extends e Implements que deberia ser rellenada en este punto del codigo para que luego en typecheck.go pueda validar bien los tipos con wildcard
tambien revisa class.go

		if classDef.IsGeneric && len(args) > 0 {
			// Check if first arguments are type parameters
			numTypeParams := len(classDef.TypeParams)
			if len(args) >= numTypeParams {
				// Try to extract type arguments
				typeArgs = make([]string, 0, numTypeParams)
				for i := 0; i < numTypeParams && i < len(args); i++ {
					if typeInfo, ok := args[i].(map[string]any); ok {
						// Handle wildcard or variance-annotated types
						if isWildcard, ok := typeInfo["isWildcard"].(bool); ok {
							if isWildcard {
								kind, _ := typeInfo["kind"].(string)
								bound, _ := typeInfo["bound"].(string)
								variance, _ := typeInfo["variance"].(string)
								wildcardStr := formatWildcard(kind, bound)
								if variance != "" {
									wildcardStr = variance + " " + wildcardStr
								}
								typeArgs = append(typeArgs, wildcardStr)
							} else {
								// Variance-annotated type
								name, _ := typeInfo["name"].(string)
								variance, _ := typeInfo["variance"].(string)
								if variance != "" && name != "" {
									typeArgs = append(typeArgs, variance+" "+name)
								} else {
									typeArgs = append(typeArgs, name)
								}
							}
						} else {
							// Not a type info map, treat as regular arg
							break
						}
					} else if typeStr, ok := args[i].(string); ok && isTypeName(typeStr) {
						// Regular type name
						typeArgs = append(typeArgs, typeStr)
					} else {
						// Not a type argument, stop extraction
						break
					}
				}

				// If we extracted the right number of type params, use them
				if len(typeArgs) == numTypeParams {
					constructorArgs = args[numTypeParams:]
				} else {
					// Couldn't extract type params, treat all as constructor args
					typeArgs = nil
					constructorArgs = args
				}
			} else {
				constructorArgs = args
			}
		} else {
			constructorArgs = args
		}