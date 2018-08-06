package gin_http

import (
	"net/http"
	"web-service-users/src/model"
	"web-service-users/src/user"

	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserUseCase user.UseCase
}

func NewUserGinHTTPHandler(e *gin.Engine, uc user.UseCase) {
	handler := &UserHandler{
		UserUseCase: uc,
	}
	routeGroupV1 := e.Group("/v1")
	routeGroupV1.POST("/user", handler.Store)
	routeGroupV1.GET("/user/:user_name", handler.GetByUserName)
	routeGroupV1.DELETE("/user/:user_name", handler.DeleteByUserName)
	routeGroupV1.GET("/user", handler.FetchAllUsers)
	routeGroupV1.PUT("/user/kudos", handler.UpdateKudosQuantity)
}

type ResponseError struct {
	Message string `json:"message"`
}

type StoreRequest struct {
	UserName string `form:"userName" binding:"required" json:"user_name"`
	Name     string `form:"name" binding:"required" json:"name"`
}

func (u *UserHandler) Store(c *gin.Context) {

	var requestData StoreRequest
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	userModel := model.User{
		UserName: requestData.UserName,
		Name:     requestData.Name,
	}
	saveUser, err := u.UserUseCase.Store(&userModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("User added")
	c.JSON(http.StatusOK, saveUser)

}

func (u *UserHandler) GetByUserName(c *gin.Context) {

	userName := c.Param("user_name")
	userFound, err := u.UserUseCase.GetByUserName(userName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("User showed")
	c.JSON(http.StatusOK, userFound)

}

func (u *UserHandler) DeleteByUserName(c *gin.Context) {

	userName := c.Param("user_name")
	err := u.UserUseCase.DeleteByUserName(userName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("User deleted")
	c.JSON(http.StatusOK, nil)

}

func (u *UserHandler) FetchAllUsers(c *gin.Context) {

	pageSizeString := c.DefaultQuery("pageSize", "1")
	pageSize, err := strconv.ParseInt(pageSizeString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	numberPageString := c.DefaultQuery("numberPage", "10")
	numberPage, err := strconv.ParseInt(numberPageString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	users, err := u.UserUseCase.FetchAllUsers(pageSize, numberPage)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Fetch all users")
	c.JSON(http.StatusOK, users)

}

type UpdateRequest struct {
	UserName string `form:"userName" binding:"required" json:"user_name"`
	Quantity int    `form:"quantity" binding:"required" json:"quantity"`
}

func (u *UserHandler) UpdateKudosQuantity(c *gin.Context) {

	var requestData UpdateRequest
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	err = u.UserUseCase.UpdateQuantityKudos(requestData.UserName, requestData.Quantity)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Quantity kudos updated")
	c.JSON(http.StatusOK, nil)

}
