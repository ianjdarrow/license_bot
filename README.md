# goals

updating open source licenses in large organizations with many repos & contributors is very, very difficult. let's see how much we can automate it.

# vision

your project wants to adopt a new open source license. this script:

- scrubs every repo in your organization, identifying all those that don't currently used the desired license
- assembles a list of every contributor whose consent is required to fully adopt the new license for all existing code
- creates a tracking issue where folks can consent to the license change
- submits a single PR for every single repo to delete existing license files and add the new ones
- tracks progress & identifies people to follow up with for consent for the license change

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
