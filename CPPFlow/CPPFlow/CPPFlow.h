//
// Created by sunbo on 2021/1/19.
//

#ifndef CPPFLOW_CPPFLOW_H
#define CPPFLOW_CPPFLOW_H

#include <memory>
#include <functional>
#include <string>
#include <vector>


enum class NodeType {
    Unspecific, NormalNode, IfNode, ElseNode, ForNode, ParallelNode, ElseIfNode
};

template<typename T,typename = void>
struct HasStatusCode: std::false_type {};

template<typename T>
struct HasStatusCode<T,std::void_t<decltype(std::declval<T>().statusCode)>>:std::true_type {};

template<typename T>
constexpr bool HasStatusCodeV = HasStatusCode<T>::value;

template<typename Data, typename Result>
using ICallable = std::function<std::shared_ptr<Result>(std::shared_ptr<Data>)>;

template<typename Data>
using IBoolFunc = std::function<bool(std::shared_ptr<Data>)>;

template<typename Data, typename PrepareInput, typename Result>
using IPrepareFunc = std::function<std::shared_ptr<Result>(std::shared_ptr<Data>, PrepareInput)>;

template<typename Data>
using INodeBeginLogger = std::function<void(std::string_view,std::shared_ptr<Data>)>;

template<typename Data, typename Result>
using INodeEndLogger = std::function<void(std::string_view, std::shared_ptr<Data>, std::shared_ptr<Result>)>;

template<typename Data, typename Result>
using IOnSuccessFunc = std::function<void(std::shared_ptr<Data>,std::shared_ptr<Result>)>;

template<typename Data, typename Result>
using IOnSuccessFunc = std::function<void(std::shared_ptr<Data>,std::shared_ptr<Result>)>;

template<typename Data, typename Result, std::enable_if_t<HasStatusCodeV<Result>>* =nullptr>
class CPPFlow;

template<typename Data, typename Result>
class BasicFlowNode{
    friend class CPPFlow<Data,Result>;

protected:
    NodeType nodeType = NodeType::Unspecific;
    BasicFlowNode *next = nullptr;
    std::shared_ptr<Data> data = nullptr;
    bool shouldSkip = false;
    std::shared_ptr<std::shared_ptr<Result>> parentResult = nullptr;
    INodeBeginLogger<Data> beginLogger;
    INodeEndLogger<Data,Result> endLogger;
    std::string note;

    BasicFlowNode(std::shared_ptr<Data> data,std::shared_ptr<std::shared_ptr<Result>> result, NodeType type):
            data(data),parentResult(result),nodeType(type){}

    void setParentResult(std::shared_ptr<Result> result){
        *this->parentResult = result;
    }

    std::shared_ptr<Result> getParentResult(){
        return *parentResult;
    }

    virtual void run(){
        if(shouldSkip || this->getParentResult()->statusCode != 0){
            return;
        }
        if(beginLogger!=nullptr){
            beginLogger(note,data);
        }

        auto result = implTask();
        if(result != nullptr){
            setParentResult(result);
        }

        if(endLogger!= nullptr){
            endLogger(note,data,getParentResult());
        }
    }

    virtual std::shared_ptr<Result> implTask() = 0;

};

template<typename Data, typename Result>
class IfNode final:public BasicFlowNode<Data,Result>{
    friend class CPPFlow<Data,Result>;

private:
    IBoolFunc<Data> condition;
    std::vector<ICallable<Data,Result>> functors;
protected:
    template<typename... FUNC, std::enable_if_t<(std::is_convertible_v<FUNC,ICallable<Data,Result>>&&...)>* =nullptr>
    IfNode(std::shared_ptr<Data> data, std::shared_ptr<std::shared_ptr<Result>> result, IBoolFunc<Data> condition, FUNC... functors):BasicFlowNode<Data,Result>(data,result,NodeType::IfNode){

        this->condition = condition;
        (this->functors.push_back(std::forward<FUNC>(functors)), ...);
    }

    void run() override {
        if(this->shouldSkip || this->getParentResult()->statusCode !=0){
            return;
        }

        auto res = this->implTask();
        if(res!= nullptr){
            this->setParentResult(res);
        }
    }

    std::shared_ptr<Result> implTask() override {
        if(this->condition == nullptr){
            return std::make_shared<std::shared_ptr<Result>>(std::make_shared<Result>());
        }

        if(this->condition(this->data)) {
            if(this->beginLogger!=nullptr){
                this->beginLogger(this->note,this->data);
            }

            for (auto f:functors) {
                auto res = f(nullptr);
                if(res != nullptr && res->statusCode !=0 ){
                    if(this->endLogger!= nullptr){
                        this->endLogger(this->note,this->data,res);
                    }
                    return res;
                }
            }

            for(auto current = this->next;current!=nullptr && (current->nodeType == NodeType::ElseIfNode || current->nodeType == NodeType::ElseNode);current = current->next){
                current->shouldSkip = true;
            }
            if(this->endLogger!= nullptr){
                this->endLogger(this->note,this->data,this->getParentResult());
            }
        }
        return this->getParentResult();
    }
};






template<typename Data, typename Result, std::enable_if_t<HasStatusCodeV<Result>>*>
class CPPFlow {
private:
    std::shared_ptr<Data> data;

public:

};


#endif //CPPFLOW_CPPFLOW_H
