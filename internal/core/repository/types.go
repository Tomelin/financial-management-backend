package repository

func checkFirebaseCondition(condition *string) string {
	if condition == nil {
		return "=="
	}

	validConditions := []string{"==", "!=", "<", "<=", ">", ">=", "array-contains", "array-contains-any", "in", "not-in"}
	for _, validCondition := range validConditions {
		if *condition == validCondition {
			return *condition
		}
	}

	return "=="
}
