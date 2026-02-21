import os
from mcp.server.fastmcp import FastMCP
import tree_sitter
import tree_sitter_python
import tree_sitter_go
import tree_sitter_javascript
import tree_sitter_typescript
from tree_sitter import Language, Parser

mcp = FastMCP("McpAstExplorer")

# Initialize languages and parsers
LANGUAGES = {
    'python': Language(tree_sitter_python.language()),
    'go': Language(tree_sitter_go.language()),
    'javascript': Language(tree_sitter_javascript.language()),
    'typescript': Language(tree_sitter_typescript.language_typescript()),
    'tsx': Language(tree_sitter_typescript.language_tsx())
}

PARSERS = {
    name: Parser(lang) for name, lang in LANGUAGES.items()
}

def get_language_from_ext(ext):
    ext = ext.lower()
    if ext == '.py':
        return 'python'
    elif ext == '.go':
        return 'go'
    elif ext in ['.js', '.jsx', '.cjs', '.mjs']:
        return 'javascript'
    elif ext in ['.ts', '.cts', '.mts']:
        return 'typescript'
    elif ext == '.tsx':
        return 'tsx'
    return None

def extract_nodes(node, lang_type, level=0):
    results = []

    if lang_type == 'python':
        target_types = ['class_definition', 'function_definition']
    elif lang_type == 'go':
        target_types = ['type_spec', 'function_declaration', 'method_declaration']
    else: # js/ts
        target_types = ['class_declaration', 'function_declaration', 'method_definition', 'variable_declarator', 'interface_declaration', 'type_alias_declaration']

    if node.type in target_types:
        name_node = node.child_by_field_name('name')
        if name_node:
            name = name_node.text.decode('utf8')
            if node.type == 'variable_declarator': # JS/TS specific for Arrow Functions
                value_node = node.child_by_field_name('value')
                if value_node and value_node.type == 'arrow_function':
                    results.append((level, "arrowFunc", name))
            elif node.type == 'type_spec':
                # check if it's a struct or interface
                type_node = node.child_by_field_name('type')
                type_desc = "type"
                if type_node:
                    if type_node.type == 'struct_type': type_desc = "struct"
                    elif type_node.type == 'interface_type': type_desc = "interface"
                results.append((level, type_desc, name))
            elif node.type == 'method_declaration' and lang_type == 'go':
                # try to get receiver
                receiver_node = node.child_by_field_name('receiver')
                receiver_str = ""
                if receiver_node:
                    receiver_str = receiver_node.text.decode('utf8')
                results.append((level, "method", f"{receiver_str} {name}"))
            else:
                display_type = node.type.split('_')[0] # class, function, method, interface, type
                results.append((level, display_type, name))

    # recurse
    for child in node.children:
        # Increase level only if we found something useful in this parent, simple heuristic
        inc = 1 if node.type in target_types else 0
        results.extend(extract_nodes(child, lang_type, level + inc))

    return results

@mcp.tool()
def get_project_architecture(workspace_path: str, sub_path: str = "", max_files: int = 150) -> str:
    """
    Get a structural overview (AST-based) of the project architecture.
    Extracts Classes, Functions, and Methods across Python, Go, and JS/TS files.

    Args:
        workspace_path: Absolute path to the workspace root.
        sub_path: Optional relative path to narrow down the search (e.g., 'backend/api').
        max_files: Limit the number of parsed files to avoid overwhelming the context window.
    """
    try:
        base_dir = os.path.join(workspace_path, sub_path)
        if not os.path.exists(base_dir):
            return f"❌ Path not found: {base_dir}"

        IGNORE_DIRS = {'.git', 'node_modules', 'vendor', '.venv', 'venv', 'dist', 'build', '.next', '.agent'}

        file_count = 0
        output = [f"🏗 PROJECT ARCHITECTURE AST: {sub_path or 'ROOT'}\n"]

        for root, dirs, files in os.walk(base_dir):
            # filter dirs
            dirs[:] = [d for d in dirs if d not in IGNORE_DIRS and not d.startswith('.')]

            # sort for deterministic output
            files.sort()

            for file in files:
                ext = os.path.splitext(file)[1]
                lang = get_language_from_ext(ext)
                if not lang:
                    continue

                filepath = os.path.join(root, file)
                rel_path = os.path.relpath(filepath, workspace_path)

                try:
                    with open(filepath, 'rb') as f:
                        code = f.read()
                except:
                    continue

                parser = PARSERS[lang]
                tree = parser.parse(code)

                nodes = extract_nodes(tree.root_node, lang)

                if nodes:
                    output.append(f"📄 {rel_path}")
                    for lvl, typ, name in nodes:
                        indent = "  " * (lvl + 1)
                        output.append(f"{indent}▪ [{typ}] {name}")
                    output.append("")

                file_count += 1
                if file_count >= max_files:
                    output.append(f"\n⚠️ Reached limit of {max_files} files. Use `sub_path` to focus on specific directories.")
                    return "\n".join(output)

        return "\n".join(output)

    except Exception as e:
        return f"❌ Error extracting architecture: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
