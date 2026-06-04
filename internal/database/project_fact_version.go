package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProjectFactVersion 事实历史快照（同 fact_key 更新前归档）。
type ProjectFactVersion struct {
	ID                     string    `json:"id"`
	FactID                 string    `json:"fact_id"`
	ProjectID              string    `json:"project_id"`
	FactKey                string    `json:"fact_key"`
	Category               string    `json:"category"`
	Summary                string    `json:"summary"`
	Body                   string    `json:"body"`
	Confidence             string    `json:"confidence"`
	SourceConversationID   string    `json:"source_conversation_id,omitempty"`
	SourceMessageID        string    `json:"source_message_id,omitempty"`
	Pinned                 bool      `json:"pinned"`
	RelatedVulnerabilityID string    `json:"related_vulnerability_id,omitempty"`
	ArchivedAt             time.Time `json:"archived_at"`
}

// InsertProjectFactVersion 将当前事实行快照写入版本表。
func (db *DB) InsertProjectFactVersion(f *ProjectFact) (string, error) {
	if f == nil || f.ID == "" {
		return "", fmt.Errorf("无效的事实记录")
	}
	id := uuid.New().String()
	now := time.Now()
	_, err := db.Exec(
		`INSERT INTO project_fact_versions (
			id, fact_id, project_id, fact_key, category, summary, body, confidence,
			source_conversation_id, source_message_id, pinned, related_vulnerability_id, archived_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, f.ID, f.ProjectID, f.FactKey, f.Category, f.Summary, f.Body, f.Confidence,
		nullIfEmpty(f.SourceConversationID), nullIfEmpty(f.SourceMessageID), boolToInt(f.Pinned),
		nullIfEmpty(f.RelatedVulnerabilityID), now,
	)
	if err != nil {
		return "", fmt.Errorf("归档事实版本失败: %w", err)
	}
	return id, nil
}

// GetProjectFactVersion 按版本 ID 获取快照。
func (db *DB) GetProjectFactVersion(versionID string) (*ProjectFactVersion, error) {
	row := db.QueryRow(
		`SELECT id, fact_id, project_id, fact_key, category, summary, COALESCE(body,''), confidence,
			COALESCE(source_conversation_id,''), COALESCE(source_message_id,''), pinned,
			COALESCE(related_vulnerability_id,''), archived_at
		 FROM project_fact_versions WHERE id = ?`, versionID,
	)
	return scanProjectFactVersionRow(row)
}

// ListProjectFactVersions 列出某条事实的全部历史版本（新→旧）。
func (db *DB) ListProjectFactVersions(factID string, limit int) ([]*ProjectFactVersion, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := db.Query(
		`SELECT id, fact_id, project_id, fact_key, category, summary, COALESCE(body,''), confidence,
			COALESCE(source_conversation_id,''), COALESCE(source_message_id,''), pinned,
			COALESCE(related_vulnerability_id,''), archived_at
		 FROM project_fact_versions WHERE fact_id = ? ORDER BY archived_at DESC LIMIT ?`,
		factID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ProjectFactVersion
	for rows.Next() {
		v, err := scanProjectFactVersionFromRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func projectFactContentChanged(existing, incoming *ProjectFact) bool {
	if existing == nil || incoming == nil {
		return false
	}
	mergedBody := mergeFactBodyOnUpdate(incoming.Body, existing.Body)
	inCat := stringsTrimDefault(incoming.Category, existing.Category)
	inConf := stringsTrimDefault(incoming.Confidence, existing.Confidence)
	return existing.Summary != incoming.Summary ||
		existing.Body != mergedBody ||
		existing.Category != inCat ||
		existing.Confidence != inConf
}

func stringsTrimDefault(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return strings.TrimSpace(s)
}

func scanProjectFactVersionRow(row *sql.Row) (*ProjectFactVersion, error) {
	var v ProjectFactVersion
	var pinned int
	var archivedAt string
	err := row.Scan(
		&v.ID, &v.FactID, &v.ProjectID, &v.FactKey, &v.Category, &v.Summary, &v.Body, &v.Confidence,
		&v.SourceConversationID, &v.SourceMessageID, &pinned,
		&v.RelatedVulnerabilityID, &archivedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("事实版本不存在")
		}
		return nil, err
	}
	v.Pinned = pinned != 0
	v.ArchivedAt = parseDBTime(archivedAt)
	return &v, nil
}

func scanProjectFactVersionFromRows(rows *sql.Rows) (*ProjectFactVersion, error) {
	var v ProjectFactVersion
	var pinned int
	var archivedAt string
	err := rows.Scan(
		&v.ID, &v.FactID, &v.ProjectID, &v.FactKey, &v.Category, &v.Summary, &v.Body, &v.Confidence,
		&v.SourceConversationID, &v.SourceMessageID, &pinned,
		&v.RelatedVulnerabilityID, &archivedAt,
	)
	if err != nil {
		return nil, err
	}
	v.Pinned = pinned != 0
	v.ArchivedAt = parseDBTime(archivedAt)
	return &v, nil
}
