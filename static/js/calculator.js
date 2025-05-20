// 全局状态
const state = {
    currentProblem: null,
    problemStartTime: null,
    totalProblems: 0,
    correctProblems: 0,
    totalTime: 0,
    practiceStartTime: null,
    selectedDifficulty: null,
    selectedOperations: [],
    timer: null,
    problemTimer: null
};

// DOM元素
const elements = {
    settingsPanel: document.getElementById('settings-panel'),
    practicePanel: document.getElementById('practice-panel'),
    resultsPanel: document.getElementById('results-panel'),
    restReminder: document.getElementById('rest-reminder'),
    problemText: document.getElementById('problem-text'),
    answerInput: document.getElementById('answer-input'),
    submitAnswer: document.getElementById('submit-answer'),
    nextProblem: document.getElementById('next-problem'),
    problemCount: document.getElementById('problem-count'),
    correctCount: document.getElementById('correct-count'),
    timer: document.getElementById('timer'),
    timeLimitProgress: document.getElementById('time-limit-progress'),
    resultMessage: document.getElementById('result-message'),
    difficultyButtons: document.querySelectorAll('.difficulty-btn'),
    operationCheckboxes: document.querySelectorAll('.operation-checkbox'),
    startPractice: document.getElementById('start-practice'),
    endPractice: document.getElementById('end-practice'),
    returnToSettings: document.getElementById('return-to-settings'),
    continueButton: document.getElementById('continue-practice'),
    restButton: document.getElementById('take-rest'),
    totalProblemsDisplay: document.getElementById('total-problems'),
    totalCorrectDisplay: document.getElementById('total-correct'),
    accuracyRateDisplay: document.getElementById('accuracy-rate'),
    avgTimeDisplay: document.getElementById('avg-time'),
    performanceMessage: document.getElementById('performance-message'),
    audioContext: new (window.AudioContext || window.webkitAudioContext)()
};

// 初始化事件监听器
function initializeEventListeners() {
    // 难度选择
    elements.difficultyButtons.forEach(button => {
        button.addEventListener('click', () => {
            elements.difficultyButtons.forEach(btn => btn.classList.remove('selected'));
            button.classList.add('selected');
            state.selectedDifficulty = button.dataset.level;
        });
    });

    // 开始练习
    elements.startPractice.addEventListener('click', startPractice);

    // 提交答案
    elements.submitAnswer.addEventListener('click', submitAnswer);
    elements.answerInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            submitAnswer();
        }
    });

    // 下一题
    elements.nextProblem.addEventListener('click', getNextProblem);

    // 结束练习
    elements.endPractice.addEventListener('click', endPractice);

    // 返回设置
    elements.returnToSettings.addEventListener('click', returnToSettings);

    // 休息提醒
    elements.continueButton.addEventListener('click', () => {
        elements.restReminder.classList.add('hidden');
    });

    elements.restButton.addEventListener('click', () => {
        elements.restReminder.classList.add('hidden');
        endPractice();
    });

    // 键盘快捷键支持
    document.addEventListener('keydown', handleKeyboardShortcuts);
}

// 键盘快捷键处理
function handleKeyboardShortcuts(e) {
    // ESC键 - 结束练习
    if (e.key === 'Escape' && !elements.practicePanel.classList.contains('hidden')) {
        endPractice();
    }
    
    // 空格键 - 下一题
    if (e.key === ' ' && !elements.nextProblem.classList.contains('hidden')) {
        e.preventDefault();
        getNextProblem();
    }
}

// 开始练习
async function startPractice() {
    // 验证选择
    if (!state.selectedDifficulty) {
        showToast('请选择难度级别！');
        return;
    }

    // 获取选中的运算类型
    state.selectedOperations = Array.from(elements.operationCheckboxes)
        .filter(cb => cb.checked)
        .map(cb => cb.value);

    if (state.selectedOperations.length === 0) {
        showToast('请至少选择一种运算类型！');
        return;
    }

    // 重置状态
    state.totalProblems = 0;
    state.correctProblems = 0;
    state.totalTime = 0;
    state.practiceStartTime = new Date();

    // 切换面板
    elements.settingsPanel.classList.add('hidden');
    elements.practicePanel.classList.remove('hidden');

    // 开始计时器
    startTimer();

    // 获取第一道题
    await getNextProblem();
}

