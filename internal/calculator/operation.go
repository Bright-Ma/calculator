package calculator

import "log"

// Operation 定义计算器操作的接口
type Operation interface {
	// Calculate 执行计算操作
	Calculate(a, b float64) float64
	// GetName 获取操作名称
	GetName() string
	// GetSymbol 获取操作符号
	GetSymbol() string
}

// AddOperation 加法操作
type AddOperation struct{}

func (o *AddOperation) Calculate(a, b float64) float64 {
	log.Printf("执行加法运算: %f + %f", a, b)
	result := a + b
	log.Printf("加法运算结果: %f", result)
	return result
}

func (o *AddOperation) GetName() string {
	return "加法"
}

func (o *AddOperation) GetSymbol() string {
	return "+"
}

// SubtractOperation 减法操作
type SubtractOperation struct{}

func (o *SubtractOperation) Calculate(a, b float64) float64 {
	log.Printf("执行减法运算: %f - %f", a, b)
	result := a - b
	log.Printf("减法运算结果: %f", result)
	return result
}

func (o *SubtractOperation) GetName() string {
	return "减法"
}

func (o *SubtractOperation) GetSymbol() string {
	return "-"
}

// MultiplyOperation 乘法操作
type MultiplyOperation struct{}

func (o *MultiplyOperation) Calculate(a, b float64) float64 {
	log.Printf("执行乘法运算: %f * %f", a, b)
	result := a * b
	log.Printf("乘法运算结果: %f", result)
	return result
}

func (o *MultiplyOperation) GetName() string {
	return "乘法"
}

func (o *MultiplyOperation) GetSymbol() string {
	return "*"
}

// DivideOperation 除法操作
type DivideOperation struct{}

func (o *DivideOperation) Calculate(a, b float64) float64 {
	log.Printf("执行除法运算: %f / %f", a, b)
	if b == 0 {
		log.Printf("警告: 除数为零")
		return float64(0)
	}
	result := a / b
	log.Printf("除法运算结果: %f", result)
	return result
}

func (o *DivideOperation) GetName() string {
	return "除法"
}

func (o *DivideOperation) GetSymbol() string {
	return "/"
}
