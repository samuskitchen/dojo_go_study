package user

const(

	// selectAllUserTest is a query that selects all rows in the user table
	selectAllUserTest = "SELECT id, name, surname, username, email, created_at, updated_at FROM users;"

	// selectUserByIdTest is a query that selects a row from the users table based off of the given id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectUserByIdTest = "SELECT id, name, surname, username, email, created_at, updated_at FROM users WHERE id \\= \\$1;"

	// selectUSerByUsernameTest is a query that selects a row from the users table based off of the given username.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectUSerByUsernameTest = "SELECT id, name, surname, username, email, password, created_at, updated_at FROM users WHERE username \\= \\$1;"

	// insertUserTest is a query test that inserts a new row in the user table using the values
	// for insert queries. You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	insertUserTest = "INSERT INTO users \\(name, surname, username, email, password, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) RETURNING id;"

	// updateUserTest is a query that updates a row in the users table based off of id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	updateUserTest = "UPDATE users SET name\\=\\$1, surname\\=\\$2, username\\=\\$3, updated_at\\=\\$4 WHERE id\\=\\$5;"

	// deleteUserTest is a query that deletes a row in the users table given a id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	deleteUserTest = "DELETE FROM users WHERE id\\=\\$1;"
)
