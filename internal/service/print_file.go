package service

type PrintFile struct {
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Formatter string `json:"formatter"`
}

type PrintFileResponse struct {
	Message string `json:"message"`
}

//
//func CreatePrintFile(filename, content, extension string, size int64, formatter string) *PrintFile {
//	return &PrintFile{}
//}
