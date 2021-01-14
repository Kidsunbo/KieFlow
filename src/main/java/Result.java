//public record Result(int statusCode, String msg, Exception err) {}

public class Result{
    private int statusCode = 0;
    private String msg = "";
    private Exception err = null;

    public Result(int statusCode,String msg,Exception err){
        this.statusCode = statusCode;
        this.msg = msg;
        this.err = err;
    }

    public int statusCode(){
        return statusCode;
    }

    public String msg(){
        return msg;
    }

    public Exception err(){
        return err;
    }

    @Override
    public String toString() {
        return "Result{" +
                "statusCode=" + statusCode +
                ", msg='" + msg + '\'' +
                ", err=" + err +
                '}';
    }
}