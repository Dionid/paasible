# paasible

Cli, UI & code management for `ansible`.

![paasible UI](./ui.png)

# Features

1. Run `ansible-playbook` and store results in `.json` files
1. Query and filter results using web UI
1. Separate `playbooks` into `projects` and keep track of them
1. Expand information about every `host` (deployed `playbooks`, connected `ssh_keys`, etc.)
1. Per playbbok variables validation

# Setup

1. Install `ansible-playbook`
1. `paasible init` or create 2 files in root directory manually:
    1. `paasible.yaml` – defines shared configuration (must be stored in repository)
    1. `paasible.hidden.yaml` – defines hidden configuration (must be stored in `.gitignore`)
1. Configure `paasible.hidden.yaml`
1. Edit yours `.gitignore` and add:
    1. `paasible.hidden.yaml`
    1. `db` (this folder is for local SQLite for UI)

# Commands

## `paasible ansible-playbook` – run ansible-playbook and store results

When you need just to fallback to original `ansible-playbook` command, but want to store results:

1. Run `paasible ansible-playbook playbook.yml -- (any ansible-playbook params)`
1. This will create `run_results` and store `.json` files with result of playbook run
1. You can search through this `.json` files as history or run UI to query them with advanced filters

## `paasible init` – initialization

Creates `paasible.yaml` and `paasible.hidden.yaml` files.

## `paasible ui` – web UI

1. Run `paasible ui`
1. Create first admin user.
1. Go to `http://localhost:PORT/_/` into `run_result` table and search by filters.

## `paasible run` – run paasible performances

1. Add `ssh_keys`, `hosts`, `inventories`, `projects`, `playbooks`, `performances`,
`variable_schemas` and `variables` to your `paasible.yaml` file or any other `.yaml` file,
that included into `paasible.yaml`.
1. Run `paasible performe <performance name or id>`
1. This will:
    1. Validate `performance.targets.ssh_keys` are applicable to `hosts`
    1. Validate `variables` against `variable_schemas`
    1. Create correct `inventory` file based on `hosts`, `ssh_keys` and `variables`
    1. Run `ansible-playbook` with `inventories` and `variables`
    1. Save result `.json`

# Dictionary

1. `User` – this is arbitrary field that you use to understand who ran the playbook,
    stored in `paasible.hidden.yaml`*
1. `Machine` – this is arbitrary field that you use to understand on what machine the playbook
    was has run, stored in `paasible.hidden.yaml`*
1. `Project` – combination of multiple `Playbooks` (can be stored in repository or just
    in local folder)
1. `Performance` – a set of `Playbooks` with `Hosts`, `SSH Keys`, `Inventories`, `Variables`
    and other configuration, needed to run playbooks.

* – you can name it however you want, but keep them unique per user and per machine

# Special thanks

To project called pocketbase that is used for UI and sqlite DB. It is a great project and I am using it for my own projects. You can find it here:

https://pocketbase.io

# Roadmap

## (DONE) Stage 1. Save and query ansible-playbook result

Main goal: develope `paasible ansible-playbook` that saves playbook run into local
`.json` files so they can be stored in repository and sync with SqliteDB so them can
be queried.

## Stage 2. Paasible configuration

Main goal: create `.yaml` configuration for paasible, that describes: `ssh_keys`, `hosts`, `inventories`, `projects`, `playbooks`, `performances`, `variable_schemas` etc.

1. ~~Add `inventories`, `projects`, `playbooks`, `performances` and test with `paasible run`~~
1. ~~Add `ssh_keys`, `hosts` and test with `paasible run`~~
1. ~~Make ability to include other `.yaml` files into `paasible.yaml`~~
1. Add `-c` as config path
1. Add `variable_schemas` and `variables` and test with `paasible run`
1. Add `paasible validate` command to validate `paasible.yaml` configuration
1. Add `expand <entity_name>` command to show relationships between entities

## Stage 3. UI

MG: create UI for Paasible to query and edit all entities, play performances,
edit playbooks, projects, etc.

1. Add authentication based on hidde
1. ...

## Stage 4. Remote performe

MG: edit playbooks locally, but run them on remote machines, so you can use Paasible as a remote playbook runner.

1. ...

## Stage 5. Crons

MG: create crons for performances, so you can run them periodically.

1. ...

## Stage 6. Notifications

MG: send notifications about performance and cron results, errors, etc.

1. ...

## Stage X. Terraform support

MG: Add ability to run `terraform` commands and store results in `.json` files, so you can query them in the same way as `ansible-playbook`.

1. ...

## Aidtional

1. Auto install `ansible-playbook` if it is not installed
1. ...
