package service

import (
	"api-rest-test/internal/app/repository"
	"api-rest-test/internal/models"
	"api-rest-test/pkg/logger"
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v4"
)

type IUserService interface {
	Login(context.Context, string, string, string) error
	Register(context.Context, string, string, string, string) (string, error)
	Search(context.Context, string, string) ([]map[string]interface{}, error)
}

type UserService struct {
	Log        logger.Logger
	repository repository.UserRepository
}

func NewUserService(repo *repository.UserRepository, log *logger.Log) *UserService {
	return &UserService{Log: log, repository: *repo}
}

func (us *UserService) Register(ctx context.Context, User string, Mail string, Tel string, Pwd string) (string, error) {
	//Validate data
	sErr, err := validationInfo(User, Mail, Tel, Pwd)
	if err != nil {
		err = errors.New(sErr)
		return "", err
	}
	//Validate Mail
	bOK := validationEmail(Mail)
	if bOK == false {
		err = errors.New("El correo no es valido.")
		return "", err
	}

	bOK = validationPwd(Pwd)
	if bOK == false {
		err = errors.New("Debe ser 6 a 12 caracteres al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.")
		return "", err
	}
	bOK = validationTel(Tel)
	if bOK == false {
		err = errors.New("Teléfono debe ser de 10 digitos.")
		return "", err
	}
	// Generar el JWT
	token, err := GenerateJWT(Mail, Pwd)
	if err != nil {
		return "", errors.New("Token no valido")
	}
	log.Println(token)
	//Checa si el user exist
	OK, err := us.repository.UserExist(ctx, Mail, Tel)
	if OK == true && err != nil {
		logger.GetLogger().Error("userService", "El usuario ya existe", err)
		return "", err
	}

	err = us.repository.Register(ctx, User, Mail, Tel, Pwd)
	if err != nil {
		logger.GetLogger().Error("userService", "Error register user", err)
		return "", err
	}
	return token, nil
}

func (us *UserService) Login(ctx context.Context, sToken string, Mail string, Pwd string) error {
	//Validate data
	sErr, err := validationInfoLogin(Mail, Pwd)
	if err != nil {
		err = errors.New(sErr)
		return err
	}
	//Validate Mail
	bOK := validationEmail(Mail)
	if bOK == false {
		err = errors.New("El correo no es valido.")
		return err
	}

	bOK = validationPwd(Pwd)
	if bOK == false {
		err = errors.New("Debe ser 6 a 12 caracteres al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.")
		return err
	}

	// Generar el JWT
	token, err := DecodeJWT_Data(sToken, Pwd)
	if err != nil {
		return errors.New("Token no valido")
	}
	log.Println(token)

	OK, err := us.repository.UserExist(ctx, Mail, Pwd)
	if OK == false && err != nil {
		//logger.GetLogger().Error("userService", "No fue posible auntenticar.", err)
		return err
	} else if OK == true && err == nil {
		//logger.GetLogger().Error("userService", "Login Ok", err)
		return nil
		//Continue with process
	}
	return nil
}

func (us *UserService) Search(ctx context.Context, Mail string, Tel string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	var err error
	results, err = us.repository.Search(ctx, Mail, Tel)

	if err != nil {
		logger.GetLogger().Error("CensoService", "Error executing query", err)
		return nil, err
	}
	return results, nil
}

func validationInfo(User string, Mail string, Tel string, Pwd string) (string, error) {
	var err error
	Detalle := ""

	if User == "" {
		Detalle += "Se requiere un usuario. "
		err = errors.New("Datos incompletos")
	}
	if Mail == "" {
		Detalle += "Se requiere un correo."
		err = errors.New("Datos incompletos")
	}
	if Tel == "" {
		Detalle += "Se requiere un teléfono."
		err = errors.New("Datos incompletos")
	}
	if Pwd == "" {
		Detalle += "Se requiere una contraseña."
		err = errors.New("Datos incompletos")
	}

	return Detalle, err
}

func validationInfoLogin(Mail string, Pwd string) (string, error) {
	var err error
	Detalle := ""

	if Mail == "" {
		Detalle += "Se requiere un correo."
		err = errors.New("Datos incompletos")
	}
	if Pwd == "" {
		Detalle += "Se requiere una contraseña."
		err = errors.New("Datos incompletos")
	}

	return Detalle, err
}
func validationEmail(email string) bool {
	// Usamos mail.ParseAddress para validar el formato básico del correo electrónico
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func validationPwd(pwd string) bool {
	//De 6 a 12 caractres
	//al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.
	var regex = `^[A-Za-z0-9@$&]{6,12}$`
	re := regexp.MustCompile(regex)

	// Verificamos que la contraseña cumpla con la longitud y los caracteres permitidos
	if !re.MatchString(pwd) {
		return false
	}

	// Comprobamos que la contraseña tenga al menos una minúscula, una mayúscula, un número y un carácter especial
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range pwd {
		if unicode.IsLower(ch) {
			hasLower = true
		}
		if unicode.IsUpper(ch) {
			hasUpper = true
		}
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if strings.ContainsRune("@$&", ch) {
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func validationTel(phone string) bool {
	//Long 10 numbers
	if len(phone) != 10 {
		return false
	}

	// Only numbers
	for _, char := range phone {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	return true
}

func GenerateJWT(mail string, key string) (string, error) {
	var JwtSecret = []byte(key)
	// Crear los claims del JWT
	claims := models.Claims{
		Username: mail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expired in 24 horas
			Issuer:    "HugoBH",
		},
	}

	//Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signed token with key
	signedToken, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func DecodeJWT_Data(sToken string, clave string) (*models.Claims, error) {
	var JwtSecret = []byte(clave)

	// Claim for the token
	var claims models.Claims

	// Get the token and check the signed
	token, err := jwt.ParseWithClaims(sToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma inválido: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error al validar el JWT: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token no válido")
	}

	return &claims, nil
}
