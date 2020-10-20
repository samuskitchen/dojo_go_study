package user

const(

	// selectAllUser is a query that selects all rows in the user table
	selectAllUser = "SELECT id, name, surname, username, email, created_at, updated_at FROM users;"

	// selectUserById is a query that selects a row from the users table based off of the given id.
	selectUserById = "SELECT id, name, surname, username, email, created_at, updated_at FROM users WHERE id = $1;"

	// selectUSerByUsername is a query that selects a row from the users table based off of the given username
	selectUSerByUsername = "SELECT id, name, surname, username, email, password, created_at, updated_at FROM users WHERE username = $1;"

	// insertUser is a query that inserts a new row in the user table using the values
	// given in order for name, surname, username, email, password, created_at and updated_at.
	insertUser = "INSERT INTO users (name, surname, username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"

	// updateUser is a query that updates a row in the users table based off of id.
	// The values able to be updated are name, surname, username and updated_at.
	updateUser = "UPDATE users SET name=$1, surname=$2, username=$3, updated_at=$4 WHERE id=$5;"

	// deleteUser is a query that deletes a row in the users table given a id.
	deleteUser = "DELETE FROM users WHERE id=$1;"
)
