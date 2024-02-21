package vm

type OpCode byte

type Instruction struct {
	Op  OpCode
	Val int
}

const (
	__MEM__ OpCode = iota + 0
	LOADC
	LOAD
	DECL
	ASGN
	POP
	DUP
	JMP
	JMPT
	JMPF
	CALL
	CALWT
	FUNC
	RET
	RETV
	ENVST
	ENVED
	GET
	SET
	GETWT
)

const (
	__CALC__ OpCode = iota + 0x40
	ADD
	SUB
	MUL
	DIV
	MOD
	NEG
	NOT
	EQ
	NEQ
	LT
	LE
	GT
	GE
	PRINT
)

func (op OpCode) String() string {
	switch op {
	case LOADC:
		return "LOADC"
	case LOAD:
		return "LOAD"
	case DECL:
		return "DECL"
	case ASGN:
		return "ASGN"
	case POP:
		return "POP"
	case DUP:
		return "DUP"
	case JMP:
		return "JMP"
	case JMPT:
		return "JMPT"
	case JMPF:
		return "JMPF"
	case CALL:
		return "CALL"
	case CALWT:
		return "CALWT"
	case FUNC:
		return "FUNC"
	case RET:
		return "RET"
	case RETV:
		return "RETV"
	case ENVST:
		return "ENVST"
	case ENVED:
		return "ENVED"
	case GET:
		return "GET"
	case SET:
		return "SET"
	case GETWT:
		return "GETWT"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case NEG:
		return "NEG"
	case NOT:
		return "NOT"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case LT:
		return "LT"
	case LE:
		return "LE"
	case GT:
		return "GT"
	case GE:
		return "GE"
	case PRINT:
		return "PRINT"
	default:
		return ""
	}
}
