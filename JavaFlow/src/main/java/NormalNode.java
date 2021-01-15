public class NormalNode<T> extends Checkable<T>{
    NormalNode(T data, Wrapper<Result> parentResult, ICallable<T>... functors) {
        super(data, parentResult, NodeType.NormalNode, functors);
    }

    @Override
    Result check() {
        for(var functor : this.getFs()){
            var result = functor.call(this.getData());
            if(result.err()!=null || result.statusCode()!=0){
                return result;
            }
        }
        return this.getParentResult();
    }
}
