package calculator

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
	return a + b
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
	return a - b
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
	return a * b
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
	if b == 0 {
		// 避免除以零错误，返回特殊值
		return float64(0)
	}
	return a / b
}

func (o *DivideOperation) GetName() string {
	return "除法"
}

func (o *DivideOperation) GetSymbol() string {
	return "/"
}
