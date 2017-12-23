set(WCOSA_CMD "python {{wcos-path}}/core/wcosa.py")
set(VER {{cmake-version}})
set(NAME {{cmake-project-name}})

cmake_minimum_required(VERSION ${VER})

project(${NAME} C CXX ASM)

SET(CMAKE_C_COMPILER avr-gcc)
SET(CMAKE_CXX_COMPILER avr-g++)
SET(CMAKE_CXX_FLAGS_DISTRIBUTION "{{cxx_flags}}")
SET(CMAKE_C_FLAGS_DISTRIBUTION "{{cc_flags}}")
set(CMAKE_CXX_STANDARD 11)

% def-search
{{add_definitions({{define}})}}
% end

# add search paths for all the user libraries
% lib-search
{{include-directories({{lib-path}})}}
% end

add_custom_target(
    WCOSA_BUILD ALL
    COMMAND ${WCOSA_CMD} build --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_CLEAN ALL
    COMMAND ${WCOSA_CMD} clean --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPDATE_ALL ALL
    COMMAND ${WCOSA_CMD} update --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPLOAD ALL
    COMMAND ${WCOSA_CMD} upload
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

% firmware-gen
{{add_executable({{srcs}})}}
% end
