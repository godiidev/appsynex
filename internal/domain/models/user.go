package models

imprort (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID				uint  			`gorm:"primaryKey"`
	Username		string	   		`gorm:"size:100;uniqueIndex" json:"username"`
	PasswordHash	string			`gorm:"size:255" json:"-"`
	Emai			string			`gorm:"size:100;json:"email"`
	Phone			*string			`gorm:"size:50;json:"phone"`
	LastLogin		*time.Time		`json:"last_login"`
	AccountStatus	string			`gorm:"size:50;default:active" json:"account_status"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	DeleteAt		gorm.DeletedAt	`gorm:"index" json:"deleted_at"`
	Roles			[]Role			`gorm:"many2many:user_roles" json:"roles,omitempty"`
}

type UserRole struch {
	UserID	 uint 		`gorm:"primaryKey"`
	RoleID	 uint		`gorm:"primaryKey"`
	CreateAt time.Time
	UpdateAt time.Time
}
