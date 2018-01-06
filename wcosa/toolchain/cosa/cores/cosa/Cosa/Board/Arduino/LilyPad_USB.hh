/**
 * @file Cosa/Board/Arduino/LilyPad_USB.hh
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

#ifndef COSA_BOARD_ARDUINO_LILYPAD_USB_HH
#define COSA_BOARD_ARDUINO_LILYPAD_USB_HH

/* This board is based on ATmega32U4 */
#define BOARD_ATMEGA32U4

/**
 * Compiler warning on unused varable.
 */
#if !defined(UNUSED)
#define UNUSED(x) (void) (x)
#endif

/**
 * Cosa pin symbol and hardware definitions for the ATmega32U4 based
 * SparkFun LilyPad USB board. Cosa does not use pin numbers as
 * Arduino/Wiring, instead strong data type is used (enum types) for
 * the specific pin classes; DigitalPin, AnalogPin, PWMPin, etc.
 *
 * The pin numbers for ATmega32u4 are mapped as in Arduino. The static
 * inline functions, SFR, BIT and UART, rely on compiler optimizations
 * to be reduced.
 *
 * @section Board
 * @code
 *             LilyPad Arduino USB
 *                    -----
 *                 +--|   |--+
 *             A2 /o  |   |  o\ D11
 *            A3 /o   -----   o\ D10
 *              +               +
 *           A4 |o             o| D9
 *              +               +
 *            A5 \o           o/ D3/SCL
 *            GND \o    o    o/ D2/SDA
 *                 +---------+
 *                     VCC
 * @endcode
 *
 * WARNING: The pin binding is not complete. This board does not have
 * all the defined pins. They remain to allow default parameters.
 */
class Board {
  friend class Pin;
  friend class GPIO;
  friend class UART;
private:
  /**
   * Do not allow instances. This is a static singleton; name space.
   */
  Board() {}

  /**
   * Return Special Function Register for given Arduino pin number.
   * @param[in] pin number.
   * @return special register pointer.
   */
  static volatile uint8_t* SFR(uint8_t pin)
    __attribute__((always_inline))
  {
    return (pin < 8  ? &PINB :
	    pin < 16 ? &PINC :
	    pin < 24 ? &PIND :
	    pin < 32 ? &PINE :
	               &PINF);
  }

  /**
   * Return bit position for given Arduino pin number in Special
   * Function Register.
   * @param[in] pin number.
   * @return pin bit position.
   */
  static uint8_t BIT(uint8_t pin)
    __attribute__((always_inline))
  {
    return (pin & 0x07);
  }

  /**
   * Return Pin Change Mask Register for given Arduino pin number.
   * @param[in] pin number.
   * @return pin change mask register pointer.
   */
  static volatile uint8_t* PCIMR(uint8_t pin)
    __attribute__((always_inline))
  {
    UNUSED(pin);
    return (&PCMSK0);
  }

  /**
   * Return UART Register for given Arduino serial port.
   * @param[in] port number.
   * @return UART register pointer.
   */
  static volatile uint8_t* UART(uint8_t port)
    __attribute__((always_inline))
  {
    UNUSED(port);
    return (&UCSR1A);
  }

public:
  /**
   * Initiate board ports. Default void.
   */
  static void init() {}

  /**
   * Digital pin symbols.
   */
  enum DigitalPin {
    D0 = 18,			// PD2/Not used
    D1 = 19,			// PD3/Not used
    D2 = 17,			// PD1
    D3 = 16,			// PD0
    D4 = 20,			// PD4/Not used
    D5 = 14,			// PC6/Not used
    D6 = 23,			// PD7/Not used
    D7 = 30,			// PE6/Not used
    D8 = 4,			// PB4/Not used
    D9 = 5,			// PB5
    D10 = 6,			// PB6
    D11 = 7,			// PB7
    D12 = 22,			// PD6/Not used
    D13 = 15,			// PC7/Not used
    D14 = 39,			// PF7/Not used
    D15 = 38,			// PF6/Not used
    D16 = 37,			// PF5
    D17 = 36,			// PF4
    D18 = 33,			// PF1
    D19 = 32,			// PF0
    D20 = 0,			// PB0/Not used
    D21 = 1,			// PB1/ISCP
    D22 = 2,			// PB2/ISCP
    D23 = 3,			// PB3/ISCP
    LED = 15,			// PC7
    TXLED = 0,			// Green
    RXLED = 21			// Yellow
  } __attribute__((packed));

  /**
   * Analog pin symbols (ADC channel numbers)
   */
  enum AnalogPin {
    A0 = 7,			// PF7/D14/Not used
    A1 = 6,			// PF6/D15/Not used
    A2 = 5,			// PF5/D16
    A3 = 4,			// PF4/D17
    A4 = 1,			// PF1/D18
    A5 = 0,			// PF0/D19
    A6 = 32,			// PD4/D4/Not used
    A7 = 33,			// PD6/D12/Not used
    A8 = 34,			// PD7/D6/Not used
    A9 = 35,			// PB4/D8/Not used
    A10 = 36,			// PB5/D9
    A11 = 37			// PB6/D10
  } __attribute__((packed));

  /**
   * Reference voltage; ARef pin, Vcc or internal 2V56
   */
  enum Reference {
    APIN_REFERENCE = 0,
    AVCC_REFERENCE = _BV(REFS0),
    A2V56_REFERENCE = (_BV(REFS1) | _BV(REFS0))
  } __attribute__((packed));

