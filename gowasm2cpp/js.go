// SPDX-License-Identifier: Apache-2.0

package gowasm2cpp

import (
	"os"
	"path/filepath"
	"text/template"
)

func writeJS(dir string, incpath string, namespace string) error {
	{
		f, err := os.Create(filepath.Join(dir, "js.h"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := jsHTmpl.Execute(f, struct {
			IncludeGuard string
			IncludePath  string
			Namespace    string
		}{
			IncludeGuard: includeGuard(namespace) + "_JS_H",
			IncludePath:  incpath,
			Namespace:    namespace,
		}); err != nil {
			return err
		}
	}
	{
		f, err := os.Create(filepath.Join(dir, "js.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := jsCppTmpl.Execute(f, struct {
			IncludePath string
			Namespace   string
		}{
			IncludePath: incpath,
			Namespace:   namespace,
		}); err != nil {
			return err
		}
	}
	return nil
}

var jsHTmpl = template.Must(template.New("js.h").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#ifndef {{.IncludeGuard}}
#define {{.IncludeGuard}}

#include <deque>
#include <functional>
#include <iostream>
#include <map>
#include <memory>
#include <string>
#include <vector>
#include "{{.IncludePath}}bytes.h"

namespace {{.Namespace}} {

class Object;

class Writer {
public:
  explicit Writer(std::ostream& out);
  void Write(BytesSpan bytes);

private:
  std::ostream& out_;
  // TODO: std::queue should be enough?
  std::deque<uint8_t> buf_;
};

class ArrayBuffer;

class Value {
public:
  enum class Type {
    Undefined,
    Null,
    Bool,
    Number,
    String,
    Object,
  };

  class Hash {
  public:
    std::size_t operator()(const Value& value) const;
  };

  static Value Null();
  static Value Global();
  static Value ReflectGet(Value target, const std::string& key);
  static void ReflectSet(Value target, const std::string& key, Value value);
  static void ReflectDelete(Value target, const std::string& key);
  static Value ReflectConstruct(Value target, std::vector<Value> args);
  static Value ReflectApply(Value target, Value self, std::vector<Value> args);

  Value();
  explicit Value(bool b);
  explicit Value(double num);
  explicit Value(const char* str);
  explicit Value(const std::string& str);
  explicit Value(std::shared_ptr<Object> object);
  explicit Value(const std::vector<Value>& array);

  Value(const Value& rhs);
  Value& operator=(const Value& rhs);
  bool operator==(const Value& rhs) const;

  bool IsNull() const;
  bool IsUndefined() const;
  bool IsBool() const;
  bool IsNumber() const;
  bool IsString() const;
  bool IsBytes() const;
  bool IsObject() const;
  bool IsArray() const;

  bool ToBool() const;
  double ToNumber() const;
  std::string ToString() const;
  BytesSpan ToBytes();
  Object& ToObject();
  const Object& ToObject() const;
  std::vector<Value>& ToArray();
  std::shared_ptr<ArrayBuffer> ToArrayBuffer();

  std::string Inspect() const;

private:
  static Value MakeGlobal();

  explicit Value(Type type);
  Value(Type type, double num);

  Type type_ = Type::Undefined;
  double num_value_ = 0;
  std::string str_value_;
  std::shared_ptr<Object> object_value_;
  std::shared_ptr<std::vector<Value>> array_value_;
};

class Object {
public:
  using Func = std::function<Value (Value, std::vector<Value>)>;

  virtual ~Object();
  virtual Value Get(const std::string& key);
  virtual void Set(const std::string& key, Value value);
  virtual void Delete(const std::string& key);

  virtual bool IsFunction() const { return false; }
  virtual bool IsConstructor() const { return false; }
  virtual Value Invoke(Value self, std::vector<Value> args);
  virtual Value New(std::vector<Value> args);

  virtual BytesSpan ToBytes();

  virtual std::string ToString() const = 0;
  virtual std::string Inspect() const;
};

class ArrayBuffer : public Object {
public:
  explicit ArrayBuffer(size_t size);

  size_t ByteLength() const;
  Value Get(const std::string& key) override;
  BytesSpan ToBytes() override;
  std::string ToString() const override;

private:
  std::vector<uint8_t> data_;
};

class DictionaryValues : public Object {
public:
  DictionaryValues();
  explicit DictionaryValues(const std::map<std::string, Value>& dict);
  Value Get(const std::string& key) override;
  void Set(const std::string& key, Value value) override;
  void Delete(const std::string& key) override;
  std::string ToString() const override;
  std::string Inspect() const override;

private:
  std::map<std::string, Value> dict_;
};

class Function : public Object {
public:
  explicit Function(Object::Func fn);

  bool IsFunction() const override { return true; }
  bool IsConstructor() const override { return false; }
  Value Invoke(Value self, std::vector<Value> args) override;
  std::string ToString() const override { return "(function)"; }

private:
  Object::Func fn_;
};

class Constructor : public Object {
public:
  Constructor(const std::string& name, Object::Func fn);

  bool IsFunction() const override { return true; }
  bool IsConstructor() const override { return true; }
  Value New(std::vector<Value> args) override;
  std::string ToString() const override;

private:
  std::string name_;
  Object::Func fn_;
};

}

#endif  // {{.IncludeGuard}}
`))

var jsCppTmpl = template.Must(template.New("js.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}js.h"

#include <algorithm>
#include <cassert>
#include <cstdlib>
#include <random>
#include <tuple>

namespace {{.Namespace}} {

namespace {

void panic(const std::string& msg) {
  // TODO: Can we call a Go function without registering _panic?
  auto handler = Value::Global().ToObject().Get("_panic");
  if (handler.IsUndefined()) {
    std::cerr << msg << std::endl;
    assert(false);
    std::exit(1);
    return;
  }
  handler.ToObject().Invoke(Value{}, std::vector<Value>{Value{msg}});
}

std::string JoinObjects(const std::vector<Value>& objs) {
  std::string str;
  for (int i = 0; i < objs.size(); i++) {
    str += objs[i].Inspect();
    if (i < objs.size() - 1) {
      str += ", ";
    }
  }
  return str;
}

void WriteObjects(std::ostream& out, const std::vector<Value>& objs) {
  std::vector<std::string> inspects(objs.size());
  for (int i = 0; i < objs.size(); i++) {
    out << objs[i].Inspect();
    if (i < objs.size() - 1) {
      out << ", ";
    }
  }
  out << std::endl;
}

class TypedArray : public Object {
public:
  explicit TypedArray(size_t size)
      : array_buffer_{std::make_shared<ArrayBuffer>(size)},
        length_{size} {
  }

  TypedArray(std::shared_ptr<ArrayBuffer> arrayBuffer, size_t offset, size_t length)
      : array_buffer_{arrayBuffer},
        offset_{offset},
        length_{length} {
  }

  Value Get(const std::string& key) override {
    if (key == "byteLength") {
      return Value{static_cast<double>(length_)};
    }
    if (key == "byteOffset") {
      return Value{static_cast<double>(offset_)};
    }
    if (key == "buffer") {
      return Value{array_buffer_};
    }
    return Value{};
  }

  void Set(const std::string& key, Value value) override {
    if (key == "byteLength") {
      length_ = static_cast<size_t>(value.ToNumber());
      return;
    }
    if (key == "byteOffset") {
      offset_ = static_cast<size_t>(value.ToNumber());
      return;
    }
    if (key == "buffer") {
      array_buffer_ = value.ToArrayBuffer();
      return;
    }
    panic("TypedArray::Set: invalid key: " + key);
  }

  BytesSpan ToBytes() override {
    auto bs = array_buffer_->ToBytes();
    return BytesSpan{bs.begin() + offset_, length_};
  }

  std::string ToString() const override {
    return "TypedArray";
  }

private:
  std::shared_ptr<ArrayBuffer> array_buffer_;
  size_t offset_ = 0;
  size_t length_ = 0;
};

class Uint8Array : public TypedArray {
public:
  explicit Uint8Array(size_t size)
      : TypedArray(size) {
  }

  Uint8Array(std::shared_ptr<ArrayBuffer> arrayBuffer, size_t offset, size_t length)
      : TypedArray(arrayBuffer, offset, length) {
  }

  std::string ToString() const override {
    return "Uint8Array";
  }
};

class Uint16Array : public TypedArray {
public:
  Uint16Array(std::shared_ptr<ArrayBuffer> arrayBuffer, size_t offset, size_t length)
      : TypedArray(arrayBuffer, offset*2, length*2) {
  }

  std::string ToString() const override {
    return "Uint16Array";
  }
};

class Float32Array : public TypedArray {
public:
  explicit Float32Array(size_t size)
      : TypedArray(size*4) {
  }

  Float32Array(std::shared_ptr<ArrayBuffer> arrayBuffer, size_t offset, size_t length)
      : TypedArray(arrayBuffer, offset*4, length*4) {
  }

  std::string ToString() const override {
    return "Float32Array";
  }
};

class Enosys : public Object {
public:
  explicit Enosys(const std::string& name)
    : name_(name) {
  }

  Value Get(const std::string& key) override {
    if (key == "message") {
      return Value{name_ + " not implemented"};
    }
    if (key == "code") {
      return Value{"ENOSYS"};
    }
    return Value{};
  }

  std::string ToString() const override {
    return "ENOSYS: " + name_;
  }

private:
  std::string name_;
};

class FS {
public:
  FS()
      : stdout_{std::cout},
        stderr_{std::cerr} {
  }

  Value Write(Value self, std::vector<Value> args) {
    int fd = (int)(args[0].ToNumber());
    BytesSpan buf = args[1].ToBytes();
    int offset = (int)(args[2].ToNumber());
    int length = (int)(args[3].ToNumber());
    Value position = args[4];
    Value callback = args[5];
    if (offset != 0 || length != buf.size()) {
      Value::ReflectApply(callback, Value{}, std::vector<Value>{ Value{std::make_shared<Enosys>("write")} });
      return Value{};
    }
    if (!position.IsNull()) {
      Value::ReflectApply(callback, Value{}, std::vector<Value>{ Value{std::make_shared<Enosys>("write")} });
      return Value{};
    }
    switch (fd) {
    case 1:
      stdout_.Write(buf);
      break;
    case 2:
      stderr_.Write(buf);
      break;
    default:
      Value::ReflectApply(callback, Value{}, std::vector<Value>{ Value{std::make_shared<Enosys>("write")} });
      break;
    }
    // The first argument must be null or an error. Undefined doesn't work.
    Value::ReflectApply(callback, Value{}, std::vector<Value>{ Value::Null(), Value{static_cast<double>(buf.size())} });
    return Value{};
  }

private:
  Writer stdout_;
  Writer stderr_;
};

}  // namespace

Writer::Writer(std::ostream& out)
    : out_{out} {
}

void Writer::Write(BytesSpan bytes) {
  buf_.insert(buf_.end(), bytes.begin(), bytes.end());
  for (;;) {
    auto it = std::find(buf_.begin(), buf_.end(), '\n');
    if (it == buf_.end()) {
      break;
    }
    std::string str(buf_.begin(), it);
    out_ << str << std::endl;
    ++it;
    buf_.erase(buf_.begin(), it);
  }
}

std::size_t Value::Hash::operator()(const Value& value) const {
  size_t h = 17;
  h = h * 31 + std::hash<decltype(value.type_)>()(value.type_);
  h = h * 31 + std::hash<decltype(value.num_value_)>()(value.num_value_);
  h = h * 31 + std::hash<decltype(value.str_value_)>()(value.str_value_);
  h = h * 31 + std::hash<decltype(value.object_value_)>()(value.object_value_);
  h = h * 31 + std::hash<decltype(value.array_value_)>()(value.array_value_);
  return h;
}

Value Value::Null() {
  return Value{Type::Null};
}

Value::Value() = default;

Value::Value(bool b)
    : type_{Type::Bool},
      num_value_{static_cast<double>(b)}{
}

Value::Value(double num)
    : type_{Type::Number},
      num_value_{num} {
}

Value::Value(const char* str)
    : Value{std::string(str)} {
}

Value::Value(const std::string& str)
    : type_{Type::String},
      str_value_{str} {
}

Value::Value(std::shared_ptr<Object> object)
    : type_{Type::Object},
      object_value_{object} {
}

Value::Value(const std::vector<Value>& array)
    : type_{Type::Object},
      array_value_{std::make_shared<std::vector<Value>>(array.begin(), array.end())} {
}

Value::Value(const Value& rhs) = default;

Value& Value::operator=(const Value& rhs) = default;

bool Value::operator==(const Value& rhs) const {
  return type_ == rhs.type_ &&
      num_value_ == rhs.num_value_ &&
      str_value_ == rhs.str_value_ &&
      object_value_ == rhs.object_value_ &&
      array_value_ == rhs.array_value_;
}

Value::Value(Type type)
    : type_{type} {
}

Value::Value(Type type, double num)
    : type_{type},
      num_value_{num} {
}

bool Value::IsNull() const {
  return type_ == Type::Null;
}

bool Value::IsUndefined() const {
  return type_ == Type::Undefined;
}

bool Value::IsBool() const {
  return type_ == Type::Bool;
}

bool Value::IsNumber() const {
  return type_ == Type::Number;
}

bool Value::IsString() const {
  return type_ == Type::String;
}

bool Value::IsBytes() const {
  return type_ == Type::Object && !object_value_->ToBytes().IsNull();
}

bool Value::IsObject() const {
  return type_ == Type::Object && !!object_value_;
}

bool Value::IsArray() const {
  return type_ == Type::Object && !!array_value_;
}

bool Value::ToBool() const {
  if (type_ != Type::Bool) {
    panic("Value::ToBool: the type must be Type::Bool but not: " + Inspect());
  }
  return static_cast<bool>(num_value_);
}

double Value::ToNumber() const {
  if (type_ != Type::Number) {
    panic("Value::ToNumber: the type must be Type::Number but not: " + Inspect());
  }
  return num_value_;
}

std::string Value::ToString() const {
  if (type_ != Type::String) {
    panic("Value::ToString: the type must be Type::String but not: " + Inspect());
  }
  return str_value_;
}

BytesSpan Value::ToBytes() {
  if (type_ != Type::Object) {
    panic("Value::ToBytes: the type must be Type::Object but not: " + Inspect());
  }
  BytesSpan bytes = object_value_->ToBytes();
  if (bytes.IsNull()) {
    panic("Value::ToBytes: object_value_->ToBytes() must not be null");
  }
  return bytes;
}

Object& Value::ToObject() {
  if (type_ != Type::Object) {
    panic("Value::ToObject: the type must be Type::Object but not: " + Inspect());
  }
  if (!object_value_) {
    panic("Value::ToObject: object_value_ must not be null");
  }
  return *object_value_;
}

const Object& Value::ToObject() const {
  if (type_ != Type::Object) {
    panic("Value::ToObject: the type must be Type::Object but not: " + Inspect());
  }
  if (!object_value_) {
    panic("Value::ToObject: object_value_ must not be null");
  }
  return *object_value_;
}

std::vector<Value>& Value::ToArray() {
  if (type_ != Type::Object) {
    panic("Value::ToArray: the type must be Type::Object but not: " + Inspect());
  }
  if (!array_value_) {
    panic("Value::ToArray: array_value_ must not be null");
  }
  return *array_value_;
}

std::shared_ptr<ArrayBuffer> Value::ToArrayBuffer() {
  if (type_ != Type::Object) {
    panic("Value::ToArrayBuffer: the type must be Type::Object but not: " + Inspect());
  }
  if (!object_value_) {
    panic("Value::ToArrayBuffer: object_value_ must not be null");
  }
  return std::static_pointer_cast<ArrayBuffer>(object_value_);
}

std::string Value::Inspect() const {
  switch (type_) {
  case Type::Null:
    return "null";
  case Type::Undefined:
    return "undefined";
  case Type::Bool:
    return ToBool() ? "true" : "false";
  case Type::Number:
    return std::to_string(ToNumber());
  case Type::String:
    return ToString();
  case Type::Object:
    if (IsArray()) {
      std::string str = "[";
      for (auto& v : *array_value_) {
        str += v.Inspect() + " ";
      }
      if (array_value_->size()) {
        str.resize(str.size()-1);
      }
      str += "]";
      return str;
    }
    if (IsObject()) {
      return ToObject().Inspect();
    }
    return "(object)";
  default:
    panic("invalid type: " + std::to_string(static_cast<int>(type_)));
  }
  return "";
}

Object::~Object() = default;

Value Object::Get(const std::string& key) {
  panic("Object::Get is not implemented: this: " + Inspect() + ", key: " + key);
  return Value{};
}

void Object::Set(const std::string& key, Value value) {
  panic("Object::Set is not implemented: this: " + Inspect() + ", key: " + key + ", value: " + value.Inspect());
}

void Object::Delete(const std::string& key) {
  panic("Object::Delete is not implemented: this: " + Inspect() + ", key: " + key);
}

Value Object::Invoke(Value self, std::vector<Value> args) {
  // TODO: Make this a pure virtual function?
  panic("Object::Invoke is not implemented: this: " + Inspect() + ", self: " + self.Inspect());
  return Value{};
};

Value Object::New(std::vector<Value> args) {
  // TODO: Make this a pure virtual function?
  panic("Object::New is not implemented: this: " + Inspect());
  return Value{};
};

BytesSpan Object::ToBytes() {
  return BytesSpan{};
}

std::string Object::Inspect() const {
  return ToString();
}

ArrayBuffer::ArrayBuffer(size_t size)
    : data_(size) {
}

size_t ArrayBuffer::ByteLength() const {
  return data_.size();
}

Value ArrayBuffer::Get(const std::string& key) {
  if (key == "byteLength") {
    return Value{static_cast<double>(ByteLength())};
  }
  return Value{};
}

BytesSpan ArrayBuffer::ToBytes() {
  return BytesSpan{&*data_.begin(), data_.size()};
}

std::string ArrayBuffer::ToString() const {
  return "ArrayBuffer";
}

DictionaryValues::DictionaryValues() {
}

DictionaryValues::DictionaryValues(const std::map<std::string, Value>& dict)
    : dict_{dict} {
}

Value DictionaryValues::Get(const std::string& key) {
  auto it = dict_.find(key);
  if (it == dict_.end()) {
    return Value{};
  }
  return it->second;
}

void DictionaryValues::Set(const std::string& key, Value object) {
  dict_[key] = object;
}

void DictionaryValues::Delete(const std::string& key) {
  dict_.erase(key);
}

std::string DictionaryValues::ToString() const {
  return "DictionaryValues";
}

std::string DictionaryValues::Inspect() const {
  std::string str = "{";
  for (auto& kv : dict_) {
    str += kv.first + ":" + kv.second.Inspect() + " ";
  }
  if (dict_.size()) {
    str.resize(str.size()-1);
  }
  str += "}";
  return str;
}

Function::Function(Object::Func fn)
    : fn_(fn) {
}

Value Function::Invoke(Value self, std::vector<Value> args) {
  return fn_(Value{}, args);
}

Value Value::Global() {
  static Value global = MakeGlobal();
  return global;
}

Value Value::MakeGlobal() {
  std::shared_ptr<Constructor> arr = std::make_shared<Constructor>("Array",
    [](Value self, std::vector<Value> args) -> Value {
      // TODO: Implement this.
      return Value{};
    });
  std::shared_ptr<Constructor> obj = std::make_shared<Constructor>("Object",
    [](Value self, std::vector<Value> args) -> Value {
      if (args.size() == 1) {
        panic("new Object(" + args[0].Inspect() + ") is not implemented");
      }
      return Value{std::make_shared<DictionaryValues>()};
    });

  std::shared_ptr<Constructor> arrayBuffer = std::make_shared<Constructor>("ArrayBuffer",
    [](Value self, std::vector<Value> args) -> Value {
      if (args.size() == 0) {
        panic("new ArrayBuffer() is not implemented");
      }
      if (args.size() == 1) {
        Value vlen = args[0];
        if (!vlen.IsNumber()) {
          panic("new ArrayBuffer(" + args[0].Inspect() + ") is not implemented");
        }
        size_t len = static_cast<size_t>(vlen.ToNumber());
        return Value{std::make_shared<ArrayBuffer>(len)};
      }
      panic("new ArrayBuffer with " + std::to_string(args.size()) + " args is not implemented");
      return Value{};
    });

  std::shared_ptr<Constructor> u8 = std::make_shared<Constructor>("Uint8Array",
    [](Value self, std::vector<Value> args) -> Value {
      if (args.size() == 0) {
        return Value{std::make_shared<Uint8Array>(0)};
      }
      if (args.size() == 1) {
        Value vlen = args[0];
        if (!vlen.IsNumber()) {
          panic("new Uint8Array(" + args[0].Inspect() + ") is not implemented");
        }
        size_t len = static_cast<size_t>(vlen.ToNumber());
        return Value{std::make_shared<Uint8Array>(len)};
      }
      if (args.size() == 3) {
        if (!args[0].IsObject()) {
          panic("new Uint8Array's first argument must be an ArrayBuffer but " + args[0].Inspect());
        }
        if (!args[1].IsNumber()) {
          panic("new Uint8Array's second argument must be a number but " + args[1].Inspect());
        }
        if (!args[2].IsNumber()) {
          panic("new Uint8Array's third argument must be a number but " + args[2].Inspect());
        }
        std::shared_ptr<ArrayBuffer> ab = args[0].ToArrayBuffer();
        size_t offset = static_cast<size_t>(args[1].ToNumber());
        size_t length = static_cast<size_t>(args[2].ToNumber());
        auto u8 = std::make_shared<Uint8Array>(ab, offset, length);
        return Value{u8};
      }
      panic("new Uint8Array with " + std::to_string(args.size()) + " args is not implemented");
      return Value{};
    });

  std::shared_ptr<Constructor> u16 = std::make_shared<Constructor>("Uint16Array",
    [](Value self, std::vector<Value> args) -> Value {
      if (args.size() == 3) {
        if (!args[0].IsObject()) {
          panic("new Uint16Array's first argument must be an ArrayBuffer but " + args[0].Inspect());
        }
        if (!args[1].IsNumber()) {
          panic("new Uint16Array's second argument must be a number but " + args[1].Inspect());
        }
        if (!args[2].IsNumber()) {
          panic("new Uint16Array's third argument must be a number but " + args[2].Inspect());
        }
        std::shared_ptr<ArrayBuffer> ab = args[0].ToArrayBuffer();
        size_t offset = static_cast<size_t>(args[1].ToNumber());
        size_t length = static_cast<size_t>(args[2].ToNumber());
        auto u16 = std::make_shared<Uint16Array>(ab, offset, length);
        return Value{u16};
      }
      panic("new Uint16Array with " + std::to_string(args.size()) + " args is not implemented");
      return Value{};
    });

  std::shared_ptr<Constructor> f32 = std::make_shared<Constructor>("Float32Array",
    [](Value self, std::vector<Value> args) -> Value {
      if (args.size() == 0) {
        return Value{std::make_shared<Float32Array>(0)};
      }
      if (args.size() == 3) {
        if (!args[0].IsObject()) {
          panic("new Float32Array's first argument must be an ArrayBuffer but " + args[0].Inspect());
        }
        if (!args[1].IsNumber()) {
          panic("new Float32Array's second argument must be a number but " + args[1].Inspect());
        }
        if (!args[2].IsNumber()) {
          panic("new Float32Array's third argument must be a number but " + args[2].Inspect());
        }
        std::shared_ptr<ArrayBuffer> ab = args[0].ToArrayBuffer();
        size_t offset = static_cast<size_t>(args[1].ToNumber());
        size_t length = static_cast<size_t>(args[2].ToNumber());
        auto f32 = std::make_shared<Float32Array>(ab, offset, length);
        return Value{f32};
      }
      panic("new Float32Array with " + std::to_string(args.size()) + " args is not implemented");
      return Value{};
    });

  Value getRandomValues{std::make_shared<Function>(
    [](Value self, std::vector<Value> args) -> Value {
      BytesSpan bs = args[0].ToBytes();
      // TODO: Use cryptographically strong random values instead of std::random_device.
      static std::random_device rd;
      std::uniform_int_distribution<uint8_t> dist(0, 255);
      for (size_t i = 0; i < bs.size(); i++) {
        bs[i] = dist(rd);
      }
      return Value{};
    })};
  std::shared_ptr<DictionaryValues> crypto = std::make_shared<DictionaryValues>(std::map<std::string, Value>{
    {"getRandomValues", getRandomValues},
  });

  static Value& writeObjectsToStdout = *new Value(std::make_shared<Function>(
    [](Value self, std::vector<Value> args) -> Value {
      WriteObjects(std::cout, args);
      return Value{};
    }));
  static Value& writeObjectsToStderr = *new Value(std::make_shared<Function>(
    [](Value self, std::vector<Value> args) -> Value {
      WriteObjects(std::cerr, args);
      return Value{};
    }));
  std::shared_ptr<DictionaryValues> console = std::make_shared<DictionaryValues>(std::map<std::string, Value>{
    {"error", writeObjectsToStderr},
    {"debug", writeObjectsToStderr},
    {"info", writeObjectsToStdout},
    {"log", writeObjectsToStdout},
    {"warm", writeObjectsToStderr},
  });

  std::shared_ptr<Function> fetch = std::make_shared<Function>(
    [](Value self, std::vector<Value> args) -> Value {
      // TODO: Implement this.
      return Value{};
    });

  static FS& fsimpl = *new FS();
  std::shared_ptr<DictionaryValues> fs = std::make_shared<DictionaryValues>(std::map<std::string, Value>{
    {"constants", Value{std::make_shared<DictionaryValues>(std::map<std::string, Value>{
        {"O_WRONLY", Value{-1.0}},
        {"O_RDWR", Value{-1.0}},
        {"O_CREAT", Value{-1.0}},
        {"O_TRUNC", Value{-1.0}},
        {"O_APPEND", Value{-1.0}},
        {"O_EXCL", Value{-1.0}},
      })}},
    {"write", Value{std::make_shared<Function>(
      [](Value self, std::vector<Value> args) -> Value {
        return fsimpl.Write(self, args);
      })}},
  });

  std::shared_ptr<DictionaryValues> process = std::make_shared<DictionaryValues>(std::map<std::string, Value>{
    {"pid", Value{-1.0}},
    {"ppid", Value{-1.0}},
  });

  std::shared_ptr<DictionaryValues> global = std::make_shared<DictionaryValues>(std::map<std::string, Value>{
    {"Array", Value{arr}},
    {"Object", Value{obj}},
    {"ArrayBuffer", Value{arrayBuffer}},
    {"Uint8Array", Value{u8}},
    {"Uint16Array", Value{u16}},
    {"Float32Array", Value{f32}},
    {"console", Value{console}},
    {"crypto", Value{crypto}},
    {"fetch", Value{fetch}},
    {"fs", Value{fs}},
    {"process", Value{process}},
  });

  return Value{global};
}

Value Value::ReflectGet(Value target, const std::string& key) {
  if (target.IsUndefined()) {
    panic("get on undefined (key: " + key + ") is forbidden");
    return Value{};
  }
  if (target.IsNull()) {
    panic("get on null (key: " + key + ") is forbidden");
    return Value{};
  }
  if (target.IsObject()) {
    return target.ToObject().Get(key);
  }
  if (target.IsArray()) {
    int idx = std::stoi(key);
    if (idx > 0 || (idx == 0 && key == "0")) {
      return target.ToArray()[idx];
    }
  }
  panic(target.Inspect() + "." + key + " not found");
  return Value{};
}

void Value::ReflectSet(Value target, const std::string& key, Value value) {
  if (target.IsUndefined()) {
    panic("set on undefined (key: " + key + ") is forbidden");
  }
  if (target.IsNull()) {
    panic("set on null (key: " + key + ") is forbidden");
  }
  if (target.IsObject()) {
    target.ToObject().Set(key, value);
    return;
  }
  panic(target.Inspect() + "." + key + " cannot be set");
}

void Value::ReflectDelete(Value target, const std::string& key) {
  if (target.IsUndefined()) {
    panic("delete on undefined (key: " + key + ") is forbidden");
  }
  if (target.IsNull()) {
    panic("delete on null (key: " + key + ") is forbidden");
  }
  if (target.IsObject()) {
    target.ToObject().Delete(key);
    return;
  }
  panic(target.Inspect() + "." + key + " cannot be deleted");
}

Value Value::ReflectConstruct(Value target, std::vector<Value> args) {
  if (target.IsUndefined()) {
    panic("new on undefined is forbidden");
    return Value{};
  }
  if (target.IsNull()) {
    panic("new on null is forbidden");
    return Value{};
  }
  if (target.IsObject()) {
    Object& t = target.ToObject();
    if (!t.IsConstructor()) {
      panic(t.ToString() + " is not a constructor");
      return Value{};
    }
    return t.New(args);
  }
  panic("new " + target.Inspect() + "(" + JoinObjects(args) + ") cannot be called");
  return Value{};
}

Value Value::ReflectApply(Value target, Value self, std::vector<Value> args) {
  if (target.IsUndefined()) {
    panic("apply on undefined is forbidden");
    return Value{};
  }
  if (target.IsNull()) {
    panic("apply on null is forbidden");
    return Value{};
  }
  if (target.IsObject()) {
    Object& t = target.ToObject();
    if (t.IsConstructor()) {
      panic(t.ToString() + " is a constructor");
      return Value{};
    }
    return t.Invoke(self, args);
  }
  panic(target.Inspect() + "(" + JoinObjects(args) + ") cannot be called");
  return Value{};
}

Constructor::Constructor(const std::string& name, Object::Func fn)
    : name_(name),
      fn_(fn) {
}

Value Constructor::New(std::vector<Value> args) {
  return fn_(Value{}, args);
}

std::string Constructor::ToString() const {
  return name_;
}

}
`))
