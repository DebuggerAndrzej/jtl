<div align="center" width="100%">
    <img src="assets/JTL.png" width="300">
</div>
<h2 align="center">JTL - Jira Tui Logger</h2>

JTL is a simple TUI written to assist me in daily logging in jira.

# Installation
```
 go install github.com/DebuggerAndrzej/jtl@latest
```
Requirements:
- go
- unix system

> default installation path is ~/go/bin so in order to have jtl command available this path has to be added to shell user paths

# Configuration
As of now config file path is hardcoded to `~/.config/jtl.toml` this will change in future (possibly a flag to pass conig path).

Config template:
```
username = ""
password = ""
jiraBaseUrl = ""
additionalIssues = "" # comma separated list of Issue Keys. This is an optional argument
```

# Quick features recap
Shortcuts:
- w - log hours under selected issue
- s - log hours under selected issue's scrum issue (*company/team specific*)
- e - increment selected issue's status
- E - decrement selected issue's status
- r - refresh issues list

On selected issue there is always a [glamour](https://github.com/charmbracelet/glamour) generated markdown preview of issue description.

Program keeps track of time logged in current session.
