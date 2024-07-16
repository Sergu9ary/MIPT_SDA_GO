//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	guestCount := make(map[int]int)
	var changes []int

	for _, guest := range guests {
		guestCount[guest.CheckInDate] += 1
		guestCount[guest.CheckOutDate] -= 1
		changes = append(changes, guest.CheckInDate, guest.CheckOutDate)
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i] < changes[j]
	})

	changes = remove(changes)
	var result []Load
	count := 0
	for _, date := range changes {
		guest := guestCount[date]
		count += guest
		if guest != 0 {
			result = append(result, Load{StartDate: date, GuestCount: count})
		}
	}
	return result
}

func remove(nums []int) []int {
	result := []int{}
	used := map[int]bool{}
	for i := range nums {
		if used[nums[i]] {

		} else {
			used[nums[i]] = true
			result = append(result, nums[i])
		}
	}
	return result
}
