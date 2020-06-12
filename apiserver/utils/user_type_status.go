package utils

type UserType int8
type UserStatus int8

const (
	AdminUser UserType = iota
	GeneralUser
)

const (
	Enable UserStatus = iota
	Disabled
	Deleted
)

func UserTypeToStr(userType UserType) string {
	switch userType {
	case AdminUser:
		return "admin"
	case GeneralUser:
		return "user"
	default:
		return "user"
	}
}

func StrToUserType(userType string) UserType {
	switch userType {
	case "admin":
		return AdminUser
	case "user":
		return GeneralUser
	default:
		return GeneralUser
	}
}

func UserTypeToInt(userType UserType) int {
	switch userType {
	case AdminUser:
		return 0
	case GeneralUser:
		return 1
	default:
		return 1
	}
}

func IntToUserType(userType int) UserType {
	switch userType {
	case 0:
		return AdminUser
	case 1:
		return GeneralUser
	default:
		return GeneralUser
	}
}

func UserStatusToInt(userType UserStatus) int {
	switch userType {
	case Enable:
		return 0
	case Disabled:
		return 1
	case Deleted:
		return 2
	default:
		return 0
	}
}

func IntToUserStatus(userType int) UserStatus {
	switch userType {
	case 0:
		return Enable
	case 1:
		return Disabled
	case 2:
		return Deleted
	default:
		return Enable
	}
}

func StrToUserStatus(userType string) UserStatus {
	switch userType {
	case "enable":
		return Enable
	case "disabled":
		return Disabled
	case "deleted":
		return Deleted
	default:
		return Enable
	}
}

func UserStatusToStr(userType UserStatus) string {
	switch userType {
	case Enable:
		return "enable"
	case Disabled:
		return "disabled"
	case Deleted:
		return "deleted"
	default:
		return "enable"
	}
}
