package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userRole := c.GetString("role")

	if userRole != role {
		err = fmt.Errorf("UnAuthenticated to access this resource")
		return err
	}
	return nil
}

// user can access this resource only via his token or he is an ADMIN...
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	uid := c.GetString("uid")
	userType := c.GetString("role")

	if err := CheckUserType(c, "LIBRARIAN"); err == nil {
		return nil
	}

	if uid == userId && userType == "MEMBER" {
		return nil
	}

	err = fmt.Errorf("UnAuthenticated to access this resource")
	return err
}

func MatchCustomerTypeToCid(c *gin.Context, customerId string) (err error) {
	cid := c.GetString("cid")

	if cid == customerId {
		return nil
	}

	err = fmt.Errorf("UnAuthenticated to access this resource")
	return err
}
