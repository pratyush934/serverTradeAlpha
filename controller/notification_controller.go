package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/dto"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
	"github.com/pratyush934/tradealpha/server/util"
)

func AddNotification(c echo.Context) error {

	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "userid is blank in AddNotification", nil)
	}

	var notification dto.NotificationDTO
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "not able to bind the notification /AddNotification",
			"error":   err,
		})
	}

	newNotification := models.NotificationModel{
		UserId:     userId,
		Message:    notification.Message,
		ReadStatus: notification.ReadStatus,
	}

	createNotification, err := newNotification.CreateNotification()

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to create notification in AddNotification", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notification": createNotification,
	})

}

func GetNotification(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "not able to get the userId", nil)
	}

	noticeId := c.Param("id")

	notificationByNotificationId, err := models.GetNotificationByNotificationId(noticeId)

	if err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to get the notificationById", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notification": notificationByNotificationId,
	})
}

func DeleteNotification(c echo.Context) error {
	userId := c.Get("userId").(string)

	if userId == "" {
		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "userId is not available", nil)
	}

	noticeId := c.Param("id")

	if err := models.DeleteNotificationByNId(noticeId); err != nil {
		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to delete the notification", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Notification Deleted successfully",
	})

}

//func UpdateNotification(c echo.Context) error {
//
//	userId := c.Get("userId").(string)
//
//	if userId == "" {
//		return util.NewAppError(http.StatusUnauthorized, types.StatusUnauthorized, "userId is not available in UpdateNotification", nil)
//	}
//
//	noticeId := c.Param("id")
//
//	var notification dto.NotificationDTO
//
//	if err := c.Bind(&notification); err != nil {
//		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able to bind the notification", err)
//	}
//
//	newNotice := models.NotificationModel{
//		Id:         noticeId,
//		Message:    notification.Message,
//		ReadStatus: notification.ReadStatus,
//	}
//
//	if err := models.UpdateNotification(&newNotice); err != nil {
//		return util.NewAppError(http.StatusBadRequest, types.StatusBadRequest, "not able the update notification", err)
//	}
//
//	return c.JSON(http.StatusOK, map[string]interface{}{
//		"message": "we did it",
//	})
//
//}
