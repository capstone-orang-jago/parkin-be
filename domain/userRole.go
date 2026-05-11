package domain

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RolePengguna UserRole = "pengguna"
	RolePetugas  UserRole = "petugas"
	RoleDishub   UserRole = "dishub"
)
