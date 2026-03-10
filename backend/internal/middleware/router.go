package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pizzattolog/licencas/internal/auth"
	"github.com/pizzattolog/licencas/internal/handlers"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	licencaHandler *handlers.LicencaHandler,
	dashboardHandler *handlers.DashboardHandler,
	jwtService *auth.JWTService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:3002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "pizzattolog-api"})
	})

	v1 := r.Group("/api/v1")
	{
		// Rotas públicas
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// Rotas protegidas
		protected := v1.Group("")
		protected.Use(AuthMiddleware(jwtService))
		{
			protected.POST("/auth/logout", authHandler.Logout)

			// Dashboard
			protected.GET("/dashboard", dashboardHandler.GetResumo)
			protected.GET("/alertas", dashboardHandler.GetAlertas)

			// Licenças
			licencas := protected.Group("/licencas")
			{
				licencas.GET("", licencaHandler.List)
				licencas.POST("", licencaHandler.Create)
				licencas.GET("/:id", licencaHandler.GetByID)
				licencas.PUT("/:id", licencaHandler.Update)
				licencas.DELETE("/:id", licencaHandler.Delete)
				licencas.GET("/:id/arquivo", licencaHandler.DownloadArquivo)
			}

			// Usuários (admin only)
			usuarios := protected.Group("/usuarios")
			usuarios.Use(AdminMiddleware())
			{
				usuarios.GET("", authHandler.ListUsuarios)
				usuarios.POST("", authHandler.CreateUsuario)
				usuarios.PUT("/:id", authHandler.UpdateUsuario)
				usuarios.DELETE("/:id", authHandler.DeleteUsuario)
			}
		}
	}

	return r
}

// AuthMiddleware valida JWT
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token não fornecido"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido ou expirado"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// AdminMiddleware restringe rotas a admins
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "Acesso restrito a administradores"})
			return
		}
		c.Next()
	}
}
