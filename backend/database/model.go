package database

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	// ID 使用者 ID
	ID string `json:"userId" gorm:"primaryKey"`

	// Name 使用者名稱
	Name string `json:"name"`

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

func (u *User) WebAuthnID() []byte {
	return []byte(u.Name)
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	credentials := []webauthn.Credential{}

	allUser, err := GetUsers()
	if err != nil {
		fmt.Println(err.Error())
		return credentials
	}

	for _, user := range allUser {
		if !strings.HasPrefix(user.Credential, "`") || user.Credential == "`{}`" {
			continue
		}

		s, err := strconv.Unquote(string(user.Credential))
		if err != nil {
			fmt.Println(err.Error())
			return credentials
		}
		var credsMap map[string]interface{}
		err = json.Unmarshal([]byte(s), &credsMap)
		if err != nil {
			fmt.Println(err.Error())
			return credentials
		}
		credsJson, err := json.Marshal(credsMap)
		if err != nil {
			fmt.Println(err.Error())
			return credentials
		}
		var cred webauthn.Credential
		err = json.Unmarshal(credsJson, &cred)
		if err != nil {
			fmt.Println(err.Error())
			return credentials
		}
		credentials = append(credentials, cred)
	}

	// fmt.Println("WebAuthnCredentials: ", utils.PrintJSON(credentials))

	return credentials
}

func (u *User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, credential := range u.WebAuthnCredentials() {
		descriptor := credential.Descriptor()
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}
