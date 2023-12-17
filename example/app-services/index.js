const express = require('express')
const mysql = require('mysql2');
const app = express()

const connection = mysql.createConnection({
  host     : process.env.MYSQL_HOST,
  user     : process.env.MYSQL_USER,
  password : process.env.MYSQL_PASSWORD,
  database : process.env.MYSQL_DATABASE,
})

// respond with "hello world" when a GET request is made to the homepage
app.get('/', (req, res) => {
  res.send('hello world')
})

app.get('/users', (req, res) => {
  connection.query('SELECT * FROM users', function (error, results, fields) {
    if (error) throw error
    res.send(results)
  })
})

app.post('/users', (req, res) => {
  connection.query('SELECT * FROM users', function (error, results, fields) {
    if (error) throw error
    res.send(results)
  })
})

const port = process.env.PORT || 3000
app.listen(port, () => {
    connection.connect()
    console.log(`Start server on ${port}!`);
})
