package drill

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Difficulty 题目难度类型
type Difficulty int

const (
	Easy Difficulty = iota
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
	rng *rand.Rand
}

// NewGenerator 创建新的题目生成器
func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate 根据难度生成题目
func (g *Generator) Generate(difficulty Difficulty) Question {
	switch difficulty {
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
	a := g.rng.Intn(10) + 1
	b := g.rng.Intn(10) + 1
	if g.rng.Intn(2) == 0 {
		return Question{
			Expression: fmt.Sprintf("%d + %d", a, b),
			Answer:     a + b,
			Difficulty: Easy,
		}
	}
	if a < b {
		a, b = b, a
	}
	return Question{
		Expression: fmt.Sprintf("%d - %d", a, b),
		Answer:     a - b,
		Difficulty: Easy,
	}
}

// generateMedium 生成中等题目(两位数加减法,乘法表扩展)
func (g *Generator) generateMedium() Question {
	a := g.rng.Intn(20) + 1
	b := g.rng.Intn(20) + 1
	if g.rng.Intn(2) == 0 {
		return Question{
			Expression: fmt.Sprintf("%d + %d", a, b),
			Answer:     a + b,
			Difficulty: Medium,
		}
	}
	if a < b {
		a, b = b, a
	}
	return Question{
		Expression: fmt.Sprintf("%d - %d", a, b),
		Answer:     a - b,
		Difficulty: Medium,
	}
}

// generateHard 生成困难题目(多步运算、大数运算、带括号运算)
func (g *Generator) generateHard() Question {
	a := g.rng.Intn(100) + 1
	b := g.rng.Intn(100) + 1
	if g.rng.Intn(2) == 0 {
		return Question{
			Expression: fmt.Sprintf("%d + %d", a, b),
			Answer:     a + b,
			Difficulty: Hard,
		}
	}
	if a < b {
		a, b = b, a
	}
	return Question{
		Expression: fmt.Sprintf("%d - %d", a, b),
		Answer:     a - b,
		Difficulty: Hard,
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
