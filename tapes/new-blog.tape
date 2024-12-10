
Require echo

Set Shell "bash"
Set FontSize 24
Set Width 1920
Set Height 1080
Set Shell zsh

Set TypingSpeed 1ms

Hide 
Type "dir=`mktemp -p /tmp -d markatanew-XXXX`"
Enter
Type pushd "$dir"
Enter
Type clear
Enter
Show

Set TypingSpeed 25ms

Sleep 2s
Type "pipx run markata new blog myblog"
Enter
Sleep 3.5s
Type "my blog"
Sleep 500ms
Enter
Sleep 1s
Backspace 11
Type "myblog.com"
Enter
Sleep 2.5s
Type "My Blog"
Enter
Sleep 1s
Type "my"
Sleep 1s
Backspace 2
Sleep 1s
Type "hello@myblog.com"
Enter
Sleep 1s
Type "myblog"
Sleep 500ms
Enter
Sleep 2.5s
Backspace 44
Sleep 500ms
Type "This is myblog"
Sleep 1s
Enter
Sleep 1.5s
Enter
Sleep 3.5s
