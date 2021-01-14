import java.nio.file.Watchable;
import java.util.LinkedList;

public class CheckEngine<T> {
    private T data = null;
    private LinkedList<Checkable<T>> nodes = new LinkedList<>();
    private Wrapper<Result> result = new Wrapper<>();


    public CheckEngine(T data){
        this.result.setData(new Result(0,"",null));
        this.data = data;
    }

    public CheckEngine<T> Prepare(T data){
        System.out.println("Prepare: "+data);
        this.data = data;
        return this;
    }

    public CheckEngine<T> Check(ICallable<T>... functors){
        System.out.println("Check");
        var node = new NormalNode<T>(data,result,functors);
        if(nodes.size()!=0) {
            nodes.getLast().setNext(node);
        }
        nodes.add(node);
        return this;
    }

    public CheckEngine<T> For(int times, ICallable<T>... functors){
        System.out.println("For");
        var node = new ForNode<T>(times,data,result,functors);
        if(nodes.size()!=0) {
            nodes.getLast().setNext(node);
        }
        nodes.add(node);
        return this;
    }

    public CheckEngine<T> Parallel(ICallable<T>... functors){
        System.out.println("Parallel");
        var node = new ParallelNode<T>(data,result,functors);
        if(nodes.size()!=0) {
            nodes.getLast().setNext(node);
        }
        nodes.add(node);
        return this;
    }

    public ElseCheckEngine<T> If(IBooleanFunc<T> condition, ICallable<T>... functors){
        System.out.println("If");
        var node = new IfNode<T>(data,result,condition,functors);
        if(nodes.size()!=0) {
            nodes.getLast().setNext(node);
        }
        nodes.add(node);
        return new ElseCheckEngine<T>(data,this,result,nodes);
    }

    public Result Wait(){
        for(var node : nodes){
            node.doCheck();
        }
        return result.getData();
    }

}
