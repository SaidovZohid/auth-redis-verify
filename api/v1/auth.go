package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/SaidovZohid/auth-redis-verify/api/models"
	"github.com/SaidovZohid/auth-redis-verify/pkg/email"
	passwordgenerator "github.com/SaidovZohid/auth-redis-verify/pkg/passwordGenerator"
	"github.com/SaidovZohid/auth-redis-verify/storage/repo"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func (h *handler) Register(ctx *gin.Context) {
	var (
		req models.User
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	randomCode := passwordgenerator.EncodeToString(6)
	go func() {
		email.SendEmail(h.cfg, &email.SendEmailRequest{
			To:      []string{req.Email},
			Subject: "Verification Email",
			Body: map[string]string{
				"code": randomCode,
			},
		})
	}()
	rdb.Set(context.Background(), "code", randomCode, time.Minute*2)
	rdb.Set(context.Background(), "email", req.Email, time.Minute*2)
	rdb.Set(context.Background(), "first_name", req.FirstName, time.Minute*2)

	ctx.JSON(http.StatusOK, models.ResponseOK{
		Message: "Verification code has been sent to your email address. Code is valid for 2 minutes.",
	})
}

func (h *handler) Verify(ctx *gin.Context) {
	var (
		req models.VerifyEmail
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	code, err := rdb.Get(context.Background(), "code").Result()
	email, err := rdb.Get(context.Background(), "email").Result()
	name, err := rdb.Get(context.Background(), "first_name").Result()
	if err != nil {
		ctx.JSON(http.StatusGatewayTimeout, errResponse(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusGatewayTimeout, errResponse(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusGatewayTimeout, errResponse(err))
		return
	}

	if email != req.Email {
		ctx.JSON(http.StatusBadGateway, models.ResponseError{
			Error: "Email is not correct",
		})
		return
	}
	if code != req.Code {
		ctx.JSON(http.StatusBadGateway, models.ResponseError{
			Error: "Code is not correct",
		})
		return
	}

	user, err := h.storage.User().Create(&repo.User{
		FirstName: name,
		Email:     email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.GetUser{
		ID:        int(user.ID),
		CreatedAt: user.CreatedAt,
	})
}
