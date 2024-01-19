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

Only basic types are allowed. No constant structs, maps, slices etc. Otherwise similar to var
```go
const pi = 3.1415
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

`append(slice, item1 [, item2 ..])` or `append(slice, slice2..)` is not in-place

`=` or `[:]` don't copy, but use the same var. `[::]` is like `[:]`, but limits capacity. Absurd overlapping example:
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
array2 := [size]int(slice1)  // can't use [...] or [len(slice)], fixed length ≤ slice_size needed
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

`val, ok = map[key]` ok contains bool (key is in the map or not)

can't use == !=; maps.Equal(a, b)

can't use &m["key"] nor m["key"].field !


# Structs

```go
type mydata struct {
    field1 int
    field2 string
}
v1 := mydata{
    field1: 1
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
Structs don't go along nicely with maps, can't use m[key].field; better use `map[int]*mystruct{}`

# If-then-else

```go
if condition {
} else if condition2 {
} else {
}
// var declaration is possible in condition; var is scoped to block only
if a := rand.Intn(10); a > 5 { }
```
