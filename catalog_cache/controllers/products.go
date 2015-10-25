package controllers

import (
	"github.com/FoxComm/core_services/catalog_cache/models"
	"github.com/gin-gonic/gin"
	"github.com/FoxComm/libs/utils"

	"net/http"
	"net/url"
)

func Products(c *gin.Context) {
	c.Request.ParseForm()

	solrHost := utils.GetSolrHost(c)
	if resp, err := models.Search(c.Request.Form, solrHost); err == nil {
		if resp.Results == nil {
			resp.Results = []models.Product{}
		}

		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func Product(c *gin.Context) {
	var params = url.Values{}
	params.Add("slug", c.Params.ByName("slug"))

	solrHost := utils.GetSolrHost(c)
	if resp, err := models.Search(params, solrHost); err == nil {
		if len(resp.Results) > 0 {
			c.JSON(http.StatusOK, resp.Results[0])
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
