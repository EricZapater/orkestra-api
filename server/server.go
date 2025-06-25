package server

import (
	"database/sql"
	"orkestra-api/config"
	"orkestra-api/internal/auth"
	"orkestra-api/internal/costitems"
	"orkestra-api/internal/customers"
	"orkestra-api/internal/groups"
	"orkestra-api/internal/health"
	"orkestra-api/internal/meetings"
	"orkestra-api/internal/projects"
	"orkestra-api/internal/searches"
	"orkestra-api/internal/tasks"
	"orkestra-api/internal/users"
	"orkestra-api/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		router: gin.Default(),
		cfg:    cfg,
		db:     db,
	}
}

func (s *Server) Setup() error {
	// CORS middleware
	s.router.Use(middleware.SetupCORS())
	

	// JWT middleware
	authMiddleware, err := middleware.SetupJWT(s.cfg)
	if err != nil {
		return err
	}

	// Action log middleware
	actionLogMiddleware := middleware.NewActionLogMiddleware(s.db)
	
	// Inicialitzar repositoris
	userRepo := users.NewUserRepository(s.db)
	groupRepo := groups.NewGroupRepository(s.db)
	meetingRepo := meetings.NewMeetingRepository(s.db)
	searchRepo := searches.NewSearchRepository(s.db)
	customerRepo := customers.NewCustomerRepository(s.db)
	projectRepo := projects.NewProjectRepository(s.db)
	taskRepo := tasks.NewTaskRepository(s.db)
	costItemRepo := costitems.NewCostItemRepository(s.db)


	// Inicialitzar serveis
	userService := users.NewUserService(userRepo)
	authService := auth.NewAuthService(userRepo, authMiddleware)
	groupService := groups.NewGroupService(groupRepo)
	meetingService := meetings.NewMeetingService(meetingRepo)
	searchService := searches.NewSearchService(searchRepo)
	customerService := customers.NewCustomerService(customerRepo)
	projectService := projects.NewProjectService(projectRepo, customerService)
	taskService := tasks.NewTaskService(taskRepo, userService, projectService)
	costItemService := costitems.NewCostItemService(costItemRepo, projectService)


	// Inicialitzar handlers
	userHandler := users.NewUserHandler(userService)
	authHandler := auth.NewAuthHandler(authService, authMiddleware)
	groupHandler := groups.NewGroupHandler(groupService)
	meetingHandler := meetings.NewMeetingHandler(meetingService)
	searchHandler := searches.NewSearchHandler(searchService)
	customerHandler := customers.NewCustomerHandler(customerService)
	projectHandler := projects.NewProjectHandler(projectService)
	taskHandler := tasks.NewTaskHandler(taskService)
	costItemHandler := costitems.NewCostItemHandler(costItemService)

	
	// Configurar les rutes públiques (sense autenticació)
	public := s.router.Group("/auth")
	public.Use(actionLogMiddleware.LogAction())
	public.GET("/health", health.CheckHealth)
	users.RegisterPublicRoutes(public, userHandler)
	auth.RegisterRoutes(public, authHandler, authMiddleware)


	// Configurar les rutes protegides (amb autenticació JWT)
	protected := s.router.Group("/api")
	protected.Use(authMiddleware.MiddlewareFunc())
	protected.Use(actionLogMiddleware.LogAction())
	

	// Registrar les rutes protegides
	users.RegisterRoutes(protected, userHandler)
	groups.RegisterRoutes(protected, groupHandler)
	meetings.RegisterRoutes(protected, meetingHandler)
	searches.RegisterRoutes(protected, searchHandler)
	customers.RegisterRoutes(protected, customerHandler)
	projects.RegisterRoutes(protected, projectHandler)
	tasks.RegisterRoutes(protected, taskHandler)
	costitems.RegisterRoutes(protected, costItemHandler)
	
	return nil
}

func (s *Server) Run() error {
	//return s.router.RunTLS(":" + s.cfg.ApiPort, "./certs/cert.pem", "./certs/key.pem")
	return s.router.Run(":" + s.cfg.ApiPort)
}