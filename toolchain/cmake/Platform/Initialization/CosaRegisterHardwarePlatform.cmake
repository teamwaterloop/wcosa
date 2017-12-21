write_sep()
message(STATUS "Obtaining hardware information from ${COSA_SDK_PATH}")

# Find paths to directories and files
# Cosa does not provide a `programmers.txt`
find_file(COSA_CORES_PATH
        NAMES cores
        PATHS ${COSA_SDK_PATH}
        DOC "Path to directory containing the Cosa core sources.")

find_file(COSA_VARIANTS_PATH
        NAMES variants
        PATHS ${COSA_SDK_PATH}
        DOC "Path to directory containing the Cosa variant sources.")

find_file(COSA_VARIANTS_PATH
        NAMES arduino
        PATHS ${COSA_VARIANTS_PATH}
        DOC "Path to directory containing the Cosa variant sources.")

find_file(COSA_BOOTLOADERS_PATH
        NAMES bootloaders
        PATHS ${COSA_SDK_PATH}
        Doc "Path to directory containing the Cosa bootloaders images and sources.")

find_file(COSA_BOARDS_PATH
        NAMES boards.txt
        PATHS ${COSA_SDK_PATH}
        DOC "Path to Cosa boards definition file.")

message(STATUS "Founds paths")
message(STATUS "Cores:       ${COSA_CORES_PATH}")
message(STATUS "Variants:    ${COSA_VARIANTS_PATH}")
message(STATUS "Bootloaders: ${COSA_BOOTLOADERS_PATH}")
message(STATUS "Boards:      ${COSA_BOARDS_PATH}")

if (NOT COSA_CORES_PATH OR NOT EXISTS ${COSA_CORES_PATH})
    message(FATAL_ERROR "Failed to find COSA_CORES_PATH to `cores`")
endif ()
if (NOT COSA_VARIANTS_PATH OR NOT EXISTS ${COSA_VARIANTS_PATH})
    message(FATAL_ERROR "Failed to find COSA_VARIANTS_PATH to `variants`")
endif ()
if (NOT COSA_BOOTLOADERS_PATH OR NOT EXISTS ${COSA_BOOTLOADERS_PATH})
    message(FATAL_ERROR "Failed to find COSA_BOOTLOADERS_PATH to `bootloaders`")
endif ()
if (NOT COSA_BOARDS_PATH OR NOT EXISTS ${COSA_BOARDS_PATH})
    message(FATAL_ERROR "Failed to find COSA_BOARDS_PATH to `boards.txt`")
endif ()

# Read in `boards.txt`
set(SETTINGS_LIST COSA_BOARDS)
set(SETTINGS_PATH ${COSA_BOARDS_PATH})
include(LoadArduinoPlatformSettings)

# Display example boards read
list(GET COSA_BOARDS 0 example_board_0)
list(GET COSA_BOARDS 1 example_board_1)
list(GET COSA_BOARDS 2 example_board_2)
message(STATUS "Parsed `boards.txt` (e.g. ${example_board_0}, ${example_board_1}, ${example_board_2})")
unset(example_board_0)
unset(example_board_1)
unset(example_board_2)

# Read in variant boards
file(GLOB variant_sub_dir ${COSA_VARIANTS_PATH}/*)
foreach (variant_dir ${variant_sub_dir})
    message(${variant_dir})
endforeach ()
