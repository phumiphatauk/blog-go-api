package api

import (
	"fmt"
	"net/http"

	"blog-go-api/constants"
	db "blog-go-api/db/sqlc"
	_ "blog-go-api/docs"
	"blog-go-api/token"
	"blog-go-api/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type jsonResponseWithPaginate struct {
	jsonResponse
	Total int64 `json:"total"`
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()
	return server, nil
}

// Middleware to check the size of the request body
func LimitRequestBodySize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		if err := c.Request.ParseForm(); err != nil {
			c.String(http.StatusRequestEntityTooLarge, "Request body too large")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (server *Server) setupRouter() {

	// Load Config Environment
	config_env, _ := util.LoadConfig(".")

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{config_env.URL_LOCALHOST}
	// config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	gin.SetMode(config_env.GIN_MODE)

	router := gin.Default()
	router.Use(LimitRequestBodySize(50 * 1024 * 1024))
	router.Use(cors.New(config))

	routerGroup := router.Group("/")

	// Auth
	routerGroup.GET("/api/version", server.getVersion)
	routerGroup.POST("/api/login", server.loginUser)
	routerGroup.POST("/api/tokens/renew_access", server.renewAccessToken)
	routerGroup.POST("/api/signup", server.createUser)
	routerGroup.POST("/api/forgot_password", server.forgotPassword)
	routerGroup.GET("/api/reset_password", server.resetPasswordDetail)
	routerGroup.POST("/api/reset_password", server.resetPassword)

	// User
	routerGroup.GET("/api/users", authMiddleware(*server, &[]string{constants.PermissionViewUser.Code}), server.listUser)
	routerGroup.GET("/api/users/:id", authMiddleware(*server, &[]string{constants.PermissionViewUser.Code}), server.getUser)
	routerGroup.PUT("/api/users", authMiddleware(*server, &[]string{constants.PermissionEditUser.Code}), server.updateUser)
	routerGroup.DELETE("/api/users/:id", authMiddleware(*server, &[]string{constants.PermissionEditUser.Code}), server.deleteUser)

	// Profile
	routerGroup.GET("/api/profile", authMiddleware(*server, nil), server.GetProfile)
	routerGroup.PUT("/api/profile", authMiddleware(*server, nil), server.UpdateProfile)
	routerGroup.PUT("/api/profile/change_password", authMiddleware(*server, nil), server.ChangePassword)

	// Permission
	routerGroup.GET("/api/permission_group", authMiddleware(*server, &[]string{constants.PermissionEditRole.Code}), server.GetAllPermissionGroup)

	// Role
	routerGroup.GET("/api/role", authMiddleware(*server, &[]string{constants.PermissionViewRole.Code}), server.GetAllRole)
	routerGroup.GET("/api/role/:id", authMiddleware(*server, &[]string{constants.PermissionViewRole.Code}), server.GetRoleById)
	routerGroup.POST("/api/role", authMiddleware(*server, &[]string{constants.PermissionEditRole.Code}), server.CreateRole)
	routerGroup.PUT("/api/role", authMiddleware(*server, &[]string{constants.PermissionEditRole.Code}), server.UpdateRole)
	routerGroup.DELETE("/api/role/:id", authMiddleware(*server, &[]string{constants.PermissionEditRole.Code}), server.DeleteRole)
	routerGroup.GET("/api/role/dropdownlist", authMiddleware(*server, &[]string{constants.PermissionViewUser.Code}), server.GetRoleForDropDownList)

	// Tag
	routerGroup.GET("/api/tag", authMiddleware(*server, &[]string{constants.PermissionViewTag.Code}), server.GetAllTag)
	routerGroup.GET("/api/tag/:id", authMiddleware(*server, &[]string{constants.PermissionViewTag.Code}), server.GetTagById)
	routerGroup.POST("/api/tag", authMiddleware(*server, &[]string{constants.PermissionEditTag.Code}), server.CreateTag)
	routerGroup.PUT("/api/tag", authMiddleware(*server, &[]string{constants.PermissionEditTag.Code}), server.UpdateTag)
	routerGroup.DELETE("/api/tag/:id", authMiddleware(*server, &[]string{constants.PermissionEditTag.Code}), server.DeleteTag)

	// Blog
	routerGroup.GET("/api/blog", server.GetAllBlog)
	routerGroup.GET("/api/blog/tag", server.GetAllBlogWithTag)
	routerGroup.GET("/api/blog/:url", server.GetBlogByUrl)
	routerGroup.GET("/api/blog/id", authMiddleware(*server, &[]string{constants.PermissionViewBlog.Code}), server.GetBlogByID)
	routerGroup.POST("/api/blog", authMiddleware(*server, &[]string{constants.PermissionEditBlog.Code}), server.CreateBlog)
	routerGroup.PUT("/api/blog", authMiddleware(*server, &[]string{constants.PermissionEditBlog.Code}), server.UpdateBlog)
	routerGroup.DELETE("/api/blog/:id", authMiddleware(*server, &[]string{constants.PermissionEditBlog.Code}), server.DeleteBlog)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": true, "message": err.Error(), "data": nil}
}

func successResponse(data interface{}) gin.H {
	return gin.H{"error": false, "message": "successfully", "data": data}
}
