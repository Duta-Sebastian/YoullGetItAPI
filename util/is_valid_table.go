package util

func IsTableAllowed(table string) bool {
	validTables := []string{"job_cart", "auth_user", "cv", "question"}

	isValidTable := false

	for _, validTable := range validTables {
		if validTable == table {
			isValidTable = true
			break
		}
	}
	return isValidTable
}
