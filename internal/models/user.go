package models

import (
	"github.com/thinhlu123/shortener/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var UserDB db.Collection

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Usr string `json:"usr,omitempty" bson:"usr,omitempty"`
	Pwd string `json:"pwd,omitempty" bson:"pwd,omitempty"`

	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	FullName    string `json:"full_name,omitempty" bson:"full_name,omitempty"`
	PhoneNumber string `json:"phone,omitempty" bson:"phone_number,omitempty"`

	PaymentAcc string `json:"payment_acc,omitempty" bson:"payment_acc,omitempty"`
}

func (u *User) HashPwd() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Pwd = string(hashedPassword)
	return nil
}

func (u User) ComparePwd(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(password)); err != nil {
		return false
	}
	return true
}
