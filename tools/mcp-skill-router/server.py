import os
from mcp.server.fastmcp import FastMCP
import chromadb
from chromadb.utils import embedding_functions

# Config
DB_DIR = "/Users/datran/LearnDev/antigravity-kit/tools/mcp-skill-router/.chroma_db"

# Initialize MCP Server
mcp = FastMCP("McpSkillRouter")

@mcp.tool()
def search_skills(query: str, top_k: int = 3) -> str:
    """
    Search for relevant AI agent skills based on a semantic query.
    Use this tool when you need to find which skills are best suited for a user's request.

    Args:
        query: Semantic query to search for (e.g., 'beautiful ui design', 'fix database query', 'how to write python backend')
        top_k: Number of skills to return (default is 3)
    """
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        return "ERROR: OPENAI_API_KEY environment variable is not set."

    try:
        client = chromadb.PersistentClient(path=DB_DIR)
        emb_fn = embedding_functions.OpenAIEmbeddingFunction(
            api_key=api_key,
            model_name="text-embedding-3-small"
        )

        try:
            collection = client.get_collection(
                name="skills_collection_openai",
                embedding_function=emb_fn
            )
        except ValueError:
            return "ERROR: Collection not found. Please run `python skill_indexer.py` first to generate the Vector Index."

        results = collection.query(
            query_texts=[query],
            n_results=top_k
        )

        if not results['documents'] or not len(results['documents'][0]):
            return "❌ No relevant skills found."

        formatted_results = [f"🎯 SEMANTIC SEARCH RESULTS FOR QUERY: '{query}'"]
        for i in range(len(results['ids'][0])):
            skill_name = results['ids'][0][i]
            metadata = results['metadatas'][0][i]

            res = f"\n{i+1}. **{skill_name}**"
            res += f"\n   - Path: `{metadata.get('path', 'Unknown')}`"
            res += f"\n   - Description: {metadata.get('description', '')}"
            formatted_results.append(res)

        formatted_results.append("\n💡 ADVICE FOR AGENT: Use the `view_file` (or `view_file_outline`) tool on the Path provided above to read the skill details before starting your task.")
        return "\n".join(formatted_results)

    except Exception as e:
        return f"❌ Error during search: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
