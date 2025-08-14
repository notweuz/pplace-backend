package response

type SheetInfoDto struct {
	Version string         `json:"version"`
	Size    map[string]int `json:"size"`
}
