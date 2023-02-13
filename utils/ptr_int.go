package utils

func IntPtrToUintPtr(i *int) *uint {
	if i == nil {
		return nil
	}

	return Ptr(uint(*i))
}

func UintPtrToIntPtr(i *uint) *int {
	if i == nil {
		return nil
	}

	return Ptr(int(*i))
}
