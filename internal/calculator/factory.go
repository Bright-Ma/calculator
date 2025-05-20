package calculator

import (
	"fmt"
	"log"
	"reflect"
)

// OperationFactory 使用反射机制创建操作对象的工厂
type OperationFactory struct {
	// 注册的操作类型映射
	operationTypes map[string]reflect.Type
}

// NewOperationFactory 创建一个新的操作工厂
func NewOperationFactory() *OperationFactory {
	factory := &OperationFactory{
		operationTypes: make(map[string]reflect.Type),
	}

	// 注册默认的操作类型
	factory.Register("add", reflect.TypeOf(AddOperation{}))
	factory.Register("subtract", reflect.TypeOf(SubtractOperation{}))
	factory.Register("multiply", reflect.TypeOf(MultiplyOperation{}))
	factory.Register("divide", reflect.TypeOf(DivideOperation{}))

	return factory
}

// Register 注册一个操作类型
func (f *OperationFactory) Register(name string, operationType reflect.Type) {
	f.operationTypes[name] = operationType
}

// Create 使用反射机制创建一个操作对象
func (f *OperationFactory) Create(name string) (Operation, error) {
	log.Printf("尝试创建操作: %s", name)

	operationType, exists := f.operationTypes[name]
	if !exists {
		log.Printf("错误: 未知的操作类型: %s", name)
		return nil, fmt.Errorf("未知的操作类型: %s", name)
	}

	log.Printf("找到操作类型: %s", operationType.Name())

	// 使用反射创建对象实例
	operationValue := reflect.New(operationType)
	log.Printf("使用反射创建对象: %v", operationValue.Type())

	operation, ok := operationValue.Interface().(Operation)
	if !ok {
		log.Printf("错误: 类型转换失败: %s", operationType.Name())
		return nil, fmt.Errorf("无法将类型 %s 转换为 Operation", operationType.Name())
	}

	log.Printf("成功创建操作对象: %s", operation.GetName())
	return operation, nil
}

// GetAvailableOperations 获取所有可用的操作名称
func (f *OperationFactory) GetAvailableOperations() []string {
	operations := make([]string, 0, len(f.operationTypes))
	for name := range f.operationTypes {
		operations = append(operations, name)
	}
	return operations
}
