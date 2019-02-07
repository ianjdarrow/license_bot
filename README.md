# installation

1. clone repo
2. create `.env` in project root
3. add `GITHUB_USER=<username>` and `GITHUB_TOKEN=<token>` to `.env`
4. `go install`

# usage

`license_bot --help`

currently supported:

- `license_bot -o <org_name> -a contributors` prints a list of all contributors to the org
- `license_bot -o <org_name> -a licenses` does a LICENSE COMPLIANCE CHECK across the org against the MIT/Apache-2 dual license
