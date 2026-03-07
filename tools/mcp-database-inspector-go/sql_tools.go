package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func handleSQLListTables(connStr string) (string, error) {
	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	migrator := db.Migrator()

	tables, err := migrator.GetTables()
	if err != nil {
		return "", fmt.Errorf("error getting tables: %v", err)
	}

	var result strings.Builder
	result.WriteString("🗄 DATABASE ENTITIES\n\n")
	result.WriteString("### Tables & Views:\n")
	for _, t := range tables {
		result.WriteString(fmt.Sprintf("- %s\n", t))
	}

	if len(tables) == 0 {
		return "❌ No tables found in this database.", nil
	}

	return result.String(), nil
}

func handleSQLInspectSchema(connStr string, tableName string) (string, error) {
	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	migrator := db.Migrator()

	if !migrator.HasTable(tableName) {
		return fmt.Sprintf("❌ Table '%s' does not exist in the database.", tableName), nil
	}

	colTypes, err := migrator.ColumnTypes(tableName)
	if err != nil {
		return "", fmt.Errorf("error getting columns: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📋 SCHEMA FOR TABLE: `%s`\n\n", tableName))
	result.WriteString("### Columns:\n")

	for _, col := range colTypes {
		colType, _ := col.ColumnType()
		nullable := "NOT NULL"
		if nullRef, ok := col.Nullable(); ok && nullRef {
			nullable = "NULL"
		}

		prefix := "   "
		if pk, ok := col.PrimaryKey(); ok && pk {
			prefix = "🔑 "
		}

		result.WriteString(fmt.Sprintf("%s`%s` : **%s** (%s)\n", prefix, col.Name(), colType, nullable))
	}

	return result.String(), nil
}

func handleSQLExplainQuery(connStr string, query string) (string, error) {
	forbidden := regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|ALTER|TRUNCATE|GRANT|REVOKE|CREATE)\b`)
	if forbidden.MatchString(query) {
		return "❌ SECURITY BLOCK: explain_query is only for SELECT statements.", nil
	}

	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	dbType := db.Dialector.Name()
	explainPrefix := "EXPLAIN "
	if dbType == "postgres" {
		explainPrefix = "EXPLAIN (ANALYZE, VERBOSE, BUFFERS) "
	} else if dbType == "mysql" {
		explainPrefix = "EXPLAIN format=json "
	}

	fullQuery := explainPrefix + query

	sqlDB, err := db.DB()
	if err != nil {
		return "", err
	}

	rows, err := sqlDB.Query(fullQuery)
	if err != nil {
		return "", fmt.Errorf("error explaining query: %v", err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	var result strings.Builder
	result.WriteString(fmt.Sprintf("🔍 **QUERY PLAN (%s)**\n\n```text\n", strings.ToUpper(dbType)))

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return "", err
		}

		for i, col := range columns {
			if b, ok := col.([]byte); ok {
				result.WriteString(string(b))
			} else {
				result.WriteString(fmt.Sprintf("%v", col))
			}
			if i < len(columns)-1 {
				result.WriteString(" | ")
			}
		}
		result.WriteString("\n")
	}
	result.WriteString("```\n")

	return result.String(), nil
}

func handleSQLGetTableSample(connStr string, tableName string) (string, error) {
	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	migrator := db.Migrator()
	if !migrator.HasTable(tableName) {
		return fmt.Sprintf("❌ Table '%s' does not exist.", tableName), nil
	}

	colTypes, err := migrator.ColumnTypes(tableName)
	if err != nil {
		return "", fmt.Errorf("error getting columns: %v", err)
	}

	var schemaMd strings.Builder
	schemaMd.WriteString(fmt.Sprintf("### 📊 Table Schema: `%s`\n\n", tableName))
	schemaMd.WriteString("| Column | Type | Nullable | PK |\n")
	schemaMd.WriteString("| :--- | :--- | :--- | :--- |\n")

	for _, col := range colTypes {
		colType, _ := col.ColumnType()
		nullable := "NOT NULL"
		if nullRef, ok := col.Nullable(); ok && nullRef {
			nullable = "NULL"
		}
		isPk := ""
		if pk, ok := col.PrimaryKey(); ok && pk {
			isPk = "✅"
		}

		schemaMd.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n", col.Name(), colType, nullable, isPk))
	}

	// Fetch sample
	sqlDB, err := db.DB()
	if err != nil {
		return schemaMd.String(), nil
	}

	sampleQuery := fmt.Sprintf("SELECT * FROM %s LIMIT 5", tableName)
	rows, err := sqlDB.Query(sampleQuery)
	if err != nil {
		schemaMd.WriteString(fmt.Sprintf("\n*Error fetching sample data: %v*", err))
		return schemaMd.String(), nil
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	if len(cols) == 0 {
		schemaMd.WriteString("\n*No sample data found (table is empty).*")
		return schemaMd.String(), nil
	}

	var allRows []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err == nil {
			rowMap := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				if b, ok := (*val).([]byte); ok {
					rowMap[colName] = string(b)
				} else if t, ok := (*val).(time.Time); ok {
					rowMap[colName] = t.Format(time.RFC3339)
				} else {
					rowMap[colName] = *val
				}
			}
			allRows = append(allRows, rowMap)
		}
	}

	if len(allRows) == 0 {
		schemaMd.WriteString("\n*No sample data found (table is empty).*")
	} else {
		schemaMd.WriteString("\n### 📝 Sample Data (First 5 Rows)\n\n")
		schemaMd.WriteString("| " + strings.Join(cols, " | ") + " |\n")

		dashes := make([]string, len(cols))
		for i := range dashes {
			dashes[i] = "---"
		}
		schemaMd.WriteString("| " + strings.Join(dashes, " | ") + " |\n")

		for _, rowMap := range allRows {
			vals := make([]string, len(cols))
			for i, colName := range cols {
				v := rowMap[colName]
				jsonBytes, _ := json.Marshal(v)
				vals[i] = string(jsonBytes)
			}
			schemaMd.WriteString("| " + strings.Join(vals, " | ") + " |\n")
		}
	}

	return schemaMd.String(), nil
}

