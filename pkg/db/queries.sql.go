// Code generated by sqlc. DO NOT EDIT.
// source: queries.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (
    provider,
    job_name,
    job_id,
    url,
    started,
    finished,
    cluster_version,
    cluster_name,
    cluster_id,
    multi_az,
    channel,
    environment,
    region,
    numb_worker_nodes,
    network_provider,
    image_content_source,
    install_config,
    hibernate_after_use,
    reused,
    result
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING id
`

type CreateJobParams struct {
	Provider           string    `json:"provider"`
	JobName            string    `json:"job_name"`
	JobID              string    `json:"job_id"`
	Url                string    `json:"url"`
	Started            time.Time `json:"started"`
	Finished           time.Time `json:"finished"`
	ClusterVersion     string    `json:"cluster_version"`
	ClusterName        string    `json:"cluster_name"`
	ClusterID          string    `json:"cluster_id"`
	MultiAz            string    `json:"multi_az"`
	Channel            string    `json:"channel"`
	Environment        string    `json:"environment"`
	Region             string    `json:"region"`
	NumbWorkerNodes    int32     `json:"numb_worker_nodes"`
	NetworkProvider    string    `json:"network_provider"`
	ImageContentSource string    `json:"image_content_source"`
	InstallConfig      string    `json:"install_config"`
	HibernateAfterUse  bool      `json:"hibernate_after_use"`
	Reused             bool      `json:"reused"`
	Result             JobResult `json:"result"`
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (int64, error) {
	row := q.queryRow(ctx, q.createJobStmt, createJob,
		arg.Provider,
		arg.JobName,
		arg.JobID,
		arg.Url,
		arg.Started,
		arg.Finished,
		arg.ClusterVersion,
		arg.ClusterName,
		arg.ClusterID,
		arg.MultiAz,
		arg.Channel,
		arg.Environment,
		arg.Region,
		arg.NumbWorkerNodes,
		arg.NetworkProvider,
		arg.ImageContentSource,
		arg.InstallConfig,
		arg.HibernateAfterUse,
		arg.Reused,
		arg.Result,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createTestcase = `-- name: CreateTestcase :one
INSERT INTO testcases (
    job_id,
    result,
    name,
    duration,
    error,
    stdout,
    stderr
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`

type CreateTestcaseParams struct {
	JobID    int64           `json:"job_id"`
	Result   TestResult      `json:"result"`
	Name     string          `json:"name"`
	Duration pgtype.Interval `json:"duration"`
	Error    string          `json:"error"`
	Stdout   string          `json:"stdout"`
	Stderr   string          `json:"stderr"`
}

func (q *Queries) CreateTestcase(ctx context.Context, arg CreateTestcaseParams) (int64, error) {
	row := q.queryRow(ctx, q.createTestcaseStmt, createTestcase,
		arg.JobID,
		arg.Result,
		arg.Name,
		arg.Duration,
		arg.Error,
		arg.Stdout,
		arg.Stderr,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getJob = `-- name: GetJob :one
SELECT id, provider, job_name, job_id, url, started, finished, duration, cluster_version, cluster_name, cluster_id, multi_az, channel, environment, region, numb_worker_nodes, network_provider, image_content_source, install_config, hibernate_after_use, reused, result
FROM jobs
WHERE jobs.id = $1
`

func (q *Queries) GetJob(ctx context.Context, id int64) (Job, error) {
	row := q.queryRow(ctx, q.getJobStmt, getJob, id)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.Provider,
		&i.JobName,
		&i.JobID,
		&i.Url,
		&i.Started,
		&i.Finished,
		&i.Duration,
		&i.ClusterVersion,
		&i.ClusterName,
		&i.ClusterID,
		&i.MultiAz,
		&i.Channel,
		&i.Environment,
		&i.Region,
		&i.NumbWorkerNodes,
		&i.NetworkProvider,
		&i.ImageContentSource,
		&i.InstallConfig,
		&i.HibernateAfterUse,
		&i.Reused,
		&i.Result,
	)
	return i, err
}

const getTestcase = `-- name: GetTestcase :one
SELECT id, job_id, result, name, duration, error, stdout, stderr
FROM testcases
WHERE testcases.id = $1
`

func (q *Queries) GetTestcase(ctx context.Context, id int64) (Testcase, error) {
	row := q.queryRow(ctx, q.getTestcaseStmt, getTestcase, id)
	var i Testcase
	err := row.Scan(
		&i.ID,
		&i.JobID,
		&i.Result,
		&i.Name,
		&i.Duration,
		&i.Error,
		&i.Stdout,
		&i.Stderr,
	)
	return i, err
}

const getTestcaseForJob = `-- name: GetTestcaseForJob :many
SELECT id, job_id, result, name, duration, error, stdout, stderr
from testcases
WHERE testcases.job_id = $1
ORDER BY testcases.id
`

func (q *Queries) GetTestcaseForJob(ctx context.Context, jobID int64) ([]Testcase, error) {
	rows, err := q.query(ctx, q.getTestcaseForJobStmt, getTestcaseForJob, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Testcase
	for rows.Next() {
		var i Testcase
		if err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.Result,
			&i.Name,
			&i.Duration,
			&i.Error,
			&i.Stdout,
			&i.Stderr,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listJobs = `-- name: ListJobs :many
SELECT id, provider, job_name, job_id, url, started, finished, duration, cluster_version, cluster_name, cluster_id, multi_az, channel, environment, region, numb_worker_nodes, network_provider, image_content_source, install_config, hibernate_after_use, reused, result
FROM jobs
ORDER BY id
`

func (q *Queries) ListJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.query(ctx, q.listJobsStmt, listJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.Provider,
			&i.JobName,
			&i.JobID,
			&i.Url,
			&i.Started,
			&i.Finished,
			&i.Duration,
			&i.ClusterVersion,
			&i.ClusterName,
			&i.ClusterID,
			&i.MultiAz,
			&i.Channel,
			&i.Environment,
			&i.Region,
			&i.NumbWorkerNodes,
			&i.NetworkProvider,
			&i.ImageContentSource,
			&i.InstallConfig,
			&i.HibernateAfterUse,
			&i.Reused,
			&i.Result,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTestcases = `-- name: ListTestcases :many
SELECT id, job_id, result, name, duration, error, stdout, stderr
FROM testcases
ORDER BY id
`

func (q *Queries) ListTestcases(ctx context.Context) ([]Testcase, error) {
	rows, err := q.query(ctx, q.listTestcasesStmt, listTestcases)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Testcase
	for rows.Next() {
		var i Testcase
		if err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.Result,
			&i.Name,
			&i.Duration,
			&i.Error,
			&i.Stdout,
			&i.Stderr,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}