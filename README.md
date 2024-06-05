# DEALLS

This dating apps built using Go with Echo framework and Mongo DB consisting these function:
1. Sign up
2. Sign in
3. View Other Profile
4. Get Premium
5. Like
6. Pass 

HOW TO RUN ON LOCAL (without docker)
  1. Pull the repository
  2. Rename .env.example to .env
  3. insert your Mongo URI and JWT_SECRET (up to your choice)
  4. hit "go get" in terminal
  5. hit "go run main.go" in terminal

HOW TO RUN ON DOCKER
  1. Pull the repository
  2. Rename .env.example to .env
  3. insert your Mongo URI, JWT_SECRET (up to your choice), ClearView to decide how long non premium user can reset their ten per day limit
  4. edit YOUR MONGO URI in docker-compose.yml
  5. hit "docker compose up" in terminal 


You can sign up by providing some data about yourself (beware of existing username, yours should be unique). After successful signup, hit the sign in API using your username and password.
If you're the authorized user, then it will give you the token. Use your token to scroll up other profile which you'll be taking action (Like or Pass). 
You can only view 10 different profiles a day, unless you become premium member by hitting the Get Premium API and provide your token. Have Fun!
