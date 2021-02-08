package store

// TypeStore : 存储类型(表示文件存到哪里)
type TypeStore int

const (
	_ TypeStore = iota
	// LocalStore : 节点本地
	LocalStore
	// CephStore : Ceph集群
	CephStore
	// OSSStore : 阿里OSS
	OSSStore
	// MixStore : 混合(Ceph及OSS)
	MixStore
	// AllStore : 所有类型的存储都存一份数据
	AllStore
)
