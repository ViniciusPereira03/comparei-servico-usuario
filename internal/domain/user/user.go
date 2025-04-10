package user

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Status      int    `json:"status"`
	Photo       string `json:"photo"`
	RayDistance int    `json:"ray_distance"`
	Level       int    `json:"level"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
	DeletedAt   string `json:"deleted_at"`
}
