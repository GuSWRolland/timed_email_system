package models

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"strconv"
	"strings"
	"github.com/go-gomail/gomail"
	"os"
	"bufio"
	"io"
	"time"
	"github.com/astaxie/beego/toolbox"
	"io/ioutil"
	"runtime"
)

func getCollectionNameByFrameName(frameName string) string{
	var collectionName string

	if frameName == "Caffe" {
		collectionName = "caffe"
	}
	if frameName == "MxNet" {
		collectionName = "mxnet"
	}
	if frameName == "Caffe2" {
		collectionName = "caffe2"
	}
	if frameName == "PyTorch" {
		collectionName = "pytorch"
	}
	if frameName == "BigDL" {
		collectionName = "big_dl"
	}
	if frameName == "Tensorflow" {
		collectionName = "tensorflow_sh"
	}
	if frameName == "Chainer" {
		collectionName = "chainer"
	}
	if frameName == "PaddlePaddle" {
		collectionName = "paddlepaddle"
	}

	return collectionName
}

func judgeIsRemoveFile(fileName string) {

	//taskFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0755)
	taskContent, err := ioutil.ReadFile(fileName)
	if err != nil{
		log.Println(err)
		return
	}
	if string(taskContent) == ""{
		clearTaskFile(fileName)
	}
	return
}

func replaceOldContent(content string, times string, strToken string) string {
	newContentElements := strings.Split(content, "#")
	newContent := newContentElements[0] + "#" + newContentElements[1] + "#" + newContentElements[2] + "#" + newContentElements[3]  + "#Time:" + times + strToken
	return newContent
}

func judgeSameTabTime(filename string, framework string, tab string, date string, times string) bool {

	var strToken string
	systemName := GetSystemName()
	if systemName != "windows" {
		strToken = "\n"
	} else {
		strToken = "\r\n"
	}
	//timeLayOut := "2006-01-02 15:04:05"
	timeLayOut := "15:04:05"

	timeUse, _ := time.Parse(timeLayOut, times)

	taskFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
		//return false
	}
	defer taskFile.Close()

	reader := bufio.NewReader(taskFile)

	for {
		taskContent, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		taskFrameAttr := strings.Split(taskContent, "#")[0]
		taskFrame := strings.Split(taskFrameAttr, ":")[1]
		fmt.Println("taskFrame ----->", taskFrame)

		taskTabAttr := strings.Split(taskContent, "#")[3]
		taskTab := strings.Split(taskTabAttr, ":")[1]

		fmt.Println("taskTab ------>", taskTab)

		taskDateAttr := strings.Split(taskContent, "#")[1]
		taskDate := strings.Split(taskDateAttr, ":")[1]

		fmt.Println("taskDate ----->", taskDate)

		taskTimeAttr := strings.Split(taskContent, "#")[4]
		//taskTimeTmp := strings.Split(taskTimeAttr, ":")[1]

		fmt.Println("taskTimeAttr --->", taskTimeAttr)

		tmpUploadElement := strings.Split(taskTimeAttr, ":")
		if len(tmpUploadElement) < 4{
			return false
		}

		tmpTimeHour := strings.Split(taskTimeAttr, ":")[1]
		tmpTimeMin := strings.Split(taskTimeAttr, ":")[2]
		tmpTimeSec := strings.Split(taskTimeAttr, ":")[3]
		tmpTimeSecUse := strings.Split(tmpTimeSec, strToken)[0]
		fmt.Println("tmpTimeSec ------->", tmpTimeSec)
		tmpTime := tmpTimeHour + ":" + tmpTimeMin + ":" + tmpTimeSecUse
		//taskTime := strings.Split(taskTimeTmp, strToken)[0]
		fmt.Println("tmpTime ------ >", tmpTime)

		//fmt.Println("taskTime ------>", tmpTime)
		runTime, _ := time.Parse(timeLayOut, tmpTime)
		if taskFrame == framework && taskTab == tab && taskDate == date {
			if timeUse.After(runTime) {
				newContent := replaceOldContent(taskContent, times, strToken)
				oldContent := taskContent
				RewriteFileContent(newContent, oldContent, filename)
			}
			return false
		}
	}
	return true
}

