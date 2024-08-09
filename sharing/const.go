package main

import (
	"fmt"
	"time"
)

const (
	UserTypePremium = "premium"

	CacheKeyActiveUsersByTypeFmt = "users:%s"
	CacheTtlActiveUserByType     = 1 * time.Minute
)

func getCacheKeyActiveUsersByType(userType string) string {
	return fmt.Sprintf(CacheKeyActiveUsersByTypeFmt, userType)
}
