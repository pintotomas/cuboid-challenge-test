package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchCuboid(cuboidID string) (*models.Cuboid, error) {
	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not Found")
		} else {
			return nil, r.Error
		}
	}
	return &cuboid, nil
}

func GetCuboid(c *gin.Context) {

	cuboidID := c.Param("cuboidID")

	cuboid, err := FetchCuboid(cuboidID)

	if err != nil {
		if err.Error() == "not Found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, cuboid)
}

func ListCuboids(c *gin.Context) {
	var cuboids []models.Cuboid
	if r := db.CONN.Find(&cuboids); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, cuboids)
}

func CreateCuboid(c *gin.Context) {
	var cuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := c.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboid := models.Cuboid{
		Width:  cuboidInput.Width,
		Height: cuboidInput.Height,
		Depth:  cuboidInput.Depth,
		BagID:  cuboidInput.BagID,
	}

	bag, err := FetchBag(strconv.FormatUint(uint64(cuboidInput.BagID), 10))
	if err != nil {
		if err.Error() == "not Found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if bag.Disabled == true {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return
	}
	if cuboid.PayloadVolume() > bag.AvailableVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return
	}

	if r := db.CONN.Create(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusCreated, &cuboid)
}

func UpdateCuboid(c *gin.Context) {
	var updateCuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := c.BindJSON(&updateCuboidInput); err != nil {
		return
	}

	cuboidID := c.Param("cuboidID")
	cuboid, err := FetchCuboid(cuboidID)
	if err != nil {
		if err.Error() == "not Found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	bag, err := FetchBag(strconv.FormatUint(uint64(updateCuboidInput.BagID), 10))
	if err != nil {
		if err.Error() == "not Found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	cuboid.Depth = updateCuboidInput.Depth
	cuboid.Height = updateCuboidInput.Height
	cuboid.Width = updateCuboidInput.Width
	cuboid.BagID = updateCuboidInput.BagID

	if bag.Disabled == true {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return
	}
	if cuboid.PayloadVolume() > bag.AvailableVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return
	}

	if r := db.CONN.Save(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, &cuboid)

}

func DeleteCuboid(c *gin.Context) {
	cuboidID := c.Param("cuboidID")

	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&cuboid); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
