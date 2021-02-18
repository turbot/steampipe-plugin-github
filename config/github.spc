connection "github" {
  plugin    = "github"   

  # The Github plugin uses a personal access token to authenticate 
  # to the Github APIs  (it looks like `3b99b12218f63bcd702ad90d345975ef6c62f7d8`).
  # You must create a Personal Access Token
  # (https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) 
  # and assign the following scopes:
  #  - `repo` (all)
  #  - `read:org`
  #  - `read:user`
  #  - `user:email`
  #token  = "YOUR_TOKEN_HERE"           
}