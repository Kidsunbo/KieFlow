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

Because of Go's lack of generic prior to version 2, a Python3 script is provided for you to generate the code you want.

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

## Flow with sub-flow
```go
_ = flow.Prepare(InputParam{}, PrepareData).
    Do(SimpleFunc1).
    Do(SimpleFunc4).Do(SimpleFunc5).Do(SimpleFunc6).
    IfSubPath(CondTrue,
        NewFlow().Do(Fun2WithData).Do(Fun1WithData).IfSubPath(CondFalse, 
        	NewFlow().Do(SimpleFunc4).If(CondTrue, Fun2WithData).Else(SimpleFunc4).Do(SimpleFunc5)).
        ElseIfSubPath(CondTrue, NewFlow().Do(SimpleFunc5).If(CondTrue, Fun2WithData)).
        Else(SimpleFunc4).Do(SimpleFunc5)).
    ElseSubPath(NewFlow().Do(SimpleFunc4).Do(SimpleFunc5).Do(SimpleFunc6)).
    Do(SimpleFunc5).
    Wait()
```

# API

| Name | Signature|  Note|
|------|:------:|--------|
|Normal Flow| `Do`|Register some functors to run later. This is the most common API in GoFlow|
|If Flow| `If` |Register the condition and some functors, it's the same as the if expression but in functional form|
|ElseIf Flow| `ElseIf`| The same with `If`, but it only appears behind `If` and `ElseIf`|
|Else Flow| `Else`| The same with `ElseIf` but it does not need condition|
|SubPath If Flow| `IfSubPath` | It's exactly the same with `If`, while it accepts a sub-flow as parameters. It will `Attach` the data from the parent flow. Notice that he logger will not be inherited|
|SubPath ElseIf Flow| `ElseIfSubPath` | It's exactly the same with `Elseif` while a sub-flow is expected. The same notation of `IfSubPath` is still applied here|
|SubPath Else Flow| `ElseSubPath` |It's exactly the same with `Else` while a sub-flow is expected. The same notation of `IfSubPath` is still applied here|
|For Flow| `For` | Register some functors and run them for several times. The first parameter is the times that user expects these functors run |
|Parallel Flow| `Parallel` | Run some functors in parallel. Only if all the functors finish with or without success, this node will end and return the result if there is any |
|Prepare Flow| `Prepare` | Given some input parameters and a prepare function follows the interface `IPrepareFunc`. GoFlow will use this function to prepare all the data and stored in the flow.
|Attach Function| `Attach`| Attach the result and data from flow A to flow B. All the change to either of A and B will be seen by the other flow |
|Inherit Function| `Inherit` | Attach the result and data from another flow like `Attach` does, and use the same `OnSuccess` and `OnFail` handler of that flow|
|Note Function| `SetNote` | Set the not to a certain node and the note can be accessed from Logger|
|Begin Logger| `SetBeginLogger`| Set the begin logger to a certain node. The parameter must implement `INodeBeginLogger` interface |
|End Logger| `SetEndLogger`| Set the end logger to a certain node. The parameter must implement `INodeEndLogger` interface |
|Global Begin Logger| `SetGlobalBeginLogger`| Set the begin logger to all the nodes which do not have a begin logger |
|Global End Logger| `SetGlobalEndLogger`| Set the end logger to all the nodes which do not have a end logger |
|On Success| `OnSuccess`| A function that will run only if the flow exits successfully. It must follow the interface `IOnSuccessFunc` |
|On Fail| `OnFail`| A function that will run only if the flow fails to exit successfully. It must follow the interface `IOnFailFunc` |
|Wait The Result| `Wait` | Run all the registered nodes and give out result to the caller |




# Thanks

Thank me:)