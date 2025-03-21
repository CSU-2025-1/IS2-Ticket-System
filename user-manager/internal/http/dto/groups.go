package dto

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type CreateGroupResponse struct {
	UUID string `json:"uuid"`
}

type AddUsersToGroupRequest struct {
	Users []string `json:"users"`
}

type RemoveUsersFromGroupRequest struct {
	Users []string `json:"users"`
}
