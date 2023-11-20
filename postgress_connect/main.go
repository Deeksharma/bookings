package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	// connect to the
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "host=localhost port=5432 dbname=postgres user=deekshasharma password=") // returns the connection of pools, we can configure the number
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to postgres, err:= %v", err))
	}
	defer conn.Close(ctx)
	log.Println("Connected to database!")

	// test the connection
	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to ping to postgres, err:= %v", err))
	}
	log.Println("Pinged to database!")

	// read all the data
	err = getAllRows(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	// push some data

	query := `insert into users (name, id, age) values ($1, $2, $3)` // can write in multiple lines
	//_, err = conn.Exec(ctx, query, "Jack", 9, "Brown")             // safe practice

	if err != nil {
		log.Fatal(err)
	}

	// someone can perform sql injection if they send name as `;drop users;` - sql injection attack - and that'll cause the whole table to be deleted
	// using it like the above function, attacker cannot perform sql injection because it should be used as args only

	// read the data again
	err = getAllRows(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	// update a row
	stmt := `update users set name=$1 where id = $2`
	_, err = conn.Exec(ctx, stmt, "kylie", 2) // safe practice

	if err != nil {
		log.Fatal(err)
	}

	// get the data again
	err = getAllRows(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	// get row by id
	var name, age string
	var id int
	query = `select name, id, age from users where id= $1`
	row := conn.QueryRow(ctx, query, 1)
	err = row.Scan(&name, &id, &age)

	log.Println("query row retruns: ", id, name, age)

	// delete a row

	query = `delete from users where id = $1` // can write in multiple lines
	_, err = conn.Exec(ctx, query, 2)         // safe practice

	if err != nil {
		log.Fatal(err)
	}

	// get the data again
	err = getAllRows(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(ctx context.Context, conn *pgx.Conn) error {
	rows, err := conn.Query(ctx, "select * from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close() // close connection and rows in defer always otherwise you'll go out of db resources

	var name, age string
	var id int
	for rows.Next() {
		err := rows.Scan(&name, &id, &age) //  here the order matters
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Record is: ", id, name, age)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("error scanning rows", err)
	}
	fmt.Println("-----------------------------")
	return nil
}
