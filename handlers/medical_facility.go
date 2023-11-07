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

// API for fetching the medical facility details related to a user, if exists
func (apiCfg *ApiCfg) GetMedicalFacilityDetails(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	dbMedicalFacility, err := apiCfg.DB.GetMedicalFacilityByUserId(c, dbUser.ID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Error while fetching facility details or facility does not exists")
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseMedicalFacilityToMedicalFacility(dbMedicalFacility))
}

// API for adding facility
func (apiCfg *ApiCfg) AddMedicalFacility(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	type Parameters struct {
		Type         string  `json:"type" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		Description  string  `json:"description" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		MobileNumber string  `json:"mobile_number" binding:"required"`
		Charges      float64 `json:"charges" binding:"required"`
		Address      string  `json:"address" binding:"required"`
		Location     struct {
			Lat float64 `json:"lat" binding:"required"`
			Lng float64 `json:"lng" binding:"required"`
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

	_, err = apiCfg.DB.CreateMedicalFacility(c, database.CreateMedicalFacilityParams{
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
		UserID:        dbUser.ID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error caught while creating medical facility: %v", err))
		return
	}

	SuccessResponse(c, http.StatusCreated, "Medical Facility Added Successfully!", params)
}
