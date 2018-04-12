set(VER 3.0.0)
set(NAME {{project-name}})

set(CMAKE_TOOLCHAIN_FILE "{{toolchain-path}}")

cmake_minimum_required(VERSION ${VER})
project(${NAME} C CXX ASM)

file(GLOB_RECURSE LIB_SRC_FILES "../src/*.cpp" "../src/*.cc" "../src/*.c")
generate_arduino_library({{project-name}}
	SRCS ${LIB_SRC_FILES}
    BOARD {{board}})
target_compile_definitions({{project-name}} PRIVATE __AVR_{{framework}}__ {{lib-flags}})
include_directories("../include")
{{link-library}}
file(GLOB_RECURSE SRC_FILES "../test/*.cpp" "/../test/*.cc" "../test/*.c")
generate_arduino_firmware({{target-name}}
    SRCS ${SRC_FILES}
    LIBS {{project-name}}
    BOARD {{board}})
target_compile_definitions(test PRIVATE __AVR_{{framework}}__ {{target-flags}})