func handleSQLRunReadQuery(connStr string, query string, limit int, offset int) (string, error) {
	if limit > 2000 {
		limit = 2000
	}

	forbidden := regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|ALTER|TRUNCATE|GRANT|REVOKE|CREATE)\b`)
	if forbidden.MatchString(query) {
		cleanQuery := strings.TrimSpace(strings.ToUpper(query))
		if !strings.HasPrefix(cleanQuery, "SELECT") {
			return "❌ SECURITY BLOCK: Only SELECT/Read-only queries are allowed in this tool.", nil
		}
	}

	cleanQuery := strings.TrimSpace(query)
	upperQuery := strings.ToUpper(cleanQuery)
	if !strings.Contains(upperQuery, "LIMIT") && strings.HasPrefix(upperQuery, "SELECT") {
		cleanQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", cleanQuery, limit, offset)
	}

	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return "", err
	}

	rows, err := sqlDB.Query(cleanQuery)
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err), nil
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var allRows []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err == nil {
			rowMap := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				if b, ok := (*val).([]byte); ok {
					rowMap[colName] = string(b)
				} else if t, ok := (*val).(time.Time); ok {
					rowMap[colName] = t.Format(time.RFC3339)
				} else {
					rowMap[colName] = *val
				}
			}
			allRows = append(allRows, rowMap)
		}
	}

	if len(allRows) == 0 {
		return "✅ Query executed. 0 rows returned.", nil
	}

	jsonBytes, err := json.MarshalIndent(allRows, "", "  ")
	if err != nil {
		return "", err
	}

	suffix := fmt.Sprintf("\n*(Showing %d rows starting at offset %d)*", len(allRows), offset)
	return fmt.Sprintf("✅ QUERY RESULTS\n```json\n%s\n```%s", string(jsonBytes), suffix), nil
}

func handleSQLRunWriteQuery(connStr string, query string, confirm bool) (string, error) {
	if !confirm {
		return fmt.Sprintf("⚠️  CONFIRMATION REQUIRED: You are about to execute a WRITE/DML operation:\n\n`%s`\n\nPlease confirm if you want to proceed. Set 'confirm=True' only after user approval.", query), nil
	}

	db, err := getGormDB(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to db: %v", err)
	}

	// split by ; simple implementation
	statements := strings.Split(query, ";")

	totalRowsAffected := int64(0)

	sqlDB, err := db.DB()
	if err != nil {
		return "", err
	}

	tx, err := sqlDB.Begin()
	if err != nil {
		return "", err
	}

	executed := 0
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		res, err := tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			errorMsg := strings.ToLower(err.Error())
			if strings.Contains(errorMsg, "not exist") {
				return fmt.Sprintf("❌ Table or column does not exist: %v", err), nil
			}
			if strings.Contains(errorMsg, "duplicate key") || strings.Contains(errorMsg, "integrity") {
				return fmt.Sprintf("❌ Integrity Error (Constraint violation): %v", err), nil
			}
			return fmt.Sprintf("❌ DB Error: %v", err), nil
		}

		rowsAff, _ := res.RowsAffected()
		if rowsAff > 0 {
			totalRowsAffected += rowsAff
		}
		executed++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Sprintf("❌ Transaction Commit Error: %v", err), nil
	}

	if executed == 0 {
		return "❌ No valid SQL statements found.", nil
	}

	return fmt.Sprintf("✅ WRITE SUCCESS\nExecuted %d statement(s).\nTotal rows affected: %d", executed, totalRowsAffected), nil
}