func GetSystemName() string {
	systemName := runtime.GOOS
	return systemName
}

func writeTaskFile(content string, fileName string) {

	// So, Only send content without token

	taskFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	defer taskFile.Close()
	var strToken string
	systemName := GetSystemName()
	if systemName != "windows" {
		strToken = "\n"
	} else {
		strToken = "\r\n"
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	taskFile.Seek(0, 2)
	taskFile.WriteString(content + strToken)
}

func clearkTaskFileOnTask(fileName, taskName string){
	clearTaskFile(fileName)
	toolbox.DeleteTask(taskName)
}

func RewriteFileContent(newContent string, oldContent string, fileName string) {

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fileContent := string(buf)

	useContent := strings.Replace(fileContent, oldContent, newContent, 1)
	ioutil.WriteFile(fileName, []byte(useContent), 0755)

}

func judgeIsSendEmail(fileName string, taskName string) bool {
	taskFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0755)
	var strToken string
	systemName := GetSystemName()
	if systemName != "windows" {
		strToken = "\n"
	} else {
		strToken = "\r\n"
	}

	var taskNameUse = taskName + strToken

	if err != nil {
		log.Fatal(err)
		return false
	}
	defer taskFile.Close()

	reader := bufio.NewReader(taskFile)

	for {
		taskContent, err := reader.ReadString('\n')

		if taskContent == taskNameUse {

			return false
		}
		if err == io.EOF {
			//fmt.Println("No Same Name")
			return true
		}
	}
}

func getFrameNameByCollectionName(collectionName string) string{
	var frameName string

	if collectionName == "caffe" {
		frameName = "Caffe"
	}
	if collectionName == "mxnet" {
		frameName = "MxNet"
	}
	if collectionName == "caffe2" {
		frameName = "Caffe2"
	}
	if collectionName == "pytorch" {
		frameName = "PyTorch"
	}
	if collectionName == "big_dl" {
		frameName = "BigDL"
	}
	if collectionName == "tensorflow_sh" {
		frameName = "Tensorflow"
	}
	if collectionName == "chainer" {
		frameName = "Chainer"
	}
	if collectionName == "paddlepaddle" {
		frameName = "PaddlePaddle"
	}

	return frameName
}

