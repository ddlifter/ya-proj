package main

import "sort"

func merge(nums1 []int, m int, nums2 []int, n int) {
	index := 0
	for i := m; i < n+m; i++ {
		nums1[i] = nums2[index]
		index++
	}

	sort.Ints(nums1)
}
