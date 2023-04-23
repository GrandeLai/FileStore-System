// Code generated by goctl. DO NOT EDIT.
package types

type FileListRequest struct {
	ParentId string `json:"parent_id"`     //查询的文件夹id
	Page     int64  `json:"page,optional"` //查询的第几页
	Size     int64  `json:"size,optional"` //每页页数
}

type FileListResponse struct {
	List  []*File `json:"list"`
	Count int64   `json:"count"`
}

type File struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type FolderListRequest struct {
	ParentId string `json:"parent_id"`
	Page     int64  `json:"page,optional"` //查询的第几页
	Size     int64  `json:"size,optional"` //每页页数
}

type FolderListResponse struct {
	List  []*Folder `json:"list"`
	Count int64     `json:"count"`
}

type Folder struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
}

type FileNameUpdateRequest struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
}

type FileNameUpdateResponse struct {
}

type FolderCreateRequest struct {
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
}

type FolderCreateResponse struct {
	Id string `json:"id"`
}

type StoreDeleteRequest struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
}

type StoreDeleteResponse struct {
}

type FileMoveRequest struct {
	Id             string `json:"id"`
	ParentId       string `json:"parent_id"`
	TargetParentId string `json:"target_parent_id"`
}

type FileMoveResponse struct {
}

type FileUploadRequest struct {
}

type FileUploadResponse struct {
	Id string `json:"id"`
}

type FileUploadByChunkRequest struct {
}

type FileUploadByChunkResponse struct {
	Id string `json:"id"`
}

type FileDownloadRequest struct {
	Id string `json:"id"`
}