package meta

//File Meta information struct
type Meta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]*Meta

func init() {
	fileMetas = make(map[string]*Meta)
}

//add *Meta
func UpdateFileMeta(fileMeta *Meta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

//get *Meta
func GetFileMeta(fileSha1 string) (fileMeta *Meta) {
	fileMeta = fileMetas[fileSha1]
	return
}

//remove *Meta PS: Thread-safe operation, map is Non-thread safe
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
