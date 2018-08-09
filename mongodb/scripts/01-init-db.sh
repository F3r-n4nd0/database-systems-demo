# !/bin/bash

if ! [[ -a /data/db/mydb-initialized ]]; then
	mongod --shutdown \
	&& mongod --fork --dbpath /data/db --logpath /var/log/mongodb/mongod.log \
	&& mongo <<-EOF	
		use users;		
		db.createUser({ 
			user: "web_service_user",
			pwd: "web_service_pass", 
			roles: [ "readWrite" ] 
		});
	EOF
	
	mongod --shutdown \
		&& mongod --auth --fork --dbpath /data/db --logpath /var/log/mongodb/mongod.log --replSet rs0 \
		&& mongo <<-EOF
		use admin;
		db.auth("$MONGO_INITDB_ROOT_USERNAME", "$MONGO_INITDB_ROOT_PASSWORD");
		rs.initiate({
			_id: "rs0",
			members: [
				{ _id: 0, host: "localhost:27017" }
			]
		});
	EOF
	
	touch /data/db/mydb-initialized
fi