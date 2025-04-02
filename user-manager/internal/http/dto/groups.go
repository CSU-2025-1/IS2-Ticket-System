package dto

type Group struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateGroupResponse struct {
	UUID string `json:"uuid"`
}

type AddUsersToGroupRequest struct {
	Users []string `json:"users" binding:"required"`
}

type RemoveUsersFromGroupRequest struct {
	Users []string `json:"users" binding:"required"`
}
