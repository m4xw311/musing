Qwen3 Code:
Implemented the outline of the application based on the information provided in REQUIREMENTS.md
 - Added an autocompletion script generation option without being asked
   - Necesssity of this would need to be evaulated
   - Probably picked from some example of a commandline application implemented in go
Has some issues with tool use.
 - Fails to give tool argument but does right after it errors.
 - Stops too soon and passes back to user after tool use. Proceeds when nudged.
 - Tool call with improper json
Makes minor mistakes in go import block
Using TODO comments in file to steer agent behavior works but
  - Some comments are being ignored as the comment gets removed by the agent when it implements the code
    - Look into tool for partial file replace rather than full file overwrite
    - Or Do a pre-setp where the TODO comments are turned into a TODO list in separate md file
Correctly pushed back when my ask was not exactly typical.
  - I asked if the RSS/Atom feeds should be having only new posts. Initially a change to only have last 30 days was made but then the change was reverted and I was told it was not a good idea.
