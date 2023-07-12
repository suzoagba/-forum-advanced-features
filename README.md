# forum-image-upload

This project consists in creating a web forum that allows :

- communication between users (posts, comments, likes/dislikes);
- associating categories to posts (for logged-in users when creating a new post);
- liking and disliking posts and comments (logged-in users);
- filtering posts (logged-in users).

### Storing the Data

In order to store the data in this forum (like users, posts, comments, etc.) the database library 
[SQLite](https://www.sqlite.org/index.html) is used.

SELECT, CREATE and INSERT queries are used.

### Authentication

The client is able to register as a new user on the forum, by inputting their credentials. 
A login session is created to access the forum and be able to add posts and comments.

Cookies are used to allow each user to have only one opened session. Each of these sessions contain an 
expiration date (24h). It is up to you to decide how long the cookie stays "alive". UUID is used as a session ID.

Instructions for user registration:
- An email is required:
  - When the email is already taken, an error response is returned;
- Username is required:
  - When the username is already taken, an error response is returned;
- Password is required:
  - The password is encrypted when stored.

### Communication

In order for users to communicate between each other, they are able to create posts and comments.

- Only registered users are able to create posts and comments;
- When registered users are creating a post they can associate one or more categories (tags) to it;
- The implementation and choice of the categories (tags) was up to the developers;
- The posts and comments are visible to all users (registered or not);
- Non-registered users are only able to see posts and comments.

### Likes and Dislikes

Only registered users are able to like or dislike posts and comments.

The number of likes and dislikes are visible by all users (registered or not).

### Filter

A filter mechanism has been implemented, that will allow users to filter the displayed posts by:

- categories (tags);
- created posts;
- liked posts.

The last two are only available for registered users and must refer to the logged-in user.

### Image Upload

In forum image upload, registered users have the possibility to create a post containing an image as well as text.

- When viewing the post, users and guests can see the image associated to it.
In this project JPG, JPEG, PNG and GIF types are handled.

The max size of the images to load is 20 mb. If there is an attempt to load an image greater than 20 mb, 
an error message will the user that the image is too big.

### Authentication

The goal of this project was to implement new ways of authentication. You are able to register and to login 
using Google and GitHub authentication tools.

To use the new ways of authentication, register an OAuth app at 
[Google](https://developers.google.com/identity/protocols/oauth2/web-server) and 
[GitHub](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app).
Redirect URI-s are:
- Google: `http://localhost:8080/oauth/google`
- GitHub: `http://localhost:8080/oauth/github`

<!-- Input your `client ID`-s and `secrets` to `oauth/clientInfo.go`. -->

If you log in via Google or GitHub and a user is already registered by the same email, the previously registered
username will be used. Otherwise, the username will be based on your email.

### Advanced Features

Users are notified (page 'Notifications') when their posts are:
- liked/disliked (also comments)
- commented

At page 'My Activity' the user can track their own activity:
- Shows the user created posts
- Shows where the user left a like or a dislike
- Shows where and what the user has been commenting. For this, the comment will have to be shown, as well as the post commented

The user can Edit/Remove their posts and comments.

### Moderation

The forum moderation is based on a moderation system. It presents a moderator that, depending on the access level of a user, 
approves posted messages before they become publicly visible.

The filtering can be done depending on the categories of the post being sorted by irrelevant, obscene, illegal or insulting.

There are a total of 4 types of users:
- Guests
  - These are unregistered-users that can neither post, comment, like nor dislike a post. 
  They only have the permission to see those posts, comments, likes or dislikes.
- Users
  - These are the users that will be able to create, comment, like or dislike posts.
- Moderators (user: `user`, password: `psw`)
  - Moderators, as explained above, are users that have a granted access to special functions:
    - They are able to monitor the content in the forum by deleting or reporting post to the admin.
  - To create a moderator the user should request an admin for that role.
- Administrators (user: `admin`, password: `psw`)
  - Users that manage the technical details required for running the forum. This user is able to:
    - Promote or demote a normal user to, or from a moderator user.
    - Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report.
    - Delete posts and comments
    - Manage the categories, by being able to create and delete them.

### Docker

For the forum project Docker is used.

How to:
- Build the Docker image by running the following command: `./docker/build.sh`.
- Once the image is built, you can run a container based on the image using the following command: `./docker/run.sh`.
- The container will start, and your Go application will be accessible at 
[`http://localhost:8080`](http://localhost:8080) in your web browser.
- To stop and remove the image, run the following command: `./docker/stop.sh`.

Make sure you have Docker installed and running on your machine before building and running the Docker image.

### Allowed Packages

- All standard Go packages are allowed;
- sqlite3;
- bcrypt;
- UUID;

No frontend libraries or frameworks like React, Angular, Vue etc. have been used.

### Audit

Questions can be found [here]:
- [Advanced Features](https://github.com/01-edu/public/blob/master/subjects/forum/advanced-features/audit.md).
- [Moderation](https://github.com/01-edu/public/blob/master/subjects/forum/moderation/audit.md).

## Developers
- Willem Kuningas / *thinkpad*
- Samuel Uzoagba / *suzoagba*