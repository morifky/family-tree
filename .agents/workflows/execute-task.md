---
description: How to automate code generation for a Linear task
---

# Linear Task Execution Workflow

When working on a Linear issue (e.g., MOR-9), follow this workflow to automate branch creation, code implementation, testing, atomic commits, and pull requests.

## 1. Checkout a New Branch

- For the given Linear task, get its `gitBranchName` using the Linear MCP tool or create one manually following the convention `MOR-<ID>-<short-description>`.
- Use the `run_command` tool to create and checkout the new branch:
  ```bash
  git checkout -b <branch_name>
  ```

## 2. Generate Code iteratively

- Implement the requested feature or bug fix.
- Make sure to review the technical specification and project rules before making changes to ensure everything stays compliant.

## 3. Verify the Code (Run Tests / Linters)

- Before committing any code, you MUST ensure that your code compiles and passes tests.
- You should run tests for the changes you just made:

  ```bash
  # For Go backend
  make test || go test ./...

  # For Svelte frontend (if applicable)
  npm run check
  ```

- If the tests/build fail, debug and fix the code before proceeding.

## 4. Atomic Commits

- Follow the **Atomic Commit Principle**: each commit should be a single, logical change. Do not dump all code into a single large commit.
- Use Conventional Commits formatting (e.g., `feat:`, `fix:`, `refactor:`, `test:`, `chore:`).
- Use `run_command` to stage files and create commits:
  ```bash
  git add <specific-files>
  git commit -m "feat(module): description of change"
  ```
- Repeat Steps 2-4 until the task is completely implemented.

## 5. Create Pull Request

- Once the task is fully implemented and all code is committed:
  ```bash
  # Push the branch to the remote repository
  git push -u origin HEAD
  ```
- Create a Pull Request against the main branch. You can use the GitHub MCP server `mcp_github-mcp-server_create_pull_request` (if available) or use the GitHub CLI (`gh`):
  ```bash
  gh pr create --title "MOR-XX: Title of the Linear Task" --body "Resolves MOR-XX."
  ```

## 6. PR Verification

- The CI pipeline on the Pull Request will run unit tests and linters. Review any failing status checks and address them immediately by following steps 2-4 and pushing the fixes to the same branch.
