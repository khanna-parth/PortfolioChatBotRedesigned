package transfer

type DocumentRequest struct {
	UserID string `json:"userID"`
	Document string `json:"docPath"`
}

type DocumentResponse struct {
	Message string `json:"message"`
}

type DocumentListResponse struct {
	Docs[]string `json:"documents"`
}

type DocumentDeleteRequest struct {
	Document string `json:"document"`
}

type DocumentList struct {
	Docs []string `json:"documents"`
}

func CreateDocumentList() *DocumentList {
	return &DocumentList{
		Docs: []string{},
	}
}