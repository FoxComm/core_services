# Integration of router and social analytics 

## Current integration design

1. Router receives a request to `social_analytics` 
1. Router creates `site_activity` 
1. If site activity is created request is sent to `social_analytics`, result is not checked

## Improved version

1. Router receives a request to `social_analytics` 
1. Router creates `site_activity`
1. If site activity is created request is sent to `social_analytics`
1. Response is checked, if response is error or timeout request payload is saved to Mongo
1. Goroutine checks Mongo. If there are failed requests, it tries to send them again
1. When request is sent, goroutine deletes it from Mongo





