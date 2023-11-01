provider "bar" {
  user     = "Alice"
  password = "abc123"
}
resource "baz" "quz" {
  a = 1
  b = 2
}
