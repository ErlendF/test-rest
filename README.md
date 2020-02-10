# Test-Rest

This is a REST API for a very simple message board. It supports and uses a MySQL or PostgreSQL database for persistant storage (set by the -d, --database flag, defaults to mysql).

## Endpoints


### /post
*GET/POST*

- A GET requests retrieves all posts stored in the database and related comments.
- A POST request makes a new post. The body of the POST request need to be in the following format: ```{"content": "This is an example post."}```


### /comment
*POST*

- Used for making new comments on a post. The body of the POST request need to be in the following format: ```{"post": 1, "content": "This is an example comment on post 1."}```
