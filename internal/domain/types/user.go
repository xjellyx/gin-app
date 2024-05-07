package types

// UserGender 用户性别
type UserGender uint // user gender
const (
	GenderUnknown UserGender = iota + 1
	GenderMale               // male
	GenderFemale             // female
)

// UserStatus 用户状态
type UserStatus uint // user status
const (
	StatusNormal  UserStatus = iota + 1 // 正常
	StatusLocked                        // 锁定
	StatusFreeze                        // 冻结
	StatusDeleted                       // deleted

)
