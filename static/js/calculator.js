document.addEventListener('DOMContentLoaded', function() {
    const num1Input = document.getElementById('num1');
    const num2Input = document.getElementById('num2');
    const operationSelect = document.getElementById('operation');
    const calculateButton = document.getElementById('calculate');
    const clearButton = document.getElementById('clear');
    const resultDiv = document.getElementById('result');

    // 计算按钮点击事件
    calculateButton.addEventListener('click', async function() {
        const num1 = parseFloat(num1Input.value);
        const num2 = parseFloat(num2Input.value);
        const operation = operationSelect.value;

        // 输入验证
        if (isNaN(num1) || isNaN(num2)) {
            showError('请输入有效的数字');
            return;
        }

        // 除法时检查除数是否为0
        if (operation === 'divide' && num2 === 0) {
            showError('除数不能为0');
            return;
        }

        try {
            const response = await fetch('/api/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    operation: operation,
                    a: num1,
                    b: num2
                })
            });

            const data = await response.json();

            if (data.error) {
                showError(data.error);
            } else {
                showResult(data.result);
            }
        } catch (error) {
            showError('计算请求失败');
            console.error('计算错误:', error);
        }
    });

    // 清除按钮点击事件
    clearButton.addEventListener('click', function() {
        num1Input.value = '';
        num2Input.value = '';
        resultDiv.textContent = '0';
        resultDiv.classList.remove('error');
    });

    // 显示错误信息
    function showError(message) {
        resultDiv.textContent = message;
        resultDiv.classList.add('error');
    }

    // 显示计算结果
    function showResult(result) {
        resultDiv.textContent = result;
        resultDiv.classList.remove('error');
    }

    // 输入框按回车时触发计算
    function handleEnterKey(event) {
        if (event.key === 'Enter') {
            calculateButton.click();
        }
    }

    num1Input.addEventListener('keypress', handleEnterKey);
    num2Input.addEventListener('keypress', handleEnterKey);
});