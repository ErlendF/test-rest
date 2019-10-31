This is a REST API for a very simple message board. It uses a PostgreSQL database for persistant storage.

The API provides the following endpoints:


 **"/":** Only accepts GET requests. Provides a simple success message with a timestamp of when it was last changed. Used to check the availability of the API and test CI/CD when changing the Golang application.

 **"/post":** Accepts GET and POST requests. A GET requests retrieves all posts stored in the PostgreSQL database and related comments. A POST request makes a new post. The body of the POST request need to be in the following format:
 ```{"content": "This is an example post."}```


 **"/comment":** Only accepts POST requests. Used for making new comments on a post. The body of the POST request need to be in the following format: 
 ```{"post": 1, "content": "This is an example comment on post 1."}```