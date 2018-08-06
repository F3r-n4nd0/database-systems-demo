package gin_http

import (
	"net/http"
	"strconv"
	"web-service-kudos/src/kudos"

	"web-service-kudos/src/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type kudosHandler struct {
	KudosUseCase kudos.UseCase
}

func NewKudosGinHTTPHandler(e *gin.Engine, kc kudos.UseCase) {
	handler := &kudosHandler{
		KudosUseCase: kc,
	}
	routeGroupV1 := e.Group("/v1")
	routeGroupV1.POST("/kudos", handler.Store)
	routeGroupV1.GET("/kudos/:id", handler.GetByID)
	routeGroupV1.DELETE("/kudos/:id", handler.DeleteByID)
	routeGroupV1.GET("/kudos", handler.FetchAllKudos)
	routeGroupV1.GET("/quantity/kudos/:user_name", handler.GetQuantityByUserName)
}

type ResponseError struct {
	Message string `json:"message"`
}

type StoreRequest struct {
	FromUserName string `form:"fromUserName" binding:"required" json:"from_user_name"`
	ToUserName   string `form:"toUserName" binding:"required" json:"to_user_name"`
	Message      string `form:"message" binding:"required" json:"message"`
}

func (u *kudosHandler) Store(c *gin.Context) {

	var requestData StoreRequest
	err := c.Bind(&requestData)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	kudosModel := model.Kudos{
		FromUserName: requestData.FromUserName,
		ToUserName:   requestData.ToUserName,
		Message:      requestData.Message,
	}
	saveKudos, err := u.KudosUseCase.Store(&kudosModel)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Kudos added")
	c.JSON(http.StatusOK, saveKudos)

}

func (u *kudosHandler) GetByID(c *gin.Context) {

	id := c.Param("id")
	kudosFound, err := u.KudosUseCase.GetByID(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Kudos showed")
	c.JSON(http.StatusOK, kudosFound)

}

func (u *kudosHandler) DeleteByID(c *gin.Context) {

	id := c.Param("id")
	err := u.KudosUseCase.DeleteByID(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Kudos deleted")
	c.JSON(http.StatusOK, nil)

}

func (u *kudosHandler) FetchAllKudos(c *gin.Context) {

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
	allKudos, err := u.KudosUseCase.FetchAllKudos(pageSize, numberPage)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Fetch all kudos")
	c.JSON(http.StatusOK, allKudos)

}

func (u *kudosHandler) GetQuantityByUserName(c *gin.Context) {

	userName := c.Param("user_name")
	quantity, err := u.KudosUseCase.GetQuantityByUserName(userName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	log.Info("Get kudos quantity")
	c.JSON(http.StatusOK, gin.H{
		"quantity": quantity,
	})

}
