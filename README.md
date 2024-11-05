# reddit_integration

Integrations Tech Assessment

## Overview

A critical function of the business relies on programmatically retrieving documents/data from external sources,
i.e., APIs and Web Portals on behalf of multiple customers.
These use cases present a number of challenges:

- Enforcing data access control for multiple customers
- Reliance on 3rd party sites/services that could change, without warning, or be unavailable
- Throttling & rate limiting
- Error handling
- Ensuring all new documents are downloaded
- Ensuring no documents are downloaded more than once
- Scaling as the number of documents that need to be downloaded continues to grow

## Task

We have a product that collects reddit.com posts and categorizes them. Depending on information about the post found in
it’s metadata available through an API, we route it to a particular internal team. The task:

- Since https://reddit.com can return any of their subreddit pages as JSON by appending .json to the URL, please use
  this as the “API” to integrate with. E.g. the computerscience subreddit page can be accessed as JSON like
  this: https://www.reddit.com/r/computerscience.json
- Have the integration you build take a list of subreddits (e.g. computerscience in the above URL), connect to each, and
  save the title and author for each post listed on the page.
- Each run of the app should save the list of posts found for each subreddit with their title and author , along with
  the subreddit URL and date of the run.
- Certain subreddits have far more posts than others, for example https://www.reddit.com/r/pics/. How would your design
  account for keeping the local list of categorized posts up-to-date with these much larger volume subreddits that
  change a lot?
- Use a design pattern that will allow your integration to be easily extensible to support multiple 3rd party APIs in
  the future. E.g. imagine we wanted to also pull titles and authors from the LinkedIn API
- Use a SQL database of your choice to maintain the saved data for each run.

## How to run

1. Clone the repository
2. Create a Reddit app to get the client id and client secret
    - Go to https://www.reddit.com/prefs/apps
    - Click on "Create App" or "Create another app..."
    - Select "script" as the app type
    - Fill in the required fields and click on "Create app"
    - Under the app name, you will find the client id
    - The client secret will be under the "secret" field
      ![image](https://github.com/user-attachments/assets/f68ee8b4-7983-45d8-accb-ba1c48c49894)

3. Create a .env file following the example:
    ```
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_HOST=postgres_db
    DB_PORT=5432
    DB_NAME=reddit_integration
    DB_SSLMODE=disable
    
    JWT_KEY=secret
    
    REDDIT_CLIENT_ID=your_reddit_client_id
    REDDIT_CLIENT_SECRET=your_reddit_client_secret
    REDDIT_USERNAME=your_reddit_username
    REDDIT_PASSWORD=your_reddit_password
    ```
4. Run `docker compose up` in the root directory
    - Two containers will be created, one for the app and another for the database
    - The app will be available at `localhost:8080`
       ```shell
      go_app       | [GIN-debug] Listening and serving HTTP on :8080
      ```
5. Run the following CURL command to get the access token:
    ```shell
    curl --location 'localhost:8080/login' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "customer_id": "afacf9fa-7516-468f-b048-ac4c0562aa3f",
    "email": "user1@test.com",
    "password": "123"
    }'
    ```
6. Run the following CURL command to sync the posts:
    ```shell
   curl --location --request POST 'localhost:8080/v1/reddit/posts/sync?subreddits=computerscience%2Cpics' \
   --header 'Authorization: {{token}}'
    ```
7. You can see the log outputs in the Docker console and the data will be in the database
   ```shell
   go_app       | [GIN] 2024/09/26 - 18:13:22 | 200 |     8.71779ms |      172.19.0.1 | POST     "/login"
   go_app       | [GIN] 2024/09/26 - 18:13:28 | 200 |      52.689µs |      172.19.0.1 | POST     "/v1/reddit/posts/sync?subreddits=computerscience,pics"
   go_app       | 2024/09/26 18:13:30 Fetched 100 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:30 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:32 Fetched 100 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:32 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:33 Fetched 100 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:34 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:35 Fetched 100 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:36 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:36 Fetched 100 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:38 Fetched 11 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:38 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:38 Successfully saved 511 posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:40 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:42 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:44 Fetched 100 posts for subreddit pics
   go_app       | 2024/09/26 18:13:45 Fetched 33 posts for subreddit pics
   go_app       | 2024/09/26 18:13:45 Successfully saved 833 posts for subreddit pics
   go_app       | [GIN] 2024/09/26 - 18:13:51 | 200 |       70.73µs |      172.19.0.1 | POST     "/v1/reddit/posts/sync?subreddits=computerscience,pics"
   go_app       | 2024/09/26 18:13:53 No new posts for subreddit computerscience
   go_app       | 2024/09/26 18:13:53 No new posts for subreddit pics
   ```
## Design

The application is designed to be scalable and extensible. The application is divided into the following layers:

- **API Layer**: This layer is responsible for handling the HTTP requests and responses. It uses the Gin framework to
  handle the requests and responses. The API layer is responsible for validating the input, authenticating the user,
  and calling the service layer to perform the post synchronization.
- **Service Layer**: This layer is responsible for the business logic of the application. It contains the logic for
  synchronizing the posts from Reddit and saving them to the database. The service layer is designed to be extensible
  so that it can support multiple 3rd party APIs in the future.
- **Repository Layer**: This layer is responsible for interacting with the database. It contains the logic for saving and
  retrieving the posts from the database. The repository layer is designed to be extensible so that it can support
  different types of databases.
- **Model Layer**: This layer contains the data models used by the application. It defines the structure of the data
  stored in the database.

## Notes

- Data access control for multiple customers
  > The application is designed to support multiple customers by using JWT
    tokens for authentication. Each user is assigned a unique customer ID, roles, and permissions. The application enforces
    data access control by checking the user's permissions before allowing them to access the data.
- Throttling & rate limiting
  > The application uses a rate limiter middleware to limit the number of requests that can be made to the API. This
    helps prevent abuse and ensures that the application can handle a large number of requests without being
    overwhelmed.
- Error handling
  > The application uses structured error handling to handle errors that occur during the execution of the application.
    Errors are logged and returned to the client with an appropriate HTTP status code and error message. This also helps
    ensure that the application is robust and can recover from errors gracefully.
- Ensuring all new documents are downloaded and no documents are downloaded more than once
  > The application uses a synchronization mechanism to ensure that all new posts are downloaded from Reddit. The
    application keeps track of the last synchronization time for each subreddit and only downloads new posts that have
    been added since the last synchronization.
- Scaling as the number of documents that need to be downloaded continues to grow
  > The application is designed to be scalable by using a microservices architecture. The application can be deployed
    on multiple servers and scaled horizontally to handle a large number of requests. The application is also designed
    to be extensible so that it can support multiple 3rd party APIs in the future.
- Further improvements
  > The application can be further improved by adding support for more 3rd party APIs, adding more features such as
    filtering and sorting the posts, and improving the error handling and logging, improving some specific parts of
    the code like the password encryption and comparison, reading the customer IDs from a repository, and adding a new 
    layer for use cases only, so we can have multiple different entry points for the same use case, for example, 
    we can trigger the sync by a scheduled routine from a crontab or similar, and also improving the async operations to .
    by using channels in the existing goroutines, so we can start storing data as soon they are fetched from the API.   
