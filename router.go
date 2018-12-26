package router

import (
	"github.com/astaxie/beego"
	"xxx/controllers"
)


beego.Router("/api/pushDataToXxx/?:collection/?:category/?:device/?:date", &controllers.DlFrameworkController{}, "get:PushDataToXxxHandler")


