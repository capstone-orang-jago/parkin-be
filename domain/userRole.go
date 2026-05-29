package domain

type UserRole string

const (
	RoleAdmin    UserRole = "superadmin"
	RolePengguna UserRole = "pengguna"
	RolePetugas  UserRole = "petugas"
	RoleDishub   UserRole = "dishub"
)
