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
var x int  # type declared, "zero" value used: 0 for integers
var x = 1  # type inferred, int as the most generic one
var x uint32 = 1  # both type and value set explicitely
var x1, x2 string = "one", "two"  # more vars of the same type
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
x1, x2 := 10, "something"
x := byte(10)  # explicit type
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
var x = [5]int{1, 4:9}  # sparse, [1, 0, 0, 0, 9]
var x = [...]int{1, 3:8}  # special syntax: [...] means compile-time size calculation
var x [2][3]int   # two-dimensional array [[0, 0, 0], [0, 0, 0]]
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
s1 := make([]int, initialLength [, initialCapacity])  # another way to create a slice
clear(s1)  # clears elements, keeps length and capacity
copy(dest, src) #  copies up to LENGTH, returns copied length. Only slices, not arrays.
```
