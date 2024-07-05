package services

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres golang driver // NOTE: uses only init functions and nothing else (_)
)

func InitDB() *sql.DB {
	db := connectDB()
	defer db.Close()
	if err := createTables(db); err != nil {
		log.Fatalf("database initialization failed, [err] %#v", err)
	}
	return db
}

func connectDB() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// open connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Panicf("failed to open db connection")
		// NOTE: panic allows deferred func, log.Fatal call os.Exit doesnot allow defer func run
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Database connection Success!")

	return db // return connection
}

func createTables(db *sql.DB) error {
	startQuery := `CREATE TABLE IF NOT EXISTS `
	tableName := `Users`
	endQuery := `(
  		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
 		github_username VARCHAR(50) UNIQUE,
  		leetcode_username VARCHAR(50) UNIQUE);`

	_, err := db.Exec(startQuery + tableName + endQuery)
	if err != nil {
		log.Println("failed to create table", tableName, "error:", err)
		panic(err)
	}

	log.Println("tables create Succcess!")
	return nil
}

// func InsertUser(name, ip string) error {
// 	db := connectDB()
// 	insertQuery := fmt.Sprintf("INSERT INTO Users(name, ip) VALUES(%s, %s)", name, ip)
// 	_, err := db.Exec(insertQuery)
// 	return err
// }

// // insert post into db, returns inserted postID on success
// func InsertPost(post models.Post) int64 {
// 	db := connectDB()
// 	defer db.Close()

// 	sqlQuery := `INSERT INTO posts (title, body, created_at) VALUES ($1, $2, $3) RETURNING post_id`

// 	var postID int64

// 	err := db.QueryRow(sqlQuery, post.Title, post.Body, post.CreatedAt).Scan(&postID)
// 	if err != nil {
// 		log.Fatal("failed to execute post insert query", err)
// 	}

// 	fmt.Println("insert single record, post_id", postID)
// 	return postID
// }

// func GetPost(id int64) (models.Post, error) {
// 	db := connectDB()
// 	defer db.Close()

// 	var post models.Post

// 	sqlQuery := `SELECT * FROM posts WHERE post_id=$1`

// 	row := db.QueryRow(sqlQuery, id)

// 	// unmarshal row object into post struct
// 	err := row.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt)

// 	switch err {
// 	case sql.ErrNoRows:
// 		fmt.Println("No rows returned from DB")
// 		return post, nil
// 	case nil:
// 		return post, nil
// 	default:
// 		log.Fatal("Unable to scan rows in DB")
// 	}
// 	return post, nil
// }

// func GetAllPosts() ([]models.Post, error) {
// 	db := connectDB()
// 	defer db.Close()

// 	var posts []models.Post

// 	sqlQuery := `SELECT * FROM posts`

// 	rows, err := db.Query(sqlQuery)
// 	if err != nil {
// 		log.Fatalf("Unable to execute the query. %v", err)
// 	}

// 	defer rows.Close()

// 	// iterate over the rows
// 	for rows.Next() {
// 		var post models.Post

// 		// unmarshal the row object to post
// 		err = rows.Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt)

// 		if err != nil {
// 			log.Fatalf("Unable to scan the row. %v", err)
// 		}

// 		// append the post in the posts slice
// 		posts = append(posts, post)
// 	}

// 	return posts, nil
// }

// func UpdatePost(id int64, post models.Post) int64 {
// 	db := connectDB()
// 	defer db.Close()

// 	sqlQuery := `UPDATE posts SET title=$2, body=$3 WHERE post_id=$1`

// 	res, err := db.Exec(sqlQuery, id, post.Title, post.Body)
// 	if err != nil {
// 		log.Fatalf("Unable to execute the query. %v", err)
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		log.Fatalf("Error while checking the affected rows. %v", err)
// 	}

// 	fmt.Printf("Total rows/record affected %v", rowsAffected)

// 	return rowsAffected
// }

// func DeletePost(id int64) int64 {
// 	db := connectDB()
// 	defer db.Close()

// 	sqlQuery := `DELETE FROM posts WHERE post_id=$1`

// 	res, err := db.Exec(sqlQuery, id)
// 	if err != nil {
// 		log.Fatalf("Unable to execute the query. %v", err)
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		log.Fatalf("Error while checking the affected rows. %v", err)
// 	}

// 	fmt.Printf("Total rows/record affected %v", rowsAffected)

// 	return rowsAffected
// }
