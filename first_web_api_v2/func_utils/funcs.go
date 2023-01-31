package funcutils

func Fibo(terms int) []int {
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
func SortArray(arr []int) []int {
	//start till end
	if len(arr) < 2 {
		return arr
	}
	return quickSort(arr, 0, len(arr)-1)
}

func GetDups(arr []int) map[int]int {
	dict := map[int]int{}
	for _, elem := range arr {
		if dict[elem] != 0 {
			dict[elem]++
		} else {
			dict[elem] = 1
		}
	}

	return dict
}

func ElimDups(arr []int) []int {
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

func GetLongestSubstring(input string) string {
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

func GetMissingNr(arr []int) int {
	//pt ca incepem de la zero, n va fi len(arr)
	//daca incepeam de la 1, n era len(arr)+1
	n := len(arr)
	totalSum := n * (n + 1) / 2
	sum := 0
	for _, val := range arr {
		sum += val
	}

	return totalSum - sum
}
