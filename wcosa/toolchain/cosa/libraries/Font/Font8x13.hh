/**
 * @file Font8x13.hh
 * @version 1.0
 *
 * @author contributed by jeditekunum
 *
 * @section License
 * Copyright (C) 2015, Mikael Patel
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
 * This file is part of the Arduino Che Cosa project.
 */

#ifndef COSA_FONT8X13_HH
#define COSA_FONT8X13_HH

#include "Font.hh"

/**
 * Font size 8x13.
 */
class Font8x13 : public Font {
public:
  /**
   * Construct font (8x13) singleton.
   */
  Font8x13() :
    Font(width, height, first, last, bitmap, compression_type)
  {}

private:
  static const uint8_t bitmap[] PROGMEM;
  static const uint8_t compression_type;
  static const uint8_t width;
  static const uint8_t height;
  static const uint8_t first;
  static const uint8_t last;
} font8x13;

#include "Data/8x13.h"

#endif