// 获取下一道题
async function getNextProblem() {
    try {
        const response = await fetch('/api/problem/new', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                difficulty: state.selectedDifficulty,
                operations: state.selectedOperations
            })
        });

        const data = await response.json();
        if (data.error) {
            throw new Error(data.error);
        }

        // 更新状态
        state.currentProblem = data.problem;
        state.problemStartTime = new Date();

        // 更新界面
        elements.problemText.textContent = data.problem.expression;
        elements.problemText.classList.add('animate__animated', 'animate__fadeIn');
        elements.answerInput.value = '';
        elements.answerInput.disabled = false;
        elements.submitAnswer.disabled = false;
        elements.nextProblem.classList.add('hidden');
        elements.resultMessage.textContent = '';
        elements.resultMessage.className = 'result-message';

        // 启动题目计时器
        startProblemTimer(data.problem.timeLimit);

        // 更新统计
        elements.problemCount.textContent = state.totalProblems + 1;

        // 聚焦到输入框
        elements.answerInput.focus();

    } catch (error) {
        console.error('获取题目失败:', error);
        showToast('获取题目失败，请重试！');
    }
}

// 提交答案
async function submitAnswer() {
    if (!state.currentProblem) return;

    const answer = parseFloat(elements.answerInput.value);
    if (isNaN(answer)) {
        showToast('请输入有效的数字！');
        return;
    }

    const timeSpent = Math.round((new Date() - state.problemStartTime) / 1000);

    try {
        const response = await fetch('/api/problem/answer', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                problemId: state.currentProblem.id,
                answer: answer,
                timeSpent: timeSpent
            })
        });

        const data = await response.json();
        if (data.error) {
            throw new Error(data.error);
        }

        // 停止题目计时器
        if (state.problemTimer) {
            clearInterval(state.problemTimer);
        }

        // 更新统计
        state.totalProblems++;
        state.totalTime += timeSpent;

        if (data.correct) {
            state.correctProblems++;
            elements.correctCount.textContent = state.correctProblems;
            showCorrectAnimation();
            playCorrectSound();
        } else {
            showWrongAnimation(data.correctAnswer);
            playWrongSound();
        }

        // 禁用输入
        elements.answerInput.disabled = true;
        elements.submitAnswer.disabled = true;
        elements.nextProblem.classList.remove('hidden');

        // 检查是否需要休息
        if (data.needRest) {
            showRestReminder();
        }

    } catch (error) {
        console.error('提交答案失败:', error);
        showToast('提交答案失败，请重试！');
    }
}

// 显示正确动画
function showCorrectAnimation() {
    elements.resultMessage.textContent = '答对了！👍';
    elements.resultMessage.className = 'result-message correct';
    elements.problemText.classList.add('correct-animation');
    setTimeout(() => {
        elements.problemText.classList.remove('correct-animation');
    }, 500);
}

// 显示错误动画
function showWrongAnimation(correctAnswer) {
    elements.resultMessage.textContent = `答错了！正确答案是: ${correctAnswer}`;
    elements.resultMessage.className = 'result-message wrong';
    elements.problemText.classList.add('wrong-animation');
    setTimeout(() => {
        elements.problemText.classList.remove('wrong-animation');
    }, 500);
}

// 显示休息提醒
function showRestReminder() {
    elements.restReminder.classList.remove('hidden');
}

// 开始计时器
function startTimer() {
    if (state.timer) {
        clearInterval(state.timer);
    }

    state.timer = setInterval(() => {
        const totalSeconds = Math.round((new Date() - state.practiceStartTime) / 1000);
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        elements.timer.textContent = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
    }, 1000);
}

// 开始题目计时器
function startProblemTimer(timeLimit) {
    if (state.problemTimer) {
        clearInterval(state.problemTimer);
    }

    const startTime = new Date();
    elements.timeLimitProgress.style.width = '100%';

    state.problemTimer = setInterval(() => {
        const elapsed = (new Date() - startTime) / 1000;
        const remaining = timeLimit - elapsed;
        const percentage = (remaining / timeLimit) * 100;

        if (percentage <= 0) {
            clearInterval(state.problemTimer);
            elements.timeLimitProgress.style.width = '0%';
            submitAnswer();
        } else {
            elements.timeLimitProgress.style.width = `${percentage}%`;
            // 改变颜色提示
            if (percentage < 30) {
                elements.timeLimitProgress.style.backgroundColor = '#dc3545';
            } else if (percentage < 60) {
                elements.timeLimitProgress.style.backgroundColor = '#ffc107';
            }
        }
    }, 100);
}

