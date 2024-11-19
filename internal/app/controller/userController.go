package controller

import (
	"api-rest-test/internal/app/service"
	"api-rest-test/internal/constants"
	"api-rest-test/internal/models"
	"bytes"

	"api-rest-test/pkg/logger"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Log     logger.Logger
	Service service.IUserService
}

func NewUserController(userService *service.UserService, log *logger.Log) *UserController {
	return &UserController{Log: log,
		Service: userService,
	}
}

func (uc *UserController) Register(ctx *fiber.Ctx) error {
	var err error
	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Category Content-Type header is not application/json.",
		})
	}

	var datReg = models.DataReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//      return err
	//}

	err = eJson.Decode(&datReg)
	if err != nil {
		return err
	}
	sToken, err := uc.Service.Register(ctx.Context(), datReg.User, datReg.Mail, datReg.Tel, datReg.Pwd)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"Estatus": "Error",
			"Msg":     err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"Msg":   "OK",
		"Token": sToken,
	})

}

func (uc *UserController) Login(ctx *fiber.Ctx) error {
	var err error
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header missing",
		})
	}

	// El encabezado debe tener el formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	sToken := parts[1]

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Category Content-Type header is not application/json.",
		})
	}

	var datReg = models.DataReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//      return err
	//}

	err = eJson.Decode(&datReg)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": "Login",
			"msg":   err.Error(),
		})

	}

	err = uc.Service.Login(ctx.Context(), sToken, datReg.Mail, datReg.Pwd)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": "Login",
			"msg":   err.Error(),
		})

	}
	return ctx.Status(201).JSON(fiber.Map{
		"Estatus": "OK",
		"Msg":     "Login ok, token valid",
	})
}

func (uc *UserController) HealthCheck(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodGet {
		return ctx.Status(fiber.StatusMethodNotAllowed).SendString("Method Not Allowed")
	}
	healthCheck := models.NewHealthCheck("1", constants.HealthPass)
	return ctx.Status(fiber.StatusOK).JSON(healthCheck)
}
