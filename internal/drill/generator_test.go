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

func TestGenerateEasy(t *testing.T) {
	g := NewGenerator()

	for i := 0; i < 100; i++ { // 多次测试确保随机性
		q := g.generateEasy()

		// 验证答案在0-18之间(10+10=20, 但减法最小为0)
		if q.Answer < 0 || q.Answer > 18 {
			t.Errorf("简单题目答案超出范围: %d", q.Answer)
		}
	}
}

func TestGenerateMedium(t *testing.T) {
	g := NewGenerator()

	for i := 0; i < 100; i++ {
		q := g.generateMedium()

		// 验证答案在合理范围内
		if q.Answer < 0 || q.Answer > 400 { // 20*20=400
			t.Errorf("中等题目答案超出范围: %d", q.Answer)
		}
	}
}

func TestGenerateHard(t *testing.T) {
	g := NewGenerator()

	for i := 0; i < 100; i++ {
		q := g.generateHard()

		// 验证答案在合理范围内
		if q.Answer < 0 || q.Answer > 10000 { // 100*100=10000
			t.Errorf("困难题目答案超出范围: %d", q.Answer)
		}
	}
}
