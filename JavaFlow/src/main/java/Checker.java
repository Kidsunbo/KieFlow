import java.util.LinkedList;

public class Checker<T> {

    public Checker() {

    }

    public CheckEngine<T> Prepare(T data) {
        System.out.println("Prepare");
        return new CheckEngine<>(data);
    }
}