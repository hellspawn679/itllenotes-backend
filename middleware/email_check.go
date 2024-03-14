package middleware
import (
	"regexp"
	"github.com/jinzhu/gorm"
	"github.com/nekonotes/database"
)
func EmailCheck(email string,db *gorm.DB) string {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    if !emailRegex.MatchString(email){
		return "invalid email address"
	}
	var users []Database.User
	db.Find(&users, "email = ?", email)
	for _, user := range users {
		if user.Email == email {
			return "email is already in use already exists"
		}
	}
	return "done" 
}
//testing code for email check
// for _, email := range []string{
//     "good@exmaple.com",
//     "bad-example",
// } {
//     fmt.Printf("%18s valid: %t\n", email, valid(email))
// }
// fmt.Println(isEmailValid("test44@gmail.com"))         // true 
// fmt.Println(isEmailValid("bad-email"))               // false
// fmt.Println(isEmailValid("test44$@gmail.com"))      // false
// fmt.Println(isEmailValid("test-email.com"))        // false
// fmt.Println(isEmailValid("test+email@test.com"))  // true