# @ralph-plan-review

Verify that markdown files in a folder for planning a task follow expected patterns.

## Instructions

Given the folder path provided by the user:

1. **Scan the folder** - Read all markdown files inside the folder

2. **Identify affected source directories** - For each changed file path in the plan, determine which source directory it belongs to.

3. **Spawn verification subagents** - Spawn a `ralph-reviewer` subagent for structure verification that will:
   - Read `docs/ralph.md`
   - Read files from <folder-path>
   - Verify that the files follow the structure defined in `docs/ralph.md`
   - Report any structural violations or missing sections

4. **Aggregate results** - Collect all subagent reports and only present a summary of:
   - Any pattern violations found
   - Recommendations for fixes

5. **Ask to apply recommendations** - Ask the user if they would like to apply the recommendations to their plan files.

6. **Suggest ralph loop execution** - After user replying and there is nothing else to be done, suggest the user to run ralph for the feature in a new terminal: `./ralph.sh <feature-name>`. Also suggest tmux command for it `tmux new-session -d -s <feature-name> './ralph.sh <feature-name>'`

### Variables

- `<feature-name>` is the folder name where the feature was created. If it was created in `docs/design/foo` it should be `foo`
