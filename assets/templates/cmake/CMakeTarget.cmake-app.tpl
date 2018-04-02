# include and build libraries
include("{{target-name}}Libs.cmake")

# include and build the main framework
include("{{target-name}}Framework.cmake")




# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "${WCOSA_PATH}/toolchain/cmake/CosaToolchain.cmake")

project(${NAME} C CXX ASM)

# add search paths for all the user libraries and build them
% lib-search
{{include_directories("{{lib-path}}")}}
{{generate_arduino_library({{name}}\n\tSRCS {{srcs}}\n\tHDRS {{hdrs}}\n\tBOARD {{board}})}}
{{target_compile_definitions({{name}} PRIVATE __AVR_Cosa__ {{custom-definitions}})}}
% end

file(GLOB_RECURSE SRC_FILES "../../src/*.cpp" "../../src/*.cc" "../../src/*.c")

# create the firmware
% firmware-gen
{{generate_arduino_firmware({{name}}\n\tSRCS ${SRC_FILES}\n\tLIBS {{libs}}\n\tPORT {{port}}\n\tBOARD {{board}})}}
{{target_compile_definitions({{name}} PRIVATE __AVR_Cosa__ {{custom-definitions}})}}
% end
