### CustomFunction
* GetRandomNumber(`最小值`,`最大值`) : 取隨機整數
* GetRandomString(`指定長度`) : 生成隨機字串
* SlceStringInsert(`目標字串切片`,`指定index`,`欲存入字串`) : 將字串存入切片的指定index
* SlceIntInsert(`目標int切片`,`指定index`,`欲存入int`) : 將int元素存入切片的指定index
* ConvertJsonToMap(`json檔案路徑`) : 將json檔轉換為map
### Logger
* GetLogger(`logger名稱`,`是否使用色碼`) : 取得logger
* 設定檔參數說明 : 
    * 設定檔存放位置 : `根目錄/config/logger.json`
    * key = logger name(未設定時，使用`預設設定`)
    * params : 
        * TimestampFormat : 
        時間顯示格式,`2006-01-02 15:04:05.000`
        * TimeUnit : 時間設定單位
            * `second`  : 秒
            * `minute`  : 分
            * `hour`    : 時
            * `day`     : 天 
        * RotationType : 檔案分割方式
            * `time` : 時間間隔分割,使用參數
                * `WithMaxAge`
                * `WithRotationTime`
            * `size` : 檔案大小分割,使用參數
                * `WithRotationCount`
                * `WithRotationSize_MB`
        * WithMaxAge : 存活時限(單位參照`TimeUnit`設定)
        * WithRotationTime : 分割時間間隔(單位參照`TimeUnit`設定)
        * WithRotationCount : 最高留存檔案數量
        * WithRotationSize_MB : 每個檔案最大size(單位:MB)
        * LogDir : log檔存放位置

* 預設設定
    * TimestampFormat : `2006-01-02 15:04:05.000`
    * TimeUnit : `hour`
    * RotationType : `time`
    * WithMaxAge : `120`
    * WithRotationTime : `24`
    * WithRotationCount : `5`
    * WithRotationSize_MB : `100`
    * LogDir : `logs/`

* 取得logger
```js
import (
	"gitlab.inlive7.com/golang/utils"
)

// @param loggerName string logger名稱
// @param disableColor bool 是否使用色碼
logger := utils.GetLogger("logger name", true)
```

* 設定檔範例請參考`example/logger.json`

### HashMaker
* CreateSHAHash(`目標字串`, `SHA類別`) : 生成SHA Hash

### GoogleCloudStorageHandler
* func CreateClient() (Storage_client, error): 建立Storage連線, 回傳 Storage_client{Ctx:context.Context, Client *storage.Client}
}
```
s, err := utils.CreateClient()
	if err != nil {
		log.Fatal(err)
	}
```

* func (s *Storage_client) Close(): 結束連線
* func (s *Storage_client) DownloadFile(bucket, location, object string) error: 下載檔案到當前目錄
```
err = s.DownloadFile("srs-dev-record-file", "./12_12331001_259033858.mp4", "12_12331001_259033858.mp4")
if err != nil {
    log.Fatal(err)
}
```
* func (s *Storage_client) UploadFile(bucket, location, object string) error: 上傳檔案
```
err = s.UploadFile("srs-dev-record-file", "./12_12331001_259033858-1.mp4", "12_12331001_259033858-1.mp4")
if err != nil {
    log.Fatal(err)
}
```
* func (s *Storage_client) DeleteFile(bucket, object string) error: 刪除檔案
```
err = s.DeleteFile("srs-dev-record-file", "12_12331001_259033858-1.mp4")
if err != nil {
    log.Fatal(err)
}
```
* func (s *Storage_client) ListFilesWithPrefix(bucket, prefix string) ([]ObjectAttrs, error): 查詢, 列出bucket物件清單, 包含檔名, 大小(Bytes), 新增時間 
```
fileList, err := s.ListFilesWithPrefix("srs-dev-record-file", "20211202")
if err != nil {
	log.Fatal(err)
}
fmt.Println(fileList)
```
* func (s *Storage_client) CopyFile(dstBucket, srcBucket, srcObject string) error: 複製檔案, 檔名-copy
```
err = s.CopyFile("srs-dev-record-file", "srs-dev-record-file-2", "12_12331001_259033858.mp4")
if err != nil {
    log.Fatal(err)
}
```

* func (s *Storage_client) MoveFile(bucket, object, target string) error: 移動檔案(重新命名), 
```
err = s.MoveFile("srs-dev-record-file", "12_12331001_259033858.mp4", "12_12331001_259033858-1.mp4")
if err != nil {
    log.Fatal(err)
}
```
