package path

import (
	"DStorage/config/store"
)

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "/data/storage/"
	// TempPartRootDir : 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/data/storage_part/"
	// CephRootDir : Ceph的存储路径prefix
	CephRootDir = "/ceph"
	// OSSRootDir : OSS的存储路径prefix
	OSSRootDir = "oss/"
	// CurrentStoreType : 设置当前文件的存储类型
	// 设置当前文件的存储类型
	DefaultStoreType = store.OSSStore
)
