import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedList;

abstract public class Checkable<T> {
    private final LinkedList<ICallable<T>> fs = new LinkedList<>();
    private NodeType nodeType = NodeType.NotSpecified;
    private Checkable<T> next = null;
    private T data = null;
    private boolean shouldSkip = false;
    private Wrapper<Result> parentResult = null;

    public LinkedList<ICallable<T>> getFs() {
        return fs;
    }

    public NodeType getNodeType() {
        return nodeType;
    }


    public T getData() {
        return data;
    }

    public Checkable<T> getNext() {
        return next;
    }

    public void setNext(Checkable<T> next) {
        this.next = next;
    }

    public boolean isShouldSkip() {
        return shouldSkip;
    }

    public void setShouldSkip(boolean shouldSkip) {
        this.shouldSkip = shouldSkip;
    }

    public Result getParentResult() {
        return parentResult.getData();
    }

    public void setParentResult(Result parentResult) {
        this.parentResult.setData(parentResult);
    }

    public void doCheck(){
        if(this.isShouldSkip() || this.getParentResult().err()!=null || this.getParentResult().statusCode()!=0){
            return;
        }
        var result = check();
        this.setParentResult(result);
    }

    abstract Result check();

    Checkable(T data, Wrapper<Result> parentResult,NodeType type,ICallable<T>... functors){
        this.data = data;
        fs.addAll(Arrays.asList(functors));
        nodeType = type;
        this.parentResult = parentResult;
    }
}
