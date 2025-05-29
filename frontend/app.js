document.addEventListener('DOMContentLoaded', () => {
    // 登录注册相关元素
    const authContainer = document.getElementById('auth-container');
    const appContainer = document.getElementById('app');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const loginUsername = document.getElementById('login-username');
    const loginPassword = document.getElementById('login-password');
    const registerUsername = document.getElementById('register-username');
    const registerPassword = document.getElementById('register-password');
    const registerConfirmPassword = document.getElementById('register-confirm-password');
    const showRegister = document.getElementById('show-register');
    const showLogin = document.getElementById('show-login');
    const logoutBtn = document.getElementById('logout-btn');
    const usernameDisplay = document.getElementById('username-display');

    // 口算题相关元素
    const difficultySelect = document.getElementById('difficulty');
    const generateBtn = document.getElementById('generate');
    const questionContainer = document.getElementById('question-container');
    const questionElement = document.getElementById('question');
    const answerInput = document.getElementById('answer');
    const submitBtn = document.getElementById('submit');
    const resultElement = document.getElementById('result');

    // 统计相关元素
    const statsContainer = document.getElementById('stats-container');
    const totalQuestionsElement = document.getElementById('total-questions');
    const easyQuestionsElement = document.getElementById('easy-questions');
    const mediumQuestionsElement = document.getElementById('medium-questions');
    const hardQuestionsElement = document.getElementById('hard-questions');
    const totalAttemptsElement = document.getElementById('total-attempts');
    const correctAnswersElement = document.getElementById('correct-answers');
    const accuracyElement = document.getElementById('accuracy');

    // 热度排行榜相关元素
    const hotRankBtn = document.getElementById('hot-rank-btn');
    const hotRankModal = document.getElementById('hot-rank-modal');
    const closeBtn = document.querySelector('.close-btn');
    const rankTabs = document.querySelectorAll('.tab-btn');
    const rankListBody = document.getElementById('rank-list-body');

    let currentQuestion = null;
    let currentUser = null;

    // 热度排行榜相关函数
    let currentRankType = 'hourly'; // 默认显示小时榜

    function showHotRankModal() {
        hotRankModal.style.display = 'block';
        loadRankingData(currentRankType);
    }

    function closeHotRankModal() {
        hotRankModal.style.display = 'none';
    }

    function switchTab(type) {
        currentRankType = type;
        // 更新按钮状态
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        event.target.classList.add('active');
        // 加载对应类型的排行榜数据
        loadRankingData(type);
    }

    async function loadRankingData(type) {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('请先登录');
            }

            const data = await apiRequest(`/api/drill/rankings?type=${type}`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            // 清空现有数据
            const tbody = document.getElementById('rankingTableBody');
            tbody.innerHTML = '';

            // 添加新数据
            data.rankings.forEach(rank => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${rank.rank}</td>
                    <td>${rank.username}</td>
                    <td style="text-align: center;">${rank.hot_score.toFixed(1)}</td>
                `;
                tbody.appendChild(row);
            });
        } catch (error) {
            console.error('获取排行榜失败:', error);
            alert(error.message || '获取排行榜失败，请重试');
        }
    }

    // 表单切换
    showRegister.addEventListener('click', (e) => {
        e.preventDefault();
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
    });

    showLogin.addEventListener('click', (e) => {
        e.preventDefault();
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    });

    // 通用API请求函数
    async function apiRequest(url, options = {}) {
        try {
            // 确保 URL 以 / 开头
            const apiUrl = url.startsWith('/') ? url : `/${url}`;
            console.log('发送请求:', apiUrl, options);
            
            const defaultHeaders = {
                'Content-Type': 'application/json'
            };
            
            // 添加认证头
            const token = localStorage.getItem('token');
            if (token) {
                defaultHeaders['Authorization'] = `Bearer ${token}`;
            }

            const response = await fetch(apiUrl, {
                ...options,
                headers: {
                    ...defaultHeaders,
                    ...(options.headers || {})
                }
            });

            console.log('收到响应:', response.status, response.statusText);
            const data = await response.json();
            console.log('响应数据:', data);

            if (!response.ok) {
                // 处理401未授权
                if (response.status === 401) {
                    clearAuthData();
                    appContainer.style.display = 'none';
                    authContainer.style.display = 'block';
                    throw new Error('会话已过期，请重新登录');
                }
                // 直接抛出服务器返回的错误信息
                throw new Error(data.error || `请求失败: ${response.status}`);
            }

            return data;
        } catch (error) {
            console.error('请求错误:', error);
            // 如果是网络错误，才显示网络错误提示
            if (error.name === 'TypeError' && error.message.includes('Failed to fetch')) {
                throw new Error('网络错误，请检查网络连接');
            }
            // 其他错误直接抛出
            throw error;
        }
    }

    // 登录功能
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = loginUsername.value;
        const password = loginPassword.value;

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });

            const data = await response.json();
            
            // 检查响应状态
            if (response.ok && data.token) {
                // 保存token
                localStorage.setItem('token', data.token);
                // 保存用户名
                localStorage.setItem('username', username);
                // 显示用户名
                usernameDisplay.textContent = username;
                // 切换到主界面
                authContainer.style.display = 'none';
                appContainer.style.display = 'block';
                // 清空表单
                loginForm.reset();
            }
        } catch (error) {
            console.error('登录错误:', error);
        }
    });

    // 注册功能
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = registerUsername.value.trim();
        const password = registerPassword.value.trim();
        const confirmPassword = registerConfirmPassword.value.trim();
        const role = document.getElementById('register-role').value;

        // 移除之前的错误提示
        const oldError = registerForm.querySelector('.error-message');
        if (oldError) {
            oldError.remove();
        }

        if (!username || !password) {
            const errorMessage = document.createElement('div');
            errorMessage.textContent = '请输入用户名和密码';
            errorMessage.className = 'error-message';
            errorMessage.style.color = '#e74c3c';
            errorMessage.style.fontSize = '0.9em';
            errorMessage.style.marginTop = '5px';
            errorMessage.style.textAlign = 'center';
            registerForm.appendChild(errorMessage);
            return;
        }

        if (password !== confirmPassword) {
            const errorMessage = document.createElement('div');
            errorMessage.textContent = '两次输入的密码不一致';
            errorMessage.className = 'error-message';
            errorMessage.style.color = '#e74c3c';
            errorMessage.style.fontSize = '0.9em';
            errorMessage.style.marginTop = '5px';
            errorMessage.style.textAlign = 'center';
            registerForm.appendChild(errorMessage);
            return;
        }

        try {
            await apiRequest('/api/auth/register', {
                method: 'POST',
                body: JSON.stringify({ username, password, role })
            });

            // 注册成功后切换到登录表单
            registerForm.style.display = 'none';
            loginForm.style.display = 'block';
            registerUsername.value = '';
            registerPassword.value = '';
            registerConfirmPassword.value = '';
            
            // 在登录表单上方显示注册成功提示
            const successMessage = document.createElement('div');
            successMessage.textContent = '注册成功，请登录';
            successMessage.style.color = '#2ecc71';
            successMessage.style.textAlign = 'center';
            successMessage.style.marginBottom = '15px';
            loginForm.insertBefore(successMessage, loginForm.firstChild);
            
            // 3秒后移除提示
            setTimeout(() => {
                successMessage.remove();
            }, 3000);

        } catch (error) {
            console.error('注册错误:', error);
            const errorMessage = document.createElement('div');
            errorMessage.textContent = error.message; // 直接显示错误信息
            errorMessage.className = 'error-message';
            errorMessage.style.color = '#e74c3c';
            errorMessage.style.fontSize = '0.9em';
            errorMessage.style.marginTop = '5px';
            errorMessage.style.textAlign = 'center';
            registerForm.appendChild(errorMessage);
        }
    });

    // 登出功能
    logoutBtn.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/auth/logout', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                },
            });

            if (response.ok) {
                // 清除本地存储
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                // 切换回登录界面
                appContainer.style.display = 'none';
                authContainer.style.display = 'block';
                loginForm.style.display = 'block';
                registerForm.style.display = 'none';
            }
        } catch (error) {
            console.error('退出错误:', error);
            alert('退出失败，请重试');
        }
    });

    // 生成题目
    generateBtn.addEventListener('click', async () => {
        try {
            const difficulty = difficultySelect.value;
            const data = await apiRequest(`/api/drill/question?difficulty=${difficulty}`, {
                method: 'GET'
            });

            currentQuestion = data;
            questionElement.textContent = data.question;
            questionContainer.classList.remove('hidden');
            answerInput.value = '';
            resultElement.textContent = '';
        } catch (error) {
            console.error('获取题目错误:', error);
            alert(error.message);
        }
    });

    // 获取统计数据
    async function updateStats() {
        try {
            const stats = await apiRequest('/api/history/stats', {
                method: 'GET'
            });

            // 更新统计数据显示
            totalQuestionsElement.textContent = stats.total_questions;
            easyQuestionsElement.textContent = stats.easy_questions;
            mediumQuestionsElement.textContent = stats.medium_questions;
            hardQuestionsElement.textContent = stats.hard_questions;
            totalAttemptsElement.textContent = stats.total_attempts;
            correctAnswersElement.textContent = stats.correct_answers;
            accuracyElement.textContent = stats.accuracy.toFixed(1) + '%';
        } catch (error) {
            console.error('获取统计数据错误:', error);
        }
    }

    // 检查本地存储中是否有token，实现自动登录
    async function checkAuthStatus() {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        
        if (token && username) {
            try {
                currentUser = {
                    username,
                    token
                };
                usernameDisplay.textContent = username;
                authContainer.style.display = 'none';
                appContainer.style.display = 'block';
            } catch (error) {
                console.error('自动登录错误:', error);
                clearAuthData();
            }
        }
    }

    // 在页面加载时获取统计数据
    checkAuthStatus().then(() => {
        if (localStorage.getItem('token')) {
            updateStats();
        }
    });

    // 在提交答案后更新统计数据
    submitBtn.addEventListener('click', async () => {
        if (!currentQuestion) {
            alert('请先生成题目');
            return;
        }

        const answer = parseInt(answerInput.value);
        if (isNaN(answer)) {
            alert('请输入有效的答案');
            return;
        }

        try {
            const data = await apiRequest('/api/drill/answer', {
                method: 'POST',
                body: JSON.stringify({
                    question_id: currentQuestion.id,
                    answer: answer,
                    question: currentQuestion.question,
                    difficulty: currentQuestion.difficulty
                })
            });

            resultElement.textContent = data.message;
            if (data.correct) {
                resultElement.className = 'correct';
            } else {
                resultElement.className = 'incorrect';
            }

            // 更新统计数据
            updateStats();
        } catch (error) {
            console.error('提交答案错误:', error);
            alert(error.message);
        }
    });

    // 清除认证数据
    function clearAuthData() {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        currentUser = null;
    }

    // 显示热度排行榜
    hotRankBtn.addEventListener('click', showHotRankModal);

    // 关闭热度排行榜
    closeBtn.addEventListener('click', closeHotRankModal);

    // 点击模态框外部关闭
    window.addEventListener('click', (e) => {
        if (e.target === hotRankModal) {
            closeHotRankModal();
        }
    });

    // 切换排行榜标签
    rankTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            // 移除所有标签的active类
            rankTabs.forEach(t => t.classList.remove('active'));
            // 添加当前标签的active类
            tab.classList.add('active');
            // 加载对应类型的排行榜
            loadRankingData(tab.dataset.type);
        });
    });
});