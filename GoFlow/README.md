# GoFlow

---------

![Go](https://github.com/Kidsunbo/KieFlow/workflows/Go/badge.svg)


This is the flow project for Golang. Because of lack of virtual method, some code is duplicated.

# Requirement

To use this library, one should follow some rules.
#### 1. The `prepare` function should implement `IPrepareFunc`

input is the wrapper for passing arguments to the function due to lack of perfect forwarding. The programmer should take the
responsibility for the logistic and the error handling.

**Do not return**

#### 2. The `condition` function for `If` and `ElseIf` should implement `IBoolFunc`

#### 3. The `BeginLogger` should implement `INodeBeginLogger`

#### 4. The `EndLogger` should implement `INodeEndLogger`

#### 5. The `OnSuccess` should implement `IOnSuccessFunc`

#### 6. The `OnFail` should implement `IOnFailFunc`

#### 7. All the method used as a task should implement `ICallable`

# Example

## Simple Workflow

```go
_ = NewFlow().
    Do(Func1).
    Prepare(Prepare,PrepareTest{}).
    If(CondTrue,Func2,Func3).
    Else(Func4,Func5).
    Do(Func4,Func5,Func6).
    If(CondFalse,Func2,Func3).
    ElseIf(CondFalse,Func6,Func5).
    Else(Func1,Func2).
    For(3,Func9).
    Parallel(Func9,Func2,Func1,Func6).
    If(CondTrue,Func2).
    Wait()
```
