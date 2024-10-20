# Git activity Importer (Gitlab -> Github)
A tool to migrate your commit activity from GitLab to GitHub. 
This tool fetches commit history from private GitLab repositories you've contributed to and imports it into a specified GitHub repository, recreating your activity chart.
e the activity chart. 

## Usage

### 1. Automatic Imports (Recommended)
This approach will automatically keep your activity up to date. The program is being run daily at midnight UTC.
It imports your latest commits and automatically pushes them to specified repository.

To do that follow these steps:
1. **Fork this repository** to your GitHub account.
2. **Create an empty repository** in your GitHub profile where the commits will be pushed.
3. **Configure repository secrets** in your forked repository:
   - Go to your forked repository settings.
   - Under **Security**, navigate to **Secrets and variables > Actions**.
     ![Repository Secrets Configuration](assets/image.png)
   - Add the following secrets:


        | Secret Name       | Description                                                            |
        | ----------------- | ---------------------------------------------------------------------- |
        | `BASE_URL`        | URL of your GitLab instance (e.g., `https://gitlab.com`)               |
        | `COMMITER_EMAIL`  | Email associated with your GitHub profile                              |
        | `COMMITER_NAME`   | Your full name as it appears on GitHub                                 |
        | `GITLAB_TOKEN`    | GitLab personal access token (read permissions only)                   |
        | `ORIGIN_TOKEN`    | GitHub personal access token (with write permissions for auto-push)    |
        | `ORIGIN_REPO_URL` | HTTPS URL of your GitHub repository (ensure it has a `.git` extension) |

Once these variables are saved in your Repository secrets, your commits will be automatically updated every day.

### 2. Manual Imports
If you prefer to run the importer manually:
1. **Download the latest release** of the tool.
2. Set up the same environment variables on your local machine (e.g., `export BASE_URL=https://gitlab.com`).
3. Run the tool locally whenever you want to sync your activity.

### Important Notes:
- **GitLab permissions:** The tool only requires read access to your GitLab repositories.
- **GitHub permissions:** Your GitHub token must have write access to the destination repository for automatic pushes.

[![Go Report Card](https://github.com/furmanp/gitlab-activity-importer)](https://goreportcard.com/report/github.com/furmanp/gitlab-activity-importer)
