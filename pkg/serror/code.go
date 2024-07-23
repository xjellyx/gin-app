package serror

type ErrorCode string

const (
	ErrInvalidInput            ErrorCode = "INVALID_INPUT"
	ErrUnauthorized            ErrorCode = "UNAUTHORIZED"
	ErrInvalidToken            ErrorCode = "INVALID_TOKEN"
	ErrCodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrEmailAlredayInUse       ErrorCode = "EMAIL_ALREADY_IN_USE"
	ErrUsernameAlredayInUse    ErrorCode = "USERNAME_ALREADY_IN_USE"
	ErrPhoneNumberAlredayInUse ErrorCode = "PHONE_NUMBER_ALREADY_IN_USE"
	ErrUserInactivate          ErrorCode = "USER_INACTIVATE"
	ErrUserFreeze              ErrorCode = "USER_FREEZE"
	ErrUserDeleted             ErrorCode = "USER_DELETED"
	ErrUserStatusAbnormal      ErrorCode = "USER_STATUS_ABNORMAL"
	ErrUserRecordNotFound      ErrorCode = "USER_RECORD_NOT_FOUND"
	ErrRecordNotFound          ErrorCode = "RECORD_NOT_FOUND"
	ErrCasbinRemoveFail        ErrorCode = "CASBIN_REMOVE_FAIL"
	ErrCasbinAddFail           ErrorCode = "CASBIN_ADD_FAIL"
)
