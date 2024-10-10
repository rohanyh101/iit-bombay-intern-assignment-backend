# Library Management System

## Description

The Library Management System is a web application designed to streamline the management of library resources, including books, members, and transactions. This system enables librarians to efficiently manage book inventory, handle member registrations, and track borrowing and returning of books. With a user-friendly interface, both librarians and members can interact seamlessly with the system, ensuring a smooth and organized library experience.

## Database Diagram

![Screenshot (140)](https://github.com/user-attachments/assets/f1066531-3d21-4483-9efb-2d551d636752)


## Features

- **User Authentication**: Secure login for both librarians and members.
- **Book Management**: Add, update, delete, and view book details.
- **Member Management**: Register new members, update details, and manage memberships.
- **Borrowing System**: Track book borrowings and returns, including due dates.
- **Role-Based Access**: Different functionalities based on user roles (Librarian/Members).
- **Responsive Design**: Accessible from various devices and screen sizes.

## Technologies Used

- **Frontend**: HTML, CSS, JavaScript
- **Backend**: Golang (Gin Framework)
- **Database**: MongoDB for data storage
- **Authentication**: JWT for secure user authentication

## Getting Started

To run this project locally, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/rohanyh101/iit-bombay-intern-assignment-backend.git
   cd library-management-system
   ```
   
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start the backend server:
   ```bash
   go run main.go
   ```

### System logs (GIN),
   ```bash
   Connected to MongoDB!
Connected to MongoDB!
Connected to MongoDB!
Connected to MongoDB!
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /users/signup             --> github.com/roh4nyh/iit_bombay/routes.AuthRoutes.UserSignUp.func1 (2 handlers)
[GIN-debug] POST   /users/login              --> github.com/roh4nyh/iit_bombay/routes.AuthRoutes.UserLogIn.func2 (2 handlers)
[GIN-debug] GET    /health                   --> main.main.func1 (2 handlers)
[GIN-debug] POST   /librarian/books          --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.AddBook.func3 (4 handlers)
[GIN-debug] GET    /librarian/books          --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetBooks.func4 (4 handlers)
[GIN-debug] GET    /librarian/books/:isbn    --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetBook.func5 (4 handlers)
[GIN-debug] PUT    /librarian/books/:isbn    --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.UpdateBook.func6 (4 handlers)
[GIN-debug] DELETE /librarian/books/:isbn    --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.DeleteBook.func7 (4 handlers)
[GIN-debug] GET    /librarian/users          --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetUsers.func8 (4 handlers)
[GIN-debug] POST   /librarian/users          --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.AddUser.func9 (4 handlers)
[GIN-debug] GET    /librarian/users/:user_id --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetUser.func10 (4 handlers)
[GIN-debug] PUT    /librarian/users/:user_id --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.UpdateUser.func11 (4 handlers)
[GIN-debug] DELETE /librarian/users/:user_id --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.DeActivateUser.func12 (4 handlers)
[GIN-debug] DELETE /librarian/users/:user_id/force --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.DeleteUser.func13 (4 handlers)
[GIN-debug] GET    /librarian/users/active   --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetActiveUsers.func14 (4 handlers)
[GIN-debug] GET    /librarian/users/deleted  --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetNonActiveUsers.func15 (4 handlers)
[GIN-debug] GET    /librarian/users/:user_id/history --> github.com/roh4nyh/iit_bombay/routes.LibrarianRoutes.GetTransactionHistory.func16 (4 handlers)
[GIN-debug] GET    /member/books             --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.GetBooks.func3 (4 handlers)
[GIN-debug] GET    /member/books/:isbn       --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.GetBook.func4 (4 handlers)
[GIN-debug] POST   /member/books/borrow/:isbn --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.BorrowBook.func5 (4 handlers)
[GIN-debug] PUT    /member/books/return/:isbn --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.ReturnBook.func6 (4 handlers)
[GIN-debug] GET    /member/books/borrowed    --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.BorrowedBooks.func7 (4 handlers)
[GIN-debug] DELETE /member/account           --> github.com/roh4nyh/iit_bombay/routes.MemberRoutes.DeActivateMember.func8 (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
   ```

## Example CRUD operations:
## HEALTH CHECK
1.  check the health of API => `GET    /health`
   ```bash
   #request
   curl --location --request GET 'localhost:8080/health' \
   --header 'Content-Type: application/json'

   #response
   {"success":"iit bombay server is up and running..."}
  ```

## AUTHENTICATION / AUTHORIZATION ROUTES

1. **Sign up => `POST   /users/signup`**
 ```bash
   #request
   curl --location --request POST 'http://localhost:8080/users/signup' \
   --header 'Content-Type: application/json' \
   --data-raw '{ "username": "charan", "password": "****", "role": "LIBRARIAN" }'

   #response
   {"InsertedID":"6707ad2a047fb29cf8d72c8c"}
  ```

2. **Login => `POST   /users/login`**
 ```bash
   #request
   curl --location --request POST 'http://localhost:8080/users/login' \
   --header 'Content-Type: application/json' \
   --data-raw '{ "username": "charan", "password": "1212" }'

   #response
   {
  "id": "6707ad2a047fb29cf8d72c8c",
  "username": "charan",
  "password": "$2a$15$meupWOA8h7UUyntYo/Xujem1iJsrJ0p4Hzo7ui0hSp7ndlLAQF6c.",
  "role": "LIBRARIAN",
  "is_active": false,
  "token": "token",
  "created_at": "2024-10-10T10:32:10Z",
  "updated_at": "2024-10-10T10:39:30Z",
  "user_id": "6707ad2a047fb29cf8d72c8c"
}
  ```

## LIBRARIAN ROUTES

1. **get all books => `GET    /librarian/books`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/books' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer token'

  #response
  [
  {
    "id": "67053e06bf3020c63f51ffe7",
    "isbn": "978-0062315007",
    "title": "The Alchemist",
    "author": "Paulo Coelho",
    "status": "AVAILABLE",
    "qty": 1,
    "created_at": "2024-10-08T14:13:26.679Z",
    "updated_at": "2024-10-08T14:13:26.679Z"
  }
]
```

2. **get a single books => `GET    /librarian/books/:isbn`**
```bash
  #request
  curl --location --request GET 'localhost:8080/librarian/books/978-0062323421' `
  --header 'Content-Type: application/json' `
  --header 'Authorization: Bearer <token>'

  #response
 {
  "id": "6705fda074919c7ad793d97c",
  "isbn": "978-0062323421",
  "title": "the monk who sold his ferrari",
  "author": "Robin Sharma",
  "status": "AVAILABLE",
  "qty": 1,
  "created_at": "2024-10-09T03:50:56.337Z",
  "updated_at": "2024-10-09T03:50:56.337Z"
}
```

3. **insert a book => `POST   /librarian/books`**
```bash
  #request
  curl --location --request POST 'http://localhost:8080/librarian/books' \
 --header 'Content-Type: application/json' \
 --data-raw '{ "title": "the almanic of naval ravikant", "isbn": "345-0062535002", "author": "Naval Ravikant", "status": "AVAILABLE", "qty": 1 }' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "book added successfully"
}
```

3. **update book => `PUT    /librarian/books/:isbn`**
```bash
  #request
  curl --location --request PUT 'http://localhost:8080/librarian/books/345-0062535002' \
 --header 'Content-Type: application/json' \
 --data-raw '{ "title": "The almanic of naval ravikant" }' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "book updated successfully"
}
```

4. **delete book => `DELETE /librarian/books/:isbn`**
```bash
  #request
  curl --location --request DELETE 'http://localhost:8080/librarian/books/345-0062535002' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "book deleted successfully"
}
```

5. **get all users => `GET    /librarian/users`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/users' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
[
  {
    "id": "67044140fa1b0a3bc72192f4",
    "username": "manish k",
    "password": "$2a$15$DE5zXiwB8NYAsD71oeGVwepjCDzl8EIXE/LyDpMdqQo0CPpFD4VPO",
    "role": "MEMBER",
    "is_active": false,
    "token": "<token>",
    "created_at": "2024-10-07T20:14:56.462Z",
    "updated_at": "2024-10-08T22:27:38.183Z",
    "user_id": "67044140fa1b0a3bc72192f4"
  },
  {
    "id": "6704ef4cdc19cd768dcedd51",
    "username": "manish",
    "password": "$2a$15$EKRqYXNI.ZqSy5fIokBXf.wSPd75MHRVA9QGgJ6zfLNmcUAjXLlAi",
    "role": "MEMBER",
    "is_active": false,
    "token": "<token>",
    "created_at": "2024-10-08T08:37:32.863Z",
    "updated_at": "2024-10-08T21:55:40Z",
    "user_id": "6704ef4cdc19cd768dcedd51"
  },
]
```

6. **get a single user => `GET    /librarian/users/:user_id`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/users/67044140fa1b0a3bc72192f4' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "id": "6704ef4cdc19cd768dcedd51",
  "username": "manish",
  "password": "$2a$15$EKRqYXNI.ZqSy5fIokBXf.wSPd75MHRVA9QGgJ6zfLNmcUAjXLlAi",
  "role": "MEMBER",
  "is_active": false,
  "token": "<token>",
  "created_at": "2024-10-08T08:37:32.863Z",
  "updated_at": "2024-10-08T21:55:40Z",
  "user_id": "6704ef4cdc19cd768dcedd51"
}
```
7. **add User => `POST   /librarian/users`**
```bash
  #request
  curl --location --request POST 'http://localhost:8080/librarian/users' \
 --header 'Content-Type: application/json' \
 --data-raw '{ "username": "manish", "password": "1212", "role": "MEMBER" }' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "InsertedID":"6707ad2a047fb29cf8d72c8c"
}
```

8. **update User => `PUT    /librarian/users/:user_id`**
```bash
  #request
  curl --location --request PUT 'http://localhost:8080/librarian/users/6704ef4cdc19cd768dcedd51' \
 --header 'Content-Type: application/json' \
 --data-raw '{ "is_active": true }' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "user updated successfully"
}
```

8. **de-activate User => `DELETE   /librarian/users/:user_id`**
```bash
  #request
  curl --location --request DELETE 'http://localhost:8080/librarian/users/67044140fa1b0a3bc72192f4' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "user de-activated successfully"
}
```

9. **force delete User => `DELETE   /librarian/users/:user_id/force`**
```bash
  #request
  curl --location --request DELETE 'http://localhost:8080/librarian/users/67044140fa1b0a3bc72192f4/force' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "user deleted successfully"
}
```

10. **get all active users => `GET    /librarian/users/active`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/users/active' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
[]
```

11. **force delete User => `GET    /librarian/users/deleted`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/users/deleted' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
[
  {
    "id": "67044140fa1b0a3bc72192f4",
    "username": "manish k",
    "password": "$2a$15$DE5zXiwB8NYAsD71oeGVwepjCDzl8EIXE/LyDpMdqQo0CPpFD4VPO",
    "role": "MEMBER",
    "is_active": false,
    "token": "<token>",
    "created_at": "2024-10-07T20:14:56.462Z",
    "updated_at": "2024-10-08T22:27:38.183Z",
    "user_id": "67044140fa1b0a3bc72192f4"
  },
  {
    "id": "6704ef4cdc19cd768dcedd51",
    "username": "manish",
    "password": "$2a$15$EKRqYXNI.ZqSy5fIokBXf.wSPd75MHRVA9QGgJ6zfLNmcUAjXLlAi",
    "role": "MEMBER",
    "is_active": false,
    "token": "<token>",
    "created_at": "2024-10-08T08:37:32.863Z",
    "updated_at": "2024-10-08T21:55:40Z",
    "user_id": "6704ef4cdc19cd768dcedd51"
  },
]
```

12. **get transactions of a single user => `GET    /librarian/users/:user_id/history`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/librarian/users/6704f441a734f8fa83d37008/borrowed' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
[
  {
    "id": "6705b1e065b3e400ac9e9aa4",
    "user_id": "6704f441a734f8fa83d37008",
    "book_id": "67051ed789ae4508c45c4f20",
    "borrowed_at": "2024-10-08T22:27:44.327Z",
    "returned_at": "2024-10-08T23:07:52.193Z",
    "status": "RETURNED"
  },
  {
    "id": "6705bb22f781daa0c372e223",
    "user_id": "6704f441a734f8fa83d37008",
    "book_id": "67051ed789ae4508c45c4f20",
    "borrowed_at": "2024-10-08T23:07:14.124Z",
    "returned_at": "2024-10-08T23:12:15.628Z",
    "status": "RETURNED"
  },
]
```

## MEMBER ROUTES

1. **force delete User => `DELETE   /librarian/users/:user_id/force`**
```bash
  #request
  curl --location --request GET 'http://localhost:8080/member/books' \
 --header 'Content-Type: application/json' \
 --header 'Authorization: Bearer <token>'

  #response
{
  "message": "user deleted successfully"
}
```
