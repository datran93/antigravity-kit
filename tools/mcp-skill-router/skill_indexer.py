import os
import glob
import re
import yaml
import chromadb
from chromadb.utils import embedding_functions
from typing import Dict

# Config
SKILLS_DIR = "/Users/datran/LearnDev/antigravity-kit/.agent/skills"
DB_DIR = "/Users/datran/LearnDev/antigravity-kit/tools/mcp-skill-router/.chroma_db"

def parse_skill_file(filepath: str) -> Dict:
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Extract frontmatter
    match = re.search(r'^---\n(.*?)\n---', content, re.DOTALL)
    metadata = {}
    if match:
        try:
            metadata = yaml.safe_load(match.group(1))
        except Exception:
            pass

    # Extract preview
    text_content = re.sub(r'^---\n(.*?)\n---', '', content, flags=re.DOTALL).strip()
    preview = text_content[:500]

    skill_name = os.path.basename(os.path.dirname(filepath))
    description = metadata.get('description', '')
    tags = metadata.get('tags', [])
    if isinstance(tags, list):
        tags = ', '.join(tags)

    search_text = f"Skill: {skill_name}\nTags: {tags}\nDescription: {description}\n\nPreview: {preview}"

    desc_clean = str(description).replace('\n', ' ')[:150] if description else ''

    return {
        "id": skill_name,
        "text": search_text,
        "metadata": {
            "name": skill_name,
            "description": desc_clean,
            "path": filepath
        }
    }

def build_index():
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        print("❌ ERROR: OPENAI_API_KEY environment variable is not set.")
        return

    print("🔄 Initializing ChromaDB with OpenAI Embeddings (text-embedding-3-small)...")
    client = chromadb.PersistentClient(path=DB_DIR)

    emb_fn = embedding_functions.OpenAIEmbeddingFunction(
        api_key=api_key,
        model_name="text-embedding-3-small"
    )

    collection = client.get_or_create_collection(
        name="skills_collection_openai",
        embedding_function=emb_fn
    )

    print(f"📖 Scanning SKILL.md files in {SKILLS_DIR}...")
    skill_files = glob.glob(os.path.join(SKILLS_DIR, "*/SKILL.md"))

    docs = []
    ids = []
    metadatas = []

    for filepath in skill_files:
        parsed = parse_skill_file(filepath)
        docs.append(parsed["text"])
        ids.append(parsed["id"])
        metadatas.append(parsed["metadata"])

    print(f"✅ Found {len(docs)} skills. Starting vector DB ingestion...")

    batch_size = 100
    for i in range(0, len(docs), batch_size):
        end = min(i + batch_size, len(docs))
        print(f"⏳ Processing batch {i+1} to {end}...")
        try:
            collection.upsert(
                documents=docs[i:end],
                ids=ids[i:end],
                metadatas=metadatas[i:end]
            )
        except Exception as e:
            print(f"❌ Error at batch {i+1}-{end}: {str(e)}")

    print("🎉 Finished indexing all skills with OpenAI!")

if __name__ == "__main__":
    build_index()
