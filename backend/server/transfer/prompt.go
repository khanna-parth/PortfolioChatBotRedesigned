package transfer

type UserPromptRequest struct {
	Prompt string `json:"prompt"`
	DocsSelected string `json:"docsSelected"`
}