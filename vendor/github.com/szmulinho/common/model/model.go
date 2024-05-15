package model

type Doctor struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Drug struct {
	DrugID int64  `json:"drug_id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
}

type Drugs []string

type Prescription struct {
	PreID      int64  `json:"pre_id" gorm:"primaryKey;autoIncrement"`
	Drugs      Drugs  `gorm:"type:text[]" json:"drugs"`
	Patient    string `json:"patient"`
	Expiration string `json:"expiration"`
}

type Opinion struct {
	ID      int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Login   string `json:"login"`
	Rating  int    `json:"rating" gorm:"column:rating"`
	Comment string `json:"comment"`
}

type Order struct {
	ID      int64  `json:"order_id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Items   string `json:"items"`
	Price   string `json:"price"`
}

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Login    string `gorm:"unique" json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type GithubUser struct {
	ID          int64  `gorm:"unique_index"`
	Login       string `json:"username"`
	AvatarUrl   string `json:"avatar_url"`
	HtmlUrl     string `json:"html_url"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AccessToken string `json:"-"`
}
