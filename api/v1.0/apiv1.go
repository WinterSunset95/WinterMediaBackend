package apiv1

import (
	"github.com/WinterSunset95/WinterMediaBackend/api/v1.0/movies"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		movies.ApplyRoutes(v1)
	}
}
