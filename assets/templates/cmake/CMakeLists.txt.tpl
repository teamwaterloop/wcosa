# define toolchain path, project name and version
set(TOOLCHAIN_PATH "{{toolchain-path}}")
set(VER {{cmake-version}})
set(NAME {{project-name}})

# include the file for the target we are building
include({{target-path}})

cmake_minimum_required(VERSION ${VER})
project(${NAME} C CXX ASM)
