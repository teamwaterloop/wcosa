set(VER 3.0.0)
set(NAME {{project-name}})

set(CMAKE_TOOLCHAIN_FILE "{{toolchain-path}}")

cmake_minimum_required(VERSION ${VER})
project(${NAME} C CXX ASM)

file(GLOB_RECURSE SRC_FILES "../src/*.cpp" "../src/*.cc" "../src/*.c")

# create the firmware
generate_arduino_firmware({{target-name}}
    SRCS ${SRC_FILES}
    BOARD {{board}})
target_compile_definitions({{target-name}} PRIVATE __AVR_{{framework}}__ {{target-flags}})