func SendEmailAutomatically(collectionName, dateNow, judgeFileName, fileName, strToken, taskName string){
	frameName := getFrameNameByCollectionName(collectionName)
	taskFile, err := os.OpenFile(judgeFileName, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil{
		log.Fatal(err)
	}
	defer taskFile.Close()
	reader := bufio.NewReader(taskFile)
	for {
		taskContent, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		taskFrame := strings.Split(taskContent, "_")[0]
		if taskFrame != "" && taskFrame != strToken && taskFrame == frameName {
			collectionNameUse := getCollectionNameByFrameName(taskFrame)
			sendDataUpdatedEmail(collectionNameUse, dateNow, fileName)
		}

		if taskFrame == ""{
			break
		}
	}
	toolbox.DeleteTask(taskName)
}

func SetEmailOrder(collectionName string) {

	var fileName string
	var judgeFileName string
	var fileNameToken string
	//var strToken string
	var strToken string
	systemName := GetSystemName()
	if systemName != "windows" {
		strToken = "\n"
	} else {
		strToken = "\r\n"
	}
	if systemName != "windows" {
		fileName = "./task"
		judgeFileName = "./judge"
		//strToken = "\n"
	} else {
		fileName = ".//task"
		judgeFileName = ".//judge"
		//strToken = "\r\n"
	}
	timeLayOut := "2006-01-02 15:04:05"

	dateNow := GetNowDate()

	timeNow := GetNowTime()
	frameName := getFrameNameByCollectionName(collectionName)

	taskName := frameName + "_" + dateNow

	nowUse, _ := time.Parse(timeLayOut, dateNow+" "+timeNow)
	earlyPoint, _ := time.Parse(timeLayOut, dateNow+" "+"9:00:00")
	pointOne, _ := time.Parse(timeLayOut, dateNow+" "+"12:00:00")
	pointTwo, _ := time.Parse(timeLayOut, dateNow+" "+"15:00:00")
	pointThree, _ := time.Parse(timeLayOut, dateNow+" "+"17:00:00")
	pointFour, _ := time.Parse(timeLayOut, dateNow+" "+"23:59:59")

	if nowUse.Before(pointOne) && nowUse.After(earlyPoint) {
		fileNameToken = "_12.log"

		judgeFileName += fileNameToken

		taskName += "_12"
		//taskName += strToken

		onlyFlag := judgeIsSendEmail(judgeFileName, taskName)
		if onlyFlag != true {
			return
		}
		writeTaskFile(taskName, judgeFileName)
		fileName += fileNameToken

		taskCount := judgeNumOfTask(judgeFileName, strToken, taskName)
		hourString := " 12"
		timePart := generateTimePart(taskCount, hourString)
		datePart := " * * *"
		newSpec := timePart + datePart
		taskNameUse := taskName + "_12"


		//tk := toolbox.NewTask(taskName, newSpec, func() error { sendNormalEmail(collectionName, dateNow, fileName); return nil })
		tk := toolbox.NewTask(taskNameUse, newSpec, func() error { SendEmailAutomatically(collectionName, dateNow, judgeFileName,fileName, strToken, taskNameUse); return nil })

		if taskCount == 0{
			clearName := "clearJudge_12"
			clearTime := "0 05 12"
			clearDuration := " * * *"
			clearSpec := clearTime + clearDuration
			tkClearJudge := toolbox.NewTask(clearName, clearSpec, func() error { clearkTaskFileOnTask(judgeFileName, clearName); return nil })
			toolbox.AddTask(clearName, tkClearJudge)
		}

		toolbox.AddTask(taskNameUse, tk)
		toolbox.StopTask()
		toolbox.StartTask()


	} else if nowUse.After(pointOne) && nowUse.Before(pointTwo) {
		fileNameToken = "_15.log"

		judgeFileName += fileNameToken

		taskName += "_15"

		onlyFlag := judgeIsSendEmail(judgeFileName, taskName)
		if onlyFlag != true {
			return
		}
		writeTaskFile(taskName, judgeFileName)
		fileName += fileNameToken

		taskCount := judgeNumOfTask(judgeFileName, strToken, taskName)

		hourString := " 15"
		timePart := generateTimePart(taskCount, hourString)
		datePart := " * * *"
		newSpec := timePart + datePart
		taskNameUse := taskName + "_15"

		tk := toolbox.NewTask(taskNameUse, newSpec, func() error { SendEmailAutomatically(collectionName, dateNow, judgeFileName,fileName, strToken, taskNameUse); return nil })

		if taskCount == 0{
			clearName := "clearJudge_15"
			clearTime := "0 05 15"
			clearDuration := " * * *"
			clearSpec := clearTime + clearDuration
			tkClearJudge := toolbox.NewTask(clearName, clearSpec, func() error { clearkTaskFileOnTask(judgeFileName, clearName); return nil })
			toolbox.AddTask(clearName, tkClearJudge)
		}

		toolbox.AddTask(taskNameUse, tk)
		toolbox.StopTask()
		toolbox.StartTask()


	} else if nowUse.After(pointTwo) && nowUse.Before(pointThree) {

		fileNameToken = "_17.log"

		judgeFileName += fileNameToken

		taskName += "_17"
		//taskName += strToken

		onlyFlag := judgeIsSendEmail(judgeFileName, taskName)
		if onlyFlag != true {
			return
		}
		writeTaskFile(taskName, judgeFileName)
		fileName += fileNameToken

		taskCount := judgeNumOfTask(judgeFileName, strToken, taskName)

		hourString := " 17"
		timePart := generateTimePart(taskCount, hourString)
		datePart := " * * *"
		newSpec := timePart + datePart
		taskNameUse := taskName + "_17"

		tk := toolbox.NewTask(taskNameUse, newSpec, func() error { SendEmailAutomatically(collectionName, dateNow, judgeFileName,fileName, strToken, taskNameUse); return nil })

		if taskCount == 0{
			clearName := "clearJudge_17"
			clearTime := "0 05 17"
			clearDuration := " * * *"
			clearSpec := clearTime + clearDuration
			tkClearJudge := toolbox.NewTask(clearName, clearSpec, func() error { clearkTaskFileOnTask(judgeFileName, clearName); return nil })
			toolbox.AddTask(clearName, tkClearJudge)
		}

		toolbox.AddTask(taskNameUse, tk)
		toolbox.StopTask()
		toolbox.StartTask()
	} else if nowUse.Before(earlyPoint) {
		fileNameToken = "_9.log"

		judgeFileName += fileNameToken

		taskName += "_9"

		onlyFlag := judgeIsSendEmail(judgeFileName, taskName)
		if onlyFlag != true {
			return
		}
		writeTaskFile(taskName, judgeFileName)
		fileName += fileNameToken

		taskCount := judgeNumOfTask(judgeFileName, strToken, taskName)

		hourString := " 9"
		timePart := generateTimePart(taskCount, hourString)
		datePart := " * * *"
		newSpec := timePart + datePart
		taskNameUse := taskName + "_9"

		tk := toolbox.NewTask(taskNameUse, newSpec, func() error { SendEmailAutomatically(collectionName, dateNow, judgeFileName,fileName, strToken, taskNameUse); return nil })

		if taskCount == 0{
			clearName := "clearJudge_9"
			clearTime := "0 05 9"
			clearDuration := " * * *"
			clearSpec := clearTime + clearDuration
			tkClearJudge := toolbox.NewTask(clearName, clearSpec, func() error { clearkTaskFileOnTask(judgeFileName, clearName); return nil })
			toolbox.AddTask(clearName, tkClearJudge)
		}
		toolbox.AddTask(taskNameUse, tk)
		toolbox.StopTask()
		toolbox.StartTask()

	} else if nowUse.After(pointThree) && nowUse.Before(pointFour) {
		fileNameToken = "_9.log"

		judgeFileName += fileNameToken

		taskName += "_9"
		//taskName += strToken

		onlyFlag := judgeIsSendEmail(judgeFileName, taskName)
		if onlyFlag != true {
			return
		}
		writeTaskFile(taskName, judgeFileName)
		fileName += fileNameToken

		taskCount := judgeNumOfTask(judgeFileName, strToken, taskName)
		hourString := " 9"
		timePart := generateTimePart(taskCount, hourString)
		datePart := " * * *"
		newSpec := timePart + datePart
		taskNameUse := taskName + "_9"

		tk := toolbox.NewTask(taskNameUse, newSpec, func() error { SendEmailAutomatically(collectionName, dateNow, judgeFileName,fileName, strToken, taskNameUse); return nil })

		if taskCount == 0{
			clearName := "clearJudge_9"
			clearTime := "0 05 9"
			clearDuration := " * * *"
			clearSpec := clearTime + clearDuration
			tkClearJudge := toolbox.NewTask(clearName, clearSpec, func() error { clearkTaskFileOnTask(judgeFileName, clearName); return nil })
			toolbox.AddTask(clearName, tkClearJudge)
		}

		toolbox.AddTask(taskNameUse, tk)
		toolbox.StopTask()
		toolbox.StartTask()
	}
}

func generateTimePart(count int, hourString string) string{

	var minString string

	minNum := count * 2

	if minNum == 0{
		minString = "00"
	}else if minNum == 2{
		minString = "01"
	}else{
		minString = strconv.Itoa(minNum)
	}

	timePart := "0 " + minString + hourString

	return timePart
}

func judgeNumOfTask(fileName, strToken, taskName string) int{

	var tmpCount int

	var taskNameUse = taskName + strToken

	taskFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0755)

	if err != nil {
		log.Fatal(err)
	}
	defer taskFile.Close()

	reader := bufio.NewReader(taskFile)

	for {
		taskContent, err := reader.ReadString('\n')

		if taskContent == strToken{
			continue
		}

		if taskContent == taskNameUse {
			return tmpCount
		}

		tmpCount += 1

		if err == io.EOF {
			return 0
		}
	}
	return tmpCount
}


