package drill

import (
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	g := NewGenerator()

	tests := []struct {
		name       string
		difficulty Difficulty
	}{
		{"简单题目", Easy},
		{"中等题目", Medium},
		{"困难题目", Hard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question := g.Generate(tt.difficulty)

			// 验证表达式不为空
			if question.Expression == "" {
				t.Errorf("表达式为空")
			}

			// 验证难度级别匹配
			if question.Difficulty != tt.difficulty {
				t.Errorf("难度不匹配, 期望 %d, 实际 %d", tt.difficulty, question.Difficulty)
			}

			// 简单验证表达式格式
			if len(question.Expression) < 5 { // 如 "1 + 2" 是最短长度
				t.Errorf("表达式格式不正确: %s", question.Expression)
			}
		})
	}
}
