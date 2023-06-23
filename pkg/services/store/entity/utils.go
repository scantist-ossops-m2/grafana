package entity

// 😬😬😬😬 HACK HACK HACK... avoid wire circular dependency for now
var WireCircularDependencyHack EntityStoreService

// The admin request is a superset of write request features
func ToAdminWriteEntityRequest(req *WriteEntityRequest) *AdminWriteEntityRequest {
	return &AdminWriteEntityRequest{
		GRN:             req.GRN,
		Meta:            req.Meta,
		Body:            req.Body,
		Status:          req.Status,
		Folder:          req.Folder,
		Comment:         req.Comment,
		PreviousVersion: req.PreviousVersion,
	}
}
