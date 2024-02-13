function sendExpression() {
    const mathExprInput = document.getElementById('mathExpr');
    const mathExpr = mathExprInput.value;
    const data = { MathExpr: mathExpr, Result: "", Status: "не отправлено" };

    fetch('http://localhost:8000/api/go/expressions', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Ответ от сервера:', data);
        mathExprInput.value = ''; // Очистка поля ввода
        updateExpressionsList(); // Обновляем список выражений на странице
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
}

function updateExpressionsList() {
    fetch('http://localhost:8000/api/go/expressions') // Замените '/expressions' на фактический URL вашего эндпоинта
        .then(response => response.json())
        .then(data => {
            const expressionsList = document.getElementById('expressions-list');
            expressionsList.innerHTML = ''; // Очищаем список перед обновлением

            data.forEach(expression => {
                const expressionString = expression.MathExpr;

                const expressionItem = document.createElement('div');
                expressionItem.textContent = expressionString;

                if (expression.Status === "отправлено на вычисление") {
                    expressionItem.innerHTML += ' - отправлено на вычисление';
                }

                const deleteButton = document.createElement('button');
                deleteButton.textContent = 'Удалить';
                deleteButton.addEventListener('click', () => deleteExpression(expression.ID));
                expressionItem.appendChild(deleteButton);

                expressionsList.appendChild(expressionItem);
            });
        });
}

function deleteExpression(expressionId) {
    fetch('http://localhost:8000/api/go/expressions/${expressionId}', {
        method: 'DELETE'
    })
    .then(() => updateExpressionsList())
    .catch(error => console.error('Ошибка при удалении:', error));
}

function sendToAgent() {
    fetch('http://localhost:8000/api/go/expressions/agent', {
        method: 'GET'
    })
    .then(() => updateExpressionsList())
    .catch(error => console.error('Ошибка при отправке на вычисление:', error));
}
