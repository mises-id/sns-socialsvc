package enum

import "github.com/mises-id/sns-socialsvc/lib/codes"

type FileType uint32

const (
	ImageFile FileType = iota
	VideoFile
)

var (
	fileTypeMap = map[FileType]string{
		ImageFile: "image",
		VideoFile: "video",
	}
	fileTypeStringMap = map[string]FileType{}
)

func init() {
	for key, val := range fileTypeMap {
		fileTypeStringMap[val] = key
	}
}

func (tp FileType) String() string {
	return fileTypeMap[tp]
}

func FileTypeFromString(tp string) (FileType, error) {
	fileType, ok := fileTypeStringMap[tp]
	if !ok {
		return ImageFile, codes.ErrInvalidArgument.Newf("invalid file type: %s", tp)
	}
	return fileType, nil
}
