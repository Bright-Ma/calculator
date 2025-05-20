package calculator

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 题目难度级别
const (
	Easy   = "easy"   // 简单：个位数运算
	Medium = "medium" // 中等：两位数运算
	Hard   = "hard"   // 困难：多位数或多项运算
)

// Problem 表示一道口算题
type Problem struct {
	ID           string    `json:"id"`          // 题目ID
	Expression   string    `json:"expression"`  // 表达式文本
	Answer       float64   `json:"answer"`      // 正确答案
	Difficulty   string    `json:"difficulty"`  // 难度级别
	Operations   []string  `json:"operations"`  // 包含的运算类型
	TimeLimit    int       `json:"timeLimit"`   // 时间限制(秒)
	GeneratedAt  time.Time `json:"generatedAt"` // 生成时间
	InternalExpr string    `json:"-"`           // 内部表达式(用于计算)
}

// ProblemGenerator 题目生成器
type ProblemGenerator struct {
	factory      *OperationFactory
	random       *rand.Rand
	operations   []string
	difficulties []string
}

// NewProblemGenerator 创建新的题目生成器
func NewProblemGenerator(factory *OperationFactory) *ProblemGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &ProblemGenerator{
		factory:      factory,
		random:       rand.New(source),
		operations:   []string{"add", "subtract", "multiply", "divide"},
		difficulties: []string{Easy, Medium, Hard},
	}
}

// GenerateProblem 生成一道口算题
func (g *ProblemGenerator) GenerateProblem(difficulty string, operations []string, timeLimit int) (*Problem, error) {
	log.Printf("生成题目: 难度=%s, 运算类型=%v, 时间限制=%d秒", difficulty, operations, timeLimit)

	// 如果未指定运算类型，使用所有可用的运算
	if len(operations) == 0 {
		operations = g.operations
	}

	// 如果未指定难度，随机选择一个难度
	if difficulty == "" {
		difficulty = g.difficulties[g.random.Intn(len(g.difficulties))]
	}

	// 如果未指定时间限制，根据难度设置默认值
	if timeLimit <= 0 {
		switch difficulty {
		case Easy:
			timeLimit = 10
		case Medium:
			timeLimit = 20
		case Hard:
			timeLimit = 30
		default:
			timeLimit = 15
		}
	}

	// 根据难度生成题目
	var problem *Problem
	var err error

	switch difficulty {
	case Easy:
		problem, err = g.generateEasyProblem(operations)
	case Medium:
		problem, err = g.generateMediumProblem(operations)
	case Hard:
		problem, err = g.generateHardProblem(operations)
	default:
		problem, err = g.generateEasyProblem(operations)
	}

	if err != nil {
		return nil, err
	}

	// 设置题目属性
	problem.Difficulty = difficulty
	problem.Operations = operations
	problem.TimeLimit = timeLimit
	problem.GeneratedAt = time.Now()
	problem.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	log.Printf("题目生成成功: %s = %f", problem.Expression, problem.Answer)
	return problem, nil
}

// generateEasyProblem 生成简单题目(个位数运算)
func (g *ProblemGenerator) generateEasyProblem(operations []string) (*Problem, error) {
	// 随机选择一种运算
	opType := operations[g.random.Intn(len(operations))]

	// 生成两个个位数
	a := g.random.Intn(9) + 1 // 1-9
	b := g.random.Intn(9) + 1 // 1-9

	// 对于除法，确保能整除
	if opType == "divide" {
		a = b * (g.random.Intn(9) + 1)
	}

	return g.createProblem(opType, a, b)
}

// generateMediumProblem 生成中等难度题目(两位数运算)
func (g *ProblemGenerator) generateMediumProblem(operations []string) (*Problem, error) {
	// 随机选择一种运算
	opType := operations[g.random.Intn(len(operations))]

	// 生成两个两位数
	a := g.random.Intn(90) + 10 // 10-99
	b := g.random.Intn(90) + 10 // 10-99

	// 对于除法，确保能整除或结果简单
	if opType == "divide" {
		if g.random.Intn(2) == 0 {
			// 整除
			a = b * (g.random.Intn(9) + 1)
		} else {
			// 简单结果(保留一位小数)
			b = 2 * (g.random.Intn(5) + 1) // 2, 4, 6, 8, 10
			a = b*(g.random.Intn(9)+1) + b/2
		}
	}

	return g.createProblem(opType, a, b)
}

