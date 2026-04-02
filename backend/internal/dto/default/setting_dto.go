package dto

type SettingResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Category    string `json:"category"`
	FieldType   string `json:"field_type"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type UpdateSettingsRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}
