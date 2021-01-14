public class Main {

    public static void main(String[] args) {

        var f1 = new Functor1();
        var f2 = new Functor2();
        var f3 = new Functor3();
        var f4 = new Functor4();
        var f5 = new Functor5();
        var cond1false = new Condition1();
        var cond2true = new Condition2();

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

        System.out.println(result);
    }
}

record Data(String name, int age){}


class Functor1 implements ICallable<Data>{

    @Override
    public Result call(Data data) {
        System.out.println("Functor1: "+data);
        return new Result(0,"",null);
    }
}

class Functor2 implements ICallable<Data>{

    @Override
    public Result call(Data data) {
        System.out.println("Functor2: "+data);
        return new Result(0,"",null);
    }
}

class Functor3 implements ICallable<Data>{

    @Override
    public Result call(Data data) {
        System.out.println("Functor3: "+data);
        return new Result(10000,"something wrong",null);
    }
}

class Functor4 implements ICallable<Data>{

    @Override
    public Result call(Data data) {
        System.out.println("Functor4: "+data);
        return new Result(0,"",null);
    }
}

class Functor5 implements ICallable<Data>{

    @Override
    public Result call(Data data) {
        System.out.println("Functor5: "+data);
        return new Result(0,"",null);
    }
}


class Condition1 implements IBooleanFunc<Data>{

    @Override
    public boolean func(Data data) {
        return false;
    }
}

class Condition2 implements IBooleanFunc<Data>{

    @Override
    public boolean func(Data data) {
        return true;
    }
}