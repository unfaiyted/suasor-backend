package types

// Config holds database configuration
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

// Configuration represents the complete application configuration
// @Description Complete application configuration settings
type Configuration struct {
	// App contains core application settings
	App struct {
		Name          string `json:"name" mapstructure:"name" example:"suasor" binding:"required"`
		Environment   string `json:"environment" mapstructure:"environment" example:"development" binding:"required,oneof=development staging production"`
		AppURL        string `json:"appURL" mapstructure:"appURL" example:"http://localhost:3000" binding:"required,url"`
		APIBaseURL    string `json:"apiBaseURL" mapstructure:"apiBaseURL" example:"http://localhost:8080" binding:"required,url"`
		LogLevel      string `json:"logLevel" mapstructure:"logLevel" example:"info" binding:"required,oneof=debug info warn error"`
		MaxPageSize   int    `json:"maxPageSize" mapstructure:"maxPageSize" example:"100" binding:"required,min=1,max=1000"`
		AvatarPath    string `json:"avatarPath" mapstructure:"avatarPath" example:"./uploads/avatars" binding:"required"`
		MaxAvatarSize int    `json:"maxAvatarSize" mapstructure:"maxAvatarSize" example:"5242880" binding:"required,min=1"`
	} `json:"app"`

	// Database contains database connection settings
	Db struct {
		Host     string `json:"host" mapstructure:"url" example:"localhost" binding:"required"`
		Port     string `json:"port" mapstructure:"port" example:"5432" binding:"required"`
		Name     string `json:"name" mapstructure:"name" example:"suasor" binding:"required"`
		User     string `json:"user" mapstructure:"user" example:"postgres_user" binding:"required"`
		Password string `json:"password" mapstructure:"password" example:"yourpassword" binding:"required"`
		MaxConns int    `json:"maxConns" mapstructure:"maxConns" example:"20" binding:"required,min=1"`
		Timeout  int    `json:"timeout" mapstructure:"timeout" example:"30" binding:"required,min=1"`
	} `json:"db" mapstructure:"db"`

	// HTTP contains HTTP server configuration
	HTTP struct {
		Port             string `json:"port" mapstructure:"port" example:"8080" binding:"required"`
		ReadTimeout      int    `json:"readTimeout" mapstructure:"readTimeout" example:"30" binding:"required,min=1"`
		WriteTimeout     int    `json:"writeTimeout" mapstructure:"writeTimeout" example:"30" binding:"required,min=1"`
		IdleTimeout      int    `json:"idleTimeout" mapstructure:"idleTimeout" example:"60" binding:"required,min=1"`
		EnableSSL        bool   `json:"enableSSL" mapstructure:"enableSSL" example:"false"`
		SSLCert          string `json:"sslCert" mapstructure:"sslCert" example:"/path/to/cert.pem"`
		SSLKey           string `json:"sslKey" mapstructure:"sslKey" example:"/path/to/key.pem"`
		ProxyEnabled     bool   `json:"proxyEnabled" mapstructure:"proxyEnabled" example:"false"`
		ProxyURL         string `json:"proxyURL" mapstructure:"proxyURL" example:"http://proxy:8080"`
		RateLimitEnabled bool   `json:"rateLimitEnabled" mapstructure:"rateLimitEnabled" example:"true"`
		RequestsPerMin   int    `json:"requestsPerMin" mapstructure:"requestsPerMin" example:"100" binding:"min=0"`
	} `json:"http"`

	// Auth contains authentication settings
	Auth struct {
		EnableLocal     bool     `json:"enableLocal" mapstructure:"enableLocal" example:"true"`
		SessionTimeout  int      `json:"sessionTimeout" mapstructure:"sessionTimeout" example:"60" binding:"required,min=1"`
		Enable2FA       bool     `json:"enable2FA" mapstructure:"enable2FA" example:"false"`
		JWTSecret       string   `json:"jwtSecret" mapstructure:"jwtSecret" example:"your-secret-key" binding:"required"`
		TokenExpiration int      `json:"tokenExpiration" mapstructure:"tokenExpiration" example:"24" binding:"required,min=1"`
		AllowedOrigins  []string `json:"allowedOrigins" mapstructure:"allowedOrigins" example:"http://localhost:3000"`
		// New fields to add
		AccessExpiryMinutes int    `json:"accessExpiryMinutes" mapstructure:"accessExpiryMinutes" example:"15" binding:"required,min=1"`
		RefreshExpiryDays   int    `json:"refreshExpiryDays" mapstructure:"refreshExpiryDays" example:"7" binding:"required,min=1"`
		TokenIssuer         string `json:"tokenIssuer" mapstructure:"tokenIssuer" example:"suasor-api" binding:"required"`
		TokenAudience       string `json:"tokenAudience" mapstructure:"tokenAudience" example:"suasor-client" binding:"required"`
	} `json:"auth"`
}
