import java.util.ArrayList;
import java.util.LinkedList;

public class IfNode<T> extends Checkable<T>{
    private IBooleanFunc<T> condition = null;


    IfNode(T data, Wrapper<Result> parentResult, IBooleanFunc<T> condition, ICallable<T>... functors) {
        super(data, parentResult, NodeType.IfNode, functors);
        this.condition = condition;
    }

    @Override
    Result check() {
        if(condition==null){
            throw new RuntimeException("Please give condition");
        }

        if(condition.func(this.getData())){
            for(var functor : this.getFs()){
                var result = functor.call(this.getData());
                if(result.err()!=null || result.statusCode()!=0){
                    return result;
                }
            }
            Checkable<T> current = this.getNext();
            while(current!=null && (current.getNodeType() == NodeType.ElseIfNode || current.getNodeType() == NodeType.ElseNode)){
                current.setShouldSkip(true);
                current = current.getNext();
            }
        }
        return this.getParentResult();
    }
}
