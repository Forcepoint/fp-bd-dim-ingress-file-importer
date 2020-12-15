package structs

type IncomingItem struct {
	EventType string `json:"event_type"`
	RemoteIp  string `json:"remote_ip"`
	Domain    string `json:"domain"`
	Timestamp string `json:"timestamp"`
}

type ProcessedItemsWrapper struct {
	Items []ProcessedItem `json:"items"`
}

type ProcessedItem struct {
	Source        string `json:"source"`
	ServiceName   string `json:"service_name"`
	Type          string `json:"type"`
	Value         string `json:"value"`
	Safe          bool   `json:"safe"`
	UpdateBatchId uint   `json:"batch_number"`
}

type ConfigObject struct {
	Values struct {
		UpdaterSchedule string `json:"updater_schedule"`
		Lists           string `json:"lists"`
	} `json:"values"`
}
