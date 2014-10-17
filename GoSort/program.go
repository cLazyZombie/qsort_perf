package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Implements sort.Interface for sort.Sort() and sort.IsSorted()
type Int32Slice []int32

func (s Int32Slice) Len() int {
	return len(s)
}

func (s Int32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Int32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	n := 1000000
	a := make(Int32Slice, n)

	threeway := flag.Bool("threeway", false, "")
	builtin := flag.Bool("builtin", false, "")
	flag.Parse()

	iter := 12
	trim := 1
	elapsed := make(sort.IntSlice, iter)
	var err error

	for i := 0; i < iter; i++ {
		elapsed[i], err = run(a, rand.Int63(), *threeway, *builtin)
		if err != nil {
			panic(err)
		}
	}

	elapsed.Sort()

	sum := 0
	for i := trim; i < iter-trim; i++ {
		sum += elapsed[i]
	}

	mean := float64(sum) / float64(iter-trim*2)
	fmt.Printf("%f ms elapsed\n", mean)
}

func run(a Int32Slice, seed int64, threeway bool, builtin bool) (int, error) {
	rand.Seed(seed)
	for i := 0; i < len(a); i++ {
		a[i] = rand.Int31n(math.MaxInt16) + 1
	}

	var begin time.Time
	if threeway {
		begin = time.Now()
		quicksortThreeway(a)
	} else if builtin {
		begin = time.Now()
		sort.Sort(a)
	} else {
		begin = time.Now()
		quicksort(a)
	}

	elapsed := int(time.Since(begin) / time.Millisecond)

	if !sort.IsSorted(a) {
		return 0, fmt.Errorf("Not sorted")
	}

	return elapsed, nil
}

func quicksortThreeway(a Int32Slice) {
	start := int32(0)
	end := int32(len(a))
	if end < 2 {
		return
	}

	p := a[start]
	l := start
	m := l + 1
	i := m

	for i < end {
		if a[i] < p {
			a[l] = a[i]
			l++
			a[i] = a[m]
			a[m] = p
			m++
		} else if a[i] == p {
			a[i] = a[m]
			a[m] = p
			m++
		}
		i++
	}

	quicksortThreeway(a[:l])
	quicksortThreeway(a[m:])
}

func quicksort(a Int32Slice) {
	start := 0
	end := len(a)
	if end < 2 {
		return
	}

	p := a[end/2]
	l := start
	r := end - 1

	for l <= r {
		if a[l] < p {
			l++
			continue
		}

		if a[r] > p {
			r--
			continue
		}

		a[l], a[r] = a[r], a[l]
		l++
		r--
	}

	quicksort(a[:r+1])
	quicksort(a[r+1:])
}
