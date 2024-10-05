package api

import (
	apiv1 "github.com/WinterSunset95/WinterMediaBackend/api/v1.0"
	"github.com/gin-gonic/gin"
)

func ApiRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
