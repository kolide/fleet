package errors

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

type KolideError struct {
	Err            error
	StatusCode     int
	PublicMessage  string
	PrivateMessage string
}

func (e *KolideError) Error() string {
	return e.PublicMessage
}

func NewFromError(err error, status int, publicMessage string) *KolideError {
	return &KolideError{
		Err:            err,
		StatusCode:     status,
		PublicMessage:  publicMessage,
		PrivateMessage: err.Error(),
	}
}

const StatusUnprocessableEntity = 422

func ReturnError(c *gin.Context, err error) {
	switch typedErr := err.(type) {
	case *KolideError:
		c.JSON(typedErr.StatusCode,
			gin.H{"message": typedErr.PublicMessage})
		logrus.WithError(typedErr.Err).Debug(typedErr.PrivateMessage)

	case validator.ValidationErrors:
		errors := make([](map[string]string), 0, len(typedErr))
		for _, fieldErr := range typedErr {
			m := make(map[string]string)
			m["field"] = fieldErr.Name
			m["code"] = "invalid"
			m["message"] = fieldErr.Tag
			errors = append(errors, m)
		}

		c.JSON(StatusUnprocessableEntity,
			gin.H{"message": "Validation error",
				"errors": errors})
		logrus.WithError(typedErr).Debug("Validation error")

	default:
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "Unspecified error"})
		logrus.WithError(typedErr).Debug("Unspecified error")
	}
}
