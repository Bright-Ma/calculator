// 显示错误信息
function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    errorElement.textContent = message;
    errorElement.classList.add('show');
}

// 清除错误信息
function clearError(elementId) {
    const errorElement = document.getElementById(elementId);
    errorElement.textContent = '';
    errorElement.classList.remove('show');
}

// 清除所有错误信息
function clearAllErrors() {
    const errorElements = document.querySelectorAll('.error-message');
    errorElements.forEach(element => {
        element.textContent = '';
        element.classList.remove('show');
    });
}

// 登录表单提交
document.getElementById('login').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearAllErrors();

    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 401) {
                showError('login-error', '用户名或密码错误');
            } else {
                showError('login-error', data.error || '登录失败，请稍后重试');
            }
            return;
        }

        // 登录成功
        localStorage.setItem('token', data.token);
        localStorage.setItem('username', data.username);
        localStorage.setItem('role', data.role);
        
        document.getElementById('auth-container').style.display = 'none';
        document.getElementById('app').style.display = 'block';
        document.getElementById('username-display').textContent = data.username;
    } catch (error) {
        showError('login-error', '网络错误，请稍后重试');
    }
});

// 注册表单提交
document.getElementById('register').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearAllErrors();

    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;
    const role = document.getElementById('register-role').value;

    // 验证密码
    if (password !== confirmPassword) {
        showError('register-confirm-password-error', '两次输入的密码不一致');
        return;
    }

    try {
        const response = await fetch('/api/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password, role }),
        });

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 400 && data.error.includes('用户名已存在')) {
                showError('register-username-error', '用户名已存在');
            } else {
                showError('register-error', data.error || '注册失败，请稍后重试');
            }
            return;
        }

        // 注册成功，切换到登录界面
        document.getElementById('register-form').style.display = 'none';
        document.getElementById('login-form').style.display = 'block';
        showError('login-error', '注册成功，请登录');
    } catch (error) {
        showError('register-error', '网络错误，请稍后重试');
    }
});

// 切换登录/注册表单
document.getElementById('show-register').addEventListener('click', (e) => {
    e.preventDefault();
    clearAllErrors();
    document.getElementById('login-form').style.display = 'none';
    document.getElementById('register-form').style.display = 'block';
});

document.getElementById('show-login').addEventListener('click', (e) => {
    e.preventDefault();
    clearAllErrors();
    document.getElementById('register-form').style.display = 'none';
    document.getElementById('login-form').style.display = 'block';
});

// 退出登录
document.getElementById('logout-btn').addEventListener('click', () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('role');
    document.getElementById('app').style.display = 'none';
    document.getElementById('auth-container').style.display = 'block';
    document.getElementById('login-form').style.display = 'block';
    document.getElementById('register-form').style.display = 'none';
});

// 检查登录状态
function checkAuth() {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    if (token && username) {
        document.getElementById('auth-container').style.display = 'none';
        document.getElementById('app').style.display = 'block';
        document.getElementById('username-display').textContent = username;
    }
}

// 页面加载时检查登录状态
checkAuth(); 