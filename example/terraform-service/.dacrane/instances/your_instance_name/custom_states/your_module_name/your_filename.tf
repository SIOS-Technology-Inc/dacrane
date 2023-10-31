provider "bar" {
  user     = "Alice"
  password = "abc123"
}
resource "baz" "quz" {
  b = 2
  a = 1
}
