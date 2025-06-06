package repository

import (
	"context"
	"fmt"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RewildingPhotoRepository struct{}

// @Summary Rewilding
// @Description Retrieve all rewilding
// @ID rewilding
// @Produce json
// @Tags Rewilding
// @Success 200 {object} []models.Rewilding
// @Failure 400 {object} structs.Message
// @Router /rewilding [get]
func (r RewildingPhotoRepository) Retrieve(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var Rewilding models.Rewilding
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding.RewildingPhotos)
}

func (r RewildingPhotoRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	photosId, _ := primitive.ObjectIDFromHex(c.Param("photosId"))
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_photos._id", Value: photosId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	idx := 0
	for p, v := range Rewilding.RewildingPhotos {
		if v.RewildingPhotosID == photosId {
			idx = p
			break
		}
	}

	// 直接使用 RewildingPhotosPath 進行重定向
	if Rewilding.RewildingPhotos[idx].RewildingPhotosPath != "" {
		c.Redirect(http.StatusFound, Rewilding.RewildingPhotos[idx].RewildingPhotosPath)
		return
	}

	// 如果沒有路徑，返回錯誤
	c.JSON(http.StatusNotFound, gin.H{"message": "Image URL not found"})
}

func (r RewildingPhotoRepository) FromGoogle(c *gin.Context) {
	reference := c.Param("referenceId")
	place := c.Param("placesId")

	photo := helpers.GooglePlacePhoto(c, fmt.Sprintf("places/%s/photos/%s", place, reference))
	if photo == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusSeeOther, photo.PhotoUri)
}
