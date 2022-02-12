# Go-Caching-Server
POC Caching Server written in Go and uses MonogDB as backend database and Redis as a caching database.

**Built on Twitch**: https://www.twitch.tv/notdankenoughq

---

![Architecture Overview](https://github.com/M-Faheem-Khan/Go-Caching-Server/blob/main/images/overview.PNG)


The server exposes a single HTTP restful route on port `8001`, `/user/{id}`. When ever a `GET` request is made to the server first the server checks if a user data exists in cache (REDIS) where `id` is used for look up. If the user details are found in the cache the user details are returned to the user. Otherwise, the server attempts to fetch the user details from the database (MongoDB) and inserts into cache (REDIS) and returns the user details as JSON. The user data is stored in cache for 1 hour after which it is expired(removed from cache).

```JSON
// Sample User Object
{
    "id":1000,
    "firstName":"Eli",
    "lastName":"Treleaven",
    "email":"etreleavenrr@slate.com",
    "ipAddress":"99.150.232.20",
    "employer":"Riffwire"
}
```

![Cache_Flow](https://github.com/M-Faheem-Khan/Go-Caching-Server/blob/main/images/cache_details.PNG)

<!-- EOF -->
