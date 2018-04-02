#include "Brother/brother.h"
#include "Cosa/Output.h"

OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
}

void loop() {
    ledPin.on();
    delay(50);
    ledPin.off();
    delay(500);
    int y = number();
}
