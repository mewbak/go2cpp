// SPDX-License-Identifier: Apache-2.0

package gowasm2cpp

import (
	"os"
	"path/filepath"
	"text/template"
)

type wasmData struct {
	Offset int
	Data   []byte
}

func writeMem(dir string, incpath string, namespace string, initPageNum int, data []wasmData) error {
	const pageSize = 64 * 1024

	{
		f, err := os.Create(filepath.Join(dir, "mem.h"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := memHTmpl.Execute(f, struct {
			IncludeGuard string
			IncludePath  string
			Namespace    string
			PageSize     int
		}{
			IncludeGuard: includeGuard(namespace) + "_MEM_H",
			IncludePath:  incpath,
			Namespace:    namespace,
			PageSize:     pageSize,
		}); err != nil {
			return err
		}
	}
	{
		f, err := os.Create(filepath.Join(dir, "mem.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		var flatten []byte
		for _, d := range data {
			flatten = append(flatten, d.Data...)
		}

		if err := memCppTmpl.Execute(f, struct {
			IncludePath string
			Namespace   string
			InitPageNum int
			Data        []wasmData
			FlattenData []byte
		}{
			IncludePath: incpath,
			Namespace:   namespace,
			InitPageNum: initPageNum,
			Data:        data,
			FlattenData: flatten,
		}); err != nil {
			return err
		}
	}
	return nil
}

var memHTmpl = template.Must(template.New("mem.h").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#ifndef {{.IncludeGuard}}
#define {{.IncludeGuard}}

#include "{{.IncludePath}}bytes.h"

#include <cstdint>
#include <string>
#include <vector>

namespace {{.Namespace}} {

class Mem {
public:
  static constexpr int32_t kPageSize = {{.PageSize}};

  Mem();
  ~Mem();

  int32_t GetSize() const;
  int32_t Grow(int32_t delta);

  inline int8_t LoadInt8(int32_t addr) const {
    return static_cast<int8_t>(*(bytes_ + addr));
  }

  inline uint8_t LoadUint8(int32_t addr) const {
    return *(bytes_ + addr);
  }

  inline int16_t LoadInt16(int32_t addr) const {
    return *(reinterpret_cast<const int16_t*>(bytes_ + addr));
  }

  inline uint16_t LoadUint16(int32_t addr) const {
    return *(reinterpret_cast<const uint16_t*>(bytes_ + addr));
  }

  inline int32_t LoadInt32(int32_t addr) const {
    return *(reinterpret_cast<const int32_t*>(bytes_ + addr));
  }

  inline uint32_t LoadUint32(int32_t addr) const {
    return *(reinterpret_cast<const uint32_t*>(bytes_ + addr));
  }

  inline int64_t LoadInt64(int32_t addr) const {
    return *(reinterpret_cast<const int64_t*>(bytes_ + addr));
  }

  inline float LoadFloat32(int32_t addr) const {
    return *(reinterpret_cast<const float*>(bytes_ + addr));
  }

  inline double LoadFloat64(int32_t addr) const {
    return *(reinterpret_cast<const double*>(bytes_ + addr));
  }

  void StoreInt8(int32_t addr, int8_t val) {
    *(bytes_ + addr) = static_cast<uint8_t>(val);
  }

  inline void StoreInt16(int32_t addr, int16_t val) {
    *(reinterpret_cast<int16_t*>(bytes_ + addr)) = val;
  }

  inline void StoreInt32(int32_t addr, int32_t val) {
    *(reinterpret_cast<int32_t*>(bytes_ + addr)) = val;
  }

  inline void StoreInt64(int32_t addr, int64_t val) {
    *(reinterpret_cast<int64_t*>(bytes_ + addr)) = val;
  }

  inline void StoreFloat32(int32_t addr, float val) {
    *(reinterpret_cast<float*>(bytes_ + addr)) = val;
  }

  inline void StoreFloat64(int32_t addr, double val) {
    *(reinterpret_cast<double*>(bytes_ + addr)) = val;
  }

  void StoreBytes(int32_t addr, const std::vector<uint8_t>& bytes);

  BytesSpan LoadSlice(int32_t addr);
  BytesSpan LoadSliceDirectly(int64_t array, int32_t len);
  std::string LoadString(int32_t addr) const;

  int Memcmp(int32_t a, int32_t b, int32_t len);
  int32_t Memchr(int32_t ptr, int32_t ch, int32_t count);
  void Memmove(int32_t dst, int32_t src, int32_t count);
  void Memset(int32_t dst, uint8_t ch, int32_t count);

private:
  Mem(const Mem&) = delete;
  Mem& operator=(const Mem&) = delete;

  uint8_t* bytes_;
  size_t size_ = 0;
};

}

#endif  // {{.IncludeGuard}}
`))

var memCppTmpl = template.Must(template.New("mem.cpp").Funcs(template.FuncMap{
	"needsNewLine": func(x int) bool {
		return (x+1)%16 == 0
	},
}).Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}mem.h"

#include <algorithm>
#include <cstring>

namespace {{.Namespace}} {

namespace {

// 2GB. 4GB seems too big on some machines.
constexpr size_t kMaxMemorySize = 2ull * 1024ull * 1024ull * 1024ull;

const uint8_t initial_data_[] = {
  {{range $index, $value := .FlattenData}}{{$value}}, {{if needsNewLine $index}}
  {{end}}{{end}}
};

struct WasmData {
  int32_t offset;
  int32_t length;
};

const WasmData initial_data_info_[] = {
  {{range $index, $value := .Data}}{ {{$value.Offset}}, {{len $value.Data}} }, {{if needsNewLine $index}}
  {{end}}{{end}}
};

}

Mem::Mem()
    : size_({{.InitPageNum}} * kPageSize) {
  bytes_ = reinterpret_cast<uint8_t*>(std::calloc(1, kMaxMemorySize));
  constexpr int32_t info_size = sizeof(initial_data_info_) / sizeof(initial_data_info_[0]);
  int32_t src_offset = 0;
  for (int32_t i = 0; i < info_size; i++) {
    WasmData info = initial_data_info_[i];
    std::memcpy(bytes_ + info.offset, initial_data_ + src_offset, info.length);
    src_offset += info.length;
  }
}

Mem::~Mem() {
  std::free(bytes_);
}

int32_t Mem::GetSize() const {
  return size_ / kPageSize;
}

int32_t Mem::Grow(int32_t delta) {
  int prev_page_num = GetSize();
  size_ = std::min(static_cast<size_t>((prev_page_num + delta) * kPageSize), kMaxMemorySize);
  return prev_page_num;
}

void Mem::StoreBytes(int32_t addr, const std::vector<uint8_t>& src) {
  std::memcpy(bytes_ + addr, &(*src.begin()), src.size());
}

BytesSpan Mem::LoadSlice(int32_t addr) {
  int64_t array = LoadInt64(addr);
  int64_t len = LoadInt64(addr + 8);
  return BytesSpan{&*(bytes_ + array), static_cast<BytesSpan::size_type>(len)};
}

BytesSpan Mem::LoadSliceDirectly(int64_t array, int32_t len) {
  return BytesSpan{&*(bytes_ + array), static_cast<BytesSpan::size_type>(len)};
}

std::string Mem::LoadString(int32_t addr) const {
  int64_t saddr = LoadInt64(addr);
  int64_t len = LoadInt64(addr + 8);
  return std::string{bytes_ + saddr, bytes_ + saddr + len};
}

int Mem::Memcmp(int32_t a, int32_t b, int32_t len) {
  return std::memcmp(bytes_ + a, bytes_ + b, len);
}

int32_t Mem::Memchr(int32_t ptr, int32_t ch, int32_t count) {
  void* result = std::memchr(bytes_ + ptr, ch, count);
  if (!result) {
    return 0;
  }
  return static_cast<int32_t>(reinterpret_cast<uint8_t*>(result) - bytes_);
}

void Mem::Memmove(int32_t dst, int32_t src, int32_t count) {
  std::memmove(bytes_ + dst, bytes_ + src, count);
}

void Mem::Memset(int32_t dst, uint8_t ch, int32_t count) {
  std::memset(bytes_ + dst, ch, count);
}

}
`))
