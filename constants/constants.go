package constants

type permission struct {
	ID   int
	Code string
	Name string
}

var PermissionViewDashboardAnalytics = permission{
	ID:   1,
	Code: "view_dashboard_analytics",
	Name: "View Dashboard Analytics",
}

var PermissionViewUser = permission{
	ID:   2,
	Code: "view_user",
	Name: "View User",
}

var PermissionEditUser = permission{
	ID:   3,
	Code: "edit_user",
	Name: "Edit User",
}

var PermissionViewRole = permission{
	ID:   4,
	Code: "view_role",
	Name: "View Role",
}

var PermissionEditRole = permission{
	ID:   5,
	Code: "edit_role",
	Name: "Edit Role",
}

var PermissionViewBlog = permission{
	ID:   6,
	Code: "view_blog",
	Name: "View Blog",
}

var PermissionEditBlog = permission{
	ID:   7,
	Code: "edit_blog",
	Name: "Edit Blog",
}

var PermissionViewTag = permission{
	ID:   8,
	Code: "view_tag",
	Name: "View Tag",
}

var PermissionEditTag = permission{
	ID:   9,
	Code: "edit_tag",
	Name: "Edit Tag",
}
