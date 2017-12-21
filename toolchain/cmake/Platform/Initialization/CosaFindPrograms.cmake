write_sep()

# Find examples and libraries
find_file(COSA_EXAMPLES_PATH
        NAMES examples
        PATHS ${COSA_SDK_PATH}
        DOC "Path to directory containing Cosa built-in examples.")

find_file(COSA_LIBRARIES_PATH
        NAMES libraries
        PATHS ${COSA_SDK_PATH}
        DOC "Path to directory containing Cosa libraries")

info("Found paths")
info("Examples:  ${COSA_EXAMPLES_PATH}")
info("Libraries: ${COSA_LIBRARIES_PATH}")

if (NOT COSA_EXAMPLES_PATH OR NOT EXISTS ${COSA_EXAMPLES_PATH})
    fatal("Failed to find COSA_EXAMPLES_PATH to `examples`")
endif ()
if (NOT COSA_LIBRARIES_PATH OR NOT EXISTS ${COSA_LIBRARIES_PATH})
    fatal("Failed to find COSA_LIBRARIES_PATH to `libraries`")
endif ()

#===================================================#
# Search paths for `avrdude`
# Keep a log of lists here for posterity
#
# MacOS (commandline)
# - /usr/local/bin/avrdude
#
#===================================================#

# If ARDUINO_SDK_PATH is provided, search there first
if (EXISTS ${ARDUINO_SDK_PATH})
    find_program(COSA_AVRDUDE_PROGRAM
            NAMES avrdude
            PATHS ${ARDUINO_SDK_PATH}
            PATH_SUFFIXES hardware/tools hardware/tools/avr/bin
            NO_DEFAULT_PATH)
endif ()

# Search known paths first
set(cosa_avrdude_known_paths
        /usr/bin
        /usr/local/bin
        /usr/local/Cellar/avrdude/6.3/bin)
find_program(COSA_AVRDUDE_PROGRAM
        NAMES avrdude
        PATHS ${cosa_avrdude_known_paths}
        DOC "Path to avrdude programmer binary.")

# Search through environment PATH
find_program(COSA_AVRDUDE_PROGRAM
        NAMES avrdude
        DOC "Path to avrdude programmer binary.")

info("avrdude:   ${COSA_AVRDUDE_PROGRAM}")
if (NOT COSA_AVRDUDE_PROGRAM OR NOT EXISTS ${COSA_AVRDUDE_PROGRAM})
    fatal("Unable to find path to `avrdude`")
endif ()

#===================================================#
# Search paths for `avr-size`
# Keep a log of lists here for posterity
#
# MacOS (commandline)
# - /usr/local/bin/avr-size
#
#===================================================#

if (EXISTS ${ARDUINO_SDK_PATH})
    find_program(COSA_AVR_SIZE_PROGRAM
            names avr-size
            PATHS ${ARDUINO_SDK_PATH}
            PATH_SUFFIXES hardware/tools hardware/tools/avr/bin
            NO_DEFAULT_PATH)
endif ()

set(cosa_avr_size_known_paths
        /usr/bin
        /usr/local/bin
        /usr/local/Cellar/avr-binutils/2.29/bin)

find_program(COSA_AVR_SIZE_PROGRAM
        names avr-size
        PATHS ${cosa_avr_size_known_paths}
        DOC "Path to avr-size program binary.")

find_program(COSA_AVR_SIZE_PROGRAM
        names avr-size
        DOC "Path to avr-size program binary.")

info("avr-size:  ${COSA_AVR_SIZE_PROGRAM}")
if (NOT COSA_AVR_SIZE_PROGRAM OR NOT EXISTS ${COSA_AVR_SIZE_PROGRAM})
    fatal("Unable to find path to `avr-size`")
endif ()