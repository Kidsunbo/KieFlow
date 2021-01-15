import java.util.LinkedList;

public class ElseNode<T> extends Checkable<T>{
    ElseNode(T data, Wrapper<Result> parentResult, ICallable<T>... functors) {
        super(data, parentResult,NodeType.ElseNode, functors);
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
