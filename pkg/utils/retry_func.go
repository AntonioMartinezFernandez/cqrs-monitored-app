package utils

func RetryFunc(fn func() (any, error), timesToRetry int) (any, error) {
	var (
		err    error
		result any
	)

	for i := 1; i <= timesToRetry; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
	}

	return nil, err
}
