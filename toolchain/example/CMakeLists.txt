generate_arduino_firmware(blink
        SRCS main.cpp
        ARDLIBS FAT16 ManchesterCodec exlib
        BOARD mega2560)

add_library(GUCCI libraries/exlib/exlib.h libraries/exlib/exlib.cpp)
target_link_libraries(blink GUCCI)
