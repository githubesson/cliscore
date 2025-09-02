package models

type SearchRequest struct {
	Terms    []string `json:"terms"`
	Types    []string `json:"types"`
	Wildcard bool     `json:"wildcard"`
	Source   string   `json:"source"`
	Operator *string  `json:"operator,omitempty"`
}

type CountRequest struct {
	Terms    []string `json:"terms"`
	Types    []string `json:"types"`
	Wildcard bool     `json:"wildcard"`
	Source   string   `json:"source"`
	Operator *string  `json:"operator,omitempty"`
}

type SearchResponse struct {
	Results map[string]interface{} `json:"results"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}

type DetailedCountResponse struct {
	Counts     map[string]interface{} `json:"counts"`
	TotalCount int64                  `json:"total_count"`
	Took       int64                  `json:"took"`
}

type ValidationRequest struct {
	APIKey string `json:"apiKey"`
}

type ValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	UserID  string `json:"userId,omitempty"`
	Role    string `json:"role,omitempty"`
}

type Config struct {
	BaseURL string
	APIKey  string
}

type MachineInfoRequest struct {
	UUID string `json:"uuid"`
}

type NormalizedMachineInfo struct {
	OperatingSystem string `json:"operatingSystem,omitempty"`
	OSVersion       string `json:"osVersion,omitempty"`
	Architecture    string `json:"architecture,omitempty"`
	Language        string `json:"language,omitempty"`
	TimeZone        string `json:"timeZone,omitempty"`

	HWID             string   `json:"hwid,omitempty"`
	ComputerName     string   `json:"computerName,omitempty"`
	UserName         string   `json:"userName,omitempty"`
	RAMSize          string   `json:"ramSize,omitempty"`
	CPUName          string   `json:"cpuName,omitempty"`
	CPUVendor        string   `json:"cpuVendor,omitempty"`
	CPUCores         string   `json:"cpuCores,omitempty"`
	CPUThreads       string   `json:"cpuThreads,omitempty"`
	GPUs             []string `json:"gpus,omitempty"`
	ScreenResolution string   `json:"screenResolution,omitempty"`

	IPAddress string `json:"ipAddress,omitempty"`
	Country   string `json:"country,omitempty"`
	Location  string `json:"location,omitempty"`
	ZipCode   string `json:"zipCode,omitempty"`

	AntiViruses     []string `json:"antiViruses,omitempty"`
	ProcessElevated string   `json:"processElevated,omitempty"`

	FilePath    string `json:"filePath,omitempty"`
	InstallDate string `json:"installDate,omitempty"`
	LogDate     string `json:"logDate,omitempty"`

	BuildID         string   `json:"buildId,omitempty"`
	Domain          string   `json:"domain,omitempty"`
	Hostname        string   `json:"hostname,omitempty"`
	NetBIOS         string   `json:"netBIOS,omitempty"`
	KeyboardLayouts []string `json:"keyboardLayouts,omitempty"`

	MachineID     string `json:"machineId,omitempty"`
	ProductKey    string `json:"productKey,omitempty"`
	AdminGroup    string `json:"adminGroup,omitempty"`
	Integrity     string `json:"integrity,omitempty"`
	WallpaperHash string `json:"wallpaperHash,omitempty"`

	CountryCode     string            `json:"countryCode,omitempty"`
	CountryName     string            `json:"countryName,omitempty"`
	DataInformation string            `json:"dataInformation,omitempty"`
	ParsedDataInfo  map[string]string `json:"parsedDataInfo,omitempty"`

	LocalTime     string              `json:"localTime,omitempty"`
	UTC           string              `json:"utc,omitempty"`
	IsLaptop      string              `json:"isLaptop,omitempty"`
	RunningPath   string              `json:"runningPath,omitempty"`
	ProcessCount  string              `json:"processCount,omitempty"`
	ProcessList   []string            `json:"processList,omitempty"`
	InstalledApps []string            `json:"installedApps,omitempty"`
	Monitors      []map[string]string `json:"monitors,omitempty"`

	FileType   string   `json:"fileType,omitempty"`
	SourceInfo string   `json:"sourceInfo,omitempty"`
	FileTree   []string `json:"fileTree,omitempty"`
}

type MachineInfoResponse struct {
	Data  *NormalizedMachineInfo `json:"data,omitempty"`
	Error string                 `json:"error,omitempty"`
}

type DownloadRequest struct {
	UUID string `json:"uuid"`
	File string `json:"file,omitempty"`
}

type DownloadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FilePath string `json:"filePath,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

type CreditsResponse struct {
	Credits int64  `json:"credits"`
	Message string `json:"message,omitempty"`
}

type ApiKeyValidation struct {
	ApiKey string `json:"apiKey"`
}
