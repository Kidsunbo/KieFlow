#JavaFlow

------

![Java CI with Gradle](https://github.com/Kidsunbo/JavaFlow/workflows/Java%20CI%20with%20Gradle/badge.svg)

Java flow is a library that run bunch of functions composed by the code.

# Example

```java
var result =  new Checker<Data>()
        .Prepare(new Data("Tom is Tom",10))
        .Check(f1)
        .Check(f2)
        .If(cond2true,f1,f2)
        .Else(f3,f4)
        .If(cond1false,f1,f2)
        .ElseIf(cond1false,f2,f5)
        .ElseIf(cond2true,f4,f5)
        .Check(f1)
        .If(cond2true,f1,f2)
        .Else(f3,f4)
        .Check(f3)
        .For(3,f4)
        .Parallel(f1,f2,f4)
        .Check(f1)
        .Wait();
```
The `condition variable` used by `If`, `ElseIf` and `Else` should implements `IBooleanFunc`. The `f<num>` should implement `ICallable`.

`Data` and `Result` can be designed by yourself.