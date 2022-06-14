package selection

// Пакет selection отвечает за оптимальный выбор столиков
// для бронирования

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	"github.com/AscaroLabs/EateryGlass/pkg/structures"
	"gonum.org/v1/gonum/stat/combin"
)

// SelectTables выбирает столики(или их комбинации), которые подходят по времени и по вместимости
func SelectTables(db *sql.DB, volume,
	appropriateTime string) (map[string][]structures.Table, error) {

	// Парсим время из запроса согласно стандарту RFC3339 e.g. "2006-01-02T15:04:05Z"
	requiredTime, err := time.Parse(time.RFC3339, appropriateTime)
	if err != nil {
		return nil, err
	}

	// Проверяем время на корректность (рестораны работают с 9:00 до 23:00, т.к. бронь длится 2 часа, то
	// забронировать столик после 20:59 не выйдет)
	if requiredTime.Hour() < 9 || requiredTime.Hour() >= 21 {
		return nil, errors.New("Wrong time!")
	}
	// Конвертируем вместимость из строки в чило
	capacity, err := strconv.Atoi(volume)
	if err != nil {
		return nil, err
	}

	// Получаем из БД столики, которые подходят нам по времени
	tables, err := storage.GetTablesByTime(db, requiredTime)
	if err != nil {
		return nil, err
	}

	// Хеш-таблица tablesByRestaurants делит доступные столики по идентификаторам ресторанов
	var tablesByRestaurants map[string][]structures.Table
	tablesByRestaurants = make(map[string][]structures.Table)
	for _, t := range tables {
		tablesByRestaurants[t.Restaurant_id] = append(tablesByRestaurants[t.Restaurant_id], t)
	}

	// Для каждого ресторана ищем оптимальный подходящий вариант
	for restaurant, restaurantTables := range tablesByRestaurants {
		tablesByRestaurants[restaurant] = findOpt(restaurantTables, capacity)
	}

	// Возвращаем полученную хеш-таблицу
	return tablesByRestaurants, nil
}

// findOpt ищет оптимальный вариант из возможных
func findOpt(tables []structures.Table, requiredCapacity int) []structures.Table {

	//	Задача аналогична задаче о дискретном рюкзаке,
	// 	которая обычно решается методом динамисекого
	// 	программирования

	//	Однако, мы знаем, что количество столиков в ресторане
	// 	небольшое число, так что, чтобы не усложнять,
	// 	воспользуемя перебором по подмножествам столиков
	//	(асимптотика 2^(кол-во столиков))

	//	l - количесво столиков(размернось входного массива)
	l := len(tables)
	if l == 0 {
		// Если на вход подали пустое множество столов, то подходящих нет
		return nil
	}

	// goodSubSet определяет структуру для "хорошего" подмножества столиков,
	// т.е. подмножества общая вместимость которых >= необходимой
	type goodSubSet struct {
		subSetIndexes []int // Индексы элементов в исх. массиве tables
		capacitySum   int   // Общая вместимость подмножества
	}

	// все хорошие подмножества столиков
	var goodSubSets []goodSubSet

	// i - количество элементов в рассматриваемом подмножестве (i = 1, ..., l)
	for i := 1; i-1 < l; i++ {

		// Создаем новый генератор сочетаний из C(l, i)
		gen := combin.NewCombinationGenerator(l, i)

		// Проходимся во всевозможным подмножествам из i элементов
		for gen.Next() {

			// Индексы элементов, которые входят в рассматриваемое подмножество
			subSetIndexes := gen.Combination(nil)

			// Суммируем вместимость всех столиков в подмножестве
			var capacitySum int
			for _, v := range subSetIndexes {
				capacitySum += tables[v].Capacity
			}

			// Если подмножество оказалось хорошим, то добавляем его в goodSubSets
			if capacitySum >= requiredCapacity {
				goodSubSets = append(goodSubSets, goodSubSet{subSetIndexes, capacitySum})
			}
		}
	}

	// Если вообще нет хороших подмножеств, то возвращаем nil
	if len(goodSubSets) == 0 {
		return nil
	}

	// Ищем подмножество, которое наиболее близко по вместимости к необходимому
	var minCapValue int
	minCapValue = 100000000 // аналог +бесконечности
	var minCapIndex int
	for i, val := range goodSubSets {
		if val.capacitySum < minCapValue {
			minCapValue = val.capacitySum
			minCapIndex = i
		}
	}

	// Итоговый набор столиков для ресторана
	var result []structures.Table
	for _, i := range goodSubSets[minCapIndex].subSetIndexes {
		result = append(result, tables[i])
	}
	return result
}
