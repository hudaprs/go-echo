package helpers

import "golang.org/x/crypto/bcrypt"

// @description	Hash password
// @param 		password.(string)
// @return 		string, error
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// @description	Hash password
// @param 		password.(string), hashedPassword.(string)
// @return 		string, error
func PasswordCompare(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
