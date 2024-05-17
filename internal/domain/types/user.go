package types

// UserGender 用户性别
type UserGender uint // user gender
const (
	GenderMale   UserGender = iota + 1 // male
	GenderFemale                       // female
	GenderUnknown
)

// UserStatus 用户状态
type UserStatus uint // user status
const (
	UserStatusNormal UserStatus = iota + 1 // 正常
	UserStatusLocked                       // 禁用
)
