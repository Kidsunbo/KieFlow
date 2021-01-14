import java.util.ArrayList;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;


public class ParallelNode<T> extends Checkable<T> {
    ParallelNode(T data, Wrapper<Result> parentResult, ICallable<T>... functors) {
        super(data, parentResult, NodeType.ParallelNode, functors);
    }

    @Override
    Result check() {

        var threadPool = Executors.newFixedThreadPool(this.getFs().size(),r -> {
            Thread t = Executors.defaultThreadFactory().newThread(r);
            t.setDaemon(true);
            return t;
        });
        var results = new ArrayList<Future<Result>>();
        for(var functor : this.getFs()){
            results.add(threadPool.submit(()->functor.call(this.getData())));
        }
        try {
            for(var result : results){
                var r = result.get();
                if (r.statusCode() != 0 || r.err() != null) {
                    return r;
                }
            }
        }catch (Exception e){
            threadPool.shutdown();
            e.printStackTrace();
            return new Result(10001,"interrupted",null);
        }
        return this.getParentResult();

    }
}
