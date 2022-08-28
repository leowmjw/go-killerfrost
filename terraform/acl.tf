locals {
  acl = [
    // Looks like bug in module; data always out of sync // { name = "foo", inherit = true, login = true, membership = ["s2admin"], validity = "2022-12-31 00:00:00+00", connection_limit = -1, createrole = false },
    { name = "foo", inherit = true, login = true, membership = ["s2admin"], validity = "infinity", connection_limit = 10, createrole = true },
    { name = "bar", inherit = true, login = true, membership = ["s2read"], validity = "infinity", connection_limit = -1, createrole = false },
    { name = "backend", inherit = true, login = true, membership = [], validity = "infinity", connection_limit = -1, createrole = false },
  ]
}

// Use the Users module here ..
#########################################
# Create the users inside the database
#########################################
module "create_users" {

  source  = "jparnaudeau/database-admin/postgresql//create-users"
  version = "2.0.3"

  # set the provider
  providers = {
    postgresql = postgresql.pgmgm
  }

  # targetted rds
  pgadmin_user = "leow"
  dbhost       = "localhost"
  dbport       = 5432

  # input parameters for creating users inside database
  db_users = local.acl

  # set passwords
  passwords = { for user in local.acl : user.name => "password" }

}

