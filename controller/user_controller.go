package controller

import (
	"net/http"
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
