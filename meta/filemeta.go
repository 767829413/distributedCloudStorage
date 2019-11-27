package meta

//File Meta information struct
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]*FileMeta

func init() {
	fileMetas = make(map[string]*FileMeta)
}

//add *FileMeta
func UpdateFileMeta(fileMeta *FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}
//get *FileMeta
func GetFileMeta(fileSha1 string) (fileMeta *FileMeta) {
	fileMeta = fileMetas[fileSha1]
	return
}
