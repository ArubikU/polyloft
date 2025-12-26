# Generics

*Documentation for Generics is coming soon.*

Generics allow you to write flexible, reusable code that works with any data type.

```pf
class Box<T>:
    let value: T
    
    Box(value: T):
        this.value = value
    end
end
```
