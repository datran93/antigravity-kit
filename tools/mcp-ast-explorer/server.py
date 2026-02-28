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

def get_node_text(node, code):
    if not node: return ""
    return code[node.start_byte:node.end_byte].decode('utf8')

def get_docstring(node, lang_type, code):
    if lang_type == 'python':
        block = node.child_by_field_name('body')
        if block and len(block.children) > 0:
            first_child = block.children[0]
            if first_child.type == 'expression_statement':
                str_node = first_child.children[0]
                if str_node.type == 'string':
                    return get_node_text(str_node, code).strip('"\' \n')
    else:
        # Check node's own siblings first, then parent's siblings (for Go type_spec/JS var_dec)
        doc = []

        curr = node
        while curr:
            prev = curr.prev_sibling
            found_comment = False
            while prev and prev.type in ['comment', 'line_comment', 'block_comment']:
                doc.insert(0, get_node_text(prev, code).strip('/* \n'))
                prev = prev.prev_sibling
                found_comment = True

            if found_comment: break
            # if no comment found, try parent (once)
            if curr == node:
                curr = node.parent
                # filter out generic parents like 'source_file' or 'program'
                if curr and curr.type in ['source_file', 'program', 'module']:
                    break
            else:
                break

        if doc:
            return "\n".join(doc).strip()
    return ""

def extract_nodes(node, lang_type, code, level=0):
    results = []

    if lang_type == 'python':
        target_types = ['class_definition', 'function_definition']
    elif lang_type == 'go':
        target_types = ['type_spec', 'function_declaration', 'method_declaration']
    else: # js/ts
        target_types = ['class_declaration', 'function_declaration', 'method_definition', 'variable_declarator', 'interface_declaration', 'type_alias_declaration']

    if node.type in target_types:
        name_node = node.child_by_field_name('name')
        if not name_node and node.type == 'variable_declarator':
             # some JS arrow functions
             name_node = node.child_by_field_name('name')

        if name_node:
            name = get_node_text(name_node, code)
            doc = get_docstring(node, lang_type, code)

            signature = ""
            display_type = node.type.split('_')[0]

            if lang_type == 'python' and node.type == 'function_definition':
                params = node.child_by_field_name('parameters')
                ret = node.child_by_field_name('return_type')
                sig = get_node_text(params, code)
                if ret: sig += f" -> {get_node_text(ret, code)}"
                signature = sig
            elif lang_type == 'go':
                if node.type == 'method_declaration':
                    receiver = node.child_by_field_name('receiver')
                    params = node.child_by_field_name('parameters')
                    result = node.child_by_field_name('result')
                    sig = f"{get_node_text(receiver, code)} {get_node_text(params, code)}"
                    if result: sig += f" {get_node_text(result, code)}"
                    signature = sig
                    display_type = "method"
                elif node.type == 'function_declaration':
                    params = node.child_by_field_name('parameters')
                    result = node.child_by_field_name('result')
                    sig = get_node_text(params, code)
                    if result: sig += f" {get_node_text(result, code)}"
                    signature = sig
            elif lang_type in ['javascript', 'typescript', 'tsx']:
                if node.type in ['function_declaration', 'method_definition']:
                    params = node.child_by_field_name('parameters') or node.child_by_field_name('formal_parameters')
                    ret = node.child_by_field_name('return_type')
                    sig = get_node_text(params, code)
                    if ret: sig += get_node_text(ret, code)
                    signature = sig
                elif node.type == 'variable_declarator':
                    value_node = node.child_by_field_name('value')
                    if value_node and value_node.type == 'arrow_function':
                        params = value_node.child_by_field_name('parameters')
                        ret = value_node.child_by_field_name('return_type')
                        sig = get_node_text(params, code)
                        if ret: sig += get_node_text(ret, code)
                        signature = sig
                        display_type = "arrowFunc"

            results.append({
                "level": level,
                "type": display_type,
                "name": name,
                "signature": signature,
                "doc": doc
            })

    # recurse
    for child in node.children:
        inc = 1 if node.type in target_types else 0
        results.extend(extract_nodes(child, lang_type, code, level + inc))

    return results

