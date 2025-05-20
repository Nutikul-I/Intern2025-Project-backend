package handler

type PermissionHandler interface {
	// สามารถเพิ่ม method ตามที่ต้องการได้ในอนาคต
	CheckPermission(roleID int, module string) bool
}

type permissionHandler struct {
	// สามารถเพิ่ม field ตามที่ต้องการได้
}

func NewPermissionHandler() PermissionHandler {
	return &permissionHandler{}
}

func (h *permissionHandler) CheckPermission(roleID int, module string) bool {
	// ตัวอย่างการตรวจสอบสิทธิ์
	// ในอนาคตสามารถเพิ่มการตรวจสอบจากฐานข้อมูลได้
	return true
}
