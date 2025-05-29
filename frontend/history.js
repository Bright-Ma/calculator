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
        return data;
    } catch (error) {
        console.error('Error:', error);
        return [];
    }
}

// 获取统计数据
async function fetchStats() {
    try {
        const response = await fetch('/api/history/stats', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) throw new Error('获取统计数据失败');
        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        return null;
    }
}

// 渲染贡献日历
function renderContributionCalendar(history) {
    const months = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'];
    const days = ['日', '一', '二', '三', '四', '五', '六'];
    
    // 渲染星期标签
    const daysContainer = document.querySelector('.calendar-days');
    daysContainer.innerHTML = days.map(day => `<div>${day}</div>`).join('');
    
    // 计算最近一年的日期
    const today = new Date();
    const oneYearAgo = new Date(today);
    oneYearAgo.setFullYear(today.getFullYear() - 1);
    
    // 创建日期映射
    const dateMap = new Map();
    history.forEach(record => {
        const date = new Date(record.created_at).toISOString().split('T')[0];
        dateMap.set(date, (dateMap.get(date) || 0) + 1);
    });
    
    // 渲染日历格子
    const gridContainer = document.querySelector('.calendar-grid');
    gridContainer.innerHTML = '';
    
    // 计算每个月的起始位置
    const monthPositions = new Map();
    let currentMonth = -1;
    let currentYear = -1;
    
    // 计算需要显示的月份
    const visibleMonths = new Set();
    for (let i = 0; i < 53; i++) { // 53周
        for (let j = 0; j < 7; j++) { // 7天
            const date = new Date(oneYearAgo);
            date.setDate(date.getDate() + (i * 7 + j));
            visibleMonths.add(date.getMonth());
        }
    }
    
    // 渲染日历格子
    for (let i = 0; i < 53; i++) { // 53周
        const column = document.createElement('div');
        column.className = 'calendar-column';
        
        for (let j = 0; j < 7; j++) { // 7天
            const date = new Date(oneYearAgo);
            date.setDate(date.getDate() + (i * 7 + j));
            
            const cell = document.createElement('div');
            cell.className = 'calendar-cell';
            
            const dateStr = date.toISOString().split('T')[0];
            const count = dateMap.get(dateStr) || 0;
            
            // 设置贡献等级
            if (count > 0) {
                const level = Math.min(Math.ceil(count / 5), 4);
                cell.setAttribute('data-level', level);
            }
            
            // 设置工具提示
            if (count > 0) {
                cell.setAttribute('data-tooltip', `${dateStr}: ${count} 次练习`);
            }
            
            // 点击事件
            cell.addEventListener('click', () => {
                if (count > 0) {
                    window.location.href = `history-detail.html?date=${dateStr}`;
                }
            });
            
            column.appendChild(cell);
            
            // 记录月份位置
            const month = date.getMonth();
            const year = date.getFullYear();
            if (month !== currentMonth || year !== currentYear) {
                monthPositions.set(`${year}-${month}`, i);
                currentMonth = month;
                currentYear = year;
            }
        }
        
        gridContainer.appendChild(column);
    }
    
    // 渲染月份标签
    const monthsContainer = document.querySelector('.calendar-months');
    monthsContainer.innerHTML = '';
    
    // 按时间顺序显示月份
    const sortedMonths = Array.from(monthPositions.entries())
        .sort((a, b) => {
            const [yearA, monthA] = a[0].split('-').map(Number);
            const [yearB, monthB] = b[0].split('-').map(Number);
            return yearA === yearB ? monthA - monthB : yearA - yearB;
        });
    
    sortedMonths.forEach(([dateKey, position]) => {
        const [year, month] = dateKey.split('-').map(Number);
        const monthLabel = document.createElement('div');
        monthLabel.textContent = months[month];
        monthLabel.style.gridColumn = position + 1;
        monthsContainer.appendChild(monthLabel);
    });
}

// 更新统计数据
function updateStats(stats) {
    if (!stats) return;
    
    document.getElementById('total-questions').textContent = stats.total_questions;
    document.getElementById('total-attempts').textContent = stats.total_attempts;
    document.getElementById('correct-answers').textContent = stats.correct_answers;
    document.getElementById('correct-rate').textContent = `${stats.accuracy.toFixed(2)}%`;
    document.getElementById('easy-questions').textContent = stats.easy_questions;
    document.getElementById('medium-questions').textContent = stats.medium_questions;
    document.getElementById('hard-questions').textContent = stats.hard_questions;
}

// 退出登录
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    window.location.href = 'login.html';
}

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', async () => {
    checkAuth();
    
    const history = await fetchHistory();
    const stats = await fetchStats();
    
    renderContributionCalendar(history);
    updateStats(stats);
    
    // 退出登录
    document.getElementById('logout-btn').addEventListener('click', logout);
}); 