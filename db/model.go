package db

type BaseModel struct {
	ID int64 `db:"id,primary_key,auto_increment"`
}
