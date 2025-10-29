Defines a inmutable variable whose value cannot be changed after its initial assignment. or reasing or edit their content if its a object or array.
```pf
final MAX_USERS = 100
MAX_USERS = 200  // This will raise an error
```
In this example, `MAX_USERS` is defined as a final variable and cannot be reassigned later in the code. Attempting to change its value will result in an error.
```pf
final CONFIG = {
    "host": "localhost",
    "port": 8080
}
CONFIG.port = 9090  // This will raise an error
```
In this example, `CONFIG` is defined as a final variable holding an object. Attempting to modify any of its properties will result in an error.
```pf
class Point:
    var x
    var y

    Point(x, y):
        this.x = x
        this.y = y
    end
    def move(new_x, new_y):
        this.x = new_x
        this.y = new_y
    end
end
final ORIGIN = Point(0, 0)
ORIGIN.move(5, 5)  // This will raise an error
```
In this example, `ORIGIN` is defined as a final variable holding an instance of the `Point` class. Attempting to call methods that modify its state will result in an error.