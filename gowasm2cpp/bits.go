// SPDX-License-Identifier: Apache-2.0

package gowasm2cpp

import (
	"os"
	"path/filepath"
	"text/template"
)

func writeBits(dir string, namespace string) error {
	{
		f, err := os.Create(filepath.Join(dir, "bits.h"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := bitsHTmpl.Execute(f, struct {
			IncludeGuard string
			Namespace    string
		}{
			IncludeGuard: includeGuard(namespace) + "_BITS_H",
			Namespace:    namespace,
		}); err != nil {
			return err
		}
	}
	{
		f, err := os.Create(filepath.Join(dir, "bits.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := bitsCppTmpl.Execute(f, struct {
			Namespace string
		}{
			Namespace: namespace,
		}); err != nil {
			return err
		}
	}
	return nil
}

var bitsHTmpl = template.Must(template.New("bits.h").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#ifndef {{.IncludeGuard}}
#define {{.IncludeGuard}}

#include <cstdint>

namespace {{.Namespace}} {

class Bits {
public:
  static int32_t LeadingZeros(uint32_t x);
  static int32_t LeadingZeros(uint64_t x);
  static int32_t TailingZeros(uint32_t x);
  static int32_t TailingZeros(uint64_t x);
  static int32_t OnesCount(uint32_t x);
  static int32_t OnesCount(uint64_t x);
  static uint32_t RotateLeft(uint32_t x, int32_t k);
  static uint64_t RotateLeft(uint64_t x, int32_t k);

private:
  static int32_t Len(uint32_t x);
  static int32_t Len(uint64_t x);
};

}

#endif  // {{.IncludeGuard}}
`))

var bitsCppTmpl = template.Must(template.New("bits.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "autogen/bits.h"

namespace {

static uint8_t pop8tab[] = {
  0x00, 0x01, 0x01, 0x02, 0x01, 0x02, 0x02, 0x03, 0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04,
  0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04, 0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
  0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04, 0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04, 0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06, 0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
  0x01, 0x02, 0x02, 0x03, 0x02, 0x03, 0x03, 0x04, 0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06, 0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
  0x02, 0x03, 0x03, 0x04, 0x03, 0x04, 0x04, 0x05, 0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06,
  0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06, 0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
  0x03, 0x04, 0x04, 0x05, 0x04, 0x05, 0x05, 0x06, 0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07,
  0x04, 0x05, 0x05, 0x06, 0x05, 0x06, 0x06, 0x07, 0x05, 0x06, 0x06, 0x07, 0x06, 0x07, 0x07, 0x08,
};

static uint8_t len8tab[] = {
  0x00, 0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
  0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
  0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
  0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
  0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
  0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
  0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
  0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
  0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
};

const uint32_t deBruijn32 = 0x077CB531;

static uint8_t deBruijn32tab[] = {
  0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
  31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9,
};

const uint64_t deBruijn64 = 0x03f79d71b4ca8b09;

static uint8_t deBruijn64tab[] = {
  0, 1, 56, 2, 57, 49, 28, 3, 61, 58, 42, 50, 38, 29, 17, 4,
  62, 47, 59, 36, 45, 43, 51, 22, 53, 39, 33, 30, 24, 18, 12, 5,
  63, 55, 48, 27, 60, 41, 37, 16, 46, 35, 44, 21, 52, 32, 23, 11,
  54, 26, 40, 15, 34, 20, 31, 10, 25, 14, 19, 9, 13, 8, 7, 6,
};

}

namespace {{.Namespace}} {

// The implementation is copied from the Go standard package math/bits, which is under BSD-style license.

int32_t Bits::LeadingZeros(uint32_t x) {
  return 32 - Len(x);
}

int32_t Bits::LeadingZeros(uint64_t x) {
  return 64 - Len(x);
}

int32_t Bits::TailingZeros(uint32_t x) {
  if (x == 0) {
    return 32;
  }
  return (int32_t)deBruijn32tab[(x&-x)*deBruijn32>>(32-5)];
}

int32_t Bits::TailingZeros(uint64_t x) {
  if (x == 0) {
    return 64;
  }
  return (int32_t)deBruijn64tab[(x&(uint64_t)(-(int64_t)x))*deBruijn64>>(64-6)];
}

int32_t Bits::OnesCount(uint32_t x) {
  return (int32_t)(pop8tab[x>>24] + pop8tab[x>>16&0xff] + pop8tab[x>>8&0xff] + pop8tab[x&0xff]);
}

int32_t Bits::OnesCount(uint64_t x) {
  const uint64_t m0 = 0x5555555555555555ul;
  const uint64_t m1 = 0x3333333333333333ul;
  const uint64_t m2 = 0x0f0f0f0f0f0f0f0ful;
  const uint64_t m  = 0xfffffffffffffffful;

  x = ((x>>1)&(m0&m)) + (x&(m0&m));
  x = ((x>>2)&(m1&m)) + (x&(m1&m));
  x = ((x>>4) + x) & (m2 & m);
  x += x >> 8;
  x += x >> 16;
  x += x >> 32;
  return (int32_t)(x) & ((1<<7) - 1);
}

uint32_t Bits::RotateLeft(uint32_t x, int32_t k) {
  const int32_t n = 32;
  int32_t s = k & (n - 1);
  return x<<s | x>>(n-s);
}

uint64_t Bits::RotateLeft(uint64_t x, int32_t k) {
  const int32_t n = 64;
  int32_t s = k & (n - 1);
  return x<<s | x>>(n-s);
}

int32_t Bits::Len(uint32_t x) {
  int32_t n = 0;
  if (x >= 1<<16) {
    x >>= 16;
    n = 16;
  }
  if (x >= 1<<8) {
    x >>= 8;
    n += 8;
  }
  return n + (int32_t)len8tab[x];
}

int32_t Bits::Len(uint64_t x) {
  int32_t n = 0;
  if (x >= 1ul<<32) {
    x >>= 32;
    n = 32;
  }
  if (x >= 1ul<<16) {
    x >>= 16;
    n += 16;
  }
  if (x >= 1ul<<8) {
    x >>= 8;
    n += 8;
  }
  return n + (int32_t)len8tab[x];
}

}
`))
