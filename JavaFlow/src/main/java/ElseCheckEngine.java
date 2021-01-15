import java.util.LinkedList;

public class ElseCheckEngine<T> {
    private T data = null;
    private CheckEngine<T> invoker = null;
    private Wrapper<Result> result = null;
    private LinkedList<Checkable<T>> nodes = null;

    public ElseCheckEngine(T data,CheckEngine<T> invoker, Wrapper<Result> result,LinkedList<Checkable<T>> nodes){
        this.data = data;
        this.invoker = invoker;
        this.result = result;
        this.nodes = nodes;
    }

    public CheckEngine<T> Check(ICallable<T>... functors){
        System.out.println("Check");
        var node = new NormalNode<T>(data,result,functors);
        nodes.getLast().setNext(node);
        nodes.add(node);
        return invoker;
    }

    public CheckEngine<T> For(int times, ICallable<T>... functors){
        System.out.println("For");
        var node = new ForNode<T>(times,data,result,functors);
        nodes.getLast().setNext(node);
        nodes.add(node);
        return invoker;
    }

    public CheckEngine<T> Parallel(ICallable<T>... functors){
        System.out.println("Parallel");
        var node = new ParallelNode<T>(data,result,functors);
        nodes.getLast().setNext(node);
        nodes.add(node);
        return invoker;
    }

    public ElseCheckEngine<T> If(IBooleanFunc<T> condition, ICallable<T>... functors){
        System.out.println("If");
        var node = new IfNode<T>(data,result,condition,functors);
        if(nodes.size()!=0) {
            nodes.getLast().setNext(node);
        }
        nodes.add(node);
        return this;
    }

    public ElseCheckEngine<T> ElseIf(IBooleanFunc<T> condition, ICallable<T>... functors){
        System.out.println("ElseIf");
        var node = new ElseIfNode<T>(data,result,condition,functors);
        nodes.getLast().setNext(node);
        nodes.add(node);
        return this;
    }


    public CheckEngine<T> Else(ICallable<T>... functors){
        System.out.println("Else");
        var node = new ElseNode<T>(data,result,functors);
        nodes.getLast().setNext(node);
        nodes.add(node);
        return invoker;
    }

    public Result Wait(){
        for(var node : nodes){
            node.doCheck();
        }
        return result.getData();
    }

}
