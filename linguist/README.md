# Polyloft Linguist Integration

This directory contains files for registering Polyloft (`.pf`) as a programming language with [GitHub Linguist](https://github.com/github/linguist).

## Quick Start

1. **Read the Implementation Guide**: Start with [`IMPLEMENTATION_GUIDE.md`](./IMPLEMENTATION_GUIDE.md) for complete instructions
2. **Review Files**: Check `languages.yml`, the grammar, and samples
3. **Submit to Linguist**: Follow the guide to submit a PR to the linguist repository

## What's Included

- **`IMPLEMENTATION_GUIDE.md`**: Complete step-by-step guide for Linguist submission
- **`languages.yml`**: Language definition entry for Linguist
- **`grammars/polyloft.tmLanguage.json`**: TextMate grammar for syntax highlighting
- **`samples/`**: Representative Polyloft code samples

## Why Register with Linguist?

Registering Polyloft with GitHub Linguist enables:

- ✅ **Syntax highlighting** for `.pf` files on GitHub
- ✅ **Language statistics** in repository insights
- ✅ **Language detection** for Polyloft projects
- ✅ **Search filters** for finding Polyloft code
- ✅ **Community visibility** for the language

## Requirements

Before submitting to Linguist:

1. `.pf` extension should be in use in **at least 200 unique repositories** on GitHub
2. Provide **2-3 sample files** showing real-world code (not "Hello World")
3. Grammar must have an **acceptable open-source license** (MIT, Apache, BSD, etc.)

## Language Information

- **Name**: Polyloft
- **Extension**: `.pf`
- **Type**: Programming language
- **Color**: `#E53935` (Red Loft - Polyloft brand color)
- **Scope**: `source.polyloft`
- **Website**: https://github.com/ArubikU/polyloft

## Sample Features

The TextMate grammar supports:

- **Comments**: `// line comments` and `/* block comments */`
- **Keywords**: `let`, `def`, `class`, `if`, `elif`, `for`, `loop`, `where`, `return`, `break`, `continue`, `end`
- **Modifiers**: `public`, `private`, `protected`, `static`, `abstract`, `sealed`
- **String interpolation**: `"Hello #{name}"`
- **Classes**: `class Point:` ... `end`
- **Functions**: `def fibonacci(n):` ... `end`
- **Type annotations**: `let x: Int = 5`, `name: String`, `def foo() -> Int:`
- **Generic types**: `Array<Int>`, `Map<String, Float>`
- **Built-in functions**: `println()`, `print()`, `len()`, `range()`, `int()`, `float()`, `str()`, `bool()`
- **Built-in classes**: `Sys`, `Array`, `String`, `Map`, etc.
- **Logical operators**: `&&`, `||`, `!` (note: `and`, `or`, `not` do not exist)
- **Ranges**: `0...10` (only `...`, not `..`)
- **Where clause**: `for n in numbers where n > 5:`

## Next Steps

1. **Local Testing**: Test the grammar in a TextMate-compatible editor
2. **Review Samples**: Ensure samples are representative and well-documented
3. **Follow the Guide**: Complete steps in `IMPLEMENTATION_GUIDE.md`
4. **Submit PR**: Open a pull request to github/linguist

## Contributing

Found an issue with the grammar or samples? Open an issue or PR in the [Polyloft repository](https://github.com/ArubikU/polyloft).

## License

All files in this directory are part of the Polyloft project and are licensed under the MIT License.
