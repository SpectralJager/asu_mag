import random

# Генерация всех возможных 6-битных чисел
all_numbers = [format(i, '06b') for i in range(64)]

# Выбор случайных 60 чисел из 64
selected_numbers = random.sample(all_numbers, 60)

# Замена 0 на 0.1 и 1 на 0.9
matrix = [[0.9 if bit == '1' else 0.1 for bit in number] for number in selected_numbers]

# Вывод матрицы
for row in matrix:
    print(row)