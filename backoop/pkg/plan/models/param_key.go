package models

type ParamKey string

const (
	DriverNameParamKey       = ParamKey("driver-name")
	BackupStorageDeviceName  = ParamKey("backup-storage-device-name")
	BackupStorageModel       = ParamKey("backup-storage-access-model")
	BackupStorageModelParams = ParamKey("backup-storage-model-params")
	BackupPath               = ParamKey("backup-path")
)
