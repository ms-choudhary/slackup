# slackup

A highly queryable backup of messages in slack channels, which are exposed via webserver over REST API.


### Why?

If you've free slack account, there is a limit to the history of messages
saved by default. It's always good to have backup of the messages in that case. This
project stores all messages in elasticsearch, and are highly queryable via REST
API. 

To be added, a cli client and frontend service to display messages.
