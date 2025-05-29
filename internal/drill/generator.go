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

// generateEasy 生成简单题目(10以内加减法，2-5的乘法)
func (g *Generator) generateEasy() Question {
	// 随机选择运算类型：0-加法，1-减法，2-乘法
	opType := g.rng.Intn(3)

	switch opType {
	case 0: // 加法
		a := g.rng.Intn(10) + 1 // 1-10
		b := g.rng.Intn(10) + 1 // 1-10
		return Question{
			Expression: fmt.Sprintf("%d + %d", a, b),
			Answer:     a + b,
			Difficulty: Easy,
		}
	case 1: // 减法
		a := g.rng.Intn(10) + 1 // 1-10
		b := g.rng.Intn(10) + 1 // 1-10
		if a < b {
			a, b = b, a
		}
		return Question{
			Expression: fmt.Sprintf("%d - %d", a, b),
			Answer:     a - b,
			Difficulty: Easy,
		}
	case 2: // 乘法
		a := g.rng.Intn(4) + 2 // 2-5
		b := g.rng.Intn(4) + 2 // 2-5
		return Question{
			Expression: fmt.Sprintf("%d × %d", a, b),
			Answer:     a * b,
			Difficulty: Easy,
		}
	default:
		return g.generateEasy()
	}
}

// generateMedium 生成中等题目(两位数加减法，乘法表扩展，简单除法)
func (g *Generator) generateMedium() Question {
	// 随机选择运算类型：0-加法，1-减法，2-乘法，3-除法，4-混合运算
	opType := g.rng.Intn(5)

	switch opType {
	case 0: // 加法
		a := g.rng.Intn(50) + 1 // 1-50
		b := g.rng.Intn(50) + 1 // 1-50
		return Question{
			Expression: fmt.Sprintf("%d + %d", a, b),
			Answer:     a + b,
			Difficulty: Medium,
		}
	case 1: // 减法
		a := g.rng.Intn(50) + 51 // 51-100
		b := g.rng.Intn(50) + 1  // 1-50
		return Question{
			Expression: fmt.Sprintf("%d - %d", a, b),
			Answer:     a - b,
			Difficulty: Medium,
		}
	case 2: // 乘法
		a := g.rng.Intn(9) + 2 // 2-10
		b := g.rng.Intn(9) + 2 // 2-10
		return Question{
			Expression: fmt.Sprintf("%d × %d", a, b),
			Answer:     a * b,
			Difficulty: Medium,
		}
	case 3: // 除法
		b := g.rng.Intn(9) + 2 // 2-10
		c := g.rng.Intn(9) + 2 // 2-10
		a := b * c             // 确保能整除
		return Question{
			Expression: fmt.Sprintf("%d ÷ %d", a, b),
			Answer:     c,
			Difficulty: Medium,
		}
	case 4: // 混合运算（先乘除后加减）
		a := g.rng.Intn(20) + 1 // 1-20
		b := g.rng.Intn(9) + 2  // 2-10
		c := g.rng.Intn(9) + 2  // 2-10
		op := []string{"+", "-"}[g.rng.Intn(2)]
		var answer int
		if op == "+" {
			answer = a + (b * c)
		} else {
			answer = a - (b * c)
		}
		return Question{
			Expression: fmt.Sprintf("%d %s %d × %d", a, op, b, c),
			Answer:     answer,
			Difficulty: Medium,
		}
	default:
		return g.generateEasy()
	}
}

// generateHard 生成困难题目(多步运算、大数运算、带括号运算)
func (g *Generator) generateHard() Question {
	// 随机选择运算类型：0-多步混合运算，1-带括号运算，2-大数运算
	opType := g.rng.Intn(3)

	switch opType {
	case 0: // 多步混合运算
		a := g.rng.Intn(20) + 1 // 1-20
		b := g.rng.Intn(9) + 2  // 2-10
		c := g.rng.Intn(9) + 2  // 2-10
		d := g.rng.Intn(20) + 1 // 1-20
		op1 := []string{"+", "-"}[g.rng.Intn(2)]
		op2 := []string{"+", "-"}[g.rng.Intn(2)]

		var answer int
		if op1 == "+" {
			if op2 == "+" {
				answer = a + (b * c) + d
			} else {
				answer = a + (b * c) - d
			}
		} else {
			if op2 == "+" {
				answer = a - (b * c) + d
			} else {
				answer = a - (b * c) - d
			}
		}

		return Question{
			Expression: fmt.Sprintf("%d %s %d × %d %s %d", a, op1, b, c, op2, d),
			Answer:     answer,
			Difficulty: Hard,
		}

	case 1: // 带括号运算
		a := g.rng.Intn(20) + 1 // 1-20
		b := g.rng.Intn(9) + 2  // 2-10
		c := g.rng.Intn(9) + 2  // 2-10
		d := g.rng.Intn(9) + 2  // 2-10
		op1 := []string{"+", "-"}[g.rng.Intn(2)]
		op2 := []string{"×", "÷"}[g.rng.Intn(2)]

		var answer int
		if op2 == "×" {
			if op1 == "+" {
				answer = a + (b * c * d)
			} else {
				answer = a - (b * c * d)
			}
		} else { // 除法
			// 确保能整除
			b = g.rng.Intn(5) + 2 // 2-6
			c = g.rng.Intn(5) + 2 // 2-6
			d = g.rng.Intn(5) + 2 // 2-6
			if op1 == "+" {
				answer = a + (b * c / d)
			} else {
				answer = a - (b * c / d)
			}
		}

		return Question{
			Expression: fmt.Sprintf("%d %s (%d × %d %s %d)", a, op1, b, c, op2, d),
			Answer:     answer,
			Difficulty: Hard,
		}

	case 2: // 大数运算
		a := g.rng.Intn(50) + 51 // 51-100
		b := g.rng.Intn(9) + 2   // 2-10
		c := g.rng.Intn(9) + 2   // 2-10
		d := g.rng.Intn(20) + 1  // 1-20
		op1 := []string{"+", "-"}[g.rng.Intn(2)]
		op2 := []string{"×", "÷"}[g.rng.Intn(2)]

		var answer int
		if op2 == "×" {
			if op1 == "+" {
				answer = a + (b * c * d)
			} else {
				answer = a - (b * c * d)
			}
		} else { // 除法
			// 确保能整除
			b = g.rng.Intn(5) + 2 // 2-6
			c = g.rng.Intn(5) + 2 // 2-6
			d = g.rng.Intn(5) + 2 // 2-6
			if op1 == "+" {
				answer = a + (b * c / d)
			} else {
				answer = a - (b * c / d)
			}
		}

		return Question{
			Expression: fmt.Sprintf("%d %s %d × %d %s %d", a, op1, b, c, op2, d),
			Answer:     answer,
			Difficulty: Hard,
		}

	default:
		return g.generateMedium()
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
