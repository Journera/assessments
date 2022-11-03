# Assessments
Assessments is a collection of coding evaluation tests for developer candidates.
The tests are designed to test real coding ability in a real environment.

Your task will be to complete the implementation of the task
assigned to you, and submit a PR to the main project.

## Check out and compile the code
Using your github account, fork this repository, then clone the repo
to your workstation. You may clone it wherever you like, but below is our
example. Be sure to clone your forked repo, instead of this one:

```bash
mkdir dev
cd dev
git clone git@github.com:Journera/assessments.git
cd assessments/go
go build
```

You should see the `assessments` binary in the directory. This will have a functional,
but incorrect and/or incomplete implementation of the various tasks. Your job will
be complete the implementation and insure it runs properly. Unit tests are not
mandatory, but recommended.

## Rate Limiter
The Rate Limiter will attempt to prevent a single sending client from sending too many
messages in a given time frame. The parameters it will be given will be:
- limit - the maximum number of messages allowed per minute (600 = 10/sec)
- reject - if true then messages will be rejected if they exceed the rate, otherwise they will be delayed

### Technical considerations:
- If messages are to be delayed, blocking on the call to `Send()`, to force a slowdown, is allowed.
- Proper logging is advised, see examples from the client code.
- You can assume things like unlimited memory. If there are issues that arise during your implementation
that would require a vast amount of complexity, you can make certain assumptions like this to get the 
project complete in a reasonable amount of time. Although complete solutions are preferred.

### Examples:
- `./assessments ratelimit` this will run with the default parameters
- `./assessments ratelimit -c 2 -m 10 -s 60 -d`
  - Run with 2 clients, each sending 10 messages at 1 msg/sec, with debug logging - good dev testing params