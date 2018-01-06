/**
 * @file CosaLED.ino
 * @version 1.0
 *
 * @section License
 * Copyright (C) 2014-2015, Mikael Patel
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * @section Description
 * Demonstration of the Cosa LED class.
 *
 * This file is part of the Arduino Che Cosa project.
 */

#include "Cosa/LED.hh"
#include "Cosa/Power.hh"
#include "Cosa/Event.hh"
#include "Cosa/Watchdog.hh"

Watchdog::Scheduler scheduler;
LED builtin(&scheduler);

void setup()
{
  Power::set(SLEEP_MODE_PWR_DOWN);
  builtin.alert_mode();
  Watchdog::begin();
}

void loop()
{
  Event::service();
}
