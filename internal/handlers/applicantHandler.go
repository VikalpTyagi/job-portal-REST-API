package handlers

import (
	"encoding/json"
	"finalAssing/internal/middleware"
	"finalAssing/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *handler) FilterApplicants(c *gin.Context) {
	ctx := c.Request.Context()
	trackerId, ok := ctx.Value(middleware.TrackerIdKey).(string)
	if !ok {
		log.Error().Msg("TrackerId missing from context in Filter Apllicant")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var newApplicant []models.Applicant

	err := json.NewDecoder(c.Request.Body).Decode(&newApplicant)
	if err != nil {
		log.Error().Err(err).Str("tracker Id", trackerId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

}