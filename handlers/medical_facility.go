package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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

	description := sql.NullString{
		String: params.Description,
		Valid:  true,
	}

	_, err = apiCfg.DB.CreateMedicalFacility(c, database.CreateMedicalFacilityParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		ModifiedAt:    time.Now().UTC(),
		Type:          database.FacilityType(params.Type),
		Name:          params.Name,
		Description:   description,
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

// API for updating medical facility details
func (apiCfg *ApiCfg) UpdateMedicalFacility(c *gin.Context) {
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

	type Parameters struct {
		Type         string  `json:"type"`
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		Email        string  `json:"email"`
		MobileNumber string  `json:"mobile_number"`
		Charges      float64 `json:"charges"`
		Address      string  `json:"address"`
		Location     struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	}

	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request: %v", err.Error()))
		return
	}

	if params.Type == "" {
		params.Type = string(dbMedicalFacility.Type)
	}

	if params.Name == "" {
		params.Name = dbMedicalFacility.Name
	}

	if params.Description == "" && dbMedicalFacility.Description.String != "" {
		params.Description = dbMedicalFacility.Description.String
	}

	if params.Email == "" {
		params.Email = dbMedicalFacility.Email
	}

	if !validators.EmailValidator(params.Email) {
		ErrorResponse(c, http.StatusBadRequest, "Invalid email address")
		return
	}

	if params.MobileNumber == "" {
		params.MobileNumber = dbMedicalFacility.MobileNumber
	}

	if err := validators.MobileNumberValidator(params.MobileNumber); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid mobile number: %v", err))
		return
	}

	if params.Charges == 0 {
		chargesFloatValue, err := strconv.ParseFloat(dbMedicalFacility.Charges, 64)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
			return
		}
		params.Charges = chargesFloatValue
	}

	if params.Address == "" {
		params.Address = dbMedicalFacility.Address
	}

	if params.Location.Lat == 0 && params.Location.Lng == 0 {
		params.Location.Lat = dbMedicalFacility.Lat.(float64)
		params.Location.Lng = dbMedicalFacility.Lng.(float64)
	}

	description := sql.NullString{
		String: params.Description,
		Valid:  true,
	}

	_, err = apiCfg.DB.UpdateMedicalFacility(c, database.UpdateMedicalFacilityParams{
		ID:            dbMedicalFacility.ID,
		Type:          database.FacilityType(params.Type),
		Name:          params.Name,
		Description:   description,
		Email:         params.Email,
		MobileNumber:  params.MobileNumber,
		Charges:       fmt.Sprintf("%.2f", params.Charges),
		Address:       params.Address,
		StMakepoint:   params.Location.Lat,
		StMakepoint_2: params.Location.Lng,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Medical Facility Details Updated Successfully!", params)
}
