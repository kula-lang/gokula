package vm

import (
	"encoding/binary"
	"fmt"
	"gokula/objects"
	"io"
	"os"
	"strings"
)

const (
	MAGIC_NUMBER uint16 = 0x1701
	SEPRATOR     byte   = 0xff
	NONE         byte   = 0x80
	BOOL         byte   = 0x81
	DOUBLE       byte   = 0x82
	STRING       byte   = 0x83
)

var CompiledFileInstance *CompiledFile

type FunctionChunk struct {
	Params       []uint16
	Instructions []Instruction
}

type CompiledFile struct {
	SymbolArray []string
	Literals    []any
	Chunk       []Instruction
	Functions   []*FunctionChunk
}

func Load(path string) (*CompiledFile, error) {
	compiledFile := new(CompiledFile)

	compiledFile.SymbolArray = make([]string, 0)
	compiledFile.Literals = make([]any, 0)
	compiledFile.Chunk = make([]Instruction, 0)
	compiledFile.Functions = make([]*FunctionChunk, 0)

	var err error

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open file \"%s\"\n%s", path, err.Error())
		return nil, err
	}
	defer file.Close()

	// Read Magic Number
	var magic_number uint16
	err = binary.Read(file, binary.LittleEndian, &magic_number)
	if err != nil {
		return nil, err
	}
	if magic_number != MAGIC_NUMBER {
		return nil, fmt.Errorf("not a kulac file")
	}

	var byte_buffer uint8
	// Read Signal Table
	for {
		err := binary.Read(file, binary.LittleEndian, &byte_buffer)
		if err != nil {
			return compiledFile, err
		}
		if byte_buffer == SEPRATOR {
			break
		} else {
			bytes := make([]byte, int32(byte_buffer))
			err = binary.Read(file, binary.LittleEndian, &bytes)
			if err != nil {
				return nil, err
			}
			symbol := string(bytes)
			compiledFile.SymbolArray = append(compiledFile.SymbolArray, symbol)
		}
	}

	// Read Literals
	compiledFile.Literals = append(compiledFile.Literals, false, true, nil)
	for {
		err = binary.Read(file, binary.LittleEndian, &byte_buffer)
		if err != nil {
			return nil, err
		}
		if byte_buffer == SEPRATOR {
			break
		}

		switch byte_buffer {
		case NONE, BOOL:
		case STRING:
			var bytes_size int32
			err = binary.Read(file, binary.LittleEndian, &bytes_size)
			if err != nil {
				return nil, err
			}
			bytes := make([]byte, bytes_size)
			err = binary.Read(file, binary.LittleEndian, bytes)
			if err != nil {
				return nil, err
			}
			literal := objects.KulaString(bytes)
			compiledFile.Literals = append(compiledFile.Literals, &literal)
		case DOUBLE:
			var number objects.KulaNumber
			err = binary.Read(file, binary.LittleEndian, &number)
			if err != nil {
				return nil, err
			}
			compiledFile.Literals = append(compiledFile.Literals, number)
		default:
			return nil, fmt.Errorf("undefined literal type")
		}
	}

	// Read Instructions
	for {
		err = binary.Read(file, binary.LittleEndian, &byte_buffer)
		if err != nil {
			return nil, err
		}
		if byte_buffer == SEPRATOR {
			break
		}

		inst, err := readInstruction(byte_buffer, file)
		if err != nil {
			return nil, err
		}
		compiledFile.Chunk = append(compiledFile.Chunk, inst)
	}

	// Read Functions
	for {
		err = binary.Read(file, binary.LittleEndian, &byte_buffer)
		if err != nil {
			if err == io.EOF {
				return compiledFile, nil
			}
			return nil, err
		}

		param_size := int32(byte_buffer)
		function := new(FunctionChunk)
		function.Params = make([]uint16, param_size)
		err = binary.Read(file, binary.LittleEndian, &function.Params)
		if err != nil {
			return nil, err
		}
		function.Instructions = make([]Instruction, 0)
		for {
			err = binary.Read(file, binary.LittleEndian, &byte_buffer)
			if err != nil {
				return nil, err
			}
			if byte_buffer == SEPRATOR {
				break
			}

			inst, err := readInstruction(byte_buffer, file)
			if err != nil {
				return nil, err
			}
			function.Instructions = append(function.Instructions, inst)
		}
		compiledFile.Functions = append(compiledFile.Functions, function)
	}
}

func codeSize(op OpCode) int {
	switch op {
	case LOADC, LOAD, DECL, ASGN:
		return 16
	case JMP, JMPT, JMPF:
		return 16
	case FUNC, PRINT, CALL, CALWT:
		return 8
	default:
		return 0
	}
}

func readInstruction(byte_buffer byte, file *os.File) (Instruction, error) {
	var err error
	inst := Instruction{}
	inst.Op = OpCode(byte_buffer)
	switch codeSize(OpCode(byte_buffer)) {
	case 32:
		var uint32_buffer uint32
		err = binary.Read(file, binary.LittleEndian, &uint32_buffer)
		if err != nil {
			return inst, err
		}
		inst.Val = int(uint32_buffer)
	case 16:
		var uint16_buffer uint16
		err = binary.Read(file, binary.LittleEndian, &uint16_buffer)
		if err != nil {
			return inst, err
		}
		inst.Val = int(uint16_buffer)
	case 8:
		var uint8_buffer uint8
		err = binary.Read(file, binary.LittleEndian, &uint8_buffer)
		if err != nil {
			return inst, err
		}
		inst.Val = int(uint8_buffer)
	}
	return inst, nil
}

func (kulac *CompiledFile) String() string {
	var sb strings.Builder
	sb.Grow(1024)

	sb.WriteString("==== Symbols ====\n")
	for i, v := range kulac.SymbolArray {
		sb.WriteString(fmt.Sprintf("\t%4d\t%s\n", i, v))
	}

	sb.WriteString("==== Literal ====\n")
	for i, v := range kulac.Literals {
		sb.WriteString(fmt.Sprintf("\t%4d\t%s\n", i, v))
	}

	sb.WriteString("==== Instructions ====\n")
	for i, v := range kulac.Chunk {
		sb.WriteString(fmt.Sprintf("\t%4d\t%s\t%d\n", i, v.Op.String(), v.Val))
	}

	sb.WriteString("==== Functions ====\n")
	for fIndex, f := range kulac.Functions {
		sb.WriteString(fmt.Sprintf("---- F %d ----\n", fIndex))
		for i, v := range f.Params {
			sb.WriteString(fmt.Sprintf("\t%4d\t%4d\t%s\n", i, v, kulac.SymbolArray[v]))
		}
		sb.WriteString(fmt.Sprintf("---- I %d ----\n", fIndex))
		for i, v := range f.Instructions {
			sb.WriteString(fmt.Sprintf("\t%4d\t%s\t%4d\n", i, v.Op.String(), v.Val))
		}
	}

	return sb.String()
}
