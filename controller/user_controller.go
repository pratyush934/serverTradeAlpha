package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

func GetUsersById(c echo.Context) error {

	id := c.Get("userId").(string)

	if id == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Please add user id as it is not there", nil)
	}

	byId, err := models.GetUserById(id)

	if err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "Not able to get the user via id", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Just able to do that",
		"user":    byId,
	})

}

func GetUserByEmail(c echo.Context) error {

	email := c.Get("email").(string)

	if email == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "Please add the email", nil)
	}

	byEmail, err := models.GetUserByEmail(email)

	if err != nil {
		return util.NewAppError(http.StatusNotFound, types.StatusNotFound, "Not able to get the user via email", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Get through the email",
		"user":    byEmail,
	})
}

func DeleteUser(c echo.Context) error {

	role := c.Get("role").(int)

	if role != 2 {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not an admin", nil)
	}

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userId", nil)
	}

	if err := models.DeleteUserById(userId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the user", nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted",
		"status":  types.StatusOK,
	})
}

func GetAllUsersByAdmin(c echo.Context) error {

	role := c.Get("role").(int)

	if role != 2 {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not an admin in GetAllUsers", nil)
	}

	var limit int = 10
	var offSet int = 5

	limitStr := c.QueryParam("limit")
	offSetStr := c.QueryParam("offSet")

	limit, _ = strconv.Atoi(limitStr)
	offSet, _ = strconv.Atoi(offSetStr)

	allUsers, err := models.GetAllUsers(limit, offSet)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to fetch the user", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": allUsers,
	})
}

func GetUserPortfolios(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userId in GetUserPortfolios", nil)
	}

	portfolio, err := models.GetPortFolioByUserId(userId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the portfolio", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"portfolio": portfolio,
	})
}

func GetUserTransactions(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userId in GetUserTransaction", nil)
	}

	transactionsByUserId, err := models.GetTransactionsByUserId(userId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the transactionsById", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": transactionsByUserId,
	})
}

func GetUserNotifications(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userId", nil)
	}

	notificationByUserId, err := models.GetNotificationByUserId(userId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the notification with userid", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notificationByUserId,
	})
}

func UpdateUserVerificationStatus(c echo.Context) error {
	email := c.Get("email").(string)

	if email == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the email", nil)
	}

	if err := models.UpdateUserVerificationStatus(email, true); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to update verification", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": types.StatusOK,
	})
}

/*
GetAddresses - List user addresses
AddAddress - Create new user address
UpdateAddress - Modify existing address
DeleteAddress - Remove address
*/

func GetUserAddress(c echo.Context) error {

	userid := c.Get("userId").(string)

	if userid == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userid", nil)
	}

	address, err := models.GetAddressByUserId(userid)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the address", nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": types.StatusOK,
		"address": address,
	})
}

func AddAddress(c echo.Context) error {

	userid := c.Get("userId").(string)

	if userid == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userid", nil)
	}

	var addressModel dto.AddressModelDTO

	if err := c.Bind(&addressModel); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able bind the address model", nil)
	}

	newAddress := models.AddressModel{
		UserId:  userid,
		Street:  addressModel.StreetName,
		City:    addressModel.City,
		ZipCode: addressModel.ZipCode,
		State:   addressModel.State,
		Country: addressModel.Country,
	}

	address, err := newAddress.CreateAddress()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create address", nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": types.StatusOK,
		"address": address,
	})
}

func UpdateAddress(c echo.Context) error {

	userid := c.Get("userId").(string)

	if userid == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userid", nil)
	}

	var updateRequest struct {
		AddressID  string `json:"addressId"`
		StreetName string `json:"streetName"`
		ZipCode    string `json:"zipCode"`
		City       string `json:"city"`
		State      string `json:"state"`
		Country    string `json:"country"`
	}

	if err := c.Bind(&updateRequest); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able bind the address model", err)
	}

	if updateRequest.AddressID == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "please provide addressId", nil)
	}

	existingAddress, err := models.GetAddressByAddressId(updateRequest.AddressID)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the existing address", err)
	}

	if existingAddress.UserId != userid {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to update the userid", nil)
	}

	existingAddress.Street = updateRequest.StreetName
	existingAddress.ZipCode = updateRequest.ZipCode
	existingAddress.City = updateRequest.City
	existingAddress.State = updateRequest.State
	existingAddress.Country = updateRequest.Country
	existingAddress.UpdatedAt = time.Now()

	var updatedAdd *models.AddressModel

	if updatedAdd, err = models.UpdateAddress(existingAddress); err != nil {
		return util.NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "not able to update the address", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": types.StatusOK,
		"address": updatedAdd,
	})
}

func DeleteAddress(c echo.Context) error {
	userid := c.Get("userId").(string)

	if userid == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the userid", nil)
	}

	id := c.Param("id")
	if id == "" {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "address id not provided", nil)
	}

	if err := models.DeleteAddressByAddressId(id); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the address", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": types.StatusOK,
	})

}