// 结束练习
function endPractice() {
    // 停止计时器
    if (state.timer) {
        clearInterval(state.timer);
    }
    if (state.problemTimer) {
        clearInterval(state.problemTimer);
    }

    // 播放完成音效
    playCompleteSound();

    // 计算统计数据
    const accuracyRate = state.totalProblems === 0 ? 0 : 
        Math.round((state.correctProblems / state.totalProblems) * 100);
    const avgTime = state.totalProblems === 0 ? 0 : 
        Math.round(state.totalTime / state.totalProblems);

    // 更新结果面板
    elements.totalProblemsDisplay.textContent = state.totalProblems;
    elements.totalCorrectDisplay.textContent = state.correctProblems;
    elements.accuracyRateDisplay.textContent = `${accuracyRate}%`;
    elements.avgTimeDisplay.textContent = `${avgTime}秒`;

    // 生成表现评价
    let performanceMessage = '';
    if (accuracyRate >= 90) {
        performanceMessage = '太棒了！你是小计算天才！🌟';
    } else if (accuracyRate >= 70) {
        performanceMessage = '做得不错！继续加油！👍';
    } else if (accuracyRate >= 50) {
        performanceMessage = '还需要多加练习，你可以的！💪';
    } else {
        performanceMessage = '不要灰心，慢慢来，重在参与！🌈';
    }
    elements.performanceMessage.textContent = performanceMessage;

    // 切换面板
    elements.practicePanel.classList.add('hidden');
    elements.resultsPanel.classList.remove('hidden');
}

// 返回设置
function returnToSettings() {
    elements.resultsPanel.classList.add('hidden');
    elements.settingsPanel.classList.remove('hidden');
    // 重置选择
    elements.difficultyButtons.forEach(btn => btn.classList.remove('selected'));
    state.selectedDifficulty = null;
}

// 显示提示消息
function showToast(message) {
    const toast = document.createElement('div');
    toast.className = 'toast animate__animated animate__fadeIn';
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.classList.remove('animate__fadeIn');
        toast.classList.add('animate__fadeOut');
        setTimeout(() => {
            document.body.removeChild(toast);
        }, 500);
    }, 2000);
}

// 音效函数
function playCorrectSound() {
    const oscillator = elements.audioContext.createOscillator();
    const gainNode = elements.audioContext.createGain();
    
    oscillator.connect(gainNode);
    gainNode.connect(elements.audioContext.destination);
    
    // 正确音效：上升的音调
    oscillator.frequency.setValueAtTime(600, elements.audioContext.currentTime);
    oscillator.frequency.linearRampToValueAtTime(800, elements.audioContext.currentTime + 0.1);
    
    gainNode.gain.setValueAtTime(0.3, elements.audioContext.currentTime);
    gainNode.gain.linearRampToValueAtTime(0, elements.audioContext.currentTime + 0.2);
    
    oscillator.start();
    oscillator.stop(elements.audioContext.currentTime + 0.2);
}

function playWrongSound() {
    const oscillator = elements.audioContext.createOscillator();
    const gainNode = elements.audioContext.createGain();
    
    oscillator.connect(gainNode);
    gainNode.connect(elements.audioContext.destination);
    
    // 错误音效：下降的音调
    oscillator.frequency.setValueAtTime(400, elements.audioContext.currentTime);
    oscillator.frequency.linearRampToValueAtTime(200, elements.audioContext.currentTime + 0.2);
    
    gainNode.gain.setValueAtTime(0.3, elements.audioContext.currentTime);
    gainNode.gain.linearRampToValueAtTime(0, elements.audioContext.currentTime + 0.2);
    
    oscillator.start();
    oscillator.stop(elements.audioContext.currentTime + 0.2);
}

function playCompleteSound() {
    const oscillator = elements.audioContext.createOscillator();
    const gainNode = elements.audioContext.createGain();
    
    oscillator.connect(gainNode);
    gainNode.connect(elements.audioContext.destination);
    
    // 完成音效：欢快的音阶
    const notes = [400, 500, 600, 800];
    const noteLength = 0.1;
    
    notes.forEach((freq, index) => {
        oscillator.frequency.setValueAtTime(freq, elements.audioContext.currentTime + index * noteLength);
        gainNode.gain.setValueAtTime(0.3, elements.audioContext.currentTime + index * noteLength);
        gainNode.gain.setValueAtTime(0, elements.audioContext.currentTime + (index + 0.8) * noteLength);
    });
    
    oscillator.start();
    oscillator.stop(elements.audioContext.currentTime + notes.length * noteLength);
}

// 初始化音频上下文的函数
function initAudioContext() {
    // 某些浏览器需要用户交互后才能创建AudioContext
    if (!elements.audioContext) {
        elements.audioContext = new (window.AudioContext || window.webkitAudioContext)();
    }
    // 如果音频上下文被挂起，则恢复它
    if (elements.audioContext.state === 'suspended') {
        elements.audioContext.resume();
    }
}

// 在用户第一次交互时初始化音频上下文
document.addEventListener('click', function initAudio() {
    initAudioContext();
    document.removeEventListener('click', initAudio);
}, { once: true });

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    initializeEventListeners();
});