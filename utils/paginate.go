package utils

func Paginate[T any](slice []T, page int, pageSize int, extra int) []T {

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 15
	}

	if page <= 0 {
		page = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize + extra
	if end > len(slice) {
		end = len(slice)
	}

	return slice[start:end]
}
