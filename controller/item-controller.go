package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"basic-rest-api-jwt-mysql/dto"
	"basic-rest-api-jwt-mysql/entity"
	"basic-rest-api-jwt-mysql/helper"
	"basic-rest-api-jwt-mysql/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ItemController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type itemController struct {
	itemService service.ItemService
	jwtService  service.JWTService
}

func NewItemController(itemServ service.ItemService, jwtServ service.JWTService) ItemController {
	return &itemController{
		itemService: itemServ,
		jwtService:  jwtServ,
	}
}

func (c *itemController) All(context *gin.Context) {
	var items []entity.Item = c.itemService.All()
	res := helper.BuildResponse(true, "OK", items)
	context.JSON(http.StatusOK, res)
}

func (c *itemController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var item entity.Item = c.itemService.FindByID(id)
	if (item == entity.Item{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", item)
		context.JSON(http.StatusOK, res)
	}
}

func (c *itemController) Insert(context *gin.Context) {
	var itemCreateDTO dto.ItemCreateDTO
	errDTO := context.ShouldBind(&itemCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			itemCreateDTO.UserID = convertedUserID
		}
		result := c.itemService.Insert(itemCreateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *itemController) Update(context *gin.Context) {
	var itemUpdateDTO dto.ItemUpdateDTO
	errDTO := context.ShouldBind(&itemUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.itemService.IsAllowedToEdit(userID, itemUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			itemUpdateDTO.UserID = id
		}
		result := c.itemService.Update(itemUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *itemController) Delete(context *gin.Context) {
	var item entity.Item
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	item.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.itemService.IsAllowedToEdit(userID, item.ID) {
		c.itemService.Delete(item)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *itemController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
