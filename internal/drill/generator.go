package drill

import (
	"math/rand"
	"strconv"
	"time"
)

// Difficulty 题目难度类型
type Difficulty int

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
)

// Question 表示一道口算题
type Question struct {
	Expression string // 表达式如 "3 + 5"
	Answer     int    // 正确答案
	Difficulty Difficulty
}

// Generator 口算题生成器
type Generator struct {
	rand *rand.Rand
}

// NewGenerator 创建新的题目生成器
func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate 根据难度生成题目
func (g *Generator) Generate(d Difficulty) Question {
	switch d {
	case Easy:
		return g.generateEasy()
	case Medium:
		return g.generateMedium()
	case Hard:
		return g.generateHard()
	default:
		return g.generateEasy()
	}
}

// generateEasy 生成简单题目(10以内加减法)
func (g *Generator) generateEasy() Question {
	a := g.rand.Intn(10)
	b := g.rand.Intn(10)
	op := g.rand.Intn(2) // 0:加, 1:减

	if op == 0 {
		return Question{
			Expression: formatExpression(a, b, "+"),
			Answer:     a + b,
			Difficulty: Easy,
		}
	}

	// 确保减法结果非负
	if a < b {
		a, b = b, a
	}
	return Question{
		Expression: formatExpression(a, b, "-"),
		Answer:     a - b,
		Difficulty: Easy,
	}
}

// generateMedium 生成中等题目(两位数加减法,乘法表扩展)
func (g *Generator) generateMedium() Question {
	opType := g.rand.Intn(3) // 0:加, 1:减, 2:乘

	switch opType {
	case 0: // 加法(20-99)
		a := g.rand.Intn(80) + 20
		b := g.rand.Intn(80) + 20
		return Question{
			Expression: formatExpression(a, b, "+"),
			Answer:     a + b,
			Difficulty: Medium,
		}
	case 1: // 减法(20-99)
		a := g.rand.Intn(80) + 20
		b := g.rand.Intn(80) + 20
		if a < b {
			a, b = b, a
		}
		return Question{
			Expression: formatExpression(a, b, "-"),
			Answer:     a - b,
			Difficulty: Medium,
		}
	default: // 乘法(0-12乘法表)
		a := g.rand.Intn(13)
		b := g.rand.Intn(13)
		return Question{
			Expression: formatExpression(a, b, "×"),
			Answer:     a * b,
			Difficulty: Medium,
		}
	}
}

// generateHard 生成困难题目(多步运算、大数运算、带括号运算)
func (g *Generator) generateHard() Question {
	opType := g.rand.Intn(5) // 0:大数加法, 1:大数减法, 2:复杂乘法, 3:复杂除法, 4:多步运算

	switch opType {
	case 0: // 大数加法(100-999)
		a := g.rand.Intn(900) + 100
		b := g.rand.Intn(900) + 100
		return Question{
			Expression: formatExpression(a, b, "+"),
			Answer:     a + b,
			Difficulty: Hard,
		}
	case 1: // 大数减法(100-999)
		a := g.rand.Intn(900) + 100
		b := g.rand.Intn(900) + 100
		if a < b {
			a, b = b, a
		}
		return Question{
			Expression: formatExpression(a, b, "-"),
			Answer:     a - b,
			Difficulty: Hard,
		}
	case 2: // 复杂乘法(两位数×两位数)
		a := g.rand.Intn(90) + 10
		b := g.rand.Intn(90) + 10
		return Question{
			Expression: formatExpression(a, b, "×"),
			Answer:     a * b,
			Difficulty: Hard,
		}
	case 3: // 复杂除法(确保能整除)
		b := g.rand.Intn(12) + 1       // 除数1-12
		a := b * (g.rand.Intn(20) + 1) // 被除数是除数的1-20倍
		return Question{
			Expression: formatExpression(a, b, "÷"),
			Answer:     a / b,
			Difficulty: Hard,
		}
	default: // 多步运算(强制先乘除后加减)
		switch g.rand.Intn(3) {
		case 0: // 形式: a + b × c
			a := g.rand.Intn(10) + 1
			b := g.rand.Intn(10) + 1
			c := g.rand.Intn(10) + 1
			return Question{
				Expression: strconv.Itoa(a) + " + " + strconv.Itoa(b) + " × " + strconv.Itoa(c),
				Answer:     a + b*c,
				Difficulty: Hard,
			}
		case 1: // 形式: (a + b) × c
			a := g.rand.Intn(10) + 1
			b := g.rand.Intn(10) + 1
			c := g.rand.Intn(10) + 1
			return Question{
				Expression: "(" + strconv.Itoa(a) + " + " + strconv.Itoa(b) + ") × " + strconv.Itoa(c),
				Answer:     (a + b) * c,
				Difficulty: Hard,
			}
		default: // 形式: a × b + c × d
			a := g.rand.Intn(10) + 1
			b := g.rand.Intn(10) + 1
			c := g.rand.Intn(10) + 1
			d := g.rand.Intn(10) + 1
			return Question{
				Expression: strconv.Itoa(a) + " × " + strconv.Itoa(b) + " + " + strconv.Itoa(c) + " × " + strconv.Itoa(d),
				Answer:     a*b + c*d,
				Difficulty: Hard,
			}
		}
	}
}

// evalMultiStep 计算多步运算结果
func evalMultiStep(a, b, c int, op1, op2 string) int {
	var first int
	switch op1 {
	case "+":
		first = a + b
	case "-":
		first = a - b
	}

	switch op2 {
	case "×":
		return first * c
	case "÷":
		return first / c
	}
	return 0
}

// formatExpression 格式化表达式字符串
func formatExpression(a, b int, op string) string {
	return strconv.Itoa(a) + " " + op + " " + strconv.Itoa(b)
}
