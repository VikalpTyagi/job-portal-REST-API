package handlers

import (
	"encoding/json"
	"finalAssing/internal/auth"
	"finalAssing/internal/middleware"
	"finalAssing/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func (h *handler) AcceptApplicant(c *gin.Context) {
	ctx := c.Request.Context()
	trackerId, ok := ctx.Value(middleware.TrackerIdKey).(string)
	if !ok {
		log.Error().Msg("TrackerId missing from context in Filter Apllicant")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	_, ok = ctx.Value(auth.AuthKey).(jwt.RegisteredClaims)
	if !ok {
		log.Error().Str("Traker Id", trackerId).Msg("login first")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	var newApplicant []*models.ApplicantReq
	err := json.NewDecoder(c.Request.Body).Decode(&newApplicant)
	if err != nil {
		log.Error().Err(err).Str("tracker Id", trackerId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg2": http.StatusText(http.StatusBadRequest)})
		return
	}
	validate := validator.New()
	for _, data := range newApplicant {
		err = validate.Struct(data)
		if err != nil {
			log.Error().Err(err).Str("tracker Id", trackerId).Interface("body", newApplicant).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "All fields are mandatory"})
			return
		}
	}

	filteredData, err := h.s.FIlterApplication(ctx, newApplicant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, filteredData)
}
