package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type opcode int

const (
	// division: int(A / 2^combo operand) -> A
	adv = iota
	// bitwise XOR: B ^ lit operand -> B
	bxl
	// modulo 8: combo operand %8 -> B
	bst
	// Jump: sets IP to lit operand, unless 0
	jnz
	// bitwise XOR, ignoring operand: B ^ C -> B
	bxc
	// evals combo operand, and outputs, with comma separator
	out
	// division: int(A / 2^combo operand) -> B
	bdv
	// division: int(A / 2^combo operand) -> C
	cdv
	//
	invalidOpcode
)

func (o opcode) String() string {
	return map[opcode]string{
		adv: "adv",
		bxl: "bxl",
		bst: "bst",
		jnz: "jnz",
		bxc: "bxc",
		out: "out",
		bdv: "bdv",
		cdv: "cdv",
		//
		invalidOpcode: "invC",
	}[o]
}

type operand int

const (
	lit0 = iota // literal 0
	lit1        // literal 1
	lit2        // literal 2
	lit3        // literal 3
	regA        // combo: read reg A
	regB        // combo: read reg B
	regC        // combo: read reg B
	halt        // illegal: halt
	//
	invalidOperand
)

func (o operand) String() string {
	return map[operand]string{
		lit0: "lit0",
		lit1: "lit1",
		lit2: "lit2",
		lit3: "lit3",
		regA: "regA",
		regB: "regB",
		regC: "regC",
		halt: "halt",
		//
		invalidOperand: "invP",
	}[o]
}

type operation struct {
	opcode  opcode
	operand operand
}

func (o operation) String() string {
	return fmt.Sprintf("[%v %v]", o.opcode, o.operand)
}

type computer struct {
	// Registers, can hold any number
	A, B, C int
	// Instruction pointer, identifies position in
	// program where the next opcode will be read
	// (zero-base).
	ip int
	// The program: a sequence of 2*(3bit) opcode-operand pairs.
	program []operation
	// If the program halted this is true
	halt bool
	// The program's output
	output []int
}

func (c *computer) String() string {
	var s strings.Builder
	s.WriteString("\nComputer\n----------------\n")
	s.WriteString(fmt.Sprintf("  Registers: A:%d B:%d C:%d\n", c.A, c.B, c.C))
	s.WriteString(fmt.Sprintf("  Reg A/bin: %v\n", strconv.FormatInt(int64(c.A), 2)))
	s.WriteString(fmt.Sprintf("   InstrPtr: %d\n", c.ip))
	s.WriteString(fmt.Sprintf("    Halted?: %t\n", c.halt))
	s.WriteString(fmt.Sprintf("    Program: %v\n", c.program))
	s.WriteString(fmt.Sprintf("     Output: %v\n----------------\n", c.output))
	return s.String()
}

func initComputer(a, b, c int, program []operation) *computer {
	cpu := &computer{
		A:    a,
		B:    b,
		C:    c,
		ip:   0,
		halt: false,
	}
	cpu.program = make([]operation, len(program))
	copy(cpu.program, program)
	return cpu
}

var errHalt = errors.New("HALTED")
var errNotImpl = errors.New("TODO")

func (c *computer) nextOperation() (opcode, operand, bool) {
	if c.ip == 2*len(c.program) {
		return invalidOpcode, invalidOperand, true
	}
	op := c.program[c.ip/2]
	//log.Printf("Next operation, at ip=%d: %v", c.ip, op)
	return op.opcode, op.operand, false

}

func (c *computer) eval(opa operand) int {
	switch opa {
	case lit0:
		return 0
	case lit1:
		return 1
	case lit2:
		return 2
	case lit3:
		return 3
	case regA:
		return c.A
	case regB:
		return c.B
	case regC:
		return c.C
	}
	log.Fatalf("Unknown operand %v", opa)
	return 0
}

func intPow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func (c *computer) execute() error {
	count := 0
	for {
		count++
		//log.Printf("\n\n------------- Starting op #%d --------------\n", count)
		if c.halt {
			return errHalt
		}
		opc, opa, end := c.nextOperation()
		if end {
			break
		}
		incIp := true

		div := func(num int, opa operand) int {
			return int(float64(num) / math.Pow(2.0, float64(c.eval(opa))))
		}

		log.Printf(">> %v %v", opc, opa)

		switch opc {
		case adv:
			r := div(c.A, opa)
			log.Printf("A/2^%s -> %d/%d -> %d -> A", opa, c.A, intPow(2, c.eval(opa)), r)
			c.A = r
		case bxl:
			r := c.B ^ int(opa)
			log.Printf("B^%s -> %d^%d -> %d -> B", opa, c.B, int(opa), r)
			c.B = r
		case bst:
			r := c.eval(opa) % 8
			log.Printf("%s %%8 -> %d %%8 -> %d -> B", opa, c.eval(opa), r)
			c.B = r
		case jnz:
			if c.A == 0 {
				// does nothing
				log.Printf("A==0 -> no jump")
			} else {
				log.Printf("jump %d", int(opa))
				c.ip = int(opa)
				incIp = false
			}
		case bxc:
			// operand is ignored
			r := c.B ^ c.C
			log.Printf("B^C -> %d^%d -> %d -> B", c.B, c.C, r)
			c.B = r
		case out:
			r := c.eval(opa) % 8
			log.Printf("out %s %%8 -> %d %%8 -> %d out", opa, c.eval(opa), r)
			c.output = append(c.output, r)
		case bdv:
			r := div(c.A, opa)
			log.Printf("A/2^%s -> %d/%d -> %d -> B", opa, c.A, intPow(2, c.eval(opa)), r)
			c.B = r
		case cdv:
			r := div(c.A, opa)
			log.Printf("A/2^%s -> %d/%d -> %d -> C", opa, c.A, intPow(2, c.eval(opa)), r)
			c.C = r
		}
		if incIp {
			c.ip += 2
		}
		//log.Printf("Computer state:\n%v", c.String())
	}
	return nil
}

func (c *computer) out() string {
	var b strings.Builder
	for i, v := range c.output {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	return b.String()
}

func parseInput(in []string) (*computer, error) {
	if len(in) != 5 {
		return nil, fmt.Errorf("input: got %d lines want 5", len(in))
	}
	f1 := func(s string) int { return int(mustParseInt(s[(strings.Index(s, ":") + 2):len(s)])) }
	a, b, c := f1(in[0]), f1(in[1]), f1(in[2])
	if len(in[3]) != 0 {
		return nil, fmt.Errorf("input: got non-empty line3 (%d chars) want 0", len(in[3]))
	}
	progStr := in[4][(strings.Index(in[4], ":") + 2):len(in[4])]
	bytes := strings.Split(progStr, ",")
	if len(bytes)%2 != 0 {
		return nil, fmt.Errorf("program: got %d bytes want even number", len(bytes))
	}
	var program []operation
	for i := 0; i < len(bytes); i += 2 {
		program = append(program, operation{
			opcode:  opcode(int(mustParseInt(bytes[i]))),
			operand: operand(int(mustParseInt(bytes[i+1]))),
		})
	}
	return initComputer(a, b, c, program), nil
}
