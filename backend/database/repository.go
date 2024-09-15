package database

// CreateUser 建立使用者
func CreateUser(u *User) error {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	return Context.db.Create(u).Error
}

// GetUserByID 透過 UserID 取得使用者
func GetUserByID(id string) (*User, error) {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	var u User
	err := Context.db.Where("id = ?", id).First(&u).Error
	return &u, err
}

// GetUserByName 透過 Name 取得使用者
func GetUserByName(name string) (*User, error) {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	var u User
	err := Context.db.Where("name = ?", name).First(&u).Error
	return &u, err
}

// GetUserByChallenge 透過 Challenge 取得使用者
func GetUserByChallenge(challenge string) (*User, error) {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	var u User
	err := Context.db.Where("challenge = ?", challenge).First(&u).Error
	return &u, err
}

// GetUsers 取得所有使用者
func GetUsers() ([]User, error) {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	var users []User
	err := Context.db.Find(&users).Error
	return users, err
}

// UpdateUser 更新使用者
func UpdateUser(u *User, updateData interface{}) error {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	return Context.db.Model(u).Updates(updateData).Error
}

// DeleteUser 刪除使用者
func DeleteUser(u *User) error {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	return Context.db.Delete(u).Error
}

// DeleteUserByID 透過 UserID 刪除使用者
func DeleteUserByID(id string) error {
	Context.mu.Lock()
	defer Context.mu.Unlock()

	return Context.db.Delete(&User{ID: id}).Error
}
