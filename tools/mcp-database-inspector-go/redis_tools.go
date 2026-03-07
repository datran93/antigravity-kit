package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func handleRedisListTables(connStr string) (string, error) {
	client, err := getRedisClient(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to redis: %v", err)
	}

	ctx := context.Background()
	dbSize, err := client.DBSize(ctx).Result()
	if err != nil {
		return "", err
	}

	var cursor uint64
	keys, _, err := client.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("🗄 REDIS DATABASE (Size: %d keys)\n\n", dbSize))
	result.WriteString("### Sample Keys (up to 100):\n")
	for _, k := range keys {
		result.WriteString(fmt.Sprintf("- %s\n", k))
	}

	if len(keys) == 0 {
		return "❌ No keys found in this Redis database.", nil
	}

	return result.String(), nil
}

func handleRedisInspectSchema(connStr string, keyName string) (string, error) {
	client, err := getRedisClient(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to redis: %v", err)
	}

	ctx := context.Background()
	exists, err := client.Exists(ctx, keyName).Result()
	if err != nil {
		return "", err
	}
	if exists == 0 {
		return fmt.Sprintf("❌ Key '%s' does not exist in Redis.", keyName), nil
	}

	keyType, err := client.Type(ctx, keyName).Result()
	if err != nil {
		return "", err
	}

	ttl, err := client.TTL(ctx, keyName).Result()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📋 SCHEMA FOR KEY: `%s`\n\n", keyName))
	result.WriteString("### Metadata:\n")
	result.WriteString(fmt.Sprintf("- **Type**: %s\n", keyType))
	result.WriteString(fmt.Sprintf("- **TTL**: %v seconds\n\n", ttl.Seconds()))

	result.WriteString("### Details:\n")
	switch keyType {
	case "string":
		valLen, _ := client.StrLen(ctx, keyName).Result()
		result.WriteString(fmt.Sprintf("- String Length: %d\n", valLen))
	case "hash":
		hlen, _ := client.HLen(ctx, keyName).Result()
		result.WriteString(fmt.Sprintf("- Hash Entries: %d\n", hlen))
	case "list":
		llen, _ := client.LLen(ctx, keyName).Result()
		result.WriteString(fmt.Sprintf("- List Length: %d\n", llen))
	case "set":
		scard, _ := client.SCard(ctx, keyName).Result()
		result.WriteString(fmt.Sprintf("- Set Members: %d\n", scard))
	case "zset":
		zcard, _ := client.ZCard(ctx, keyName).Result()
		result.WriteString(fmt.Sprintf("- Sorted Set Members: %d\n", zcard))
	}

	return result.String(), nil
}

func handleRedisRunQuery(connStr string, query string, isWrite bool, confirm bool) (string, error) {
	if isWrite && !confirm {
		return fmt.Sprintf("⚠️  CONFIRMATION REQUIRED: You are about to execute a WRITE/DML operation:\n\n`%s`\n\nPlease confirm if you want to proceed. Set 'confirm=True' only after user approval.", query), nil
	}

	client, err := getRedisClient(connStr)
	if err != nil {
		return "", fmt.Errorf("error connecting to redis: %v", err)
	}

	parts := strings.Fields(query)
	if len(parts) == 0 {
		return "❌ Empty query", nil
	}

	command := strings.ToUpper(parts[0])

	if !isWrite {
		safeCommands := map[string]bool{
			"GET": true, "MGET": true, "HGET": true, "HGETALL": true, "HMGET": true,
			"HKEYS": true, "HVALS": true, "HLEN": true, "LRANGE": true, "LLEN": true,
			"LINDEX": true, "SMEMBERS": true, "SCARD": true, "SISMEMBER": true,
			"ZRANGE": true, "ZCARD": true, "ZSCORE": true, "ZREVRANGE": true,
			"TYPE": true, "TTL": true, "EXISTS": true, "SCAN": true, "INFO": true,
			"DBSIZE": true, "PING": true,
		}
		if !safeCommands[command] {
			return fmt.Sprintf("❌ SECURITY BLOCK: Command '%s' is not allowed in Read-only mode.", command), nil
		}
	}

	// Prepare arguments for redis Do
	var args []interface{}
	for _, p := range parts {
		args = append(args, p)
	}

	ctx := context.Background()
	res, err := client.Do(ctx, args...).Result()
	if err != nil {
		return fmt.Sprintf("❌ Redis Error: %v", err), nil
	}

	jsonBytes, err := json.MarshalIndent(map[string]interface{}{"result": res}, "", "  ")
	if err != nil {
		return fmt.Sprintf("✅ Result: %v", res), nil
	}

	if isWrite {
		return fmt.Sprintf("✅ REDIS WRITE SUCCESS\nCMD = `%s`\nResult:\n```json\n%s\n```", query, string(jsonBytes)), nil
	}

	return fmt.Sprintf("✅ REDIS RESULT\n```json\n%s\n```", string(jsonBytes)), nil
}
