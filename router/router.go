package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobingge/dbger/app/http/controllers/auth"
	"github.com/xiaobingge/dbger/app/http/controllers/home"
	"github.com/xiaobingge/dbger/app/http/middleware"
	"net/http"
)
// Load loads the middlewares, routes, handles.
func Load(g *gin.Engine) *gin.Engine {
	loadHealthTest(g)
	loadAPI(g)
	return g
}


// loadAPI load api part
func loadAPI(g *gin.Engine) *gin.Engine {
	// Group for api
	api := g.Group("/api")
	api.GET("/token", auth.Login)
	api.GET("/create_user", auth.Register)
	api.GET("/home",middleware.JWTAuthMiddleware(), home.Home)
	return g
}



// loadHelthTest the health check handlers
func loadHealthTest(g *gin.Engine) *gin.Engine {
	// Group for health check
	svcd := g.Group("/check")
	{
		svcd.GET("/health", func(c *gin.Context){
			c.String(http.StatusOK, "ok")
		})
	}
	return g
}