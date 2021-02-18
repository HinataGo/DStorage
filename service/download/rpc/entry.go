package rpc

import (
	"context"

	cfg "DStorage/service/download/config"
	dlProto "DStorage/service/download/proto"
)

// Download :download结构体
type Download struct{}

// DownloadEntry : 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqEntry,
	res *dlProto.RespEntry) error {

	res.Entry = cfg.DownloadEntry
	return nil
}
