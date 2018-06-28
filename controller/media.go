package controller



//func MediaCreate(c *gin.Context) {
//	var mediaDto struct {
//		Title        string `form:"title"`
//		Introduction string `form:"introduction" `
//		Categories   []int  `form:"category"`
//	}
//	uploadPath := "parent/"
//	if err := c.Bind(&mediaDto); err != nil {
//		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
//		return
//	}
//	mc, err := c.MultipartForm()
//	if err != nil {
//		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
//		return
//	}
//	m := &model.Media{
//		Title:        mediaDto.Title,
//		Introduction: mediaDto.Introduction,
//	}
//	ms := make([]*model.MediaSharpness, 0)
//	cs := make([]*model.Category, 0)
//	for i, f := range mc.File["file"] {
//		if err := c.SaveUploadedFile(f, uploadPath+f.Filename); err != nil {
//			c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
//		} else {
//			ms = append(ms, &model.MediaSharpness{
//				SharpnessId: i + 1,
//				Uri:         f.Filename,
//			})
//		}
//	}
//	for _, v := range mediaDto.Categories {
//		cs = append(cs, &model.Category{Model: model.Model{ID: v}})
//	}
//
//	if media, err := service.MediaCreate(m, ms, cs); err != nil {
//		c.JSON(http.StatusBadRequest, result(http.StatusBadRequest, nil, err.Error()))
//	} else {
//		c.JSON(http.StatusOK, result(http.StatusOK, media, ""))
//	}
//}
