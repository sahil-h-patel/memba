package auth

import (
	"os"

	csh_auth "github.com/ComputerScienceHouse/csh-auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var LoginDone = make(chan bool)
var LogoutDone = make(chan bool)

func StartServer() {
	// load in env vars
	godotenv.Load()

	// setup CSH auth wrapper
	csh := csh_auth.CSHAuth{}

	scopes := []string{"openid", "roles", "groups", "profile", "offline_access", "api"}

	// fmt.Printf("DEBUG: Client ID: [%s]\n", os.Getenv("OIDC_CLIENT_ID"))
	// fmt.Printf("DEBUG: Client secret: [%s]\n", os.Getenv("OIDC_CLIENT_SECRET"))
	// fmt.Printf("DEBUG: Secret: [%s]\n", os.Getenv("JWT_SECRET"))
	// fmt.Printf("DEBUG: State: [%s]\n", os.Getenv("STATE"))
	// fmt.Printf("DEBUG: Server Host: [%s]\n", os.Getenv("SERVER_HOST"))
	// fmt.Printf("DEBUG: Redirect URI: [%s]\n", os.Getenv("REDIRECT_URI"))
	// fmt.Printf("DEBUG: Auth URI: [%s]\n", os.Getenv("AUTH_URI"))

	csh.Init(
		os.Getenv("OIDC_CLIENT_ID"),
		os.Getenv("OIDC_CLIENT_SECRET"),
		os.Getenv("SECRET"),
		os.Getenv("STATE"),
		os.Getenv("SERVER_HOST"),
		os.Getenv("REDIRECT_URI"),
		os.Getenv("AUTH_URI"),
		scopes,
	)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/auth/", func(c *gin.Context) {
		// c.String(200, "Successfully authenticated with MEMBA! You can close this window.")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/auth/login")
	})
	r.GET("/auth/login", csh.AuthRequest)
	r.GET("/auth/redirect", func(c *gin.Context) {
		csh.AuthCallback(c)
		SaveToken(c)
		LoginDone <- true
		c.String(200, "Logged out! You can close this window.")
	})
	r.GET("/auth/logout", func(c *gin.Context) {
		tokenPath := os.Getenv("HOME") + "/.memba/token"
		os.Remove(tokenPath)
		csh.AuthLogout(c)
	})
	// r.GET("/search", csh.AuthWrapper(HandleSearch))

	r.Run(":8080")
}
