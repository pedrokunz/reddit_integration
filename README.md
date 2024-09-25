# canoe_reddit_integration
Canoe Integrations Tech Assessment

## Overview

A critical function of Canoe’s business relies on programmatically retrieving documents/data from external sources, i.e., APIs and Web Portals on behalf of multiple customers.
These use cases present a number of challenges:

 - Enforcing data access control for multiple customers
 - Reliance on 3rd party sites/services that could change, without warning, or be unavailable
 - Throttling & rate limiting
 - Error handling
 - Ensuring all new documents are downloaded
 - Ensuring no documents are downloaded more than once
 - Scaling as the number of documents that need to be downloaded continues to grow

## Task

We have a product that collects reddit.com posts and categorizes them. Depending on information about the post found in it’s metadata available through an API, we route it to a particular internal team. The task:

 - Since https://reddit.com can return any of their subreddit pages as JSON by appending .json to the URL, please use this as the “API” to integrate with. E.g. the computerscience subreddit page can be accessed as JSON like this: https://www.reddit.com/r/computerscience.json
 - Have the integration you build take a list of subreddits (e.g. computerscience in the above URL), connect to each, and save the title and author for each post listed on the page.
 - Each run of the app should save the list of posts found for each subreddit with their title and author , along with the subreddit URL and date of the run.
 - Certain subreddits have far more posts than others, for example https://www.reddit.com/r/pics/. How would your design account for keeping the local list of categorized posts up-to-date with these much larger volume subreddits that change a lot?
 - Use a design pattern that will allow your integration to be easily extensible to support multiple 3rd party APIs in the future. E.g. imagine we wanted to also pull titles and authors from the LinkedIn API
 - Use a SQL database of your choice to maintain the saved data for each run.
