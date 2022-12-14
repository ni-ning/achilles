package errcode

var (
	// 具体业务
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
	ErrorCreateTagFail  = NewError(20010002, "创建标签失败")
	ErrorUpdateTagFail  = NewError(20010003, "更新标签失败")
	ErrorDeleteTagFail  = NewError(20010004, "删除标签失败")
	ErrorCountTagFail   = NewError(20010005, "统计标签失败")

	// 注册登录认证 再调整？
	ErrorRegisterFail   = NewError(20020001, "用户注册失败")
	ErrorLoginFail      = NewError(20020002, "用户登录失败")
	ErrorGenTokenFail   = NewError(20020003, "生成Token失败")
	ErrorParseTokenFail = NewError(20020004, "验证Token失败")

	ErrorUploadFileFail = NewError(20030001, "上传文件失败")
)
