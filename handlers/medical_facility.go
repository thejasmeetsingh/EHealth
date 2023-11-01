package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
	"github.com/thejasmeetsingh/EHealth/validators"
)

func (apiCfg *ApiCfg) AddMedicalFacility(c *gin.Context, dbUser database.User) {
	type Parameters struct {
		Type         string  `json:"type" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		Description  string  `json:"description" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		MobileNumber string  `json:"mobile_number" binding:"required"`
		Charges      float64 `json:"charges" binding:"required"`
		Address      string  `json:"address" binding:"required"`
		Location     struct {
			Lat float64 `json:"lat" binding:"required,latitude"`
			Lng float64 `json:"lng" binding:"required,longitude"`
		} `json:"location" binding:"required"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request: %v", err.Error()))
		return
	}

	if err := validators.MobileNumberValidator(params.MobileNumber); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid mobile number: %v", err))
		return
	}

	dbMedicalFacility, err := apiCfg.DB.CreateMedicalFacility(c, database.CreateMedicalFacilityParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		Type:       database.FacilityType(params.Type),
		Name:       params.Name,
		Description: sql.NullString{
			String: params.Description,
			Valid:  true,
		},
		Email:         params.Email,
		MobileNumber:  params.MobileNumber,
		Charges:       fmt.Sprintf("%.2f", params.Charges),
		Address:       params.Address,
		StMakepoint:   params.Location.Lat,
		StMakepoint_2: params.Location.Lng,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while creating medical facility: %v", err))
		return
	}

	SuccessResponse(c, http.StatusCreated, "Medical Facility Added Successfully!", models.DatabaseMedicalFacilityToMedicalFacility(dbMedicalFacility))
}
