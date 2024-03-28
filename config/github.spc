connection "github" {
  plugin = "github"

  # The GitHub personal access token to authenticate to the GitHub APIs, e.g., `ghp_3b99b12218f63bcd702ad90d345975ef6c62f7d8`.
  # Please see https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token for more information.
  # Can also be set with the GITHUB_TOKEN environment variable.
  # token = "ghp_J1jzniKzVbFJNB34cJPwFPCmKeFakeToken"

  # GitHub Enterprise requires a base_url to be configured to your installation location.
  # Can also be set with the GITHUB_BASE_URL environment variable.
  # base_url = "https://github.example.com"

  # GitHub Apps to authenticate with the GitHub API. This is distinct from personal user authentication. GitHub App authentication typically involves using a private key to create a JSON Web Token (JWT), which is then exchanged for an installation access token.

  # The unique identifier for the GitHub App.
  # Can also be set with the GITHUB_APP_ID environment variable.
  # app_id = 1234232

  # The installation ID for a specific installation of the GitHub App.
  # Can also be set with the GITHUB_INSTALLATION_ID environment variable.
  # installation_id = 7439287

  # The private key of the GitHub App, used for generating JWTs for authentication.
  # Can also be set with the GITHUB_PRIVATE_KEY environment variable.
  # private_key = "/Users/ec2-home/2016-10-19.private-key.pem"

  # The installation access token that enables a GitHub App to make authenticated API requests for the app's installation on an organization or individual account.
  # Please see https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app for more information.
  # Can also be set with the GITHUB_APP_TOKEN environment variable.
  # app_token = "ghp_J1jzniKzVbFJNB34cJPwFPCmKeFakeToken"
}
