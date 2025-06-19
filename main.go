package main

import (
	"fmt"
	"log"
	"os"
)

// WASM opcodes used in add.wasm
const (
	OP_LOCAL_GET = 0x20
	OP_I32_ADD   = 0x6a
	OP_END       = 0x0b
)

func main() {
	wasm, err := os.ReadFile("add.wasm")
	if err != nil {
		log.Fatalf("failed to read add.wasm: %v", err)
	}

	// Find the code section (section id 10)
	codeSection := findSection(wasm, 10)
	if codeSection == nil {
		log.Fatal("Code section not found")
	}

	// Only one function body is expected
	funcBody := extractFunctionBody(codeSection)

	// Call the function with two arguments
	result := interpretFunction(funcBody, 3, 4)
	fmt.Printf("add(3, 4) = %d\n", result)
}

// findSection locates a section by id in the WASM binary
func findSection(wasm []byte, id byte) []byte {
	pos := 8 // skip magic + version
	for pos < len(wasm) {
		sectionId := wasm[pos]
		pos++
		sectionLen, n := readULEB(wasm[pos:])
		pos += n
		if sectionId == id {
			return wasm[pos : pos+sectionLen]
		}
		pos += sectionLen
	}
	return nil
}

// extractFunctionBody extracts the first function body from the code section
func extractFunctionBody(code []byte) []byte {
	_, n := readULEB(code)
	pos := n
	bodySize, m := readULEB(code[pos:])
	pos += m
	return code[pos : pos+bodySize]
}

// interpretFunction interprets the function body for add(a, b)
func interpretFunction(body []byte, a, b int32) int32 {
	_, n := readULEB(body) // localsCount, not used
	pos := n
	locals := []int32{a, b}
	stack := []int32{}
	for pos < len(body) {
		op := body[pos]
		pos++
		switch op {
		case OP_LOCAL_GET:
			idx, m := readULEB(body[pos:])
			pos += m
			stack = append(stack, locals[idx])
		case OP_I32_ADD:
			v2 := stack[len(stack)-1]
			v1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, v1+v2)
		case OP_END:
			return stack[len(stack)-1]
		}
	}
	return 0
}

// readULEB reads a ULEB128-encoded integer
func readULEB(b []byte) (int, int) {
	result := 0
	shift := 0
	for i, v := range b {
		result |= int(v&0x7f) << shift
		if v&0x80 == 0 {
			return result, i + 1
		}
		shift += 7
	}
	return 0, 0
}
