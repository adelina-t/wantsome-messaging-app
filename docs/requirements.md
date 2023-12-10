## SERVER

# - improve logging DONE using log levels
- Direct / private messages between users
- List chat users - http ( not enforced, could be ws if so desired )
- Multiple rooms - users can chat on different rooms 
Nice to have:
- Load chat history when connecting *** 
    - will have to preserve state - how? files, db?
- User login - bind connection to username 
    - should we change username? how? 


## CLIENT

# - accept input from stdin DONE for message every time and username only once at start
- quit / disconnect 
- config server connection params: - url etc 
    - options:
        - envvars
        - file config
- NOTE: does not have to be written in go. Web based is perfectly fine.



