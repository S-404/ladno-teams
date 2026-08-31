package handler

import (
	"ladno-teams/internal/config"
	"ladno-teams/internal/handler/validation"
	"ladno-teams/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Auth          IAuthHandler
	User          IUserHandler
	Team          ITeamHandler
	Teammate      ITeammateHandler
	Invite        IInviteHandler
	Workspace     IWorkspaceHandler
	WorkspaceRole IWorkspaceRoleHandler
	Admin         IAdminHandler
	cfg           config.Config
	services      *service.Service
}

func NewHandler(cfg config.Config, services *service.Service) *Handler {
	validate := validator.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation(validation.ValidLogin, validation.LoginValidation)
		_ = v.RegisterValidation(validation.ValidPassword, validation.PasswordValidation)
		_ = v.RegisterValidation(validation.ValidRole, validation.RoleValidation)
	}

	return &Handler{
		Auth:          NewAuthHandler(services, validate),
		User:          NewUserHandler(services, validate),
		Team:          NewTeamHandler(services, validate),
		Teammate:      NewTeammateHandler(services, validate),
		Invite:        NewInviteHandler(services),
		Workspace:     NewWorkspaceHandler(services, validate),
		WorkspaceRole: NewWorkspaceRoleHandler(services, validate),
		Admin:         NewAdminHandler(services, validate),
		cfg:           cfg,
		services:      services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/", h.HomePage)
	router.GET("/login", h.AuthLoginPage)
	router.POST("/login", h.AuthLoginSubmit)
	router.GET("/register", h.AuthRegisterPage)
	router.POST("/register", h.AuthRegisterSubmit)
	router.GET("/register/invite", h.AuthInviteRegisterPage)
	router.POST("/register/invite", h.AuthInviteRegisterSubmit)
	router.POST("/logout", h.AuthLogout)

	auth := router.Group("api/auth")
	{
		auth.POST("/login", h.Auth.Login)
		auth.POST("/logout", h.Auth.Logout)
		auth.POST("/refresh", h.Auth.RefreshAccessToken)
		auth.POST("/invite-register", h.Auth.InviteRegister)
	}

	router.GET("/api/invites/:guid", h.Invite.Get)

	api := router.Group("api", h.authMiddleware)
	{
		api.POST("/auth/accept-invite", h.Auth.AcceptInvite)

		users := api.Group("/users")
		{
			users.GET("/me", h.User.Me)
			users.PUT("/me", h.User.UpdateMe)
			users.PUT("/me/profile", h.User.UpdateProfile)
		}

		teams := api.Group("/teams")
		{
			teams.POST("/", h.Team.Create)
			teams.GET("/", h.Team.List)
			teams.GET("/:guid", h.Team.Get)
			teams.PUT("/:guid", h.Team.Update)
			teams.DELETE("/:guid", h.Team.Delete)
			teams.POST("/:guid/invites", h.Invite.Create)
			teams.GET("/:guid/teammates", h.Teammate.List)
			teams.PUT("/:guid/teammates/:user_guid", h.Teammate.Update)
			teams.DELETE("/:guid/teammates/:user_guid", h.Teammate.Delete)
			teams.POST("/:guid/workspaces", h.Workspace.Create)
			teams.GET("/:guid/workspaces", h.Workspace.ListByTeam)
		}

		workspaces := api.Group("/workspaces")
		{
			workspaces.GET("/:guid", h.Workspace.Get)
			workspaces.PUT("/:guid", h.Workspace.Update)
			workspaces.DELETE("/:guid", h.Workspace.Delete)
			workspaces.POST("/:guid/roles", h.WorkspaceRole.Create)
			workspaces.GET("/:guid/roles", h.WorkspaceRole.List)
			workspaces.PUT("/:guid/roles/:teammate_guid", h.WorkspaceRole.Update)
			workspaces.DELETE("/:guid/roles/:teammate_guid", h.WorkspaceRole.Delete)
		}
	}

	router.GET("/admin/login", h.AdminLoginPage)
	router.POST("/admin/login", h.AdminLoginSubmit)
	router.POST("/admin/logout", h.AdminLogout)

	admin := router.Group("admin", h.adminMiddleware)
	{
		admin.GET("/", h.Admin.Index)
		admin.GET("/users", h.Admin.Users)
		admin.GET("/invites", h.Admin.Invites)
		admin.GET("/teams", h.Admin.Teams)
		admin.GET("/teammates", h.Admin.Teammates)
		admin.GET("/workspaces", h.Admin.Workspaces)
		admin.GET("/workspace_roles", h.Admin.WorkspaceRoles)

		admin.GET("/modals/invite", h.Admin.ModalInvite)
		admin.GET("/modals/team", h.Admin.ModalTeam)
		admin.GET("/modals/workspace-role", h.Admin.ModalWorkspaceRole)

		admin.POST("/invites", h.Admin.CreateInvite)
		admin.POST("/teams", h.Admin.CreateTeam)
		admin.POST("/workspace_roles", h.Admin.CreateWorkspaceRole)

		admin.POST("/users/:guid/toggle-block", h.Admin.ToggleUserBlock)
		admin.POST("/users/:guid/toggle-admin", h.Admin.ToggleUserAdmin)
		admin.POST("/users/:guid/admin", h.Admin.UpdateUserAdmin)
		admin.POST("/users/:guid/blocked", h.Admin.UpdateUserBlocked)
		admin.DELETE("/users/:guid", h.Admin.DeleteUser)
		admin.DELETE("/teams/:guid", h.Admin.DeleteTeam)
		admin.DELETE("/teammates/:user_guid/:team_guid", h.Admin.DeleteTeammate)
		admin.POST("/teammates/:user_guid/:team_guid/leader", h.Admin.UpdateTeammateLeader)
		admin.DELETE("/invites/:guid", h.Admin.DeleteInvite)
		admin.DELETE("/workspaces/:guid", h.Admin.DeleteWorkspace)
		admin.DELETE("/workspace_roles/:workspace_guid/:teammate_guid", h.Admin.DeleteWorkspaceRole)
	}

	return router
}
