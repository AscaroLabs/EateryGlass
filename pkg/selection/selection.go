package selection

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat/combin"

	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	"github.com/AscaroLabs/EateryGlass/pkg/structures"
)

func SelectTables(db *sql.DB, volume, appropriate_time string) (map[string][]structures.Table, error) {

	log.Printf("\n-----* start time parsing (%v) *-----\n", appropriate_time)

	t, err := time.Parse(time.RFC3339, appropriate_time)
	if err != nil {
		return nil, err
	}

	log.Printf("\n-----* time parsed (%v) *-----\n", t)

	n, err := strconv.Atoi(volume)
	if err != nil {
		return nil, err
	}

	log.Printf("\n-----* getting tables by time  *-----\n")

	tables, err := storage.GetTablesByTime(db, t)
	if err != nil {
		return nil, err
	}
	log.Printf("\n-----* tables (%v) getted!  *-----\n", len(tables))

	var TablesByRestaurants map[string][]structures.Table
	TablesByRestaurants = make(map[string][]structures.Table)
	for _, t := range tables {
		TablesByRestaurants[t.Restaurant_id] = append(TablesByRestaurants[t.Restaurant_id], t)
	}

	for rest, t := range TablesByRestaurants {
		var ok bool
		TablesByRestaurants[rest], ok = findOpt(t, n)
		if ok != true {
			return nil, nil
		}
	}

	return TablesByRestaurants, nil
}

func findOpt(tab []structures.Table, n int) ([]structures.Table, bool) {
	// количество столов в одном ресторане маленькое, так что можно обойтись без динамики
	l := len(tab)
	if l == 0 {
		return nil, true
	}
	type goodSubSet struct {
		subSetIndexes []int
		capacity_sum  int
	}
	var goodSubSets []goodSubSet
	for i := 1; i-1 < l; i++ {
		gen := combin.NewCombinationGenerator(l, i)
		for gen.Next() {
			subSetIndexes := gen.Combination(nil)
			var capacity_sum int
			for _, v := range subSetIndexes {
				capacity_sum += tab[v].Capacity
			}
			if capacity_sum >= n {
				goodSubSets = append(goodSubSets, goodSubSet{subSetIndexes, capacity_sum})
			}
		}
	}
	if len(goodSubSets) == 0 {
		return nil, true
	}
	var minCapValue int
	minCapValue = 100000000
	var minCapIndex int
	for i, val := range goodSubSets {
		if val.capacity_sum < minCapValue {
			minCapValue = val.capacity_sum
			minCapIndex = i
		}
	}
	var result []structures.Table
	for _, i := range goodSubSets[minCapIndex].subSetIndexes {
		result = append(result, tab[i])
	}
	return result, true
}
