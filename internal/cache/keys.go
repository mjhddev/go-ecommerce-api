package cache

import "fmt"

func ProductListKey(page, limit int, search string, categoryID uint, sort string) string {
	return fmt.Sprintf(
		"products:%d:%d:%s:%d:%s",
		page,
		limit,
		search,
		categoryID,
		sort,
	)
}

func ProductKey(id uint) string {
	return fmt.Sprintf("product:%d", id)
}

func CategoryKey() string {
	return "categories"
}

func DashboardKey() string {
	return "dashboard"
}
