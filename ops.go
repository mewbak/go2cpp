// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/go-interpreter/wagon/disasm"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/operators"
)

type ReturnType int

const (
	ReturnTypeVoid ReturnType = iota
	ReturnTypeI32
	ReturnTypeI64
	ReturnTypeF32
	ReturnTypeF64
)

func (r ReturnType) CSharp() string {
	switch r {
	case ReturnTypeVoid:
		return "void"
	case ReturnTypeI32:
		return "int"
	case ReturnTypeI64:
		return "long"
	case ReturnTypeF32:
		return "float"
	case ReturnTypeF64:
		return "double"
	default:
		panic("not reached")
	}
}

func opsToCSharp(code []byte, sig *wasm.FunctionSig) ([]string, error) {
	instrs, err := disasm.Disassemble(code)
	if err != nil {
		return nil, err
	}
	return instrsToCSharp(instrs, sig)
}

func instrsToCSharp(instrs []disasm.Instr, sig *wasm.FunctionSig) ([]string, error) {
	var body []string
	var newIdx int
	var idxStack []int

	nextIdx := func() int {
		idx := newIdx
		newIdx++
		return idx
	}

	for _, instr := range instrs {
		switch instr.Op.Code {
		case operators.Unreachable:
			body = append(body, `Debug.Assert(false, "not reached");`)
		case operators.Nop:
			// Do nothing
		case operators.Block:
			// TODO: Implement this.
		case operators.Loop:
			// TODO: Implement this.
		case operators.If:
			// TODO: Implement this.
		case operators.Else:
			// TODO: Implement this.
		case operators.End:
			// TODO: Implement this.
		case operators.Drop:
			idxStack = idxStack[:len(idxStack)-1]
		case operators.Select:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]

		case operators.GetLocal:
			idx := nextIdx()
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = local%d;", idx, instr.Immediates[0]))
			idxStack = append(idxStack, idx)
		case operators.SetLocal:
			idx := idxStack[len(idxStack)-1]
			body = append(body, fmt.Sprintf("local%d = stack%d;", instr.Immediates[0], idx))
			idxStack = idxStack[:len(idxStack)-1]
		case operators.TeeLocal:
			idx := idxStack[len(idxStack)-1]
			body = append(body, fmt.Sprintf("local%d = stack%d;", instr.Immediates[0], idx))
		case operators.GetGlobal:
			idx := nextIdx()
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = global%d;", idx, instr.Immediates[0]))
			idxStack = append(idxStack, idx)
		case operators.SetGlobal:
			idx := idxStack[len(idxStack)-1]
			body = append(body, fmt.Sprintf("global%d = stack%d;", instr.Immediates[0], idx))
			idxStack = idxStack[:len(idxStack)-1]

		case operators.I32Load:
			// TODO: Implement this.
		case operators.I64Load:
			// TODO: Implement this.
		case operators.F32Load:
			// TODO: Implement this.
		case operators.F64Load:
			// TODO: Implement this.
		case operators.I32Load8s:
			// TODO: Implement this.
		case operators.I32Load8u:
			// TODO: Implement this.
		case operators.I32Load16s:
			// TODO: Implement this.
		case operators.I32Load16u:
			// TODO: Implement this.
		case operators.I64Load8s:
			// TODO: Implement this.
		case operators.I64Load8u:
			// TODO: Implement this.
		case operators.I64Load16s:
			// TODO: Implement this.
		case operators.I64Load16u:
			// TODO: Implement this.
		case operators.I64Load32s:
			// TODO: Implement this.
		case operators.I64Load32u:
			// TODO: Implement this.

		case operators.I32Store:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I64Store:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.F32Store:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.F64Store:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I32Store8:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I32Store16:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I64Store8:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I64Store16:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]
		case operators.I64Store32:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-2]

		case operators.CurrentMemory:
			idx := nextIdx()
			// TOOD: Implement this.
			idxStack = append(idxStack, idx)
		case operators.GrowMemory:
			// TOOD: Implement this.

		case operators.I32Const:
			idx := nextIdx()
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = %d;", idx, instr.Immediates[0]))
			idxStack = append(idxStack, idx)
		case operators.I64Const:
			idx := nextIdx()
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = %d;", idx, instr.Immediates[0]))
			idxStack = append(idxStack, idx)
		case operators.F32Const:
			idx := nextIdx()
			// TODO: Implement this.
			// https://docs.microsoft.com/en-us/dotnet/api/system.runtime.compilerservices.unsafe?view=netcore-3.1
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = 0 /* TODO */;", idx))
			idxStack = append(idxStack, idx)
		case operators.F64Const:
			idx := nextIdx()
			// TODO: Implement this.
			// TODO: Use a proper type.
			body = append(body, fmt.Sprintf("dynamic stack%d = 0 /* TODO */;", idx))
			idxStack = append(idxStack, idx)

		case operators.I32Eqz:
			// TODO: Implement this.
		case operators.I32Eq:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Ne:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32LtS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32LtU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32GtS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32GtU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32LeS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32LeU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32GeS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32GeU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Eqz:
			// TODO: Implement this.
		case operators.I64Eq:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Ne:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64LtS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64LtU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64GtS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64GtU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64LeS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64LeU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64GeS:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64GeU:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Eq:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Ne:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Lt:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Gt:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Le:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Ge:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Eq:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Ne:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Lt:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Gt:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Le:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Ge:
			// TODO: Implement this.
			idxStack = idxStack[:len(idxStack)-1]

		case operators.I32Clz:
			// TODO: Implement this
		case operators.I32Ctz:
			// TODO: Implement this
		case operators.I32Popcnt:
			// TODO: Implement this
		case operators.I32Add:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Sub:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Mul:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32DivS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32DivU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32RemS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32RemU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32And:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Or:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Xor:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Shl:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32ShrS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32ShrU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Rotl:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I32Rotr:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Clz:
			// TODO: Implement this
		case operators.I64Ctz:
			// TODO: Implement this
		case operators.I64Popcnt:
			// TODO: Implement this
		case operators.I64Add:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Sub:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Mul:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64DivS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64DivU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64RemS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64RemU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64And:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Or:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Xor:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Shl:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64ShrS:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64ShrU:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Rotl:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.I64Rotr:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Abs:
			// TODO: Implement this
		case operators.F32Neg:
			// TODO: Implement this
		case operators.F32Ceil:
			// TODO: Implement this
		case operators.F32Floor:
			// TODO: Implement this
		case operators.F32Trunc:
			// TODO: Implement this
		case operators.F32Nearest:
			// TODO: Implement this
		case operators.F32Sqrt:
			// TODO: Implement this
		case operators.F32Add:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Sub:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Mul:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Div:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Min:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Max:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F32Copysign:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Abs:
			// TODO: Implement this
		case operators.F64Neg:
			// TODO: Implement this
		case operators.F64Ceil:
			// TODO: Implement this
		case operators.F64Floor:
			// TODO: Implement this
		case operators.F64Trunc:
			// TODO: Implement this
		case operators.F64Nearest:
			// TODO: Implement this
		case operators.F64Sqrt:
			// TODO: Implement this
		case operators.F64Add:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Sub:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Mul:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Div:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Min:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Max:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		case operators.F64Copysign:
			// TODO: Implement this
			idxStack = idxStack[:len(idxStack)-1]
		default:
			body = append(body, fmt.Sprintf("// %v", instr.Op))
		}
	}
	switch len(sig.ReturnTypes) {
	case 0:
		// TODO: Enable this error
		/*if len(idxStack) != 0 {
			return nil, fmt.Errorf("the stack length must be 0 but %d", len(idxStack))
		}*/
	case 1:
		// TODO: The stack must be exactly 1.
		/*if len(idxStack) == 0 {
			return nil, fmt.Errorf("the stack length must be 1 but %d", len(idxStack))
		}*/
		if len(idxStack) == 0 {
			body = append(body, `throw new Exception("not reached");`)
		} else {
			body = append(body, fmt.Sprintf("return stack%d;", idxStack[0]))
		}
	}
	return body, nil
}
