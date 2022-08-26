// Local only; so no remote state; in real deployment; add remote state here!!!
// TODO: Remote state here ..

// Providers config ..
#######################################
# Define Providers pgadm & pgmgm for postgresql
#######################################
provider "postgresql" {
  alias    = "pgadm"
  host     = "localhost"
  username = "leow"
  sslmode  = "disable"
}

provider "postgresql" {
  alias    = "pgmgm"
  host     = "localhost"
  username = "leow"
  database = "myterraform"
  sslmode  = "disable"
}