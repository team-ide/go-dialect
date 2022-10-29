package worker

type DataSource interface {
	Stop()
	Read(nameList []string, dataChan chan map[string]interface{}) (err error)
	Write(sheetName string, titles []string, nameList []string, dataChan chan map[string]interface{}) (err error)
}
