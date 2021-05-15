package json

import (
	"context"

	"github.com/vitalik-malkin/go-labs/backoop/pkg/plan/models"
)

type Plan = interface {
	IsReadOnly() bool

	GetMetaProps() (models.MetaProps, error)
	SetMetaProps(props models.MetaProps) error

	GetParam(ctx context.Context, key models.ParamKey) (models.ParamVal, error)
	SetParam(ctx context.Context, key models.ParamKey, val models.ParamVal) error
}

type plan struct {
	Meta   models.MetaProps   `json:"meta"`
	Driver models.DriverProps `json:"driver"`

	BackupStorageDeviceName        models.ParamVal `json:"backupStorageDeviceName,omitempty"`
	BackupStorageAccessModel       models.ParamVal `json:"backupStorageAccessModel,omitempty"`
	BackupStorageAccessModelParams models.ParamVal `json:"backupStorageAccessModelParams,omitempty"`
	BackupPath                     models.ParamVal `json:"backupPath,omitempty"`
	Source                         models.ParamVal `json:"source,omitempty"`
}

func (p *plan) IsReadOnly() bool {
	return true
}

func (p *plan) GetMetaProps() (models.MetaProps, error) {
	return p.Meta, nil
}

func (p *plan) SetMetaProps(props models.MetaProps) error {
	p.Meta = props
	return nil
}

func (p *plan) GetParam(ctx context.Context, key models.ParamKey) (models.ParamVal, error) {

	panic("not implemented yet")
}

func (p *plan) SetParam(ctx context.Context, key models.ParamKey, val models.ParamVal) error {
	panic("not implemented yet")
}
