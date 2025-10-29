# Import Statement

The `import` statement allows you to import modules and specific symbols from other Polyloft files.

## Syntax

### Basic Import
```pf
import module.name
```

### Import Specific Symbols
```pf
import module.name { Symbol1, Symbol2 }
```

### Dotted Path Import
```pf
import package.subpackage.module
```

## Examples

### Import Entire Module
```pf
// file: math_utils.pf
class MathHelper:
    static def square(x):
        return x * x
    end
    
    static def cube(x):
        return x * x * x
    end
end

// file: main.pf
import math_utils

let result = math_utils.MathHelper.square(5)
println(result)  // 25
```

### Import Specific Classes
```pf
// file: utils.pf
class Logger:
    static def log(message):
        println("[LOG] #{message}")
    end
end

class Helper:
    static def format(text):
        return text.toUpperCase()
    end
end

// file: main.pf
import utils { Logger, Helper }

Logger.log("Application started")
let formatted = Helper.format("hello")
println(formatted)  // "HELLO"
```

### Import from Package
```pf
// file: lib/data/models.pf
class User:
    let name
    let email
    
    def init(name, email):
        this.name = name
        this.email = email
    end
end

// file: main.pf
import lib.data.models { User }

let user = User("Alice", "alice@example.com")
println(user.name)
```

### Multiple Imports
```pf
import utils { Logger }
import models { User, Product }
import helpers.string { capitalize, trim }

Logger.log("Starting application")
let user = User("Bob", "bob@example.com")
```

## Module Organization

### Project Structure
```
project/
├── main.pf
├── lib/
│   ├── utils.pf
│   ├── models.pf
│   └── helpers/
│       ├── string.pf
│       └── math.pf
└── config.pf
```

### Creating Reusable Modules
```pf
// file: lib/validation.pf
class Validator:
    static def isEmail(email):
        return email.contains("@") and email.contains(".")
    end
    
    static def isNotEmpty(text):
        return text.length() > 0
    end
    
    static def isInRange(value, min, max):
        return value >= min and value <= max
    end
end

// file: main.pf
import lib.validation { Validator }

let email = "user@example.com"
if Validator.isEmail(email):
    println("Valid email")
end
```

## Common Patterns

### Utility Library
```pf
// file: lib/utils.pf
class StringUtils:
    static def capitalize(text):
        if text.isEmpty():
            return text
        end
        return text.charAt(0).toUpperCase() + text.substring(1)
    end
    
    static def reverse(text):
        let result = ""
        for i in range(text.length() - 1, -1, -1):
            result = result + text.charAt(i)
        end
        return result
    end
end

class ArrayUtils:
    static def sum(arr):
        return arr.reduce((acc, x) => acc + x, 0)
    end
    
    static def average(arr):
        return ArrayUtils.sum(arr) / arr.length()
    end
end

// file: main.pf
import lib.utils { StringUtils, ArrayUtils }

println(StringUtils.capitalize("hello"))
println(ArrayUtils.sum([1, 2, 3, 4, 5]))
```

### Model Definitions
```pf
// file: models/user.pf
class User:
    let id
    let username
    let email
    
    def init(id, username, email):
        this.id = id
        this.username = username
        this.email = email
    end
    
    def toJSON():
        return {
            id: this.id,
            username: this.username,
            email: this.email
        }
    end
end

// file: models/post.pf
class Post:
    let id
    let title
    let content
    let authorId
    
    def init(id, title, content, authorId):
        this.id = id
        this.title = title
        this.content = content
        this.authorId = authorId
    end
end

// file: main.pf
import models.user { User }
import models.post { Post }

let user = User(1, "alice", "alice@example.com")
let post = Post(1, "Hello", "First post", user.id)
```

### Service Layer
```pf
// file: services/database.pf
class DatabaseService:
    static def connect(config):
        println("Connecting to database...")
        // Connection logic
    end
    
    static def query(sql):
        println("Executing: #{sql}")
        // Query logic
        return []
    end
end

// file: services/auth.pf
import services.database { DatabaseService }

class AuthService:
    static def login(username, password):
        let query = "SELECT * FROM users WHERE username = '#{username}'"
        let results = DatabaseService.query(query)
        // Authentication logic
        return true
    end
end

// file: main.pf
import services.auth { AuthService }

if AuthService.login("alice", "password123"):
    println("Login successful")
end
```

## Best Practices

### ✅ DO - Organize code into modules
```pf
// Good: Separate concerns
import models { User }
import services { UserService }
import utils { Validator }
```

### ✅ DO - Import only what you need
```pf
// Good: Specific imports
import utils { Logger, Helper }

// Avoid: Import everything
import utils
```

### ✅ DO - Use meaningful module names
```pf
import user.management { UserManager }
import data.validation { EmailValidator }
```

### ❌ DON'T - Create circular dependencies
```pf
// Bad: Module A imports Module B, Module B imports Module A
// file: moduleA.pf
import moduleB

// file: moduleB.pf
import moduleA  // Circular dependency!
```

### ❌ DON'T - Nest imports too deeply
```pf
// Bad: Too many levels
import very.deeply.nested.module.package.subpackage.file

// Better: Flatten structure
import modules.file
```

## Module Resolution

Polyloft resolves imports relative to:
1. The current file's directory
2. The project root
3. Installed packages

```pf
// Relative import
import ./utils { Helper }

// Project root import
import lib.models { User }

// Package import (if installed)
import external.package { Class }
```

## See Also

- [Classes](../definitions/class.md)
- [Functions](../definitions/function.md)
- [Project Structure](../guides/project-structure.md)
