- [Basic types](#basic-types)
- [Variable declaration](#variable-declaration)
- [Constants](#constants)
- [Arrays](#arrays)
- [Slices](#slices)
- [Arrays × slices conversion](#arrays-vs-slices-conversion)
- [Strings](#strings)
- [Maps](#maps)
- [Structs](#structs)
- [If-then-else](#if-then-else)
- [For cycle](#for-cycle)
- [Switch](#switch)
- [Functions](#functions)
- [Pointers](#pointers)
- [Methods](#methods)
- [Interfaces](#interfaces)
- [Errors](#errors)

# Basic types

```go
var s1 string = "interpreted\nstring"
var s2 string = `non-interpreted\nstring`
var r1 rune = '🎁'
```
rune '🎁' == '\U0001F381'; also '\xHH' and '\uHHHH'. Alias to int32

Numerical: [u]int(8,16,32,64)
int and uint are platform-dependent, == [u]int64 on 64bit
byte == uint8

Float: float32, float64 (no default float)

Complex: complex64, complex128 (no default complex)


# Variable declaration

```go
var x int  // type declared, "zero" value used: 0 for integers
var x = 1  // type inferred, int as the most generic one
var x uint32 = 1  // both type and value set explicitely
var x1, x2 string = "one", "two"  // more vars of the same type
var (
    x    int
    y        = 20
    z    int = 30
    d, e     = 40, "hello"
    f, g string
)
```

Another way, only allowed in functions (not on module-level)
```go
x1, x2 := 10, "something" // implicite type
x := byte(10)             // explicite type
```


# Constants

Only basic types are allowed. No constant structs, maps, slices etc. Otherwise
similar to var

```go
const pi = 3.1415
```

Constants and `iota` could be used to imitate enums:
```go
const (
    Spades = iota
    Clubs
    Diamonds
    Hearts
)
fmt.Println(Spades, Clubs, Diamonds, Hearts)     // 0 1 2 3
```
If `iota` is used in constant definition, the next constant uses the same
definition with iota++, unless specified otherwise. Since it starts with 0, to
prevent "first value vs uninitialized value" problem, it's recommended to
either use `iota+1` or skip the first value:
```go
const (
    READ = iota + 1
    WRITE
    EXECUTE
)
// or
const (
    _
    READ = iota
    WRITE
    EXECUTE
)
```

# Arrays

Not terribly useful. Constant size
```go
var x [3]int
var x = [3]int{1, 2, 3}
var x = [5]int{1, 4:9}    // sparse, [1, 0, 0, 0, 9]
var x = [...]int{1, 3:8}  // special syntax: [...] means compile-time size calculation
var x [2][3]int           // two-dimensional array [[0, 0, 0], [0, 0, 0]]
```


# Slices

Variable size. Have length and capacity
```go
var x []int
var x = []int{1, 2, 3}
var x [][]int
```
Uninitialized slices are nil!

can't use == !=; slices.Equal(a, b)

`append(slice, item1 [, item2 ..])` or `append(slice, slice2..)` is not
in-place

`=` or `[:]` don't copy, but use the same var. `[::]` is like `[:]`, but limits
capacity. Absurd overlapping example:

```go
x := make([]string, 0, 5)
x = append(x, "a", "b", "c", "d") // capacity 5
y := x[:2]                        // capacity 5, [a b]. Better x[:2:2]
z := x[2:]                        // capacity 3, [c d]. Better x[2:4:4]
                                  // [a b c d] [a b] [c d]
y = append(y, "i", "j", "k")      // [a b i j] [a b i j k] [i j]
x = append(x, "x")                // [a b i j x] [a b i j x] [i j]
z = append(z, "y", "w")           // [a b i j y] [a b i j y] [i j y]
```
`len()` and `cap()` return length and capacity. Cant index beyond capacity!

```go
s1 := make([]int, initialLength [, initialCapacity])  // another way to create a slice
clear(s1)       // clears elements, keeps length and capacity
copy(dest, src) // copies up to LENGTH, returns copied length. Only slices, not arrays.
```

# Arrays vs slices conversion

```go
slice1 = array1[:]
array2 := [size]int(slice1)  // can't use [...] or [len(slice)]
                             // fixed length ≤ slice_size needed
```

# Strings

[x], [x:y] count in *bytes*! [x] returns byte, [x:y] string

`len()` counts bytes!

`string(int)` creates a character; strconv has `Itoa` and `FormatInt`.

```go
var s string = "🎯hit"
var r []rune = []rune(s)
var b []byte = []byte(s)
fmt.Printf("%T %[1]v\n%[2]T %[2]v\n", r, b)

[]int32 [127919 104 105 116]
[]uint8 [240 159 142 175 104 105 116]
```

# Maps

```go
var x map[string]int        // creates a nil map, barely usable
var x = map[string]int{}    // empty map, usable
x := map[string]int{}       // same as above, empty map
x := make(map[string]int)   // make works too
x := map[string][]string{}  // map string to slice of strings
```
`len(m)` number of KVPs. Access of non-existing KVP produces "zero value"!

`delete(map, key)` deletes a key

`clear(map)` clear entire map, length=0 (unlike slices)

`val, ok = map[key]` OK-idiom: ok contains bool (key is in the map or not)

can't use == !=; maps.Equal(a, b)

## Maps magic

can't use &m["key"] nor m["key"].field! See also structs and pointers magic.


# Structs

```go
type mydata struct {
    field1 int
    field2 string
}
v1 := mydata{   // explicit field names, recommended
    field1: 1,
}
v2 := mydata{   // implicit field names
    1,
    "test"
}
```
Structs could be anynomous
```go
var v1 struct {
    f1 int
}
// or
v1 := struct { f1 int } { f1: 2 }
```
Structs don't go along nicely with maps, can't use m[key].field; better use
`map[int]*mystruct{}`

## Structs magic

Struct pointers are "magically" dereferenced when accessing fields:
```go
type st struct {
    field1 int
    field2 string
}
p := &st{1, "one"}
p.field1 = 2            // OK, *p = {2, "one"}

m := map[int]*st        // map int to pointer to struct
m[1] = &st{3, "three"}  // here it insists it's a pointer
res := m[1].field1 == 3 // OK, res = true
```

## Struct embedding/composition
Struct can contain unnamed field of other type; in this case the fields of
the embedded struct could be accessed directly (but not in construction),
and methods are promoted to struct's methods:
```go
type Point struct {
    X int
    Y int
}
type Circle struct {
    Point
    R int
}
c1 := Circle{Point: Point{3,4}, R: 2}  // can't use Circle{X: 3, Y:4, R:2}
c1.X == 3                              // true; here can access .X and .Y
```
In case of of name overlap, "inner" fields need full specification: `c1.Point.X`

See also [interfaces](#interfaces).

# If-then-else

```go
if condition {
} else if condition2 {
} else {
}
// var declaration is possible in condition; var is scoped to block only
if a := rand.Intn(10); a > 5 { }
```

# For cycle
```go
for i := 0; i < 10; i++ { .. }               // can omit any part, but still needs `;`
for i < 100 { .. }                           // like while
for { .. }                                   // infinite loop
for i, v := range []string{"a", "b" } { .. } // [0 "a"], [1 "b"]
```
Range-loop on string iterates over *runes* and i is the *byte* offset!
```go
for i, v := range "🎁x" { .. }               // [0 127873], [4 120]
```

Range-loop on maps iterates over KVPs. Order is not guaranteed.
```go
for k, v := range mymap { .. }
for _, v := range mymap { .. }  // only values
for k := range mymap { .. }     // only keys
```

Range-loop in a channel iterates over data received from it
```go
for v := range channel { .. }
```

Range loop creates a copy; if you need to change the elements, use
`for i := range slice { slice[i] = .. }`

loops can be labeled: loopname: ..
`break`, `continue` by default affect the innermost, but can use loop name


# Switch

Can have two forms, with expression and without. Switch-local var is possible,
like with `if`

```go
// with an expression
switch size:=len(word); size {
    case 1, 2: ..
    case 3:
    default: ..
}

// bare, can have arbitrary conditions
switch size:=len(word); {
    case size < 5: ..
}
```
no fall-through, unless `fallthrough` is the last command in block. `break`
breaks from case

# Functions

variables are passed by value, unless pointers (map, slice)
```go
func fname1 (param1 int, param2 str) int { .. }
func fname2 (param1, param2 int) int { .. }      // both params int
func fname3 (param1 ...int) int { .. }           // variadic param
x := fname3(1, 2, 3, 4)                          // variadic param call
x := fname3(slice1...)                           // variadic param call with a slice:
                                                 // three dots needed
func fname4(param1 int) (int, int) { .. }        // returns 2 int values
```
Function can have named return vars: they are function-scoped (and seen by
defer) and bare return returns them
```go
func fname5 () (res int, err error) {
    res = rand.Intn(10)
    return
}
```
Functions are first-class, can have anonymous functions, function types and
closures

```go
type ftype1 func(int)int   // int → int function type
f := func(int i) int {     // variable f of the same type
    return i+1
}
```

Currying is possible in the following way:
```go
func Add(x, y int) int {
    return x + y
}

func AddTo(x int) func(int) int {   // returns a function
    return func(y int) int {        // one param
        return Add(x, y)            // in in the closure
    }
}
```
Currying a *method* makes a closure under the hood

Special keyword: `defer <callable>`. Defers action till the end of the
function, but parameters are evaluated immediately. Defer can access named
return params, if any.

```go
func f() int {
    a := 1
    defer fmt.Println(a)   // executed in reverse order: 2nd
    defer func() {         // 1st
        ...
    }
}
```

# Pointers
```go
var p *string       // pointer to a string
p := &stringvar     // pointer to a string variable stringvar
*p == stringvar     // true
p := new(string)    // pointer to a string
```
`&` operator can't be applied to a basic type, no `&5` or `&"test"`!

`&` of a slice or map element is risky; changing map/slice might rearrange
elements in memory, and the pointer will point to a stale data.

## Pointers magic

pointer to a struct is automatically dereferenced, see [structs](#structs)

Also for pointer- and value-receiving functions, Go applies `*` and `&`
automatically. But for the pointer receiver, the value couldn't be a r-value:
```go
Point{1,2}.PointerRecieverMethod(arg) // ERROR
```

# Methods

Can only be declared on package level, and should be in the same package as the
type
```go
type Person struct { .. }
func (p Person) to_string() string { .. }
```
Pointer-methods vs value-methods: appear to convert automatically, but
value-methods have a copy of the value


# Interfaces

Interfaces enforce methods on types. Type doesn't explicitely state it
implements interface, compiler does the check

```go
type Stringer interface {
    String() string         // to be a stringer, a type must have a
                            // method String() string
}
```

Interfaces work with struct composition:

```go
type Point struct {
    X float64
    Y float64
}
type Circle struct {
    Point
    R float64
}
func distance (p Point) float64 {
   return math.Hypot(p.X, p.Y)
}
var p1 = Point{1,2}
var c1 = Circle{Point{1,2},3}
distance(p1)       // OK
distance(c1.Point) // OK
distance(c1)       // Error: cannot use Circle as Point value
                   // However
func (p Point) distance2() float64 {
   return math.Hypot(p.X, p.Y)
}
p1.distance2()     // OK
c1.distance2()     // OK, distance got promoted from Point
                   // Explicit interface allows assignment
type Distancer2 interface {
    distance2() float64
}
var x1 Distancer2
x1 = p1            // OK
x1 = c1            // OK
```

To sort, a container should satisfy `sort.Intrface`: `Len()`, `Swap(i, j)`,
`Less(i, j)`. So to sort points we have to define these; but if we want
different criteria? Composition:

```go
type Points []Point
func (p Points) Len() int {
    return len(s)
}
func (p Points) Swap (i, j int) {
    p[i], p[j] = p[j], p[i]
}
func (p Points) Less (i, j int) bool {
    return p[i].distance() < p[j].distance()
}

// ByX is a type; all Points methods promoted
type ByX struct {Points}
func (p ByX) Less(i, j int) bool {
    return p.Points[i].X < p.Points[j].X
}

x := Points{{1,3}, {2, 2}, {0, 3}}
sort.Sort(x)         // x = [{2 2} {0 3} {1 3}]
sort.Sort(ByX{x})    // x = [{0 3} {1 3} {2 2}]
```

struct trick works with interfaces: `sort.Reverse`
```go
// embeds anything implementing sort.Interface
type reverse struct { Interface }

// opposite of a normal Less provided by Interface
func (r reverse) Less (i, j int) bool {
    return r.Interface.Less(j, i)
}

// above is sufficient with a type cast
sort.Sort(reverse{x})

// but for function interface we need a function
func Reverse (data Interface) Interface {
    return &reverse{data}
}
sort.Sort(sort.Reverse(x))
```

`any` is an alias of `interface{}`, that is, something that might have some
methods, that is, anything.

type can be asserted and checked:
```go
var i any = 20
i.(string)           // panic: interface {} is int, not string
s, ok := i.(string)  // OK-idiom: "", false
x, ok := i.(int)     // OK-idiom: 20, true

// var j is strongly typed inside the switch
switch j := i(type) {
case nil:            // nil, j is any
case int:            // int, j is int
case string:         // string, j is string
}
```

How to make any function of given signature satisfy the interface? Example from
http module:

```go
// Handler interface
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Handler-like function (of any name)
type HandlerFunc func(ResponseWriter, *Request)

// Method authomatically available on any handler-like function, which creates
// an "alias" for the function by creating a methid simply calling named function.
// Thus any handler-like function will satisfy Handler interface
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```


# Errors

Errors are interfaces:
```go
type error interface {
    Error() string
}
```

`errors.New(s string)` produces new error. More convenient is `fmt.Errorf()`,
it supports `%w` verb for wrapping errors
```go
f, err := os.Open("filename")
if err != nil {
    return fmt.Errorf("In myfunc: %w", err)
}
```
errors.Unwrap allows to unwrap the nested error.

For more structured errors, one needs custom type:
```go
type MyError struct {
    code int
    msg  string
}

func (me MyError) Error() string {
    return me.msg
}
```
