# Pull Request Template for GitHub Linguist

Use this template when submitting Polyloft to the GitHub Linguist repository.

---

## Add support for Polyloft programming language

This PR adds support for the Polyloft programming language (`.pf` extension).

### Language Overview

Polyloft is a hybrid programming language that combines object-oriented structure with scripting flexibility and reflection. It's designed for teaching, experimentation, and creating powerful scripts without sacrificing clarity.

**Key Features:**
- Hybrid design combining Ruby expressiveness, Python simplicity, and Java performance
- Object-oriented with classes and inheritance
- String interpolation with `#{}`
- Range operators (`...`, `..`)
- Flexible type system (Int, Float, String, Array, Map, etc.)
- Built-in standard library (Sys, Array, String, Math, etc.)

**Links:**
- **Repository**: https://github.com/ArubikU/polyloft
- **Documentation**: See `docs/` in repository
- **Extension**: `.pf`
- **Type**: Programming language
- **License**: MIT

### Usage on GitHub

GitHub search showing `.pf` files in use:
```
https://github.com/search?type=code&q=NOT+is%3Afork+path%3A*.pf
```

### Samples Provided

Three representative samples are included demonstrating various Polyloft features:

1. **fibonacci.pf**: Fibonacci sequence implementation using both iterative and memoization approaches
   - Shows: functions, loops, arrays, string concatenation, comments
   - License: MIT (from Polyloft repository)
   - Source: https://github.com/ArubikU/polyloft/blob/main/algorithm_samples/fibonacci.pf

2. **quicksort.pf**: QuickSort divide-and-conquer sorting algorithm
   - Shows: recursive functions, array manipulation, conditionals, partitioning
   - License: MIT (from Polyloft repository)
   - Source: https://github.com/ArubikU/polyloft/blob/main/algorithm_samples/quicksort.pf

3. **class_example.pf**: Object-oriented programming demonstration
   - Shows: classes, constructors, methods, `this` keyword, instances
   - License: MIT (from Polyloft repository)
   - Source: https://github.com/ArubikU/polyloft/blob/main/stress_tests/test9_class.pf

All samples demonstrate real-world usage patterns, not trivial "Hello World" examples.

### TextMate Grammar

The TextMate grammar (`polyloft.tmLanguage.json`) provides syntax highlighting support for:

**Language Constructs:**
- Comments: `// line comments`
- Keywords: `let`, `def`, `class`, `if`, `for`, `while`, `return`, `end`, `in`, `range`
- Storage types: `let`, `def`, `class`
- Variables: `this`, `self`
- Constants: `true`, `false`, `nil`, `null`

**Literals:**
- Strings: Double quotes with `#{}` interpolation, single quotes
- Numbers: Integers and floating-point
- Escape sequences in strings

**Operators:**
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `and`, `or`, `not`, `&&`, `||`, `!`
- Range: `...`, `..`
- Assignment: `=`

**Built-in Functions and Classes:**
- Functions: `println`, `print`, `len`, `range`, `type`, `str`, `int`, `float`, etc.
- Classes: `Sys`, `Array`, `String`, `Math`, `Object`, `JSON`, `File`, `Http`, `Time`, `Regex`

**Grammar License**: MIT (part of the Polyloft project)
**Grammar Location**: https://github.com/ArubikU/polyloft/tree/main/linguist/grammars

### Changes Made

- [x] Added entry to `lib/linguist/languages.yml`
- [x] Added TextMate grammar to `vendor/grammars/polyloft/`
- [x] Added sample files to `samples/Polyloft/`
- [x] Generated unique `language_id` with `script/update-ids`
- [x] All samples are real-world code with clear licensing
- [x] Grammar follows TextMate conventions
- [x] Extensions in alphabetical order
- [x] Linked to GitHub search showing usage

### Language Definition

```yaml
Polyloft:
  type: programming
  color: "#E53935"
  extensions:
  - ".pf"
  tm_scope: source.polyloft
  ace_mode: text
  language_id: <GENERATED_BY_SCRIPT>
  aliases:
  - pf
  - polyloft-lang
```

**Color Rationale**: `#E53935` (Red Loft) is Polyloft's official brand color as documented in the project's branding guidelines.

### Testing

Local testing performed:
```bash
bundle exec bin/github-linguist /path/to/polyloft/repo --breakdown
bundle exec rake test
bundle exec script/cross-validation --test
```

Results: `.pf` files correctly detected as Polyloft with proper syntax highlighting.

### Additional Notes

- Polyloft is under active development with growing community adoption
- The language is used for teaching algorithms and data structures
- The repository includes extensive documentation and examples
- Future updates will be submitted as grammar improvements if needed

### Checklist

- [x] Added language entry to `languages.yml`
- [x] Provided TextMate grammar with acceptable license
- [x] Included 2-3 representative sample files
- [x] Samples are real-world code, not trivial examples
- [x] Stated license for all samples and grammar
- [x] Linked to GitHub search showing `.pf` usage
- [x] Generated unique `language_id`
- [x] Local testing completed successfully
- [x] Extensions in alphabetical order
- [x] Grammar scope matches language definition

---

**Maintainer Note**: Thank you for reviewing this PR! Polyloft aims to make programming more accessible through its hybrid design. We're excited to bring official GitHub support to the language. Please let us know if any changes are needed.
