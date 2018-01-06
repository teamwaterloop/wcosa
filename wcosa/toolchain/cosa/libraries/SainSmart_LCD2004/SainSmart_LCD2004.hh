/**
 * @file SainSmart_LCD2004.hh
 * @version 1.0
 *
 * @section License
 * Copyright (C) 2013-2015, Mikael Patel
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

#ifndef COSA_SAINSMART_LCD2004_HH
#define COSA_SAINSMART_LCD2004_HH

#include <HD44780.h>
#include <PCF8574.h>

/**
 * IO handler for HD44780 (LCD-II) Dot Matix Liquid Crystal Display
 * Controller/Driver when using the SainSmart IO expander board based
 * on PCF8574A I2C IO expander device driver.
 *
 * @section Circuit
 * @code
 *                       PCF8574A/SainSmart LCD2004
 *                       +------U-----+
 * (GND)---[ ]---------1-|A0       VCC|-16--------------(VCC)
 * (GND)---[ ]---------2-|A1       SDA|-15-----------(SDA/A4)
 * (GND)---[ ]---------3-|A2       SCL|-14-----------(SCL/A5)
 * (LCD/RS)------------4-|P0       INT|-13
 * (LCD/RW)------------5-|P1        P7|-12-----------(LCD/D7)
 * (LCD/EN)------------6-|P2        P6|-11-----------(LCD/D6)
 * (LCD/BT)------------7-|P3        P5|-10-----------(LCD/D5)
 * (GND)---------------8-|GND       P4|-9------------(LCD/D4)
 *                       +------------+
 * @endcode
 */
class SainSmart_LCD2004 : public HD44780::IO, private PCF8574A {
public:
  /**
   * Construct HD44780 IO port handler using the SainSmart LCD2004
   * I2C/TWI I/O expander with given sub-address (A0..A2).
   * @param[in] subaddr sub-address (0..7, default 7).
   */
  SainSmart_LCD2004(uint8_t subaddr = 7) : PCF8574A(subaddr), m_port() {}

  /**
   * @override{HD44780::IO}
   * Initiate TWI interface. Returns false.
   * @return bool.
   */
  virtual bool setup();

  /**
   * @override{HD44780::IO}
   * Write nibble to display using TWI interface.
   * @param[in] data (4b) to write.
   */
  virtual void write4b(uint8_t data);

  /**
   * @override{HD44780::IO}
   * Write byte (8bit) to display.
   * @param[in] data (8b) to write.
   */
  virtual void write8b(uint8_t data);

  /**
   * @override{HD44780::IO}
   * Write character buffer to display.
   * @param[in] buf pointer to buffer.
   * @param[in] size number of bytes in buffer.
   */
  virtual void write8n(const void* buf, size_t size);

  /**
   * @override{HD44780::IO}
   * Set instruction/data mode; zero for instruction,
   * non-zero for data mode.
   * @param[in] flag.
   */
  virtual void set_mode(uint8_t flag);

  /**
   * @override{HD44780::IO}
   * Set backlight on/off.
   * @param[in] flag.
   */
  virtual void set_backlight(uint8_t flag);

protected:
  /** Max size of temporary buffer for TWI message. */
  static const uint8_t TMP_MAX = 32;

  /** Expander port bit fields; little endian. */
  union port_t {
    uint8_t as_uint8;		//!< Unsigned byte access.
    struct {
      uint8_t rs:1;		//!< Command/Data select (P0).
      uint8_t rw:1;		//!< Read/Write (P1).
      uint8_t en:1;		//!< Enable pulse (P2).
      uint8_t bt:1;		//!< Back-light (P3).
      uint8_t data:4;		//!< Data port (P4..P7).
    };
    operator uint8_t()
    {
      return (as_uint8);
    }
    port_t()
    {
      as_uint8 = 0;
    }
  };
  port_t m_port;		//!< Port setting.
};

#endif
