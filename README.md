# Consul-Lock

This is a simple tool to help syncrhnoizing distributed workflows across the team.

## Usage:

```
  help [<command>...]
    Show help.

  lock
    Accquire lock from consul

  status
    Check lock status

  release
    Release the lock, uses .consul_lock_id file

  release-with-id <id>
    Release the lock with explicit session ID argument. Should not be used normally
```
