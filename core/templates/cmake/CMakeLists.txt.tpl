set(WCOSA_PATH "{{wcosa-path}}")
set(VER {{cmake-version}})
set(NAME {{cmake-project-name}})

# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "${WCOSA_PATH}/toolchain/cmake/CosaToolchain.cmake")

cmake_minimum_required(VERSION ${VER})

project(${NAME} C CXX ASM)

# add search paths for all the user libraries and build them
% lib-search
{{include-directories({{lib-path}})}}
{{generate_arduino_library({{name}}\n\tSRCS {{srcs}}\n\tHDRS {{hdrs}}\n\tBOARD {{board}})\n}}
% end

# create the firmware
% firmware-gen
{{generate_arduino_firmware({{name}}\n\tSRCS {{srcs}}\n\tARDLIBS {{cosa-libs}}\n\tLIBS {{libs}}\n\tPORT {{port}}\n\tBOARD {{board}})}}
% end
