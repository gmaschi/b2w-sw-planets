package planetmodel

type (
	CreateRequest struct {
		Name    string `json:"name" binding:"required"`
		Terrain string `json:"terrain" binding:"required"`
		Climate string `json:"climate" binding:"required"`
	}

	GetRequest struct {
		ID string `uri:"id" binding:"required"`
	}

	DeleteRequest struct {
		ID string `uri:"id" binding:"required"`
	}

	ListRequest struct {
		// TODO: implement pagination
		Name string `form:"name"`
		//PageID   int32  `form:"page_id"`
		//PageSize int32  `form:"page_size"`
		//PageID   int32  `form:"page_id" binding:"required,min=1"`
		//PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
	}
)