func GetNowDate() string {
	var tmpMon string
	var tmpDay string

	yearNow, monNow, dayNow := time.Now().Date()

	if monNow < 10 {
		tmpMon = "0"
	}
	if dayNow < 10 {
		tmpDay = "0"
	}
	dateNow := fmt.Sprintf("%d-"+tmpMon+"%d-"+tmpDay+"%d", yearNow, monNow, dayNow)
	return dateNow
}

func GetNowTime() string {
	var tmpHour string
	var tmpMin string
	var tmpSec string

	hourNow, minNow, secNow := time.Now().Clock()
	if hourNow < 10 {
		tmpHour = "0"
	}
	if minNow < 10 {
		tmpMin = "0"
	}
	if secNow < 10 {
		tmpSec = "0"
	}
	timeNow := fmt.Sprintf(tmpHour+"%d:"+tmpMin+"%d:"+tmpSec+"%d", hourNow, minNow, secNow)
	return timeNow
}


func readTaskContent(fileName, frameName string) []interface{} {

	var strToken string
	systemName := GetSystemName()
	if systemName != "windows" {
		strToken = "\n"
	} else {
		strToken = "\r\n"
	}
	taskFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0755)
	defer taskFile.Close()
	contentSlice := make([]interface{}, 0)
	nowDate := GetNowDate()

	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(taskFile)

	for {
		taskContent, err := reader.ReadString('\n')
		tmpMap := make(map[string]string)

		if taskContent == "" || taskContent == strToken{
			break
		}

		if err == io.EOF {
			break
		}
		tmpFrameworkAttr := strings.Split(taskContent, "#")[0]
		tmpFramework := strings.Split(tmpFrameworkAttr, ":")[1]
		tmpMap["Framework"] = tmpFramework

		tmpDateAttr := strings.Split(taskContent, "#")[1]
		tmpDate := strings.Split(tmpDateAttr, ":")[1]
		tmpMap["Date"] = tmpDate

		tmpCategoryAttr := strings.Split(taskContent, "#")[2]
		tmpCategory := strings.Split(tmpCategoryAttr, ":")[1]
		tmpTabName := generateTabName(tmpCategory)
		tmpMap["TabName"] = tmpTabName

		tmpUploadAttr := strings.Split(taskContent, "#")[4]
		tmpUploadElement := strings.Split(tmpUploadAttr, ":")
		if len(tmpUploadElement) >= 4{
			tmpUploadHour := strings.Split(tmpUploadAttr, ":")[1]
			tmpUploadMin := strings.Split(tmpUploadAttr, ":")[2]
			tmpUploadSec := strings.Split(tmpUploadAttr, ":")[3]
			tmpUploadTime := tmpUploadHour + ":" + tmpUploadMin + ":" + tmpUploadSec
			tmpUpload := nowDate + " " + tmpUploadTime
			if strings.Contains(tmpUpload, strToken){
				tmpUpload = strings.Split(tmpUpload, strToken)[0]
			}
			tmpMap["uploadTime"] = tmpUpload
		}
		contentSlice = append(contentSlice, tmpMap)
		if tmpFramework == frameName{
			newContent := ""
			RewriteFileContent(newContent, taskContent, fileName)
		}
	}
	return contentSlice
}

