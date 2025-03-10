package cdq

type Permissions = int

const (
	Admin Permissions = 0 // 管理员
	User  Permissions = 1 // 用户
	Guest Permissions = 2 // 游客
)
