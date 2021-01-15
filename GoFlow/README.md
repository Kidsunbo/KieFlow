# GoFlow

---------

![Go](https://github.com/Kidsunbo/KieFlow/workflows/Go/badge.svg)


This is the flow project for Golang. Because of lack of virtual method, some code is duplicated.
This project is mainly solve the problem that the code will be too tediously long after several iterations.

# Requirement

To use this library, one should follow some rules.
#### 1. The `prepare` function should implement `IPrepareFunc`

The `input` is the wrapper for passing arguments to the function due to lack of perfect forwarding. The programmer should take the
responsibility for the logistic and the error handling.

#### 2. The `condition` function for `If` and `ElseIf` should implement `IBoolFunc`

The logistic should be taken by the programmer and always return the result, a boolean value.

#### 3. The `BeginLogger` should implement `INodeBeginLogger`

You can do what ever you want to do with the result and data as you wish. If it's set, the log will be printed when a node starts.

#### 4. The `EndLogger` should implement `INodeEndLogger`

You can do what ever you want to do with the result and data as you wish. If it's set, the log will be printed when a node ends.

#### 5. The `OnSuccess` should implement `IOnSuccessFunc`

You can do what ever you want to do with the result and data as you wish. If it's set, the log will be printed when flow runs successfully.

#### 6. The `OnFail` should implement `IOnFailFunc`

You can do what ever you want to do with the result and data as you wish. If it's set, the log will be printed when flow failed to run.

#### 7. All the method used as a task should implement `ICallable`


# Usage

A Python3 script is provided for you to generate the code you want.

```shell
python3 flow.py --data DataTest 
                --result ResultTest 
                --prepare PrepareTest
                -s ~/goflow
                -o ~/workdir
                -p myflow
```
`--data`, `--result` and `--prepare` is mandatory.
`--data` is for replacing the `_Data` type, `--result` is for replacing the `_Result` type while
`--prepare` is for replacing `_PrepareInput`.

`-s` or `--source` is the directory where the template files located in.

`-o` or `--output` is the output directory for the generated file

`-p` or `--package` is the package name for the generated file

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

## Simple Logger

```go
_ = NewFlow().
    Do(Func1).SetNote("This might be a check").SetBeginLogger(BeginLogger).SetEndLogger(EndLogger).
    Do(Func2).
    Do(Func3).
    SetGlobalBeginLogger(BeginLogger).
    Wait()
```

## Success/Fail Handler
```go
_ = NewFlow().
    Do(Func1).
    Do(Func2).
    Do(Func3).
    Do(Func4).
    Do(Func5).
    Do(Func6).
    Do(Func7).
    Do(Func8).
    Do(Func9).
    Do(Func1).
    Do(Func2).
    Do(Func3).
    OnSuccess(OnSuccessHandle).
    OnFail(OnFailHandle).
    Wait()

```