// generateHardProblem 生成困难题目(多位数或多项运算)
func (g *ProblemGenerator) generateHardProblem(operations []string) (*Problem, error) {
	// 决定是生成多位数运算还是多项运算
	if g.random.Intn(2) == 0 {
		// 多位数运算
		opType := operations[g.random.Intn(len(operations))]

		// 生成三位数
		a := g.random.Intn(900) + 100 // 100-999
		b := g.random.Intn(90) + 10   // 10-99

		// 对于除法，确保结果简单
		if opType == "divide" {
			a = b * (g.random.Intn(9) + 1)
		}

		return g.createProblem(opType, a, b)
	} else {
		// 多项运算(三个数)
		if len(operations) < 2 {
			operations = g.operations
		}

		// 随机选择两种不同的运算
		opIndex1 := g.random.Intn(len(operations))
		opIndex2 := (opIndex1 + 1 + g.random.Intn(len(operations)-1)) % len(operations)

		opType1 := operations[opIndex1]
		opType2 := operations[opIndex2]

		// 生成三个数
		a := g.random.Intn(9) + 1 // 1-9
		b := g.random.Intn(9) + 1 // 1-9
		c := g.random.Intn(9) + 1 // 1-9

		// 创建第一个操作
		op1, err := g.factory.Create(opType1)
		if err != nil {
			return nil, err
		}

		// 创建第二个操作
		op2, err := g.factory.Create(opType2)
		if err != nil {
			return nil, err
		}

		// 根据运算符优先级计算结果
		var result float64
		var expression, internalExpr string

		// 如果第二个运算符是乘除，先计算后面的部分
		if (opType2 == "multiply" || opType2 == "divide") &&
			(opType1 == "add" || opType1 == "subtract") {
			intermediate := op2.Calculate(float64(b), float64(c))
			result = op1.Calculate(float64(a), intermediate)
			expression = fmt.Sprintf("%d %s %d %s %d", a, op1.GetSymbol(), b, op2.GetSymbol(), c)
			internalExpr = fmt.Sprintf("%d %s (%d %s %d)", a, op1.GetSymbol(), b, op2.GetSymbol(), c)
		} else {
			// 否则从左到右计算
			intermediate := op1.Calculate(float64(a), float64(b))
			result = op2.Calculate(intermediate, float64(c))
			expression = fmt.Sprintf("%d %s %d %s %d", a, op1.GetSymbol(), b, op2.GetSymbol(), c)
			internalExpr = fmt.Sprintf("(%d %s %d) %s %d", a, op1.GetSymbol(), b, op2.GetSymbol(), c)
		}

		return &Problem{
			Expression:   expression,
			Answer:       result,
			InternalExpr: internalExpr,
		}, nil
	}
}

// createProblem 根据操作类型和操作数创建题目
func (g *ProblemGenerator) createProblem(opType string, a, b int) (*Problem, error) {
	operation, err := g.factory.Create(opType)
	if err != nil {
		return nil, err
	}

	result := operation.Calculate(float64(a), float64(b))
	expression := fmt.Sprintf("%d %s %d", a, operation.GetSymbol(), b)

	return &Problem{
		Expression:   expression,
		Answer:       result,
		InternalExpr: expression,
	}, nil
}

// CheckAnswer 检查答案是否正确
func (p *Problem) CheckAnswer(answer float64) bool {
	// 对于除法，允许小数点后两位的误差
	if p.Operations != nil && len(p.Operations) > 0 && p.Operations[0] == "divide" {
		return fmt.Sprintf("%.2f", p.Answer) == fmt.Sprintf("%.2f", answer)
	}
	return p.Answer == answer
}
