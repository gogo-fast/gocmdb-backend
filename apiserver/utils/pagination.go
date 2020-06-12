package utils

import "fmt"

type Pagination struct {
	FirstPageNum   int
	LastPageNum    int
	TotalItemCount int
	CurrentPageNum int
	CurrentPage    *Page `json:"-"`
}

type Page struct {
	Limit  int
	Offset int
}

func NewPagination(total, currentPageNum, pageSize, paginationLength int) (*Pagination, error) {

	pagination := Pagination{
		CurrentPage: &Page{},
	}

	pagination.TotalItemCount = total

	var (
		totalPageCount int
		firstPageNum   int
		lastPageNum    int
	)

	maxPaginationLength := 10

	if paginationLength <= 0 {
		return nil, fmt.Errorf("分页长度不能为负数")
	}

	if total <= 0 {
		return nil, fmt.Errorf("数据总数不能为负数")
	}

	if pageSize < 5 {
		pageSize = 5
	} else if pageSize > 500 {
		pageSize = 500
	}

	if total%pageSize != 0 {
		totalPageCount = total/pageSize + 1
	} else {
		totalPageCount = total / pageSize
	}

	if currentPageNum < 1 {
		currentPageNum = 1
	} else if currentPageNum > totalPageCount {
		currentPageNum = totalPageCount
	}

	pagination.CurrentPageNum = currentPageNum
	pagination.CurrentPage.Offset = (currentPageNum - 1) * pageSize
	pagination.CurrentPage.Limit = pageSize

	if paginationLength > totalPageCount {
		paginationLength = totalPageCount
	}

	if paginationLength > maxPaginationLength {
		paginationLength = maxPaginationLength
	}

	paginationMiddle := paginationLength // 2
	firstPageNum = currentPageNum - paginationMiddle
	lastPageNum = currentPageNum + paginationMiddle

	if firstPageNum <= 0 {
		pagination.FirstPageNum = 1
	} else {
		pagination.FirstPageNum = firstPageNum
	}
	if lastPageNum > totalPageCount {
		pagination.LastPageNum = totalPageCount
	} else {
		pagination.LastPageNum = lastPageNum
	}

	return &pagination, nil
}
