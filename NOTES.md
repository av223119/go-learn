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