func generateDataInformation(fileName string, DashboardLink string, frameName string) string {

	var frameWork string
	var dataDate string
	var dataTab string
	var uploadTime string
	var dataInformation string

	contentSlice := readTaskContent(fileName, frameName)
	//tmpCount := 0

	for _, tmpContent := range contentSlice {
		frameWork = tmpContent.(map[string]string)["Framework"]
		if frameWork != frameName {
			continue
		}
		dataDate = tmpContent.(map[string]string)["Date"]
		dataTab = tmpContent.(map[string]string)["TabName"]
		uploadTime = tmpContent.(map[string]string)["uploadTime"]
		//tmpCount += 1
		tmpFormat := `<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color:royalblue;">[ ` + uploadTime + ` ]</span>` +
			`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span>Data Date : ` + dataDate + `</span>` +
			`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span>Tab : ` + dataTab + `</span>` +
			`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span>Link : ` + DashboardLink + `</span>` +
			`<div><br>`
		dataInformation += tmpFormat
	}
	return dataInformation
}


func sendDataUpdatedEmail(collectionName string, emailDate string, fileName string) {

	m := gomail.NewMessage()
	address := []string{}
	dashboardLink := "http://abc.123.com" + collectionName

	var frameName string
	var ccAddress string
	var toName string
	//var address_c = ""
	if collectionName == "caffe" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "Caffe"
	}
	if collectionName == "mxnet" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "MxNet"

	}
	if collectionName == "caffe2" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "Caffe2"

	}
	if collectionName == "pytorch" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "PyTorch"

	}
	if collectionName == "big_dl" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "BigDL"

	}
	if collectionName == "tensorflow_sh" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "Tensorflow"

	}
	if collectionName == "chainer" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "Chainer"

	}
	if collectionName == "paddlepaddle" {
		address = append(address, "a")
		toName = "a"
		ccAddress = "b"
		frameName = "PaddlePaddle"
	}
	//fmt.Println(ccAddress)
	dataInformation := generateDataInformation(fileName, dashboardLink, frameName)
	dataInformation = judgeInformationContent(dataInformation, fileName, frameName, dashboardLink)
	if dataInformation == ""{
		log.Println("I ran, but no data Information")
		return
	}
	for _, v := range address {
		log.Println(toName)
		//name := strings.Split(v,".")[0]
		body :=
			`<br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span>` +
				frameName +
				` Data ` +
				`has been updated to Xxx successfully,</span><div>` +
				`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span>Data Information:</span>` +
				dataInformation +
				`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;` +
				`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<div>aaa</div><br><div>Xxx</div>`
		m.SetAddressHeader("From", "abc@123.com", "abc")
		m.SetHeader("To",
			m.FormatAddress(v, "address"))
		m.SetHeader("Subject", frameName+" Xxx Dashboard Updated")
		m.SetHeader("Cc", "c", "d", "e", ccAddress)

		m.SetBody("text/html", body)
		d := gomail.NewDialer("smtp.123.com", 25, "DarkFrameMaster@xxx.com", "123,")
		if err := d.DialAndSend(m); err != nil {
			log.Println("send err", err)
			return
		}
		log.Println("done.send success")
	}
	judgeIsRemoveFile(fileName)
	return
}



