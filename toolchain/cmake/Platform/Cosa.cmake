cmake_minimum_required(VERSION 3.1.0)

if (NOT COSA_SDK_PATH)
    message(FATAL_ERROR "Error: COSA_SDK_PATH is not defined")
endif ()
if (NOT ARDUINO_CMAKE_PATH)
    message(FATAL_ERROR "Error: ARDUINO_CMAKE_PATH is not defined")
endif ()

if (COSA_SCRIPT_EXECUTED)
    return()
endif ()

# Set module paths to include `arduino-cmake` components
set(ARDUINO_CMAKE_PLATFORM_PATH ${ARDUINO_CMAKE_PATH}/Platform)
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH}
        ${ARDUINO_CMAKE_PLATFORM_PATH}
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Initialization
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/BoardFlags
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Libraries
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Targets
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Sketch
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Examples
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Extras
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Generation)

# Set module paths to include `cosa-cmake` components
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH}
        ${CMAKE_CURRENT_LIST_DIR}/Initialization
        ${CMAKE_CURRENT_LIST_DIR}/Utils
        ${CMAKE_CURRENT_LIST_DIR}/Vendor)

# Include vendored files
include(JsonParser)

# Include utilities
include(CosaOutput)

# Include external utilities
include(CMakeParseArguments)
include(VariableValidator)

# Initialization scripts
include(CosaInitializer)

# Mark configuration as complete
set(COSA_SCRIPT_EXECUTED True)

write_sep()
