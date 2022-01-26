#!/bin/sh

MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-""}
MYSQL_DATABASE="storage"
MYSQL_USER="admin"
MYSQL_PASS="admin"
# ANSI colors variables for shell scripts
    Red='\033[0;31m'
    Green='\033[0;32m'
    Yellow='\033[0;33m'
    Blue='\033[0;34m'
    Cyan='\033[0;36m'
    White='\033[0;37m'
    Grey='\033[0;37m'
    NC='\033[0m'

# Check if mysql is currently running
check_mysql_running() {
    if [ -z "$(ps -ef | grep mysql | grep -v grep)" ]; then
        echo -e "${Red}MySQL is not running.${NC}"
        start_mysql
    fi
}

# Start mysql
start_mysql() {
    echo -e "${Green}Starting mysql...${NC}"
    mysql.server start
}


# Check if mysql is installed
check_mysql_installed() {
    if [ -z "$(which mysql)" ]; then
        echo -e "${Red}MySQL is not installed.${NC}"
        install_mysql
    fi
}

# Install mysql using homebrew
install_mysql() {
   echo -e "${Green}Installing mysql...${NC}"
   brew install mysql
   check_mysql_running
}

# Create user
create_user() {
    echo -e "${Green}Creating user...${NC}"
    if [ -z "$1" ]; then
        mysql -u root -e "CREATE USER '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_PASS'"
	    mysql -u root -e "GRANT ALL PRIVILEGES ON * . * TO '$MYSQL_USER'@'localhost'"
	    mysql -u root -e "FLUSH PRIVILEGES"
    else
        mysql -u root -p$1 -e "CREATE USER '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_PASS'"
	    mysql -u root -p$1 -e "GRANT ALL PRIVILEGES ON * . * TO '$MYSQL_USER'@'localhost'"
	    mysql -u root -p$1 -e "FLUSH PRIVILEGES"
    fi
}

# Use the mysql root password to create the database
create_database() {
    echo -e "${Blue}Creating database ${MYSQL_DATABASE}${NC}"
    if [ -z "$1" ]; then
        mysql -u root -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE"
        mysql -u root $MYSQL_DATABASE < db.sql
    else
        mysql -u root -p$1 -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE"
        mysql -u root -p$1 $MYSQL_DATABASE < db.sql
    fi
    create_user $1
}

rebuild_database() {
    echo -e "${Blue}Rebuilding database ${MYSQL_DATABASE}...${NC}"
    if [  -z "$1" ]; then
        mysql -u root -e "DROP DATABASE IF EXISTS $MYSQL_DATABASE"
	    mysql -u root -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE"
	    mysql -u root $MYSQL_DATABASE < db.sql
	    mysql -u root -e "FLUSH PRIVILEGES"
    else
        mysql -u root -p$1 -e "DROP DATABASE IF EXISTS $MYSQL_DATABASE"
	    mysql -u root -p$1 -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE"
	    mysql -u root -p$1 $MYSQL_DATABASE < db.sql
	    mysql -u root -p$1 -e "FLUSH PRIVILEGES"
    fi
}

test_database_creation() {
    if [  -z "$1" ]; then
        if [ -z "$(mysql -u root -e "SHOW DATABASES LIKE '$MYSQL_DATABASE'" | grep $MYSQL_DATABASE)" ]; then
            echo -e "${Red}Database $MYSQL_DATABASE was not created, retrying${NC}"
            create_database $1
        fi
    else
        if [ -z "$(mysql -u root -p$1 -e "SHOW DATABASES LIKE '$MYSQL_DATABASE'" | grep $MYSQL_DATABASE)" ]; then
            echo -e "${Red}Database $MYSQL_DATABASE was not created, retrying${NC}"
            create_database $1
        fi
    fi
}

test_user_creation() {
    if [  -z "$1" ]; then
        if [ -z "$(mysql -u root -e "SELECT User FROM mysql.user WHERE User = '$MYSQL_USER'" | grep $MYSQL_USER)" ]; then
            echo -e "${Red}User $MYSQL_USER was not created, retrying${NC}"
            create_user $1
        fi
    else
        if [ -z "$(mysql -u root -p$1 -e "SELECT User FROM mysql.user WHERE User = '$MYSQL_USER'" | grep $MYSQL_USER)" ]; then
            echo -e "${Red}User $MYSQL_USER was not created, retrying${NC}"
            create_user $1
        fi
    fi
}

main() {
    check_mysql_installed
    check_mysql_running
    # if argument is "create" then create the database and if its "rebuild" then rebuild the database
    if [ "$1" = "create" ]; then
        create_database $2
        test_database_creation $2
        test_user_creation $2
        echo -e "${Green}Database created successfully.${NC}"
    elif [ "$1" = "rebuild" ]; then
        rebuild_database $2
        echo -e "${Green}Database rebuilt successfully.${NC}"
    else
        echo -e "${Red}Invalid argument.${NC}"
        echo -e "${Yellow}Usage: ./build_meli_database.sh [create|rebuild]${NC}"
    fi
}

main $1 $2