package main

import "testing"

type Test struct {
	first  int
	second int
	out    int
}

func TestSum(t *testing.T) {
	// набор тестов
	cases := []struct {
		// имя теста
		name string
		// значения на вход тестируемой функции
		values []byte
		// желаемый результат
		want  int
		want2 error
	}{
		// тестовые данные №1
		{
			name:   "positive values",
			values: []byte("Hello, World!"),
			want:   13,
		},
		// тестовые данные №2
		{
			name:   "mixed values",
			values: []byte("Привет, весь мир?"),
			want:   17,
		},
		{
			name:   "mixed values2",
			values: []byte("😊🌟🎉🌈🎈"),
			want:   5,
		},
		{
			name:   "fgfg",
			values: []byte("Hello\xffWorld"),
			want:   0,
		},
	}
	// перебор всех тестов
	for _, tc := range cases {
		tc := tc
		// запуск отдельного теста
		t.Run(tc.name, func(t *testing.T) {
			// тестируем функцию Sum
			got, _ := GetUTFLength(tc.values)
			// if err != nil {
			// 	t.Fatalf("Error: %v", err)
			// }
			// проверим полученное значение
			if got != tc.want {
				t.Errorf("Sum(%v) = %v; want %v", tc.values, got, tc.want)
			}
		})
	}
}
