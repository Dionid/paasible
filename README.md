# Paasible

Cli and UI for `ansible` to run playbooks and store and search results in a local sqlite DB.

[Paasible UI](./ui.png)

# Setup

1. Install `ansible-playbook`
1. Create 2 files in root directory (or run `paasible init`):
    1. `paasible.yaml` – defines project
    1. `paasible.env` – defines user
1. Configure `paasible.env`:
    1. `PAASIBLE_UI_PORT` – port for Paasible UI
    1. `PAASIBLE_USER` – user from whom will playbook be runned
    1. `PAASIBLE_MACHINE` – machine where playbook will be runned
1. Configure `paasible.yaml`:
    1. `project_name` – must be unique project name
    1. `data_folder` – where Paasible will store data
    1. `cli_version` – version of Paasible CLI
1. Edit `.gitignore`:
    1. Add `paasible.env`
    1. Add `data_folder/pb_data` (this folder is for local sqlite DB for UI and must be unqiue per machine)

# Commands

1. `paasible ansible-playbook playbook.yml` – to run playbook
1. `paasible ansible-playbook playbook.yml -- ...` – to add any `ansible-playbook` specific arguments (e.g. `-i` for inventory or `-e` for extra vars, etc.)
1. `paasible serve` – to run web ui

# How it works

1. Go through `Setup` section
1. That you can run `ansible-playbook` as usual, BUT using `paasible` command like `paasible ansible-playbook playbook.yml` in the same folder as `data_folder`
1. This will create a new folder in `data_folder` and store `.json` files with result of playbook run
1. You use this `.json` files as history of runs or you can use UI for that
1. To use UI just do `paasible serve`, it will open admin page in your browser and you must create
first admin user.
1. After that go to `http://localhost:PORT/_/` into `playbook_run_result` table and you will see all playbooks that were runned and can filter by any field

# Dictionary

1. `Paasible cli` – terminal application to
    1. Init paasible confug
    1. Run playbooks
    1. Run paasible UI
1. `User` – this is arbitrary field that you use to understand who runed the playbook*
1. `Machine` – this is arbitrary field that you use to understand on what machine the playbook was runned*
1. `Project` – this is arbitrary field that you use to understand what project the playbook was runned, there can be any playbooks number in one project*

* – you can name it however you want, this is just for you to understand where it was runned

# Roadmap

1. Add ability to customize and run playbooks from UI
1. Add centralized UI for all projects

# Special thanks

To project called pocketbase that is used for UI and sqlite DB. It is a great project and I am using it for my own projects. You can find it here:

https://pocketbase.io
