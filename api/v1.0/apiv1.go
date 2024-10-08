package apiv1

import (
	"github.com/WinterSunset95/WinterMediaBackend/api/v1.0/auth"
	"github.com/WinterSunset95/WinterMediaBackend/api/v1.0/movies"
	"github.com/WinterSunset95/WinterMediaBackend/api/v1.0/user"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		movies.ApplyRoutes(v1)
		auth.ApplyRoutes(v1)
		user.ApplyRoutes(v1)
	}
}
