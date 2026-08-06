package service

// IsAdminUser 判断用户 ID 是否在管理员白名单内。
func IsAdminUser(userID int64, adminUserIDs []int64) bool {
	if userID <= 0 {
		return false
	}
	for _, id := range adminUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}
