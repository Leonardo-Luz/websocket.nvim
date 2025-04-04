## ISSUES

* Users is disconnected if hes in insert mode and receives a msg from another
* on message is recursive (modifies the buffer and trigges textchanged, which changes the buffed)
* write on the same cursor pos as the other will crash
* erase will crash

* startwsclient crashes cuz the cursor is inside one of the lines that are modified (all the lines are lol) (donsnt make sense cuz when i rewrite all lines doesnt crash)