func generateFileNameToken(timeHour int) (string,string) {
	var fileNameToken string
	var clearTime string

	fmt.Println("timeHour is ------> ", timeHour)

	if timeHour >= 9 && timeHour < 12 {
		fileNameToken = "_12.log"
		clearTime = "0 30 12 * * *"
	} else if timeHour >= 12 && timeHour < 15 {
		fileNameToken = "_15.log"
		clearTime = "0 30 15 * * *"
	} else if timeHour >= 15 && timeHour < 17 {
		fileNameToken = "_17.log"
		clearTime = "0 30 17 * * *"
	} else {
		fileNameToken = "_9.log"
		clearTime = "0 30 9 * * *"

	}
	return fileNameToken, clearTime
}

func generateInformationFile(collectionName string, category string, date string, times string) {
	var frameName string
	var tabName string
	var fileName string

	//hostName, _ := os.Hostname()
	systemName := GetSystemName()
	if systemName != "windows" {
		fileName = "./task"
	} else {
		fileName = ".//task"
	}
	frameName = getFrameNameByCollectionName(collectionName)
	tabName = generateTabName(category)
	fileContent := "Framework:" + frameName + "#Date:" + date + "#Category:" + category + "#Tab:" + tabName + "#Time:" + times
	timeHour, _ := strconv.Atoi(strings.Split(times, ":")[0])

	fileNameToken, _ := generateFileNameToken(timeHour)
	fileName += fileNameToken

	existNewFlag := judgeSameTabTime(fileName, frameName, tabName, date, times)
	if existNewFlag {
		writeTaskFile(fileContent, fileName)
	}
}

