de.createUser({
    user: "dockerFood",
    pwd: "dockerFood",
    roles: [
        {
            role: "readWrite",
            db: "food"
        }
    ]
})