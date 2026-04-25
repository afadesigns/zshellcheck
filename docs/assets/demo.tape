Output docs/assets/demo.gif

# Theme — kitty's live palette (muted, matches the user's daily terminal)
Set Theme '{ "name": "kitty-zshaaa", "background": "#101010", "foreground": "#999993", "cursor": "#9d9eca", "selection": "#4d4d4d", "black": "#333333", "red": "#8c4665", "green": "#287373", "yellow": "#7c7c99", "blue": "#395573", "magenta": "#5e468c", "cyan": "#31658c", "white": "#899ca1", "brightBlack": "#3d3d3d", "brightRed": "#bf4d80", "brightGreen": "#53a6a6", "brightYellow": "#9e9ecb", "brightBlue": "#477ab3", "brightMagenta": "#7e62b3", "brightCyan": "#6096bf", "brightWhite": "#c0c0c0" }'

Set Shell "zsh"
Set FontFamily "FiraCode Nerd Font Mono"
Set FontSize 18
Set LineHeight 1.0
Set Width 1100
Set Height 900
Set Padding 24
Set TypingSpeed 45ms
Set CursorBlink false

# Hidden setup — load p10k via ~/.zshrc (which sources 90-zshellcheck.zsh, exposing the bare command), cd to demo dir
Hide
Type "source ~/.zshrc 2>/dev/null"
Enter
Sleep 1800ms
Type "cd ~/.config/zsh/zshellcheck/samples"
Enter
Type "clear"
Enter
Show
Sleep 800ms

# 1. Lint (severity warning to keep output focused)
Type "zshellcheck -severity warning deploy.zsh"
Sleep 300ms
Enter
Sleep 4.5s

# 2. Preview the fix (unified diff)
Type "clear"
Enter
Sleep 200ms
Type "zshellcheck -severity warning -diff deploy.zsh"
Sleep 300ms
Enter
Sleep 4.5s

# 3. Apply the fix
Type "clear"
Enter
Sleep 200ms
Type "zshellcheck -fix deploy.zsh"
Sleep 300ms
Enter
Sleep 2.5s

# 4. Re-lint — one warning remains (rm -rf guard needs human judgment)
Type "clear"
Enter
Sleep 200ms
Type "zshellcheck -severity warning deploy.zsh"
Sleep 300ms
Enter
Sleep 4s
