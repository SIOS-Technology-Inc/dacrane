const express = require('express')
const app = express()

// respond with "hello world" when a GET request is made to the homepage
app.get('/', (req, res) => {
  res.send('hello world')
})

const port = process.env.PORT || 3000
app.listen(port, () => {
    console.log(`Start server on ${port}!`);
})
