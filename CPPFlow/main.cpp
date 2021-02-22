#include <iostream>
#include "CPPFlow/CPPFlow.h"

struct Result{
    int statusCode = 0;
    std::string msg;
};


int main() {
    std::cout << "Hello, World!" << std::endl;
    auto result = std::make_shared<std::shared_ptr<Result>>(std::make_shared<Result>());



    return 0;
}
