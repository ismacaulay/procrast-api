## procrast-api

An api for managing your procrastination.

### Why

procrast is basically a simple todo application. The main reason is that it gives a small and clear set of requirements which is good for learning a new languages (rust and go).
This is part of a larger system to learn about cli, api, and everything to do with managing and deploying a microservice based system.

### Routes

```
/lists
    GET - Returns all the lists for the user
    POST - Creates a new list for the user

/lists/<id>
    GET - Returns the info for the list
    PATCH - Updates the list info
    DELETE - Deletes the list and all items associated with that list

/lists/<id>/items
    GET - Returns all the items for a list
    POST - Creates a new item in the list

/lists/<id>/items/<id>
    GET - Returns the item information
    PATCH - Updates the item information
    DELETE - Deletes the item
```

### TODO

- Implemented database functions
- Use transactions when interacting with the database when deleting lists
  - Maybe something like a context manager

```
tx := db.Trasaction(func (db Database) {
    db.Delete()
    db.Create()
})

if tx.Commit() != nil {

}
```

- Proper error messages for validation
  - They are currently just empty strings
- Middleware for checking if list exists in item endpoints
  - This will remove the duplication in the handlers
  - Can send the list id in the context?