import functools

# Global cache for AST nodes to speed up consecutive calls (like repeated search_symbol)
@functools.lru_cache(maxsize=1000)
def parse_and_extract(filepath: str, workspace_path: str, lang: str, mtime: float):
    """
    Parse and extract AST nodes from a file.
    Uses file mtime as cache key invalidation mechanism.
    """
    try:
        with open(filepath, 'rb') as f:
            code = f.read()
    except Exception:
        return None, None

    parser = PARSERS[lang]
    tree = parser.parse(code)
    nodes = extract_nodes(tree.root_node, lang, code)
    rel_path = os.path.relpath(filepath, workspace_path)
    return rel_path, nodes

@mcp.tool()
def get_project_architecture(workspace_path: str, sub_path: str = "", max_files: int = 150, include_docs: bool = False) -> str:
    """
    Get a structural overview (AST-based) of the project architecture.
    Extracts Classes, Functions, and Methods with signatures.
    """
    try:
        base_dir = os.path.join(workspace_path, sub_path)
        if not os.path.exists(base_dir):
            return f"❌ Path not found: {base_dir}"

        IGNORE_DIRS = {'.git', 'node_modules', 'vendor', '.venv', 'venv', 'dist', 'build', '.next', '.agent'}

        file_count = 0
        output = [f"🏗 PROJECT ARCHITECTURE AST: {sub_path or 'ROOT'}\n"]

        for root, dirs, files in os.walk(base_dir):
            dirs[:] = [d for d in dirs if d not in IGNORE_DIRS and not d.startswith('.')]
            files.sort()

            for file in files:
                ext = os.path.splitext(file)[1]
                lang = get_language_from_ext(ext)
                if not lang: continue

                filepath = os.path.join(root, file)

                try:
                    mtime = os.path.getmtime(filepath)
                except Exception: continue

                rel_path, nodes = parse_and_extract(filepath, workspace_path, lang, mtime)

                if nodes:
                    output.append(f"📄 {rel_path}")
                    for n in nodes:
                        indent = "  " * (n['level'] + 1)
                        sig = f" {n['signature']}" if n['signature'] else ""
                        output.append(f"{indent}▪ [{n['type']}] {n['name']}{sig}")
                        if include_docs and n['doc']:
                            # Only show first line of docstring to keep it concise
                            first_line = n['doc'].split('\n')[0]
                            output.append(f"{indent}  // {first_line}")
                    output.append("")

                file_count += 1
                if file_count >= max_files:
                    output.append(f"\n⚠️ Reached limit of {max_files} files.")
                    return "\n".join(output)

        return "\n".join(output)

    except Exception as e:
        return f"❌ Error: {str(e)}"

@mcp.tool()
def search_symbol(workspace_path: str, query: str) -> str:
    """
    Search for a class or function symbol across the project using AST.
    Useful for finding definitions quickly.
    """
    try:
        results = []
        IGNORE_DIRS = {'.git', 'node_modules', 'vendor', '.venv', 'venv', 'dist', 'build', '.next', '.agent'}

        for root, dirs, files in os.walk(workspace_path):
            dirs[:] = [d for d in dirs if d not in IGNORE_DIRS and not d.startswith('.')]
            for file in files:
                ext = os.path.splitext(file)[1]
                lang = get_language_from_ext(ext)
                if not lang: continue

                filepath = os.path.join(root, file)

                try:
                    mtime = os.path.getmtime(filepath)
                except Exception: continue

                rel_path, nodes = parse_and_extract(filepath, workspace_path, lang, mtime)

                if not nodes: continue

                for n in nodes:
                    if query.lower() in n['name'].lower():
                        results.append(f"📍 {rel_path} -> [{n['type']}] {n['name']}{n['signature']}")

        if not results:
            return f"🔍 No symbols matching '{query}' found."

        return "🔎 SYMBOL SEARCH RESULTS:\n" + "\n".join(results[:50])

    except Exception as e:
        return f"❌ Error searching symbol: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')

