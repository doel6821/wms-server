package helpers

import "wms-server/helpers/models"

// Pagination ..
func Pagination(page, limit int, totalData int64) *models.Page {

	var prevPage, nextPage int

	var pagination models.Page

	pagination.CurrentPage = page

	if page > 1 && totalData > 0 {
		prevPage = page - 1
		pagination.PreviousPage = prevPage
	}

	if int64(page*limit) < totalData {
		nextPage = page + 1
		pagination.NextPage = nextPage
	}

	return &pagination
}
