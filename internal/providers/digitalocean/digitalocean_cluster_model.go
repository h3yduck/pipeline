// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package digitalocean

import (
	"fmt"
	"time"

	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/jinzhu/gorm"
)

// DigitalOceanClusterModel is the schema for the DB.
type DigitalOceanClusterModel struct {
	ID             uint                 `gorm:"primary_key"`
	Cluster        cluster.ClusterModel `gorm:"foreignkey:ClusterID"`
	ClusterID      uint
	DigitalOceanID string

	MasterVersion string
	NodePools     []*DigitalOceanNodePoolModel `gorm:"foreignkey:ClusterID;association_foreignkey:ClusterID"`
	Tags          []*DigitalOceanTagModel      `gorm:"foreignkey:ClusterID;association_foreignkey:ClusterID"`
}

// TableName changes the default table name.
func (DigitalOceanClusterModel) TableName() string {
	return "digitalocean_clusters"
}

// BeforeCreate sets some initial values for the cluster.
func (m *DigitalOceanClusterModel) BeforeCreate() error {
	m.Cluster.Cloud = Provider
	m.Cluster.Distribution = ClusterDistributionDigitalOcean

	return nil
}

// AfterUpdate removes node pool(s) marked for deletion.
func (m *DigitalOceanClusterModel) AfterUpdate(scope *gorm.Scope) error {
	for _, nodePoolModel := range m.NodePools {
		if nodePoolModel.Delete {
			// TODO: use transaction?
			err := scope.DB().Delete(nodePoolModel).Error

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m DigitalOceanClusterModel) String() string {
	return fmt.Sprintf("%s, Master version: %s, Node pools: %s",
		m.Cluster,
		m.MasterVersion,
		m.NodePools,
	)
}

// DigitalOceanNodePoolModel is the schema for the DB.
type DigitalOceanNodePoolModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	CreatedBy uint

	ClusterID uint `gorm:"unique_index:idx_cluster_id_name"`

	Name   string `gorm:"unique_index:idx_cluster_id_name"`
	Count  int
	Size   string
	Tags   []*DigitalOceanTagModel `gorm:"foreignkey:ID;association_foreignkey:ID"`
	Delete bool                    `gorm:"-"`
}

// DigitalOceanTag is the schema for the DB.
type DigitalOceanTagModel struct {
	ID    uint `gorm:"primary_key"`
	Value string
}

// TableName changes the default table name.
func (DigitalOceanNodePoolModel) TableName() string {
	return "digitalocean_node_pools"
}

func (m DigitalOceanNodePoolModel) String() string {
	return fmt.Sprintf(
		"ID: %d, createdAt: %v, createdBy: %d, Name: %s, Count: %d",
		m.ID,
		m.CreatedAt,
		m.CreatedBy,
		m.Name,
		m.Count,
	)
}

//Save the cluster to DB
func (m DigitalOceanClusterModel) Save() error {
	db := config.DB()
	err := db.Save(&m).Error
	if err != nil {
		return err
	}
	return nil
}