document.addEventListener('DOMContentLoaded', () => {
    const difficultySelect = document.getElementById('difficulty');
    const generateBtn = document.getElementById('generate');
    const questionContainer = document.getElementById('question-container');
    const questionElement = document.getElementById('question');
    const answerInput = document.getElementById('answer');
    const submitBtn = document.getElementById('submit');
    const resultElement = document.getElementById('result');

    let currentQuestion = null;

    // 生成题目
    generateBtn.addEventListener('click', async () => {
        const difficultyStr = difficultySelect.value;
        // 将字符串难度转换为数字: easy=1, medium=2, hard=3
        const difficultyMap = {
            'easy': 1,
            'medium': 2,
            'hard': 3
        };
        const difficulty = difficultyMap[difficultyStr] || 1;
        console.log('Selected difficulty:', difficultyStr, 'Mapped to:', difficulty);
        
        try {
            const response = await fetch(`http://localhost:8080/api/questions?difficulty=${difficulty}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            console.log('API response:', data);
            
            if (!data || typeof data.question !== 'string' || typeof data.answer !== 'number') {
                throw new Error('Invalid response format');
            }
            
            currentQuestion = {
                question: data.question,
                answer: data.answer
            };
            
            questionElement.innerText = data.question;
            answerInput.value = '';
            resultElement.textContent = '';
            resultElement.className = '';
            questionContainer.classList.remove('hidden');
        } catch (error) {
            console.error('Error fetching question:', error);
            questionElement.innerText = 'Error loading question';
            resultElement.textContent = '错误: ' + error.message;
            resultElement.className = 'incorrect';
        }
    });

    // 提交答案
    submitBtn.addEventListener('click', () => {
        if (!currentQuestion) return;
        
        const userAnswer = parseInt(answerInput.value);
        if (isNaN(userAnswer)) {
            resultElement.textContent = '请输入有效的数字答案';
            resultElement.className = 'incorrect';
            return;
        }

        if (userAnswer === currentQuestion.answer) {
            resultElement.textContent = '✅ 回答正确！';
            resultElement.className = 'correct';
        } else {
            resultElement.textContent = `❌ 回答错误！正确答案是 ${currentQuestion.answer}`;
            resultElement.className = 'incorrect';
        }
    });

    // 按Enter键提交答案
    answerInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            submitBtn.click();
        }
    });
});