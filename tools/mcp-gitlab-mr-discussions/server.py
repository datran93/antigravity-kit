import os
import gitlab
from mcp.server.fastmcp import FastMCP

# Initialize FastMCP Server
mcp = FastMCP("McpGitLabMRDiscussions")

def get_gitlab_client() -> gitlab.Gitlab:
    """Helper to initialize GitLab client from environment variables."""
    token = os.environ.get("GITLAB_PRIVATE_TOKEN")
    url = os.environ.get("GITLAB_URL", "https://gitlab.com")
    if not token:
        raise ValueError("GITLAB_PRIVATE_TOKEN environment variable is not set. Please set it in your environment or in mcp_config.json env mapping.")
    return gitlab.Gitlab(url, private_token=token)

@mcp.tool()
def read_mr_discussions(project_id: str, mr_iid: int) -> str:
    """
    Read all discussions (threads) from a specific GitLab Merge Request.

    Args:
        project_id: The ID or URL-encoded path of the project.
        mr_iid: The internal ID of the merge request.
    """
    try:
        gl = get_gitlab_client()
        project = gl.projects.get(project_id)
        mr = project.mergerequests.get(mr_iid)

        # Get all discussions for this MR
        discussions = mr.discussions.list(all=True)

        output = [f"💬 DISCUSSIONS FOR MR #{mr_iid} (Project: {project_id})\n"]

        count = 0
        for disc in discussions:
            # A discussion contains a list of notes
            notes = disc.attributes.get('notes', [])
            if not notes:
                continue

            # Filter out pure system notes (like 'approved this merge request') if we want to focus on discussions
            is_system_only = all(n.get('system', False) for n in notes)
            if is_system_only:
                continue

            count += 1

            # Extract resolution status if available
            status_text = ""
            if 'resolved' in notes[0]:
                status = "✅ RESOLVED" if disc.attributes.get('resolved', False) else "❌ UNRESOLVED"
                status_text = f" | Status: {status}"

            output.append(f"--- Discussion ID: {disc.id}{status_text} ---")

            for note in notes:
                author = note.get('author', {}).get('username', 'unknown')
                body = note.get('body', '')
                created = note.get('created_at', '')
                is_system = "[SYSTEM]" if note.get('system') else ""
                note_id = note.get('id', '')

                output.append(f"[{created}] @{author} {is_system} (Note ID: {note_id}):\n{body}\n")

        if count == 0:
            return "✅ No user discussions found on this Merge Request."

        return "\n".join(output)
    except Exception as e:
        return f"❌ Error reading MR discussions: {str(e)}"

@mcp.tool()
def reply_to_mr_discussion(project_id: str, mr_iid: int, discussion_id: str, body: str) -> str:
    """
    Reply to an existing discussion thread on a GitLab Merge Request.

    Args:
        project_id: The ID or URL-encoded path of the project.
        mr_iid: The internal ID of the merge request.
        discussion_id: The ID of the discussion thread to reply to.
        body: The text content of your reply.
    """
    try:
        gl = get_gitlab_client()
        project = gl.projects.get(project_id)
        mr = project.mergerequests.get(mr_iid)

        # We need to call create note on the discussion
        # The python-gitlab way is discussion.notes.create({'body': body})
        try:
            discussion = mr.discussions.get(discussion_id)
            discussion.notes.create({'body': body})
            return f"✅ Successfully replied to discussion '{discussion_id}' on MR #{mr_iid}."
        except gitlab.exceptions.GitlabGetError:
            return f"❌ Discussion '{discussion_id}' not found."

    except Exception as e:
        return f"❌ Error replying to discussion: {str(e)}"

@mcp.tool()
def resolve_mr_discussion(project_id: str, mr_iid: int, discussion_id: str, resolve: bool = True) -> str:
    """
    Resolve or unresolve a discussion thread on a GitLab Merge Request.

    Args:
        project_id: The ID or URL-encoded path of the project.
        mr_iid: The internal ID of the merge request.
        discussion_id: The ID of the discussion thread.
        resolve: True to resolve, False to unresolve.
    """
    try:
        gl = get_gitlab_client()
        project = gl.projects.get(project_id)
        mr = project.mergerequests.get(mr_iid)

        try:
            discussion = mr.discussions.get(discussion_id)
            # In python-gitlab, we use the put method to update resolved status
            # Actually you can modify attributes and call save()

            # Check if discussion is resolvable
            notes = discussion.attributes.get('notes', [])
            if not notes or 'resolvable' not in notes[0] or not notes[0]['resolvable']:
                return f"❌ Discussion '{discussion_id}' is not resolvable."

            discussion.resolved = resolve
            discussion.save()

            action = "resolved" if resolve else "unresolved"
            return f"✅ Successfully {action} discussion '{discussion_id}' on MR #{mr_iid}."

        except gitlab.exceptions.GitlabGetError:
            return f"❌ Discussion '{discussion_id}' not found."

    except Exception as e:
        return f"❌ Error resolving discussion: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
