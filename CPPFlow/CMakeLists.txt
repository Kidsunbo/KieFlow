cmake_minimum_required(VERSION 3.17)
project(CPPFlow)

set(CMAKE_CXX_STANDARD 17)

add_subdirectory(CPPFlow)

add_executable(CPPFlow main.cpp)
target_link_libraries(CPPFlow PUBLIC cppflow)