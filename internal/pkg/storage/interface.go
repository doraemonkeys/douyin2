package storage

// // ObjectStorageSerice
type VideoStorageService[T any] interface {
	// Save object
	Save(T) (uint, error)
	// Delete object
	Delete(uint) error
	// Get object
	Get(uint) (T, error)
	// SaveUnique 保存视频，如果视频已经存在则返回已存在的视频ID和Error
	SaveUnique(T) (uint, error)

	GetURL(uint) (string, string, error)
}
