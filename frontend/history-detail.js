// 检查用户是否已登录
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }
    // 显示用户名
    const username = localStorage.getItem('username');
    document.getElementById('username-display').textContent = username;
}

// 全局变量
let allHistory = [];
let currentPage = 1;
const recordsPerPage = 10;
let currentFilters = {
    difficulty: '',
    result: '',
    date: ''
};

// 获取历史记录
async function fetchHistory() {
    try {
        const response = await fetch('/api/history', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) throw new Error('获取历史记录失败');
        const data = await response.json();
        return data.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    } catch (error) {
        console.error('Error:', error);
        return [];
    }
}

// 应用筛选
function applyFilters(history) {
    return history.filter(record => {
        if (currentFilters.difficulty && record.difficulty !== currentFilters.difficulty) return false;
        if (currentFilters.result) {
            const isCorrect = record.is_correct;
            if (currentFilters.result === 'correct' && !isCorrect) return false;
            if (currentFilters.result === 'incorrect' && isCorrect) return false;
        }
        if (currentFilters.date) {
            const recordDate = new Date(record.created_at).toISOString().split('T')[0];
            if (recordDate !== currentFilters.date) return false;
        }
        return true;
    });
}

// 显示历史记录
function displayHistory(filteredHistory) {
    const historyList = document.getElementById('history-list');
    historyList.innerHTML = '';
    
    const startIndex = (currentPage - 1) * recordsPerPage;
    const endIndex = startIndex + recordsPerPage;
    const currentPageData = filteredHistory.slice(startIndex, endIndex);
    
    currentPageData.forEach(record => {
        const historyItem = document.createElement('div');
        historyItem.className = 'history-item';
        
        const difficultyClass = {
            'easy': 'easy',
            'medium': 'medium',
            'hard': 'hard'
        }[record.difficulty] || 'medium';
        
        historyItem.innerHTML = `
            <div class="question">
                <div class="question-content">${record.question}</div>
                <div class="question-meta">
                    <span class="difficulty ${difficultyClass}">${record.difficulty}</span>
                    <span class="time-spent">用时: ${record.time_spent}秒</span>
                </div>
            </div>
            <div class="answer">
                <div>
                    <div class="label">你的答案</div>
                    <div class="value">${record.user_answer}</div>
                </div>
                <div>
                    <div class="label">正确答案</div>
                    <div class="value">${record.correct_answer}</div>
                </div>
            </div>
            <div class="result">
                <span class="status ${record.is_correct ? 'correct' : 'incorrect'}">
                    ${record.is_correct ? '正确' : '错误'}
                </span>
            </div>
            <div class="time">${new Date(record.created_at).toLocaleString()}</div>
        `;
        
        historyList.appendChild(historyItem);
    });
    
    // 更新分页信息
    updatePagination(filteredHistory.length);
}

// 更新分页控件
function updatePagination(totalRecords) {
    const totalPages = Math.ceil(totalRecords / recordsPerPage);
    const pageInfo = document.getElementById('page-info');
    const prevButton = document.getElementById('prev-page');
    const nextButton = document.getElementById('next-page');
    
    pageInfo.textContent = `第 ${currentPage} 页，共 ${totalPages} 页`;
    prevButton.disabled = currentPage === 1;
    nextButton.disabled = currentPage === totalPages;
}

// 处理筛选器变化
function handleFilterChange() {
    const difficultyFilter = document.getElementById('difficulty-filter');
    const resultFilter = document.getElementById('result-filter');
    const dateFilter = document.getElementById('date-filter');
    
    currentFilters = {
        difficulty: difficultyFilter.value,
        result: resultFilter.value,
        date: dateFilter.value
    };
    
    currentPage = 1; // 重置到第一页
    const filteredHistory = applyFilters(allHistory);
    displayHistory(filteredHistory);
}

// 重置筛选器
function resetFilters() {
    document.getElementById('difficulty-filter').value = '';
    document.getElementById('result-filter').value = '';
    document.getElementById('date-filter').value = '';
    currentFilters = {
        difficulty: '',
        result: '',
        date: ''
    };
    currentPage = 1;
    displayHistory(allHistory);
}

// 检查URL参数
function checkUrlParams() {
    const urlParams = new URLSearchParams(window.location.search);
    const date = urlParams.get('date');
    if (date) {
        document.getElementById('date-filter').value = date;
        currentFilters.date = date;
    }
}

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', async () => {
    checkAuth();
    checkUrlParams();
    
    allHistory = await fetchHistory();
    displayHistory(allHistory);
    
    // 添加事件监听器
    document.getElementById('apply-filters').addEventListener('click', handleFilterChange);
    document.getElementById('reset-filters').addEventListener('click', resetFilters);
    document.getElementById('prev-page').addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            const filteredHistory = applyFilters(allHistory);
            displayHistory(filteredHistory);
        }
    });
    document.getElementById('next-page').addEventListener('click', () => {
        const totalPages = Math.ceil(applyFilters(allHistory).length / recordsPerPage);
        if (currentPage < totalPages) {
            currentPage++;
            const filteredHistory = applyFilters(allHistory);
            displayHistory(filteredHistory);
        }
    });
    
    // 退出登录
    document.getElementById('logout-btn').addEventListener('click', () => {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        window.location.href = 'login.html';
    });
}); 