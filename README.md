# release-co-pilot

A release process helper which is semi-automatic but requires certain steps from the user/pilot

## Introduction

Not all release processes can be automated to go "live" straight from the `master` branch of your repo. For all these release processes we created a **Release co-pilot**.

### Definitions

- Release: contains multiple steps typically capturing a specific version being released (for instance called "Release version 1.0.1")
- Step: steps are executed consecutively (non-parallel)
- Task: a step that can be launched and can be waited for
- Checkpoint: a step that will wait indefinitely until a user "acknowledge" it
- Webhook Checkpoint: a step that will wait indefinitely until triggered by an incoming webhook (some external system/task letting us know when it's done)