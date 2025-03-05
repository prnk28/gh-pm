#!/bin/bash

set -e

BREW_INSTALLED=$(command -v brew)
GO_INSTALLED=$(command -v go)

function checkInstall() {
  executable=$1
  if ! [ -x "$(command -v $executable)" ]; then
    echo "$executable not installed. Please install $executable and try again."
    return 1
  fi
  return 0
}

# Check if required tools are installed
GH_INSTALLED=$(checkInstall gh)
GUM_INSTALLED=$(checkInstall gum)
FZF_INSTALLED=$(checkInstall fzf)
TASKFILE_INSTALLED=$(checkInstall task)

# Define the taskfile content
taskfile=$(cat <<'EOF'
version: '3'
includes:
  dev: 
    taskfile: .taskfiles/Dev.yml
    flatten: true
    dir: .
tasks:
  default:
    cmds:
      - task -l
    silent: true
EOF
)

# Run task with the defined taskfile
echo "$taskfile" | TASK_X_REMOTE_TASKFILES=1 task -t -

