package admin

import (
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	_ "github.com/GoAdminGroup/themes/sword"
	"github.com/gin-gonic/gin"
)

// SetupAdmin configures and returns the go-admin engine
func SetupAdmin(r *gin.Engine, databaseURL string) (*engine.Engine, error) {
	eng := engine.Default()

	cfg := config.Config{
		Env: config.EnvLocal,
		Databases: config.DatabaseList{
			"default": {
				Host:   "db",
				Port:   "5432",
				User:   "postgres",
				Pwd:    "postgres",
				Name:   "labour_thekedar",
				Driver: db.DriverPostgresql,
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: "skin-black",
	}

	template.AddComp(chartjs.NewChart())

	adminPlugin := admin.NewAdmin(nil)

	if err := eng.AddConfig(&cfg).
		AddPlugins(adminPlugin).
		Use(r); err != nil {
		return nil, err
	}

	return eng, nil
}
