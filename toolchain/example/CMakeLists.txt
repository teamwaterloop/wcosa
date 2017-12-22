generate_arduino_firmware(blink
        SRCS main.cpp
        ARDLIBS FAT16 ManchesterCodec exlib
        BOARD mega2560)
