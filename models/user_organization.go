package models

type UserOrganization struct {
	tableName struct{}

	UserId         int
	User           *User
	OrganizationId int
	Organization   *Organization
}
