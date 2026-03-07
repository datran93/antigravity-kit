package main

import (
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

var Parsers map[string]*sitter.Language

func init() {
	Parsers = map[string]*sitter.Language{
		"python":     python.GetLanguage(),
		"go":         golang.GetLanguage(),
		"javascript": javascript.GetLanguage(),
		"typescript": typescript.GetLanguage(),
		"tsx":        tsx.GetLanguage(),
	}
}

func getLanguageFromExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js", ".jsx", ".cjs", ".mjs":
		return "javascript"
	case ".ts", ".cts", ".mts":
		return "typescript"
	case ".tsx":
		return "tsx"
	}
	return ""
}

func getFamilyFromExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".js", ".jsx", ".cjs", ".mjs":
		return "javascript"
	case ".ts", ".cts", ".mts", ".tsx":
		return "typescript"
	}
	return ""
}

type NodeResult struct {
	Level     int
	Type      string
	Name      string
	Signature string
	Doc       string
}

func getNodeText(node *sitter.Node, code []byte) string {
	if node == nil {
		return ""
	}
	return string(code[node.StartByte():node.EndByte()])
}

func getDocstring(node *sitter.Node, langType string, code []byte) string {
	if langType == "python" {
		block := node.ChildByFieldName("body")
		if block != nil && int(block.NamedChildCount()) > 0 {
			firstChild := block.NamedChild(0)
			if firstChild != nil && firstChild.Type() == "expression_statement" {
				strNode := firstChild.NamedChild(0)
				if strNode != nil && strNode.Type() == "string" {
					text := getNodeText(strNode, code)
					return strings.Trim(text, "\"'\n ")
				}
			}
		}
	} else {
		var doc []string
		curr := node
		for curr != nil {
			prev := curr.PrevSibling()
			foundComment := false
			for prev != nil && (prev.Type() == "comment" || prev.Type() == "line_comment" || prev.Type() == "block_comment") {
				text := getNodeText(prev, code)
				// Trim leading slashes and stars, spaces
				text = strings.TrimLeft(text, "/* ")
				text = strings.TrimRight(text, "*/ \n")
				doc = append([]string{strings.TrimSpace(text)}, doc...) // prepend
				prev = prev.PrevSibling()
				foundComment = true
			}

			if foundComment {
				break
			}
			if curr == node {
				curr = node.Parent()
				if curr != nil && (curr.Type() == "source_file" || curr.Type() == "program" || curr.Type() == "module") {
					break
				}
			} else {
				break
			}
		}
		if len(doc) > 0 {
			return strings.Join(doc, "\n")
		}
	}
	return ""
}

func extractNodes(node *sitter.Node, langType string, code []byte, level int) []NodeResult {
	if node == nil {
		return nil
	}

	var results []NodeResult

	targetTypes := make(map[string]bool)
	if langType == "python" {
		targetTypes["class_definition"] = true
		targetTypes["function_definition"] = true
	} else if langType == "go" {
		targetTypes["type_spec"] = true
		targetTypes["function_declaration"] = true
		targetTypes["method_declaration"] = true
	} else { // js/ts
		targetTypes["class_declaration"] = true
		targetTypes["function_declaration"] = true
		targetTypes["method_definition"] = true
		targetTypes["variable_declarator"] = true
		targetTypes["interface_declaration"] = true
		targetTypes["type_alias_declaration"] = true
	}

	if targetTypes[node.Type()] {
		nameNode := node.ChildByFieldName("name")
		if nameNode == nil && node.Type() == "variable_declarator" {
			nameNode = node.ChildByFieldName("name")
		}

		if nameNode != nil {
			name := getNodeText(nameNode, code)
			doc := getDocstring(node, langType, code)

			signature := ""
			displayType := strings.Split(node.Type(), "_")[0]

			if langType == "python" && node.Type() == "function_definition" {
				params := node.ChildByFieldName("parameters")
				ret := node.ChildByFieldName("return_type")
				sig := getNodeText(params, code)
				if ret != nil {
					sig += " -> " + getNodeText(ret, code)
				}
				signature = sig
			} else if langType == "go" {
				if node.Type() == "method_declaration" {
					receiver := node.ChildByFieldName("receiver")
					params := node.ChildByFieldName("parameters")
					result := node.ChildByFieldName("result")
					sig := getNodeText(receiver, code) + " " + getNodeText(params, code)
					if result != nil {
						sig += " " + getNodeText(result, code)
					}
					signature = sig
					displayType = "method"
				} else if node.Type() == "function_declaration" {
					params := node.ChildByFieldName("parameters")
					result := node.ChildByFieldName("result")
					sig := getNodeText(params, code)
					if result != nil {
						sig += " " + getNodeText(result, code)
					}
					signature = sig
				}
			} else if langType == "javascript" || langType == "typescript" || langType == "tsx" {
				if node.Type() == "function_declaration" || node.Type() == "method_definition" {
					params := node.ChildByFieldName("parameters")
					if params == nil {
						params = node.ChildByFieldName("formal_parameters")
					}
					ret := node.ChildByFieldName("return_type")
					sig := getNodeText(params, code)
					if ret != nil {
						sig += getNodeText(ret, code)
					}
					signature = sig
				} else if node.Type() == "variable_declarator" {
					valueNode := node.ChildByFieldName("value")
					if valueNode != nil && valueNode.Type() == "arrow_function" {
						params := valueNode.ChildByFieldName("parameters")
						ret := valueNode.ChildByFieldName("return_type")
						sig := getNodeText(params, code)
						if ret != nil {
							sig += getNodeText(ret, code)
						}
						signature = sig
						displayType = "arrowFunc"
					}
				}
			}

			// Clean up signature: single line formatting
			signature = strings.ReplaceAll(signature, "\n", " ")
			signature = strings.Join(strings.Fields(signature), " ")

			results = append(results, NodeResult{
				Level:     level,
				Type:      displayType,
				Name:      name,
				Signature: signature,
				Doc:       doc,
			})
		}
	}

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		inc := 0
		if targetTypes[node.Type()] {
			inc = 1
		}
		childResults := extractNodes(child, langType, code, level+inc)
		results = append(results, childResults...)
	}

	return results
}

func parseAndExtract(filepathStr, workspacePath, lang string) (string, []NodeResult) {
	code, err := os.ReadFile(filepathStr)
	if err != nil {
		return "", nil
	}

	langPtr, ok := Parsers[lang]
	if !ok {
		return "", nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(langPtr)

	// Create a new parser tree
	tree := parser.Parse(nil, code)
	defer tree.Close() // Keep memory safe if supported

	nodes := extractNodes(tree.RootNode(), lang, code, 0)

	// Ensure uniform paths
	relPath, err := filepath.Rel(workspacePath, filepathStr)
	if err != nil {
		relPath = filepathStr
	}

	return relPath, nodes
}
