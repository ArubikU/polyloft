//Como debe funcionar los tipos y sus reglas

//1. Cada clase enum record interface y builtin tienen sus tipos
class Person:
    let name
    let age
    Person(name, age):
        this.name = name
        this.age = age
    end
end
const p = Person("Alice", 30)
print(typeOf(p)) // "Person"
print(p instanceof Person) // true
//2. Si una clase hereda de otra, al comprobar el tipo debe considerar la herencia
class Employee extends Person:
    let employeeId
    Employee(name, age, employeeId):
        super(name, age)
        this.employeeId = employeeId
    end
end
const e = Employee("Bob", 25, "E123")
print(typeOf(e)) // "Employee"
print(e instanceof Employee) // true
print(e instanceof Person) // true
//3. Las interfaces tambien sirven para comprobar tipos
interface Named:
    def getName()
end
class User implements Named:
    let name
    User(name):
        this.name = name
    end
    def getName():
        return this.name
    end
end
const u = User("Charlie")
print(typeOf(u)) // "User"
print(u instanceof User) // true
print(u instanceof Named) // true
//4. Los enums funcionan correctamente con los tipos
enum Color:
    Red
    Green
    Blue
end
const c = Color.Red
print(typeOf(c)) // "Color"
print(c instanceof Color) // true
//5. Verificar tipos al usar funciones
def greet(person: Person):
    print("Hello, " + person.name)
end
greet(p) // Funciona
greet(e) // Funciona porque Employee es un Person
greet(u) // No es un Person ni lo extiende

def printInt(n: Int):
    print("Int: " + n.toString())
end
def printFloat(f: Float):
    print("Float: " + f.toString())
end
def printNumber(n: Number):
    print("Number: " + n.toString())
end

printInt(42) // Funciona
printFloat(3.14) // Funciona
printNumber(100) // Funciona
printNumber(2.718) // Funciona
//printInt(3.14) // Error de tipo
//printFloat(42) // Error de tipo
//Number es un supertipo de Int y Float en los built-in
//6. Uso de Any o ? como tipo universal
def printAny(value: Any):
    print("Value: " + value.toString())
end
printAny(42) // Funciona

let x: Any = "Hello"
print(typeOf(x)) // "String" lo infiere correctamente
let y: Int = 10 
print(typeOf(x)) // "Int"
let z: ? = 3.14
print(typeOf(x)) // "Float"
let w: User = Person("Dave", 40) // Error de tipo, Person no es User

//7. Variadic functions con tipos

def sum(numbers: Int...):
    let total = 0
    for n in numbers:
        total += n
    end
    return total
end
print(sum(1, 2, 3, 4, 5)) // 15
//print(sum(1, 2, "3")) // Error de tipo
//8. Genéricos con restricciones de tipo
class Box<T extends Number>:
    let value: T
    Box(value: T):
        this.value = value
    end
    def getValue(): T:
        return this.value
    end
    def setValue(newValue: T):
        this.value = newValue
    end
end
const intBox = Box<Int>(100)
print(typeOf(intBox)) // "Box<Int>"
const floatBox = Box<Float>(3.14)
print(typeOf(floatBox)) // "Box<Float>"
const strBox = Box<String>("Hello") // Error de tipo, String no extiende Number
const anyBox = Box<Any>(true) // Error de tipo, Any no extiende Number
const failBox: Box<Float> = Box<Int>(10) // Error de tipo, Int no es Float
//9. Comprobación de in out
let sealedBox: Box<Int> = Box<in Int>(50) // Correcto ya que el in/out solo afecta la lectura/escritura 
let testSealedBox: Box<out Int> = Box<out Int>(50) // Correcto
let testFailSealedBox: Box<out Float> = Box<out Int>(50) // Error de tipo
let vfail: Box<in Int> = Box<out Int>(50) // Error de variancia

sealedBox.setValue(60) // Correcto
sealedBox.setValue(3.14) // Error de tipo
sealedBox.getValue() // Tira error ya que la instancia el tipo generico es in, por lo que cuando se intente leer o retornar en funciones el "T" tirara error

//10. Renombrar tipos

final type Age = Int
def celebrateBirthday(age: Age):
    print("Happy " + age.toString() + "th Birthday!")
end
celebrateBirthday(30) // Funciona
celebrateBirthday("30") // Error de tipo

final type Age = Int
final type Years = Int
let a: Age = 30
let y: Years = a // Tira error de tipo, Age y Years son tipos distintos aunque ambos sean Int



//11. Los tipos built-in deben funcionar correctamente
const arr = [1, 2, 3]
print(typeOf(arr)) // "Array"
print(arr instanceof Array) // true
const dict = {"key": "value"}
print(typeOf(dict)) // "Map"
print(dict instanceof Map) // true
const s = "Hello"
print(typeOf(s)) // "String"
print(s instanceof String) // true

//12. Nil y Null

const numero: Int | Nil = Nil
print(typeOf(numero)) // "Nil"
const texto: String | Null = Null
print(typeOf(texto)) // "Nil"
const numero2: Int = Nil // Error de tipo , Nil no es Int
const anyval: Any = Nil // Funciona, Any acepta Nil

//13. 
class PrintableBox<T extends Number & Named>:
    def printValue():
        print(this.value.getName() + " = " + this.value.toString())
    end
end
const namedNumber = Employee("Eve", 28, "E456") // Employee implementa Named y extiende Number

def process(box: Box<out Number>):
    let x = box.getValue() // ok
    box.setValue(42) // error (out → no se puede escribir)
end
//Permite validar sobrecarga de funciones basadas en tipos:
def f(x: Int): String
def f(x: String): Int
def f(x: Any): Any

//14.
def printId(id: Int | String):
    print(id)
end

def ensureNamed(x: Named & Serializable):
    ...
end

let arr = [1, 2, 3] // Array , Los array no tienen tipo generico no lo necesitan
let mixed = List(1, "2", 3) // List<String | Int>

//15. Conversion explicita de tipos
let n: Float = 3 // Se convierte a 3.0


//16. Polimorfismo de subtipos con tipos genericos
def identity<T>(value: T): T
let i = identity(42) // Int
let j = identity("hi") // String


//17. etc
def printValue(x: Int | String):
    if (x instanceof Int):
        print(x + 1)
    else:
        print(x.length)