func PushDataToXxx(collectionName string, category string, device string, dataDate string) {
	session := getSession()
	defer session.Close()
	var date string
	var dateBaseLine string
	var testData interface{}
	var resultBaseLine interface{}

	c := session.DB(performanceDataBase).C(collectionName)
	if category != "v2"{
		testData = GetNewTestData(collectionName, category, device, dataDate)
		resultBaseLine = baseLineDataOnCts(collectionName, device, category)
		if testData != nil{
			date = testData.(CommonData).Date
		}else{
			return
		}
		if resultBaseLine != nil {
			dateBaseLine = resultBaseLine.(CommonData).Date
		}
	}else{
		testData = GetNewAccuracyTestData(collectionName, category, device, dataDate)
		resultBaseLine = baseLineAccuracyDataOnCts(collectionName, device, category)
		if testData != nil{
			date = testData.(AccuracyCommonData).Date
		}else{
			return
		}
		if resultBaseLine != nil {
			dateBaseLine = resultBaseLine.(AccuracyCommonData).Date
		}
	}
	query := bson.M{"date": date, "category": category, "device": device}
	queryBaseLine := bson.M{"date": bson.M{"$regex": dateBaseLine}, "category": category + "-baseline", "device": device}
	bo, _ := c.Find(queryBaseLine).Count()
	co, _ := c.Find(query).Count()

	if testData != nil {
		if category != "v2"{
			if testData.(CommonData).Results != nil {
				if co > 0 {
					fmt.Println("update testData" + " " + date)
					c.Update(query, testData)
				} else {
					fmt.Println("insert testData" + " " + date)
					c.Insert(testData)
					timeNow := GetNowTime()
					generateInformationFile(collectionName, category, date, timeNow)
					SetEmailOrder(collectionName)
				}
			}
			if resultBaseLine != nil {
				if resultBaseLine.(CommonData).Results != nil{
					if bo > 0 {
						fmt.Println("update data baseline" + " " + dateBaseLine)
						c.Update(queryBaseLine, resultBaseLine)
					} else {
						fmt.Println("insert data baseline" + " " + dateBaseLine)
						c.Insert(resultBaseLine)
					}
				}
			}
		}else{
			if testData.(AccuracyCommonData).Results != nil {
				if co > 0 {
					fmt.Println("update testData" + " " + date)
					c.Update(query, testData)
				} else {
					fmt.Println("insert testData" + " " + date)
					c.Insert(testData)
					timeNow := GetNowTime()
					generateInformationFile(collectionName, category, date, timeNow)
					SetEmailOrder(collectionName)
				}
			}
			if resultBaseLine != nil {
				if resultBaseLine.(AccuracyCommonData).Results != nil{
					if bo > 0 {
						fmt.Println("update data baseline" + " " + dateBaseLine)
						c.Update(queryBaseLine, resultBaseLine)
					} else {
						fmt.Println("insert data baseline" + " " + dateBaseLine)
						c.Insert(resultBaseLine)
					}
				}
			}
		}
	}
}