package database

type User struct {
	// ID 使用者 ID
	ID string `json:"userId" gorm:"primaryKey"`

	// DisplayName 使用者的顯示名稱
	DisplayName string `json:"displayName"`

	// Challenge 當次進行 WebAuthn 註冊 / 驗證流程時的使用者 Challenge
	Challenge string `json:"challenge"`

	// Credential 使用者的 WebAuthn Credential
	Credential string `json:"credential"`
}

func (*User) TableName() string {
	return "user"
}
