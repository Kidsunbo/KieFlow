cmake_minimum_required(VERSION 3.17)

project(cppflow)

set(CMAKE_CXX_STANDARD 17)

add_library(cppflow CPPFlow.cpp)
target_include_directories(cppflow PUBLIC .)