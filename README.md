## Branch Protection

> Version: v1.0.0

------

## About config syntax

If you need update a branch with two reviewer on a specific branch, you can do below that:

```yaml
- branch: my-branch
  protection:
    required_pull_request_reviews:
      required_approving_review_count: 2
```

But, if you need enforce your branch with a status checks (lint, test, build), ensure that nobody can push or delete commits and ensure if at least one member team review the pr, you can do below that:

```yaml
- branch: dev
  protection:
    required_status_checks:
      strict: true
      contexts:
        - lint
        - test
        - build
    required_pull_request_reviews:
      dismissal_restrictions: null
      dismiss_stale_reviews: true
      require_code_owner_reviews: false
      required_approving_review_count: 1
    required_linear_history: false
    allow_force_pushes: false
    allow_deletions: false
```

It's the same config as that github's api. For more information you can check this [link](https://developer.github.com/v3/repos/branches/#update-branch-protection)

------

## How to use it?

By default `branch-protection` uses a config file over this path `.github/config/branch_protection.yml`. You need add this config file in your repository if you want update yours branches.

> Warning: This action required an `access token` with _full repo access_. You can found more info in the next [link](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)

> Info: Don't forget a checkout your repo before create a branch protection

```yaml
jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Update branch protection
        uses: Hatzelencio/branch-protection@v1.0.0
        env:
          GITHUB_TOKEN: ${{secrets.ACCESS_TOKEN}}
```

This is another example where we define another config path:

> Warning: Ensure that config file lives over your repository

```yaml
jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Update branch protection
        uses: Hatzelencio/branch-protection@v1.0.0
        with:
          path: .github/config/another_config.yml
        env:
          GITHUB_TOKEN: ${{secrets.ACCESS_TOKEN}}
```

But if you need create a strategy to lock/unlock branches you may do next:

```yaml
jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Strategy to lock branches
        uses: Hatzelencio/branch-protection@v1.0.0
        with:
          path: .github/config/lock_branch_config.yml
        env:
          GITHUB_TOKEN: ${{secrets.ACCESS_TOKEN}}
```

```yaml
jobs:
  job-id:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Strategy to unlock branches
        uses: Hatzelencio/branch-protection@v1.0.0
        with:
          path: .github/config/unlock_branch_config.yml
        env:
          GITHUB_TOKEN: ${{secrets.ACCESS_TOKEN}}
```

------

## Inputs

`path` Path where is store the branch protection config file. This file must be a yaml file. By default, it should be on the next path: `.github/config/branch_protection.yml`

`repository` Owner and Repository name separated by `/`. Default env variable is: `GITHUB_REPOSITORY`

------

## Troubleshooting

* Why when I call this action it fails, and it prints in stdout the next message: "... 403 Resource not accessible by integration []"

    * _Because the github token doesn't have the level of permission required. You must grant it the least full repo access._

