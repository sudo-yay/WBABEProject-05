package page

type Sort struct {
	Name      string `json:"sort_name" validate:"eq=rating|eq=order_count|eq=base_time.updated_at" form:"sort-name,default=base_time.updated_at"`
	Direction int    `json:"direction" form:"direction,default=-1"`
}