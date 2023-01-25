package main

import "fmt"

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 1. Implementati o functie care returneaza sirul lui fibonacci. Functia primeste ca parametru nr de elemente pe care le returnam.
// ex: func (5) -> return primele 5 elemente din sir.
func fibo(terms int) []int {
	x := 0
	y := 1

	if terms <= 0 {
		return []int{}
	}

	if terms == 1 {
		return []int{x}
	}

	output := []int{x, y}
	for i := 0; i < terms-2; i++ {
		z := x + y
		output = append(output, z)
		x = y
		y = z
	}

	return output
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 2. Implementati o functie de sortare a elementelor dintr-un array.
// Functia primeste ca parametru un array of int si returneaza array-ul sortat.
func partition(arr []int, start, end int) ([]int, int) {
	//alegem ultima pozitie ca si pivot
	pivot := arr[end]
	i := start
	for j := start; j < end; j++ {
		//daca
		if arr[j] < pivot {
			//swap elements
			aux := arr[i]
			arr[i] = arr[j]
			arr[j] = aux
			i++
		}
	}
	aux := arr[i]
	arr[i] = arr[end] //sau pivot
	arr[end] = aux

	return arr, i
}
func quickSort(arr []int, start, end int) []int {
	//verificam sa fie adev conditia
	if start < end {
		// asta e ocul pivotului dupa fiecare iteratie
		var p int
		//facem o iteratie pentru a gasi locului primului pivot
		arr, p = partition(arr, start, end)
		//sortam in acelasi stil partea din stanga si cea din dreapta
		arr = quickSort(arr, start, p-1)
		arr = quickSort(arr, p+1, end)
	}
	return arr
}
func sortArray(arr []int) []int {
	//start till end
	if len(arr) < 2 {
		return arr
	}
	return quickSort(arr, 0, len(arr)-1)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 3. Implementati o functie care returneaza duplicatele dintr-un array.
// Functia primeste ca parametru un array si returneaza elementele duplicate.
// ex: func( [5,4,5,5,3,1,3] -> return Elementul 5 exista de 3 ori, Elementul 3 exista de 2 ori
func getDups(arr []int) {
	dict := map[int]int{}
	for _, elem := range arr {
		if dict[elem] != 0 {
			dict[elem]++
		} else {
			dict[elem] = 1
		}
	}

	for key, value := range dict {
		if value > 1 {
			fmt.Printf("Elementul %v a aparut de %v ori\n", key, value)
		}
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 4. Implementati o functie care elimina duplicatele dintr-un array.
// Functia primeste ca parametru un array si returneaza array-ul fara duplicate.
func elimDups(arr []int) []int {
	dict := map[int]bool{}
	var set []int

	for _, elem := range arr {
		// ignore value
		// if ok is true don't add the key
		if _, ok := dict[elem]; !ok {
			// add key to the set and dict
			dict[elem] = true
			set = append(set, elem)
		}
	}

	return set
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 5. Pentru un string s, returnati cel mai lung substring palindrom al lui s.
// ex: func("asd/ana/minim") -> return minim
func expandAroundLetter(input string, poz int) (string, bool) {
	var left, right int
	substr := string(input[poz])
	if poz >= 1 {
		left = poz - 1
	}
	if poz < len(input)-1 {
		right = poz + 1
	}

	for (left >= 0) && (right <= len(input)-1) {
		if input[left] != input[right] {
			return string(input[poz]), false
		}
		substr = string(input[left]) + substr + string(input[right])
		left--
		right++
	}

	return substr, true
}

func getLongestSubstring(input string) string {
	if len(input) <= 1 {
		return input
	}

	longestSubstring := string(input[0])
	maxSubstrLength := 1

	for index := range input {
		substr, flag := expandAroundLetter(input, index)
		if flag && len(substr) > maxSubstrLength {
			maxSubstrLength = len(substr)
			longestSubstring = substr
		}
	}

	return longestSubstring
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 6. Pentru un array ce contine n numere naturale distincte din intervalul [0, n],
// returnati singurul numar din interval ce lipseste din array-ul dat.
func getMissingNr(arr []int, n int) int {
	// n poate fi dat ca param al intervalului [0, n]
	// sau poate fi doar len(arr) + 1
	totalSum := n * (n + 1) / 2
	sum := 0
	for _, val := range arr {
		sum += val
	}

	return totalSum - sum
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	// 1
	// fmt.Println(fibo(1))
	// fmt.Println(fibo(2))
	// fmt.Println(fibo(5))
	// fmt.Println(fibo(10))

	// 2
	// arr := []int{5, 8, 2, 4, 3, 6, 8, 2, 4, 10, 68, 23, 54, 21, 63}
	// fmt.Println(sortArray(arr))

	// 3
	// arr := []int{5, 5, 3, 4, 3, 8, 9, 1, 1, 5, 2, 4, 6, 8}
	// getDups(arr)

	// 4
	// arr := []int{5, 5, 3, 4, 3, 8, 9, 1, 1, 5, 2, 4, 6, 8}
	// fmt.Println(elimDups(arr))

	// 5
	// strTest1 := "abcba"
	// strTest2 := "abacdeded"
	// strTest3 := "a"
	// strTest4 := "abaaaabbaaaac"
	// fmt.Println(getLongestSubstring(strTest1))
	// fmt.Println(getLongestSubstring(strTest2))
	// fmt.Println(getLongestSubstring(strTest3))
	// fmt.Println(getLongestSubstring(strTest4))

	//6
	// in intervalul [0, 3] avem 4 elemente daca includem si 0, n = 3
	// atentie la al doilea parametru (pentru n)
	// arr := []int{0, 1, 2, 3, 4, 5, 7}
	// fmt.Println(getMissingNr(arr, len(arr)))
}
