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

find_file(COSA_BOOTLOADERS_PATH
        NAMES bootloaders
        PATHS ${COSA_SDK_PATH}
        Doc "Path to directory containing the Cosa bootloaders images and sources.")

find_file(COSA_BOARDS_PATH
        NAMES boards.txt
        PATHS ${COSA_SDK_PATH}
        DOC "Path to Cosa boards definition file.")

# Cosa variants are located under `variants/arduino`
set(COSA_VARIANTS_PATH ${COSA_VARIANTS_PATH}/arduino)

message(STATUS "Founds paths")
message(STATUS "Cores:       ${COSA_CORES_PATH}")
message(STATUS "Variants:    ${COSA_VARIANTS_PATH}")
message(STATUS "Bootloaders: ${COSA_BOOTLOADERS_PATH}")
message(STATUS "Boards:      ${COSA_BOARDS_PATH}")

if (NOT COSA_CORES_PATH OR NOT EXISTS ${COSA_CORES_PATH})
    message(FATAL_ERROR "Failed to find COSA_CORES_PATH to `cores`")
endif ()
if (NOT COSA_VARIANTS_PATH OR NOT EXISTS ${COSA_VARIANTS_PATH})
    message(FATAL_ERROR "Failed to find COSA_VARIANTS_PATH to `variants/arduino`")
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
message(STATUS "Reading variants from directory")
file(GLOB variant_sub_dir ${COSA_VARIANTS_PATH}/*)
unset(COSA_VARIANTS CACHE)
foreach (variant_dir ${variant_sub_dir})
    if (IS_DIRECTORY ${variant_dir})
        get_filename_component(variant ${variant_dir} NAME)
        set(COSA_VARIANTS ${COSA_VARIANTS} ${variant} CACHE INTERNAL "A list of registered variants")
        set(${variant}.path ${variant_dir} CACHE INTERNAL "The path to variant ${variant}")
        if (WCOSA_DEBUG)
            message(STATUS "Variant [${variant}]: ${${variant}.path}")
        endif ()
    endif ()
endforeach ()
list(LENGTH COSA_VARIANTS length_cosa_variants)
message(STATUS "Found and cached ${length_cosa_variants} variants")
unset(variant_sub_dir)
unset(length_cosa_variants)

# Read in cores
message(STATUS "Reading cores from directory")
file(GLOB cores_sub_dir ${COSA_CORES_PATH}/*)
unset(COSA_CORES CACHE)
foreach (core_dir ${cores_sub_dir})
    if (IS_DIRECTORY ${core_dir})
        get_filename_component(core ${core_dir} NAME)
        set(COSA_CORES ${COSA_CORES} ${core} CACHE INTERNAL "A list of registered cores")
        set(${core}.path ${core_dir} CACHE INTERNAL "The path to core ${core}")
        if (WCOSA_DEBUG)
            message(STATUS "Core [${core}]: ${${core}.path}")
        endif ()
    endif ()
endforeach ()
list(LENGTH COSA_CORES length_cosa_cores)
message(STATUS "Found and cached ${length_cosa_cores} cores")
unset(cores_sub_dir)
unset(length_cosa_cores)
