package pages

import (
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

var DashboardPage *gorouter.ComponentRouter

func init() {
	DashboardPage = gorouter.CreateComponentRouter()

	DashboardPage.UsePrefixWrap()
	DashboardPage.Wrap(gorouter.SimpleWrapper(views.Dashboard))
	DashboardPage.Retarget("#content")

	DashboardPage.GetComponent("/metrics", gorouter.SimpleComponent(views.Metrics)).Retarget("#dashboard-component")
	DashboardPage.GetComponent("/settings", gorouter.SimpleComponent(views.Settings)).Retarget("#dashboard-component")
	DashboardPage.GetComponent("/account", gorouter.SimpleComponent(views.Account)).Retarget("#dashboard-component")

}
