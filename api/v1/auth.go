package v1

import (
	"context"
	"encoding/json"
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
	userData, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	err = rdb.Set(context.Background(), "user_"+req.Email, string(userData), time.Minute*10).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	go func() {
		randomCode := passwordgenerator.EncodeToString(6)
		err := email.SendEmail(h.cfg, &email.SendEmailRequest{
			To:      []string{req.Email},
			Subject: "Verification Email",
			Body: map[string]string{
				"code": randomCode,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		err = rdb.Set(context.Background(), "code_"+req.Email, randomCode, time.Minute*3).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}()

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

	userData, err := rdb.Get(context.Background(), "user_"+req.Email).Result()
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	var user models.User
	err = json.Unmarshal([]byte(userData), &user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if user.Email != req.Email {
		ctx.JSON(http.StatusBadGateway, models.ResponseError{
			Error: "Email is not correct",
		})
		return
	}

	code, err := rdb.Get(context.Background(), "code_"+req.Email).Result()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ResponseError{
			Error: "Code is not correct",
		})
		return
	}
	
	if code != req.Code {
		ctx.JSON(http.StatusBadGateway, models.ResponseError{
			Error: "Code is not correct",
		})
		return
	}

	user2, err := h.storage.User().Create(&repo.User{
		FirstName: user.FirstName,
		Email:     user.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.GetUser{
		ID:        int(user2.ID),
		CreatedAt: user2.CreatedAt,
	})
}