  /**
   * PWM pin symbols; sub-set of digital pins to allow compile
   * time checking
   */
  enum PWMPin {
    PWM0 = D11,			// PB7 => OCR0A
    PWM1 = D3,			// PD0 => OCR0B
    PWM2 = D9,			// PB5 => OCR1A
    PWM3 = D10,			// PB6 => OCR1B
    PWM4 = D5,			// PC6 => OCR3A/Not used
    PWM5 = D13,			// PC7 => OCR4A/LED
    PWM6 = D6			// PD7 => OCR4D/Not used
  } __attribute__((packed));

  /**
   * External interrupt pin symbols; sub-set of digital pins
   * to allow compile time checking.
   */
  enum ExternalInterruptPin {
    EXT0 = D3,			// PD0
    EXT1 = D2,			// PD1
    EXT2 = D0,			// PD2/Not used
    EXT3 = D1			// PD3/Not used
  } __attribute__((packed));

  /**
   * Pin change interrupt (PCI) pins.
   */
  enum InterruptPin {
    PCI0 = D20,			// PB0/Not used
    PCI1 = D21,			// PB1/Not used
    PCI2 = D22,			// PB2/Not used
    PCI3 = D23,			// PB3/Not used
    PCI4 = D8,			// PB4/Not used
    PCI5 = D9,			// PB5
    PCI6 = D10,			// PB6
    PCI7 = D11			// PB7
  } __attribute__((packed));

  /**
   * Size of pin maps.
   */
  enum {
    ANALOG_PIN_MAX = 12,
    DIGITAL_PIN_MAX = 24,
    EXT_PIN_MAX = 4,
    PCI_PIN_MAX = 8,
    PWM_PIN_MAX = 7
  };

  /**
   * Pins used for TWI interface (port D, bit 0-1, D2-D3)
   */
  enum TWIPin {
    SDA = 1,			// PD1/D2
    SCL = 0			// PD0/D3
  } __attribute__((packed));

  /**
   * Pins used for SPI interface (port B, bit 0-3, ICSP)
   */
  enum SPIPin {
    SS = 0,			// PB0/D20/Not used
    SCK = 1,			// PB1/D21/ICSP
    MOSI = 2,			// PB2/D22/ICSP
    MISO = 3			// PB3/D23/ICSP
  } __attribute__((packed));

  /**
   * Auxiliary
   */
  enum {
    VBG = (_BV(MUX4) | _BV(MUX3) | _BV(MUX2) | _BV(MUX1)),
    UART_MAX = 2,
    EXT_MAX = 7,
    PCMSK_MAX = 1,
    PCINT_MAX = 8
  } __attribute__((packed));
};

/**
 * Redefined symbols to allow generic code.
 */
#define UCSZ00 UCSZ10
#define UCSZ01 UCSZ11
#define UCSZ02 UCSZ12
#define UPM00 UPM10
#define UPM01 UPM11
#define USBS0 USBS1
#define U2X0 U2X1
#define TXC0 TXC1
#define RXCIE0 RXCIE1
#define RXEN0 RXEN1
#define TXEN0 TXEN1
#define UDRE0 UDRE1
#define UDRIE0 UDRIE1
#define TXCIE0 TXCIE1
#if !defined(ADCW)
#define ADCW ADC
#endif

/**
 * Forward declare interrupt service routines to allow them as friends.
 */
extern "C" {
  void ADC_vect(void) __attribute__ ((signal));
  void ANALOG_COMP_vect(void) __attribute__ ((signal));
  void INT0_vect(void) __attribute__ ((signal));
  void INT1_vect(void) __attribute__ ((signal));
  void INT2_vect(void) __attribute__ ((signal));
  void INT3_vect(void) __attribute__ ((signal));
  void INT6_vect(void) __attribute__ ((signal));
  void PCINT0_vect(void) __attribute__ ((signal));
  void SPI_STC_vect(void) __attribute__ ((signal));
  void TIMER0_COMPA_vect(void) __attribute__ ((signal));
  void TIMER0_COMPB_vect(void) __attribute__ ((signal));
  void TIMER0_OVF_vect(void) __attribute__ ((signal));
  void TIMER1_CAPT_vect(void)  __attribute__ ((signal));
  void TIMER1_COMPA_vect(void) __attribute__ ((signal));
  void TIMER1_COMPB_vect(void) __attribute__ ((signal));
  void TIMER1_COMPC_vect(void) __attribute__ ((signal));
  void TIMER1_OVF_vect(void) __attribute__ ((signal));
  void TIMER3_CAPT_vect(void)  __attribute__ ((signal));
  void TIMER3_COMPA_vect(void) __attribute__ ((signal));
  void TIMER3_COMPB_vect(void) __attribute__ ((signal));
  void TIMER3_COMPC_vect(void) __attribute__ ((signal));
  void TIMER3_OVF_vect(void) __attribute__ ((signal));
  void TIMER4_COMPA_vect(void) __attribute__ ((signal));
  void TIMER4_COMPB_vect(void) __attribute__ ((signal));
  void TIMER4_COMPD_vect(void) __attribute__ ((signal));
  void TIMER4_FPF_vect(void) __attribute__ ((signal));
  void TIMER4_OVF_vect(void) __attribute__ ((signal));
  void TWI_vect(void) __attribute__ ((signal));
  void WDT_vect(void) __attribute__ ((signal));
  void USART1_RX_vect(void) __attribute__ ((signal));
  void USART1_TX_vect(void) __attribute__ ((signal));
  void USART1_UDRE_vect(void) __attribute__ ((signal));
  void USB_COM_vect(void) __attribute__ ((signal));
  void USB_GEN_vect(void) __attribute__ ((signal));
}
#endif
