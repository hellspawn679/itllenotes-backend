package middleware
import(
	//"fmt"
	"github.com/nekonotes/database"
	"github.com/jinzhu/gorm"
)
func Name_check(name string,db *gorm.DB) string{
	if len(name) < 3 || len(name) > 100 {
		return "Name should be between 3 to 100 characters"
	}
	var users []Database.User
	db.Find(&users, "name = ?", name)
	for _, user := range users {
		if user.Name == name {
			return "user name already exists"
		}
	}
	return "done"
}