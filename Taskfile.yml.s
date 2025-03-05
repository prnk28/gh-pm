# https://taskfile.dev

version: '3'
silent: true

includes:
  dev: 
    taskfile: .taskfiles/Dev.yml
    flatten: true
    dir: .

tasks:
  default:
    cmds:
      - task --list-all
    silent: true
