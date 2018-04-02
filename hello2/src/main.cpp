#include "Brother/brother.h"

OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
}

void loop() {
    ledPin.on();
    delay(50);
    ledPin.off();
    delay(500);
    //int g = SEVEN;
    int y = number();
}
