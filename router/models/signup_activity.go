package models

import "github.com/FoxComm/libs/spree"

//I didn't want to actually store all the fields, but I ended up storing so much, that I thought to include everything.
type SignupActivityDetails struct {
	User          spree.User
	IsSocial      bool
	SocialNetwork string
}
