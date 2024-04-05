connection "github" {
  plugin = "github"

  # The following tokens are currently supported:
  # - Personal access tokens, e.g., `ghp_3b99b12218f63bcd702ad90d345975ef6c62f7d8`
  #   https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token for more information.
  # - GitHub application installation access token, e.g., `ghs_UdmjfiKzVbFJNBsaiePwFPCmKeFakeToken`
  #   https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app for more information.
  # Can also be set with the GITHUB_TOKEN environment variable.
  # token = "ghp_3b99b12218f63bcd702ad90d345975ef6c62f7d8"

  # GitHub Enterprise requires a base_url to be configured to your installation location.
  # Can also be set with the GITHUB_BASE_URL environment variable.
  # base_url = "https://github.example.com"

  # GitHub App authentication using a private key to create a JSON Web Token (JWT), which is then exchanged for an installation access token.

  # The unique identifier for the GitHub App.
  # Can also be set with the GITHUB_APP_ID environment variable.
  # app_id = 12345678

  # The installation ID for a specific installation of the GitHub App.
  # Can also be set with the GITHUB_APP_INSTALLATION_ID environment variable.
  # installation_id = 8901234

  # The path to a private key PEM file of the GitHub App, used for generating JWTs for authentication.
  # Can also be set with the GITHUB_APP_PEM_FILE environment variable.
  # private_key = "/Users/myuser/app_private_key.pem"
}
