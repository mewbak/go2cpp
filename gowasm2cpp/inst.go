// SPDX-License-Identifier: Apache-2.0

package gowasm2cpp

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"golang.org/x/sync/errgroup"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func writeInst(dir string, incpath string, namespace string, importFuncs, funcs []*wasmFunc, exports []*wasmExport, globals []*wasmGlobal, types []*wasmType, tables [][]uint32) error {
	const groupSize = 64

	sort.Slice(funcs, func(a, b int) bool {
		return funcs[a].Wasm.Name < funcs[b].Wasm.Name
	})
	sort.Slice(exports, func(a, b int) bool {
		return exports[a].Name < exports[b].Name
	})

	var g errgroup.Group
	g.Go(func() error {
		f, err := os.Create(filepath.Join(dir, "inst.h"))
		if err != nil {
			return err
		}
		defer f.Close()

		m := 0
		for _, t := range tables {
			if m < len(t) {
				m = len(t)
			}
		}
		if err := instHTmpl.Execute(f, struct {
			IncludeGuard        string
			IncludePath         string
			Namespace           string
			ImportFuncs         []*wasmFunc
			Exports             []*wasmExport
			Funcs               []*wasmFunc
			Types               []*wasmType
			Globals             []*wasmGlobal
			NumFuncs            int
			NumTable            int
			NumMaxTableElements int
		}{
			IncludeGuard:        includeGuard(namespace) + "_INST_H",
			IncludePath:         incpath,
			Namespace:           namespace,
			ImportFuncs:         importFuncs,
			Exports:             exports,
			Funcs:               funcs,
			Types:               types,
			Globals:             globals,
			NumFuncs:            len(importFuncs) + len(funcs),
			NumTable:            len(tables),
			NumMaxTableElements: m,
		}); err != nil {
			return err
		}
		return nil
	})

	groups := map[byte][]*wasmFunc{}
	for _, f := range funcs {
		n := f.Wasm.Name
		g := n[0]
		groups[g] = append(groups[g], f)
	}

	for gp, fs := range groups {
		gp := gp
		fs := fs
		g.Go(func() error {
			f, err := os.Create(filepath.Join(dir, fmt.Sprintf("inst.funcs.%c.cpp", gp)))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := instFuncCppTmpl.Execute(f, struct {
				IncludePath string
				Namespace   string
				Funcs       []*wasmFunc
			}{
				IncludePath: incpath,
				Namespace:   namespace,
				Funcs:       fs,
			}); err != nil {
				return err
			}
			return nil
		})
	}

	// exports
	g.Go(func() error {
		f, err := os.Create(filepath.Join(dir, "inst.exports.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := instExportsCppTmpl.Execute(f, struct {
			IncludePath string
			Namespace   string
			Exports     []*wasmExport
		}{
			IncludePath: incpath,
			Namespace:   namespace,
			Exports:     exports,
		}); err != nil {
			return err
		}
		return nil
	})

	// init
	g.Go(func() error {
		f, err := os.Create(filepath.Join(dir, "inst.init.cpp"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := instInitCppTmpl.Execute(f, struct {
			IncludePath string
			Namespace   string
			ImportFuncs []*wasmFunc
			Funcs       []*wasmFunc
			Types       []*wasmType
			Tables      [][]uint32
			Globals     []*wasmGlobal
		}{
			IncludePath: incpath,
			Namespace:   namespace,
			ImportFuncs: importFuncs,
			Funcs:       funcs,
			Types:       types,
			Tables:      tables,
			Globals:     globals,
		}); err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

var instHTmpl = template.Must(template.New("inst.h").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#ifndef {{.IncludeGuard}}
#define {{.IncludeGuard}}

#include <cstdint>

namespace {{.Namespace}} {

class Mem;

class IImport {
public:
  virtual ~IImport();

{{range $value := .ImportFuncs}}{{$value.CppDecl "  " true false}}

{{end -}} };

class Inst {
public:
  Inst(Mem* mem, IImport* import);

{{range $value := .Exports}}{{$value.CppDecl "  "}}
{{end}}
private:
{{range $value := .Types}}  using Type{{.Index}} = {{.Cpp}};
{{end}}
  union Func {
{{range $value := .Types}}    Type{{.Index}} type{{.Index}}_;
{{end}}  };

{{range $value := .Funcs}}{{$value.CppDecl "  " false false}}

{{end}}  Mem* mem_;
  IImport* import_;
  Func funcs_[{{.NumFuncs}}];
  uint32_t table_[{{.NumTable}}][{{.NumMaxTableElements}}];

{{range $value := .Globals}}  {{$value.Cpp}}
{{end}}};

}

#endif  // {{.IncludeGuard}}
`))

var instFuncCppTmpl = template.Must(template.New("inst.funcs.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}inst.h"

#include <cassert>
#include <cmath>
#include "{{.IncludePath}}bits.h"
#include "{{.IncludePath}}mem.h"

namespace {{.Namespace}} {

{{range $value := .Funcs}}{{$value.CppImpl "Inst" ""}}
{{end}}}
`))

var instExportsCppTmpl = template.Must(template.New("inst.exports.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}inst.h"

namespace {{.Namespace}} {

{{range $value := .Exports}}{{$value.CppImpl ""}}
{{end}}}
`))

var instInitCppTmpl = template.Must(template.New("inst.init.cpp").Parse(`// Code generated by go2cpp. DO NOT EDIT.

#include "{{.IncludePath}}inst.h"

namespace {{.Namespace}} {

IImport::~IImport() = default;

Inst::Inst(Mem* mem, IImport* import)
    : mem_{mem},
      import_{import},
      table_{
{{range $value := .Tables}}        { {{- range $value2 := $value}}{{$value2}}, {{end}} },
{{end}}      } {
{{range $value := .ImportFuncs}}  funcs_[{{.Index}}].type0_ = nullptr;
{{end}}{{range $value := .Funcs}}  funcs_[{{.Index}}].type{{.Type.Index}}_ = &Inst::{{.Identifier}};
{{end}}}

}
`))
