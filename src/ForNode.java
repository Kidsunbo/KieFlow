public class ForNode<T> extends Checkable<T>{
    private int times = 0;

    ForNode(int times, T data, Wrapper<Result> parentResult, ICallable<T>... functors) {
        super(data, parentResult,NodeType.ForNode, functors);
        this.times = times;
    }

    @Override
    Result check() {
        for(int i=0;i<times;i++){
            for(var functor : this.getFs()){
                var result = functor.call(this.getData());
                if(result.err()!=null || result.statusCode()!=0){
                    return result;
                }
            }
        }
        return this.getParentResult();
    }
}
