#include "Cosa/Power.hh"
#include "Cosa/InputPin.hh"
#include "Cosa/OutputPin.hh"
#include "Cosa/RTT.hh"
#include "Cosa/Watchdog.hh"

// Use the built-in led
OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
    Watchdog::begin();

    Power::set(SLEEP_MODE_PWR_DOWN);
}

void loop() {
    ledPin.on();

#ifdef USE_WATCHDOG_SHUTDOWN
    Watchdog::begin(16);
#endif

    delay(1);

#ifdef USE_WATCHDOG_SHUTDOWN
    Watchdog::end();
#endif

    ledPin.off();

#ifdef USE_WATCHDOG_SHUTDOWN
    Watchdog::begin(512);
#endif

    delay(2000);

#ifdef USE_WATCHDOG_SHUTDOWN
    Watchdog::end();
#endif
}
