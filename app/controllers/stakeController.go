package controllers

import (
	"strconv"

	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/services"
	"github.com/gin-gonic/gin"
)

type StakeAPI interface {
	AddStake(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetList(c *gin.Context)
}

type stakeAPI struct{
	stakeService services.StakeholderServices
}

func NewStakeAPI(stakeRepo services.StakeholderServices) *stakeAPI{
	return &stakeAPI{stakeRepo}
}

func (s *stakeAPI) AddStake(c *gin.Context){
	var stake models.Stakeholder

	if err := c.ShouldBindJSON(&stake); err != nil{
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err := s.stakeService.Store(&stake)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success create new stakeholder",
		"data" : stake,
	})	

}

func (s *stakeAPI) Update(c *gin.Context){

	stakeID := c.Param("id")

	if stakeID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}
	
	id, err := strconv.Atoi(stakeID)
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	var newStake models.Stakeholder

	if err := c.ShouldBindJSON(&newStake); err != nil{
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	newStake.ID = id

	err = s.stakeService.Update(id, newStake)

	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success update stakeholder",
		"data" : newStake,
	})
}

func (s *stakeAPI) Delete(c *gin.Context){
	stakeID := c.Param("id")
	
	if stakeID == "" {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	id, err := strconv.Atoi(stakeID)
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	err = s.stakeService.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "success delete stake",
	})
}

func (s *stakeAPI) GetByID(c *gin.Context){
	stakeID, err := strconv.Atoi(c.Param("id"))
	if stakeID == 0 {
		c.JSON(400, gin.H{
			"message" : "data notfound",
		})
		return
	}
	if err != nil {
		c.JSON(400, gin.H{
			"message" : "invalid request body",
		})
		return
	}

	result, err := s.stakeService.GetByID(stakeID)	
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}

	c.JSON(200, result)
}

func (s *stakeAPI) GetList(c *gin.Context){

	result, err := s.stakeService.GetList()
	if err != nil {
		c.JSON(500, gin.H{
			"message" : "internal server error",
		})
		return
	}
	c.JSON(200, result)
}