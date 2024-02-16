$(document).ready(function() {
    // Функция для получения данных с бэкенда
    function getItems() {
        $.ajax({
            type: 'GET',
            url: 'http://localhost:8080/api/go/expressions',
            success: function(response) {
                // Обработка полученных данных
                response.forEach(function(item) {
                    $('#response').append('<p>Name: ' + item.MathExpr + ', Description: ' + item.Result + '</p>');
                });
            },
            error: function(xhr, status, error) {
                $('#response').text('Error getting items');
            }
        });
    }

    // Вызов функции при загрузке страницы
    getItems();

    $('#itemForm').submit(function(e) {
        e.preventDefault();
        var MathExpr = $('#MathExpr').val();
        var Result = $('#Result').val();

        $.ajax({
            type: 'POST',
            url: 'http://localhost:8080/api/go/expressions',
            data: JSON.stringify({ MathExpr: MathExpr, Result: Result }),
            contentType: 'application/json',
            success: function(response) {
                $('#response').text('Item created successfully');
                // После создания элемента обновляем список
                getItems();
            },
            error: function(xhr, status, error) {
                $('#response').text('Error creating item');
            }
        });
    });
});