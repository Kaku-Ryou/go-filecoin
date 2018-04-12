# Contributor guidelines

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to the Filecoin Project. Feel free to propose changes to this document in a pull request.

### Table Of Contents
* [Project Management](#project-management)
* [Conventions and Style](#conventions-and-style)
* [Error Handling](#error-handling)
* [The Spec](#the-spec)
* [Platform Considerations](#platform-considerations)
* [Testing Philosophy](#testing-philosophy)
* [What is the bar for inclusion in master?](#what-is-the-bar-for-inclusion-in-master)
* [Pull Requests](#pull-requests)
* [Gotchas](#gotchas)

## Project Management

* We use a modified Agile scrum method. Details are in [Agile and Scrum at Protocol Labs - Filecoin](https://docs.google.com/document/d/16PQVR5vTnbt6DtgxVIEdjc4n0TrqWEtOiWZLItQXC74/edit?ts=5ab28811#).
* We track work items via the [Project Board](https://app.zenhub.com/workspace/o/filecoin-project/go-filecoin/boards?repos=113219518). We use Zenhub because it provides agile project management tools, but integrates directly with Github.

| Stage | What is it for? |
|---|---|
| Backlog | New issues |
| Ready | Issues move here after they have good descriptions, acceptance criteria, and initial estimates. |
| This Sprint | Selected by team for current sprint. Looking for work? Start here. |
| In Progress | When you start work on an issue, move here and assign yourself. |
| Awaiting Review | For PRs and issues. Move here manually when you create a PR. |
| Blocked | If you are working on something but get stuck because of external factors. |
| Done | :tada:|

- The first and last stages are kept up-to-date automatically. 
   - New issues created in `go-filecoin` show up automatically in `Backlog`
   - When a PR is merged, the referenced issues are closed and moved to `Done`.
- All other stages are updated manually.
- All work you are doing should be tracked through an issue. If an issue doesn't exist, create one.

## Conventions and Style

There are always exceptions but generally:
 * Accessors: `Foo(), SetFoo()`
 * Predicates: `isFoo()`
 * Variable names
   * Channels: `fooCh`
   * Cids: `fooCid`
 * Test-only methods: use a parallel type
 * Logging
   * `Warning`: noteworthy but not completely unexpected
   * `Error`: a truly unexpected condition that should not happen in Real Life and that a dev should go look at
 * Protocol messages are nouns (eg, `DealQuery`, `DealResponse`) and their handlers are verbs (eg, `QueryDeal`)

## Error Handling

The overarching concern is safety of the user: do not give the user an
incorrect view of the world.

* If an error is likely a function of an input, discard the input.
* If an error could be transient, attempt to continue making progress.
* If an error appears to be permanent, panic.
* If we find we have inconsistent internal state, panic.

## The Spec

The spec must be in sync with the code. Updates to the spec should
happen before the code is merged into master. Significant changes to
the spec requires discussion and buy-in from {@nicola, @whyrusleeping,
@dignifiedquire}.

At present (Q1'18) we give ourselves some leeway with spec/code update
lag: it's ok for it to be out of sync for a short period of time while
exploring the right thing to do in the code. The important thing at
this early stage is "the spec is updated" rather than "the spec is
updated in precise lockstep with the code". That will be the policy at
some point in the future.

## Platform Considerations

Prefer to derive patterns from more established, closely related
platforms than to derive them from first principles.  For example for
things like configuration, commands, persistence, and the like we draw
heavily on patterns in IPFS. Similarly we draw on patterns in Ethereum
for message processing. Features should be informed by an awareness of
related platforms.

## Testing Philosophy
* All code must be unit tested and should hit our target coverage rate (80%). 
* We prefer to test the output/contracts, not the individual lines of code (which we expect to change significantly during early work).
* Daemon tests (integration tests that run a node and send it commands):
  * Daemon tests are not a substitute for unit tests: the foo command implementation should be unit tested in the `foo_test.go` file
  * Daemon tests should test integration, not comprehensive functionality
  * Daemon tests should validate that their responses conform to a JSON schema
  * Daemon tests for foo go into `foo_daemon_test.go` (so `foo.go` should have *both* `foo_test.go` and `foo_daemon_test.go`)


## What is the bar for inclusion in master?

Presently (Q1'18) the minimum bar is:
 * Must be unittested (80% target, some slippage ok); if PR is a command, it must be command tested
 * Code review:
   * Far-reaching/foundational changes 3/3 of {@phritz, @dignifiedquire, @whyrusleeping}; otherwise 2/3
   * If you expect a response to a CR comment make it clear that you do
 * Must lint (`go run ./build/*.go lint`)
 * CI must be green
 * Must match a subset of the spec, unless when considering changes to the spec -- see details above
 * Functionality should be _somehow_ user-visible (eg, via command line interface)
 * Should be bite-sized

Present (Q1'18) NON-requirements:
 * Doesn’t have to do all or even much of anything it will ultimately do
 * Doesn't have to be fast

Likely future requirements:
 * Integration tested
 * Respects Juan’s most important API requirements

## Pull Requests

#### Workflow
* When creating a pull request:
  * Move the corresponding issue into `Awaiting Review`
  * Assign reviewers to the pull request
  * Place pull request into `This Sprint`.
* Once something is merged to `master` close the issue and move it to `Done`. (Closing an issue moves it to `Done` automatically)
* In PRs and commit messages, reference the issue(s) you are addressing with [keywords](https://help.github.com/articles/closing-issues-using-keywords/) like `Closes #123`.

#### Style
* Always squash commits.
* Anyone is free to merge but slight preference to let the creator
    merge because they have most context. Creators, if you really care
    about the commit message, squash ahead of time and provide a nice
    commit message because someone might merge for you.

## Gotchas
* Equality
  * Don't use `==` to compare `*types.Block`; use `Block.Equals()`
  * Ditto for `*cid.Cid`
  * For `types.Message` use `types.MsgCidsEqual` (or better submit a PR to add `Message.Equals()`)
  * DO use `==` for `types.Address`, it's just an address

