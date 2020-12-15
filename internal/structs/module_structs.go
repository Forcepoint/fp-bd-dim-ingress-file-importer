package structs

type ModuleMetadata struct {
	ModuleServiceName string           `json:"module_service_name"`
	ModuleDisplayName string           `json:"module_display_name"`
	ModuleType        string           `json:"module_type"`
	ModuleDescription string           `json:"module_description"`
	InboundRoute      string           `json:"inbound_route"`
	InternalIP        string           `json:"internal_ip"`
	InternalPort      string           `json:"internal_port"`
	Configured        bool             `json:"configured"`
	Configurable      bool             `json:"configurable"`
	IconURL           string           `json:"icon_url"`
	InternalEndpoints []ModuleEndpoint `json:"internal_endpoints"`
}

type ModuleEndpoint struct {
	Secure      bool         `json:"secure"`
	Endpoint    string       `json:"endpoint"`
	HttpMethods []HttpMethod `json:"http_methods"`
}

type HttpMethod struct {
	Method string `json:"method"`
}
