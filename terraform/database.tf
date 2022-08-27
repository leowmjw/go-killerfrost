// Local variables
locals {
  db1 = {
    # parameters used for creating a database named 'myterraform' and for creating objects in the s1 schema
    db_schema_name = "s1"
    db_name        = "myterraform"
    db_admin       = "app_admin_role" # owner of the database
    extensions     = []

    # ---------------------------------- ROLES ------------------------------------------------------------------------------------
    # In this example, we create 3 roles
    # - "app_admin_role" will be the role used for creation, deletion, grant operations on objects, especially for tables.
    # - "app_write_role" for write operations. If you have a backend that insert lines into tables, it will used a user that inherits permissions from it.
    # - "app_readonly_role" for readonly operations.
    # Note : "write" role does not have the permissions to create table.
    # Note : the 'createrole' field is a boolean that provides a way to create other roles and put grants on it. Be carefull when you give this permission (privilege escalation).
    db_roles = [
      { id = "admin", role = "app_admin_role", inherit = true, login = false, validity = "infinity", privileges = ["USAGE", "CREATE"], createrole = true },
      { id = "readonly", role = "app_readonly_role", inherit = true, login = false, validity = "infinity", privileges = ["USAGE"], createrole = false },
      { id = "write", role = "app_write_role", inherit = true, login = false, validity = "infinity", privileges = ["USAGE"], createrole = false },
    ],

    # ---------------------------------- GRANT PERMISSIONS ON ROLES ------------------------------------------------------------------------------------
    # Notes :
    # the concept of "Least privilege" need to be applied here.
    # in the structure of a grant, there is the "role" and "owner_role"
    # "role" corresponds to the role on which the grants will be applied
    # "owner_role" is the role used to create grants on "role".
    # you could find the available privileges on official postgresql doc : https://www.postgresql.org/docs/13/ddl-priv.html
    # Note object_type = "type" is used only for default privileges
    db_grants = [
      # role app_admin_role : define grants to apply on db 'mydatabase', schema 'public'
      { object_type = "database", privileges = ["CREATE", "CONNECT", "TEMPORARY"], objects = [], role = "app_admin_role", owner_role = "postgres", grant_option = true },
      { object_type = "type", privileges = ["USAGE"], objects = [], role = "app_admin_role", owner_role = "postgres", grant_option = true },

      # role app_readonly_role : define grant to apply on db 'mydatabase', schema 'public'
      { object_type = "database", privileges = ["CONNECT"], objects = [], role = "app_readonly_role", owner_role = "app_admin_role", grant_option = false },
      { object_type = "type", privileges = ["USAGE"], objects = [], role = "app_readonly_role", owner_role = "app_admin_role", grant_option = true },
      { object_type = "table", privileges = ["SELECT", "REFERENCES", "TRIGGER"], objects = [], role = "app_readonly_role", owner_role = "app_admin_role", grant_option = false },
      { object_type = "sequence", privileges = ["SELECT", "USAGE"], objects = [], role = "app_readonly_role", owner_role = "app_admin_role", grant_option = false },

      # role app_write_role : define grant to apply on db 'mydatabase', schema 'public'
      { object_type = "database", privileges = ["CONNECT"], objects = [], role = "app_write_role", owner_role = "app_admin_role", grant_option = false },
      { object_type = "type", privileges = ["USAGE"], objects = [], role = "app_write_role", owner_role = "app_admin_role", grant_option = true },
      { object_type = "table", privileges = ["SELECT", "REFERENCES", "TRIGGER", "INSERT", "UPDATE", "DELETE"], objects = [], role = "app_write_role", owner_role = "app_admin_role", grant_option = false },
      { object_type = "sequence", privileges = ["SELECT", "USAGE"], objects = [], role = "app_write_role", owner_role = "app_admin_role", grant_option = false },
      { object_type = "function", privileges = ["EXECUTE"], objects = [], role = "app_write_role", owner_role = "app_admin_role", grant_option = false },
    ],
  }

  db2 = {
    # parameters used for creating a database named 'myterraform' and for creating objects in the s1 schema
    db_schema_name = "s2"
    db_name        = "myterraform"
    db_admin       = "s2admin" # owner of the database
    extensions     = []

    # ---------------------------------- ROLES ------------------------------------------------------------------------------------
    db_roles = [
      { id = "admin", role = "s2admin", inherit = true, login = false, validity = "infinity", privileges = ["USAGE", "CREATE"], createrole = true },
      { id = "readonly", role = "s2read", inherit = true, login = false, validity = "infinity", privileges = ["USAGE"], createrole = false },
      { id = "write", role = "s2write", inherit = true, login = false, validity = "infinity", privileges = ["USAGE"], createrole = false },
    ],

    # ---------------------------------- GRANT PERMISSIONS ON ROLES ------------------------------------------------------------------------------------
    db_grants = [
      # role s2read : define grant to apply on db 'myterraform', schema 's2'
      { object_type = "database", privileges = ["CONNECT"], objects = [], role = "s2read", owner_role = "s2admin", grant_option = false },
      { object_type = "type", privileges = ["USAGE"], objects = [], role = "s2read", owner_role = "s2admin", grant_option = true },
      { object_type = "table", privileges = ["SELECT", "REFERENCES", "TRIGGER"], objects = [], role = "s2read", owner_role = "s2admin", grant_option = false },
      { object_type = "sequence", privileges = ["SELECT", "USAGE"], objects = [], role = "s2read", owner_role = "s2admin", grant_option = false },

      # role s2write : define grant to apply on db 'myterraform', schema 's2'
      { object_type = "database", privileges = ["CONNECT"], objects = [], role = "s2write", owner_role = "s2admin", grant_option = false },
      { object_type = "type", privileges = ["USAGE"], objects = [], role = "s2write", owner_role = "s2admin", grant_option = true },
      { object_type = "table", privileges = ["SELECT", "REFERENCES", "TRIGGER", "INSERT", "UPDATE", "DELETE"], objects = [], role = "s2write", owner_role = "s2admin", grant_option = false },
      { object_type = "sequence", privileges = ["SELECT", "USAGE"], objects = [], role = "s2write", owner_role = "s2admin", grant_option = false },
      { object_type = "function", privileges = ["EXECUTE"], objects = [], role = "s2write", owner_role = "s2admin", grant_option = false },
    ],
  }
}

// Use Database module here ..
// New DB: myterraform
// with schema: s1
module "myterraform_s1" {

  source  = "jparnaudeau/database-admin/postgresql//create-database"
  version = "2.0.3"

  # set the provider
  providers = {
    postgresql = postgresql.pgadm
  }

  # targetted rds
  pgadmin_user = "leow"
  dbhost       = "localhost"
  dbport       = 5432

  # input parameters for creating database & objects inside database
  create_database = true
  inputs          = local.db1
}
// with schema: s2
module "myterraform_s2" {

  source  = "jparnaudeau/database-admin/postgresql//create-database"
  version = "2.0.3"

  # set the provider
  providers = {
    postgresql = postgresql.pgadm
  }

  # targetted rds
  pgadmin_user = "leow"
  dbhost       = "localhost"
  dbport       = 5432

  # input parameters for creating database & objects inside database
  create_database = false
  inputs          = local.db2
}

