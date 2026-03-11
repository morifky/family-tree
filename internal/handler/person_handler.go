package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"brayat/internal/model"
	"brayat/internal/storage"
)

// PersonHandler handles HTTP requests for Person resources
type PersonHandler struct {
	service      model.PersonService
	photoStorage storage.PhotoStorage
}

// NewPersonHandler creates a new handler serving person requests
func NewPersonHandler(svc model.PersonService, photoStorage storage.PhotoStorage) *PersonHandler {
	return &PersonHandler{service: svc, photoStorage: photoStorage}
}

// POST /sessions/:id/people
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	sessionID := c.Param("id")

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	name := c.PostForm("name")
	if name == "" {
		ErrorResponse(c, http.StatusBadRequest, "Name is required")
		return
	}

	nicknameForm := c.PostForm("nickname")
	var nickname *string
	if nicknameForm != "" {
		nickname = &nicknameForm
	}

	genderForm := c.PostForm("gender")
	if genderForm == "" {
		ErrorResponse(c, http.StatusBadRequest, "Gender is required")
		return
	}

	var photoPath *string
	file, err := c.FormFile("photo")
	if err == nil {
		// File was provided
		filename, saveErr := h.photoStorage.SavePhoto(file)
		if saveErr != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to save photo")
			return
		}
		photoPath = &filename
	}

	person, err := h.service.CreatePerson(c.Request.Context(), sessionID, name, nickname, model.Gender(genderForm), photoPath)
	if err != nil {
		// Cleanup the saved file if DB creation fails
		if photoPath != nil {
			_ = h.photoStorage.DeletePhoto(*photoPath)
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create person")
		return
	}

	CreatedResponse(c, person)
}

// PUT /people/:id
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	name := c.PostForm("name")
	if name == "" {
		ErrorResponse(c, http.StatusBadRequest, "Name is required")
		return
	}

	nicknameForm := c.PostForm("nickname")
	var nickname *string
	if nicknameForm != "" {
		nickname = &nicknameForm
	}

	genderForm := c.PostForm("gender")
	if genderForm == "" {
		ErrorResponse(c, http.StatusBadRequest, "Gender is required")
		return
	}

	var photoPath *string

	// Check if user explicitly wants to remove the photo
	removePhoto := c.PostForm("remove_photo")
	if removePhoto == "true" {
		empty := ""
		photoPath = &empty
	} else {
		file, err := c.FormFile("photo")
		if err == nil {
			// New file was provided
			filename, saveErr := h.photoStorage.SavePhoto(file)
			if saveErr != nil {
				ErrorResponse(c, http.StatusInternalServerError, "Failed to save photo")
				return
			}
			photoPath = &filename
		}
	}

	err := h.service.UpdatePerson(c.Request.Context(), id, name, nickname, model.Gender(genderForm), photoPath)
	if err != nil {
		// Cleanup if update fails and we just saved a new photo
		// (Only clean up if it's a valid new path)
		if photoPath != nil && *photoPath != "" {
			_ = h.photoStorage.DeletePhoto(*photoPath)
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update person")
		return
	}

	SuccessResponse(c, gin.H{"message": "Person updated successfully"})
}

// DELETE /people/:id
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeletePerson(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete person")
		return
	}

	SuccessResponse(c, gin.H{"message": "Person deleted successfully"})
}

// GET /sessions/:id/people
func (h *PersonHandler) GetPeople(c *gin.Context) {
	sessionID := c.Param("id")

	people, err := h.service.GetPeopleBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve people")
		return
	}

	SuccessResponse(c, people)
}
