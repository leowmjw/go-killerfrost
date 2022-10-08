// Local only; so no remote state; in real deployment; add remote state here!!!
// TODO: Remote state here ..

// Providers config ..
#######################################
# Define Providers pgadm & pgmgm for postgresql
#######################################
provider "postgresql" {
  alias    = "pgadm"
  host     = "127.0.0.1"
  username = "leow"
  sslmode  = "disable"
}

provider "postgresql" {
  alias    = "pgmgm"
  host     = "127.0.0.1"
  username = "leow"
  database = "myterraform"
  sslmode  = "disable"
}
