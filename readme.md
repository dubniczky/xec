# X-ec

X-ec (or exec) is a terminal utility to simplify writing complex commands in a simpler-to-understand YAML or JSON format.

## Example

```yaml
cmd: ls
args:
- l
- a
stdout: 
stderr:
```

## Commands

- `run [path?]` - run the command from a default location or specified on the CLI
- `show [path?]` - show the final command that would be run with `xec run`
