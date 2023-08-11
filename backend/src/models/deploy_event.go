package models

import (
	"gorm.io/gorm"
	"time"
)

type DeployEvent struct {
	ID         int64     `json:"id" gorm:"->;primaryKey;column:id"`
	Tag        string    `json:"tag" gorm:"->;column:tag"`
	TaskID     int64     `json:"task_id" gorm:"->;column:task_id"`
	Start      time.Time `json:"start" gorm:"->;column:start"`
	Job        string    `json:"job" gorm:"->;column:job"`
	Cluster    string    `json:"cluster" gorm:"->;column:cluster"`
	JobCluster string    `json:"job_cluster" gorm:"->;column:job_cluster"`
	CodeUrl    string    `json:"code_url" gorm:"->;column:code_rrl"`
	State      string    `json:"state" gorm:"->;column:state"`
	User       string    `json:"user" gorm:"->;column:user"`
	Version    string    `json:"version" gorm:"->;column:version"`
	UpdateTime time.Time `json:"update_time" gorm:"->;column:update_time"`
}

func (DeployEvent) TableName() string {
	return "deploy_event"
}

func QueryJobs(db *gorm.DB, productLines []string, start time.Time, end time.Time) []DeployEvent {
	var res []DeployEvent
	db.Select(
		[]string{"start", "job", "update_time"},
	).Where(
		"start >= ? AND start <= ? AND tag IN ?",
		start, end, productLines,
	).Find(
		&res,
	)
	return res
}

func QueryProductLines(db *gorm.DB) []string {
	var res []DeployEvent
	db.Distinct("tag").Find(&res)
	var productLines []string
	for _, item := range res {
		productLines = append(productLines, item.Tag)
	}
	return productLines
}
