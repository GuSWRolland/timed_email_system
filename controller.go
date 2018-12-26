package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"Xxx/models"
	"reflect"
	//"gopkg.in/mgo.v2/bson"
	"strconv"
	"log"
)

type DlFrameworkController struct {
	beego.Controller
	// RestController
}


func (d *DlFrameworkController) AllowCross() {
	d.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
	d.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
	d.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
	d.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	d.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	d.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
}

func (d *DlFrameworkController) PushDataToXxxHandler(){
	d.AllowCross()
	var ob interface{}
	json.Unmarshal(d.Ctx.Input.RequestBody, &ob)
	//myCategory := d.Ctx.Input.Param(":category")
	collectionName := d.Ctx.Input.Param(":collection")
	category := d.Ctx.Input.Param(":category")
	device := d.Ctx.Input.Param(":device")
	dataDate := d.Ctx.Input.Param(":date")

	fmt.Println("I am collectionName ---->", collectionName)
	fmt.Println("I am category ---->", category)
	fmt.Println("I am device ---->", device)

	models.PushDataToXxx(collectionName, category, device, dataDate)

	//if ob != nil {
	d.Data["json"] = responseOk
	//	// fmt.Println(ob)
	//	models.PushDataToHeims(collectionName, category, device)
	//	// d.Data["json"] = map[string]interface{}{"status": 200, "message": "ok!", "moreinfo": "Get the Data"}
	//} else {
	//	d.Data["json"] = map[string]interface{}{"status": 200, "message": "[FAILED]the ob is nil", "moreinfo": ""}
	//}

	d.ServeJSON()
}
